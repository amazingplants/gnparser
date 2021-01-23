package parser

import (
	"sort"

	"github.com/gnames/gnparser/entity/parsed"
)

// ToOutput converts Abstract Syntax Tree of scientific name to a
// final output object.
func (sn *scientificNameNode) ToOutput(withDetails bool) parsed.Parsed {
	res := parsed.Parsed{
		Verbatim:      sn.verbatim,
		Canonical:     sn.Canonical(),
		Virus:         sn.virus,
		VerbatimID:    sn.verbatimID,
		ParserVersion: sn.parserVersion,
	}

	if res.Canonical == nil {
		return res
	}

	res.Parsed = true
	res.ParseQuality, res.QualityWarnings = qualityWarnings(sn.warnings)
	res.Normalized = sn.Normalized()
	res.Cardinality = sn.cardinality
	res.Authorship = sn.LastAuthorship(withDetails)
	res.Hybrid = sn.hybrid
	res.Surrogate = sn.surrogate
	res.Bacteria = sn.bacteria
	res.Tail = sn.tail
	if withDetails {
		res.Details = sn.Details()
		res.Words = sn.Words()
	}
	return res
}

func qualityWarnings(ws map[parsed.Warning]struct{}) (int, []parsed.QualityWarning) {
	warns := prepareWarnings(ws)
	quality := 1
	if len(warns) > 0 {
		quality = warns[0].Quality
	}
	return quality, warns
}

func prepareWarnings(ws map[parsed.Warning]struct{}) []parsed.QualityWarning {
	res := make([]parsed.QualityWarning, len(ws))
	var i int
	for k := range ws {
		res[i] = k.NewQualityWarning()
		i++
	}

	sort.Slice(res, func(i, j int) bool {
		if res[i].Quality > res[j].Quality {
			return true
		}
		if res[i].Quality < res[j].Quality {
			return false
		}
		return res[i].Warning < res[j].Warning
	})
	return res
}
