package plugins

// PreloadDllHijackPowerPoint2007 位于 迅雷11 安装目录的 XLLiveUD.exe
type PreloadDllHijackPowerPoint2007 struct {
	SubPluginPreloadDllHijackX86
}

func (p *PreloadDllHijackPowerPoint2007) GetMainProgramName() (string) {
	return "POWERPNT"
}

func (p *PreloadDllHijackPowerPoint2007) GetDllName() (string) {
	return "PPCORE"
}

func (p *PreloadDllHijackPowerPoint2007) GetPluginName() (string) {
	return "preload_dll_hijack_PowerPoint2007"
}

func (p *PreloadDllHijackPowerPoint2007) GetDllExports() []string {
	return []string{
		"?s_max@CoordRange32@Art@@2_JB",
		"?s_max@PosCoordRange32@Art@@2_JB",
		"?s_min@CoordRange32@Art@@2_JB",
		"?s_min@PosCoordRange32@Art@@2_JB",
		"DllGetLCID",
		"_PPMain@0",
		"_ShowSplashScreen@0",
	}
}
