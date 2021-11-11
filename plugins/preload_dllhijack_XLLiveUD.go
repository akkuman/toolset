package plugins

// PreloadDllHijackXLLiveUD 位于 迅雷11 安装目录的 XLLiveUD.exe
type PreloadDllHijackXLLiveUD struct {
	SubPluginPreloadDllHijackX86
}

func (p *PreloadDllHijackXLLiveUD) GetMainProgramName() (string) {
	return "XLLiveUD"
}

func (p *PreloadDllHijackXLLiveUD) GetDllName() (string) {
	return "XLLiveUpdateAgent"
}

func (p *PreloadDllHijackXLLiveUD) GetPluginName() (string) {
	return "preload_dll_hijack_XLLiveUD"
}

func (p *PreloadDllHijackXLLiveUD) GetDllExports() []string {
	return []string{
		"??0TbcString@@QAE@ABV0@@Z",
		"??0TbcString@@QAE@ABV?$basic_string@_WU?$char_traits@_W@std@@V?$allocator@_W@2@@std@@@Z",
		"??0TbcString@@QAE@PB_W@Z",
		"??0TbcString@@QAE@XZ",
		"??0UpdateInfo@XLLiveUpdate@@QAE@$$QAU01@@Z",
		"??0UpdateInfo@XLLiveUpdate@@QAE@ABU01@@Z",
		"??0UpdateInfo@XLLiveUpdate@@QAE@XZ",
		"??0XLLiveUpdateAgent@XLLiveUpdate@@QAE@$$QAV01@@Z",
		"??0XLLiveUpdateAgent@XLLiveUpdate@@QAE@ABV01@@Z",
		"??0XLLiveUpdateAgent@XLLiveUpdate@@QAE@XZ",
		"??1TbcString@@QAE@XZ",
		"??1UpdateInfo@XLLiveUpdate@@QAE@XZ",
		"??4TbcString@@QAEAAV0@ABV0@@Z",
		"??4TbcString@@QAEAAV0@ABV?$basic_string@_WU?$char_traits@_W@std@@V?$allocator@_W@2@@std@@@Z",
		"??4TbcString@@QAEAAV0@PB_W@Z",
		"??4UpdateInfo@XLLiveUpdate@@QAEAAU01@$$QAU01@@Z",
		"??4UpdateInfo@XLLiveUpdate@@QAEAAU01@ABU01@@Z",
		"??4XLLiveUpdateAgent@XLLiveUpdate@@QAEAAV01@$$QAV01@@Z",
		"??4XLLiveUpdateAgent@XLLiveUpdate@@QAEAAV01@ABV01@@Z",
		"??8TbcString@@QBE_NABV0@@Z",
		"??8TbcString@@QBE_NABV?$basic_string@_WU?$char_traits@_W@std@@V?$allocator@_W@2@@std@@@Z",
		"??8TbcString@@QBE_NPB_W@Z",
		"??9TbcString@@QBE_NABV0@@Z",
		"??9TbcString@@QBE_NABV?$basic_string@_WU?$char_traits@_W@std@@V?$allocator@_W@2@@std@@@Z",
		"??9TbcString@@QBE_NPB_W@Z",
		"??YTbcString@@QAEAAV0@ABV0@@Z",
		"??YTbcString@@QAEAAV0@ABV?$basic_string@_WU?$char_traits@_W@std@@V?$allocator@_W@2@@std@@@Z",
		"??YTbcString@@QAEAAV0@PB_W@Z",
		"??_7XLLiveUpdateAgent@XLLiveUpdate@@6B@",
		"?Empty@TbcString@@QBE_NXZ",
		"?GetInstance@XLLiveUpdateAgent@XLLiveUpdate@@SAPAV12@XZ",
		"?Length@TbcString@@QBEHXZ",
		"?ToString@TbcString@@QBE?AV?$basic_string@DU?$char_traits@D@std@@V?$allocator@D@2@@std@@XZ",
		"?c_str@TbcString@@QBEPB_WXZ",
	}
}
