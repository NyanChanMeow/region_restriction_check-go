package medias

const (
	UA_Dalvik  = "Dalvik/2.1.0 (Linux; U; Android 9; ALP-AL00 Build/HUAWEIALP-AL00)"
	UA_Browser = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36"
)

var (
	mediasJP = map[string]func(*Media) *CheckResult{
		"PCRJP":    CheckPCRJP,
		"UMAJP":    CheckUMAJP,
		"AbemaTV":  CheckAbemaTV,
		"Kancolle": CheckKancolle,
	}

	mediasHK = map[string]func(*Media) *CheckResult{
		"BilibiliHKMCTW": CheckBilibiliHKMCTW,
	}

	mediasTW = map[string]func(*Media) *CheckResult{
		"BahamutAnime": CheckBahamutAnime,
		"BilibiliTW":   CheckBilibiliTW,
	}

	MediaFuncs = map[string]map[string]func(*Media) *CheckResult{
		"JP": mediasJP,
		"TW": mediasTW,
		"HK": mediasHK,
	}

	HumanReadableNames = map[string]string{
		"PCRJP":    "Princess Connect Re:Dive Japan",
		"UMAJP":    "Pretty Derby Japan",
		"AbemaTV":  "Abema.TV",
		"Kancolle": "Kancolle Japan",
		// HK
		"BilibiliHKMCTW": "Bilibili HongKong/Macua/Taiwan",
		// TW
		"BilibiliTW":   "Bilibili Taiwan Only",
		"BahamutAnime": "Bahamut Anime",
	}
)
