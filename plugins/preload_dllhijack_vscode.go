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

func (p *PreloadDllHijackVscode) GetIsX64Arch() (bool) {
	return false
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
