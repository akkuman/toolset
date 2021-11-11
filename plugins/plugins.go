package plugins

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"toolset/utils"
)

type BasePlugin struct {
	PluginName string
}

type RunnerIface interface {
	SetShellcdoe([]byte)
	Run() ([]byte, error)
}

func (p *BasePlugin) GetPluginDataPath() string {
	return filepath.Join(utils.GetExecutableDir(), "data", p.PluginName)
}

// getDlltoolPath 获取系统中 dlltool 的路径
func (p *BasePlugin) getDlltoolPath(x64 bool) string {
	name := "i686-w64-mingw32"
	if x64 {
		name = "x86_64-w64-mingw32"
	}
	return fmt.Sprintf("/usr/%s/bin/dlltool", name)
}

// defToExp 将 def 文件转为 exp 文件
func (p *BasePlugin) defToExp(x64 bool, defPath, expPath string) error {
	dlltoolPath := p.getDlltoolPath(x64)
	cmd := exec.Command(dlltoolPath, "--input-def", defPath, "--output-exp", expPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	err := cmd.Run()
	return err
}

// buildTmpWorkDir 新建一个临时目录作为工作目录，并把所有的所需的基础文件拷贝过去
func (p *BasePlugin) buildTmpWorkDir() (string, error) {
	tmpDir, err := ioutil.TempDir("", fmt.Sprintf("%s-*", p.PluginName))
	if err != nil {
		return "", err
	}
	err = utils.CopyDir(p.GetPluginDataPath(), tmpDir)
	if err != nil {
		return "", err
	}
	return tmpDir, nil
}
