package medias

const (
	UA_Dalvik  = "Dalvik/2.1.0 (Linux; U; Android 9; ALP-AL00 Build/HUAWEIALP-AL00)"
	UA_Browser = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36"
)

var (
	mediasGlobal = map[string]func(*Media) *CheckResult{
		"Dazn":        CheckDazn,
		"Netflix":     CheckNetflix,
		"TVBAnywhere": CheckTVBAnywhere,
		"Viu.com":     CheckViuCom,
		// "NetflixCDN": CheckNetflixCDN,
		/*
			"DisneyPlus": CheckDisneyPlus,
			"HotStar":    CheckHotStar,
			"YouTubePremium": CheckYouTubePremium,
			"YouTubeCDN": CheckYouTubeCDN,
			"PrimeVideo": PrimeVideo,
			"Tiktok": CheckTiktok,
			"iQYI": CheckiQYI,
		*/
	}

	mediasJP = map[string]func(*Media) *CheckResult{
		"PCRJP":        CheckPCRJP,
		"UMAJP":        CheckUMAJP,
		"Kancolle":     CheckKancolle,
		"KonosubaFD":   CheckKonosubaFD,
		"ProjectSekai": CheckProjectSekai,
		"AbemaTV":      CheckAbemaTV,
		"HBOGoAsia":    CheckHBOGoAsia,
		"DMM":          CheckDMM,
		"Niconico":     CheckNiconico,
		"Paravi":       CheckParavi,
		"HuluJP":       CheckHuluJP,
		"KaraokeDAM":   CheckKaraokeDAM,
		"FOD":          CheckFOD,
		"Radiko":       CheckRadiko,
		/*"Unext":        CheckUnext,
		"TVer":   CheckTVer,
		"WOWOW":  CheckWOWOW,*/
	}

	mediasHK = map[string]func(*Media) *CheckResult{
		"BilibiliHKMCTW": CheckBilibiliHKMCTW,
		"MyTVSuper":      CheckMyTVSuper,
		"ViuTV":          CheckViuTV,
		"NowE":           CheckNowE,
		"HBOGoAsia":      CheckHBOGoAsia,
	}

	mediasTW = map[string]func(*Media) *CheckResult{
		"BahamutAnime":   CheckBahamutAnime,
		"BilibiliHKMCTW": CheckBilibiliHKMCTW,
		"BilibiliTW":     CheckBilibiliTW,
		"HBOGoAsia":      CheckHBOGoAsia,
		"KKTV":           CheckKKTV,
		"LiTV":           CheckLiTV,
		"4GTV":           Check4GTV,
		"LineTV":         CheckLineTV,
		"HamiVideo":      CheckHamiVideo,
		"Catchplay":      CheckCatchplay,
		"ElevenSports":   CheckElevenSports,
	}

	MediaFuncs = map[string]map[string]func(*Media) *CheckResult{
		"Global": mediasGlobal,
		"JP":     mediasJP,
		"TW":     mediasTW,
		"HK":     mediasHK,
	}

	HumanReadableNames = map[string]string{
		// Global
		"Dazn":        "Dazn",
		"Netflix":     "Netflix",
		"NetflixCDN":  "Netflix Preferred CDN",
		"TVBAnywhere": "TVBAnywhere",
		"Viu.com":     "Viu.com",
		// JP
		"PCRJP":        "Princess Connect Re:Dive Japan",
		"UMAJP":        "Pretty Derby Japan",
		"KonosubaFD":   "Konosuba Fantastic Days",
		"ProjectSekai": "Project Sekai",
		"Kancolle":     "Kancolle Japan",
		"AbemaTV":      "Abema.TV",
		"DMM":          "DMM",
		"Niconico":     "niconico",
		"Paravi":       "Paravi",
		"HuluJP":       "Hulu Japan",
		"KaraokeDAM":   "Karaoke@DAM",
		"FOD":          "FOD (Fuji TV)",
		"Radiko":       "Radiko",
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
		"ElevenSports": "Eleven Sports",
	}
)
