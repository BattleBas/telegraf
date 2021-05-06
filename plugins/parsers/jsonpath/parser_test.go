package jsonpath

import (
	"fmt"
	"testing"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/require"
)

// *** Notes:***
// import cycle is caused trying to put influx line protocol expected output in a separate file, need influx line parser
// How to print telegraf.Metric to stdout? currently only get: file map[] map[name:John] 3600000000000
// Integration tests idea, completely separate test file that uses docker to run telegraf with a toml config
// Use: https://github.com/testcontainers/testcontainers-go
// Trying to load TOML config in unit tests a bit too complicated trying to get parser data

var DefaultTime = func() time.Time {
	return time.Unix(3600, 0)
}

const stringTypesJSON = `
{
    "explicitstringtype": "Bilbo",
    "defaultstringtype": "Baggins",
	"convertbooltostring": true,
	"convertinttostring": 1,
	"convertfloattostring": 1.1
}
`
const intTypesJSON = `
{
    "explicitinttype": 1,
    "defaultinttype": 2,
    "convertfloatoint": 3.1,
	"convertstringtoint": "4",
	"convertbooltoint": false
}
`

const floatTypesJSON = `
{
    "explicitfloattype": 1.1,
    "defaultfloattype": 2.1,
    "convertintotfloat": 3,
	"convertstringtofloat": "4.1"
}
`

const boolTypesJSON = `
{
    "explicitbooltype": true,
    "defaultbooltype": false,
	"convertinttobool": 1,
	"convertstringtobool": "false",
	"convertintstringtobool": "1"
}
`

func TestParseLineTypes(t *testing.T) {
	var tests = []struct {
		name           string
		JSONInput      string
		influxDataPath string
		configs        []Config
		expected       telegraf.Metric
	}{
		{
			name:      "Parse String types from JSON",
			JSONInput: stringTypesJSON,
			configs: []Config{
				{
					MetricName: "file",
					Fields: []FieldKeys{
						{
							Name:  "explicitstringtypeName",
							Query: "explicitstringtype",
							Type:  "string",
						},
						{
							Name:  "defaultstringtypeName",
							Query: "defaultstringtype",
						},
						{
							Name:  "convertbooltostringName",
							Query: "convertbooltostring",
							Type:  "string",
						},
						{
							Name:  "convertinttostringName",
							Query: "convertinttostring",
							Type:  "string",
						},
						{
							Name:  "convertfloattostringName",
							Query: "convertfloattostring",
							Type:  "string",
						},
					},
				},
			},
			expected: testutil.MustMetric(
				"file",
				map[string]string{},
				map[string]interface{}{
					"explicitstringtypeName":   "Bilbo",
					"defaultstringtypeName":    "Baggins",
					"convertbooltostringName":  "true",
					"convertinttostringName":   "1",
					"convertfloattostringName": "1.1",
				},
				DefaultTime(),
			),
		},
		{
			name:      "Parse int types from JSON",
			JSONInput: intTypesJSON,
			configs: []Config{
				{
					MetricName: "file",
					Fields: []FieldKeys{
						{
							Name:  "explicitinttypeName",
							Query: "explicitinttype",
							Type:  "int",
						},
						{
							Name:  "defaultinttypeName",
							Query: "defaultinttype",
						},
						{
							Name:  "convertfloatointName",
							Query: "convertfloatoint",
							Type:  "int",
						},
						{
							Name:  "convertstringtointName",
							Query: "convertstringtoint",
							Type:  "int",
						},
						{
							Name:  "convertbooltointName",
							Query: "convertbooltoint",
							Type:  "int",
						},
					},
				},
			},
			expected: testutil.MustMetric(
				"file",
				map[string]string{},
				map[string]interface{}{
					"explicitinttypeName":    1,
					"defaultinttypeName":     2,
					"convertfloatointName":   3,
					"convertstringtointName": 4,
					"convertbooltointName":   0,
				},
				DefaultTime(),
			),
		},
		{
			name:      "Parse float types from JSON",
			JSONInput: floatTypesJSON,
			configs: []Config{
				{
					MetricName: "file",
					Fields: []FieldKeys{
						{
							Name:  "explicitfloattypeName",
							Query: "explicitfloattype",
							Type:  "float",
						},
						{
							Name:  "defaultfloattypeName",
							Query: "defaultfloattype",
						},
						{
							Name:  "convertintotfloatName",
							Query: "convertintotfloat",
							Type:  "float",
						},
						{
							Name:  "convertstringtofloatName",
							Query: "convertstringtofloat",
							Type:  "float",
						},
					},
				},
			},
			expected: testutil.MustMetric(
				"file",
				map[string]string{},
				map[string]interface{}{
					"explicitfloattypeName":    1.1,
					"defaultfloattypeName":     2.1,
					"convertintotfloatName":    float64(3),
					"convertstringtofloatName": 4.1,
				},
				DefaultTime(),
			),
		},
		{
			name:      "Parse bool types from JSON",
			JSONInput: boolTypesJSON,
			configs: []Config{
				{
					MetricName: "file",
					Fields: []FieldKeys{
						{
							Name:  "explicitbooltypeName",
							Query: "explicitbooltype",
							Type:  "bool",
						},
						{
							Name:  "defaultbooltypeName",
							Query: "defaultbooltype",
						},
						{
							Name:  "convertinttoboolName",
							Query: "convertinttobool",
							Type:  "bool",
						},
						{
							Name:  "convertstringtoboolName",
							Query: "convertstringtobool",
							Type:  "bool",
						},
						{
							Name:  "convertintstringtoboolName",
							Query: "convertintstringtobool",
							Type:  "bool",
						},
					},
				},
			},
			expected: testutil.MustMetric(
				"file",
				map[string]string{},
				map[string]interface{}{
					"explicitbooltypeName":       true,
					"defaultbooltypeName":        false,
					"convertinttoboolName":       true,
					"convertstringtoboolName":    false,
					"convertintstringtoboolName": true,
				},
				DefaultTime(),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &Parser{
				Configs:  tt.configs,
				Log:      testutil.Logger{Name: "parsers.jsonpath"},
				TimeFunc: DefaultTime,
			}

			actual, err := parser.ParseLine(tt.JSONInput)
			require.NoError(t, err)

			fmt.Println(actual)

			testutil.RequireMetricEqual(t, tt.expected, actual)
		})
	}
}

func TestParse(t *testing.T) {
	var tests = []struct {
		name           string
		JSONInput      string
		influxDataPath string
		configs        []Config
		expected       telegraf.Metric
	}{
		{
			name:      "Parse Multiple JSON types",
			JSONInput: stringTypesJSON,
			configs: []Config{
				{
					MetricName: "file",
					Fields: []FieldKeys{
						{
							Name:  "explicitstringtype",
							Query: "explicitstringtype",
						},
					},
				},
			},
			expected: testutil.MustMetric(
				"file",
				map[string]string{},
				map[string]interface{}{
					"explicitstringtype": "Bilbo",
				},
				DefaultTime(),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &Parser{
				Configs:  tt.configs,
				Log:      testutil.Logger{Name: "parsers.jsonpath"},
				TimeFunc: DefaultTime,
			}

			actual, err := parser.Parse([]byte(tt.JSONInput))
			require.NoError(t, err)

			for _, m := range actual {
				testutil.RequireMetricEqual(t, tt.expected, m)
			}
		})
	}
}

func TestInvalidJSON(t *testing.T) {
	invalidJSON := `
	{
		"name": "John",
	}
	`
	parser := &Parser{
		Configs:  []Config{},
		Log:      testutil.Logger{Name: "parsers.jsonpath"},
		TimeFunc: DefaultTime,
	}
	_, err := parser.ParseLine(invalidJSON)
	require.Error(t, err)
	_, err = parser.Parse([]byte(invalidJSON))
	require.Error(t, err)
}