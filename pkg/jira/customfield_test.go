package jira

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitCustomArray(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string yields single empty value",
			input:    "",
			expected: []string{""},
		},
		{
			name:     "single value",
			input:    "one",
			expected: []string{"one"},
		},
		{
			name:     "multi value splits on comma and trims",
			input:    "a, b ,c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "escaped comma is kept literal in a single value",
			input:    `Matching (CAI\, other)`,
			expected: []string{"Matching (CAI, other)"},
		},
		{
			name:     "mix of escaped and unescaped commas",
			input:    `a\,b, c, d\,e`,
			expected: []string{"a,b", "c", "d,e"},
		},
		{
			name:     "trailing backslash is preserved",
			input:    `a\`,
			expected: []string{`a\`},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, splitCustomArray(tc.input))
		})
	}
}

func TestConstructCustomFields(t *testing.T) {
	arrayOptionField := IssueTypeField{Name: "BL Product", Key: "customfield_10010"}
	arrayOptionField.Schema.DataType = customFieldFormatArray
	arrayOptionField.Schema.Items = customFieldFormatOption

	arrayStringField := IssueTypeField{Name: "Labels Multi", Key: "customfield_10011"}
	arrayStringField.Schema.DataType = customFieldFormatArray

	singleOptionField := IssueTypeField{Name: "Single Option", Key: "customfield_10012"}
	singleOptionField.Schema.DataType = customFieldFormatOption

	configured := []IssueTypeField{arrayOptionField, arrayStringField, singleOptionField}

	cases := []struct {
		name     string
		fields   map[string]string
		expected customField
	}{
		{
			name:   "array of options, multi value",
			fields: map[string]string{"bl-product": "a,b"},
			expected: customField{
				"customfield_10010": []customFieldTypeOption{{Value: "a"}, {Value: "b"}},
			},
		},
		{
			name:   "array of options, escaped comma stays one option",
			fields: map[string]string{"bl-product": `Matching (CAI\, other)`},
			expected: customField{
				"customfield_10010": []customFieldTypeOption{{Value: "Matching (CAI, other)"}},
			},
		},
		{
			name:   "array of plain strings, escaped comma stays one value",
			fields: map[string]string{"labels-multi": `x\,y,z`},
			expected: customField{
				"customfield_10011": []string{"x,y", "z"},
			},
		},
		{
			name:   "single option field with a comma is unaffected by split",
			fields: map[string]string{"single-option": "Matching (CAI, other)"},
			expected: customField{
				"customfield_10012": customFieldTypeOption{Value: "Matching (CAI, other)"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data := createRequest{}
			constructCustomFields(tc.fields, configured, &data)
			assert.Equal(t, tc.expected, data.Fields.M.customFields)
		})
	}
}
