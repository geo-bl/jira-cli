package jira

import "strings"

const (
	customFieldFormatOption  = "option"
	customFieldFormatArray   = "array"
	customFieldFormatNumber  = "number"
	customFieldFormatProject = "project"
)

type customField map[string]interface{}

type customFieldTypeNumber float64

type customFieldTypeNumberSet struct {
	Set customFieldTypeNumber `json:"set"`
}

type customFieldTypeStringSet struct {
	Set string `json:"set"`
}

type customFieldTypeOption struct {
	Value string `json:"value"`
}

type customFieldTypeOptionSet struct {
	Set customFieldTypeOption `json:"set"`
}

type customFieldTypeOptionAddRemove struct {
	Add    *customFieldTypeOption `json:"add,omitempty"`
	Remove *customFieldTypeOption `json:"remove,omitempty"`
}

type customFieldTypeProject struct {
	Value string `json:"key"`
}

type customFieldTypeProjectSet struct {
	Set customFieldTypeProject `json:"set"`
}

// splitCustomArray splits a multi-value custom-field input on unescaped commas.
// A comma escaped as `\,` is treated as a literal comma within a single value,
// which lets a single option value legitimately contain a comma
// (e.g. `Matching (CAI\, other)` -> one value `Matching (CAI, other)`).
// Each resulting value is trimmed of surrounding whitespace.
func splitCustomArray(val string) []string {
	var (
		out []string
		cur strings.Builder
	)
	for i := 0; i < len(val); i++ {
		if val[i] == '\\' && i+1 < len(val) && val[i+1] == ',' {
			cur.WriteByte(',')
			i++
			continue
		}
		if val[i] == ',' {
			out = append(out, strings.TrimSpace(cur.String()))
			cur.Reset()
			continue
		}
		cur.WriteByte(val[i])
	}
	out = append(out, strings.TrimSpace(cur.String()))

	return out
}
