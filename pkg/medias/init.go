package medias

var (
	mediasJP = map[string]func(*Media) *CheckResult{
		"PCRJP":   CheckPCRJP,
		"UMAJP":   CheckUMAJP,
		"AbemaTV": CheckAbemaTV,
	}

	mediasTW = map[string]func(*Media) *CheckResult{
		"BahamutAnime": CheckBahamutAnime,
	}

	MediaFuncs = map[string]map[string]func(*Media) *CheckResult{
		"JP": mediasJP,
		"TW": mediasTW,
	}
)
