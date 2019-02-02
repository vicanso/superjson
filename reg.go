package superjson

import "regexp"

var (
	cutWordsReg = regexp.MustCompile(`[A-Z_ ]+`)
	omitReg     = regexp.MustCompile(`[_ ]`)
)
