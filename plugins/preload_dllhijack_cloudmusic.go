package plugins

// PreloadDllHijackCloudmusic 位于网易云音乐安装目录的 cloudmusic_reporter.exe
type PreloadDllHijackCloudmusic struct {
}

func (p *PreloadDllHijackCloudmusic) GetMainProgramName() (string) {
	return "inno_updater"
}

func (p *PreloadDllHijackCloudmusic) GetDllName() (string) {
	return "vcruntime140"
}

func (p *PreloadDllHijackCloudmusic) GetPluginName() (string) {
	return "dll_hijack_cloudmusic"
}

func (p *PreloadDllHijackCloudmusic) GetIsX64Arch() (bool) {
	return false
}

func (p *PreloadDllHijackCloudmusic) GetDllExports() []string {
	return []string{
		"_except_handler4_common",
		"memset",
		"memmove",
		"memcmp",
		"memcpy",
		"__std_type_info_destroy_list",
	}
}
