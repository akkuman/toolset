package plugins

import (
	"path/filepath"
	"toolset/utils"
)

type BasePlugin struct {
	PluginName string
}

func (p *BasePlugin) GetPluginDataPath() string {
	return filepath.Join(utils.GetExecutableDir(), "data", p.PluginName)
}
