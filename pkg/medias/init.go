package medias

const (
	UA_Dalvik  = "Dalvik/2.1.0 (Linux; U; Android 9; ALP-AL00 Build/HUAWEIALP-AL00)"
	UA_Browser = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36"
)

var (
	mediasJP = map[string]func(*Media) *CheckResult{
		"PCRJP":     CheckPCRJP,
		"UMAJP":     CheckUMAJP,
		"AbemaTV":   CheckAbemaTV,
		"Kancolle":  CheckKancolle,
		"HBOGoAsia": CheckHBOGoAsia,
	}

	mediasHK = map[string]func(*Media) *CheckResult{
		"BilibiliHKMCTW": CheckBilibiliHKMCTW,
		"MyTVSuper":      CheckMyTVSuper,
		"ViuTV":          CheckViuTV,
		"NowE":           CheckNowE,
		"HBOGoAsia":      CheckHBOGoAsia,
	}

	mediasTW = map[string]func(*Media) *CheckResult{
		"BahamutAnime": CheckBahamutAnime,
		"BilibiliTW":   CheckBilibiliTW,
		"HBOGoAsia":    CheckHBOGoAsia,
		"KKTV":         CheckKKTV,
		"LiTV":         CheckLiTV,
		"4GTV":         Check4GTV,
		"LineTV":       CheckLineTV,
		"HamiVideo":    CheckHamiVideo,
		"Catchplay":    CheckCatchplay,
		//"ElevenSports": CheckElevenSports,
	}

	MediaFuncs = map[string]map[string]func(*Media) *CheckResult{
		"JP": mediasJP,
		"TW": mediasTW,
		"HK": mediasHK,
	}

	HumanReadableNames = map[string]string{
		// JP
		"PCRJP":    "Princess Connect Re:Dive Japan",
		"UMAJP":    "Pretty Derby Japan",
		"AbemaTV":  "Abema.TV",
		"Kancolle": "Kancolle Japan",
		// HK
		"BilibiliHKMCTW": "Bilibili HongKong/Macua/Taiwan",
		"ViuTV":          "Viu.TV",
		"NowE":           "Now E",
		"HBOGoAsia":      "HBO Go Asia",
		"MyTVSuper":      "MyTVSuper",
		// TW
		"BilibiliTW":   "Bilibili Taiwan Only",
		"BahamutAnime": "Bahamut Anime",
		"KKTV":         "KKTV",
		"LiTV":         "LiTV",
		"4GTV":         "4GTV",
		"LineTV":       "LineTV",
		"HamiVideo":    "HamiVideo",
		"Catchplay":    "Catchplay",
	}
)
