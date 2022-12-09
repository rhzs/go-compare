package compare

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsEquivalent(t *testing.T) {
	scenarios := []struct {
		desc        string
		inExpected  string
		inActual    string
		outExpected string
		outActual   string
		outSame     bool
		expectErr   bool
	}{
		{
			desc:        "happy path - empty",
			inExpected:  "{}",
			inActual:    "{}",
			outExpected: "",
			outActual:   "",
			outSame:     true,
			expectErr:   false,
		},
		{
			desc: "happy path - same but different order",
			inExpected: `
{
	"A": "1",
	"B": "2"
}`,
			inActual: `
{
	"B": "2",
	"A": "1"
}`,
			outExpected: "",
			outActual:   "",
			outSame:     true,
			expectErr:   false,
		},
		{
			desc: "sad path - item missing",
			inExpected: `
{
	"A": "1",
	"B": "2",
	"C": "3"
}`,
			inActual: `
{
	"B": "2",
	"A": "1"
}`,
			outExpected: "{\n\t\"A\": \"1\",\n\t\"B\": \"2\",\n\t\"C\": \"3\"\n}",
			outActual:   "{\n\t\"A\": \"1\",\n\t\"B\": \"2\"\n}",
			outSame:     false,
			expectErr:   false,
		},
		{
			desc: "sad path - sub item missing field",
			inExpected: `
{
	"A": "1",
	"B": "2",
	"C": {
		"FU": 111,
		"BAR": 222
	}
}`,
			inActual: `
{
	"A": "1",
	"B": "2",
	"C": {
		"FU": 111
	}
}`,
			outExpected: "{\n\t\"BAR\": 222,\n\t\"FU\": 111\n}",
			outActual:   "{\n\t\"FU\": 111\n}",
			outSame:     false,
			expectErr:   false,
		},
		{
			desc: "sad path - something invalid - actual",
			inExpected: `
{
	"A": "1",
	"B": "2",
	"C": {
		"FU": 111,
		"BAR": 222
	}
}`,
			inActual:    `something invalid`,
			outExpected: "",
			outActual:   "",
			outSame:     false,
			expectErr:   true,
		},
		{
			desc:       "sad path - something invalid - expected",
			inExpected: `something invalid`,
			inActual: `
{
	"A": "1",
	"B": "2",
	"C": {
		"FU": 111,
		"BAR": 222
	}
}`,
			outExpected: "",
			outActual:   "",
			outSame:     false,
			expectErr:   true,
		},
	}

	for _, s := range scenarios {
		scenario := s
		t.Run(scenario.desc, func(t *testing.T) {
			// inputs

			// call object under test
			resultExpected, resultActual, resultSame, resultErr := IsEquivalent([]byte(scenario.inExpected), []byte(scenario.inActual))

			// validation
			require.Equal(t, scenario.expectErr, resultErr != nil, "expected error: %t, err: '%s'", scenario.expectErr, resultErr)
			assert.Equal(t, scenario.outExpected, resultExpected)
			assert.Equal(t, scenario.outActual, resultActual)
			assert.Equal(t, scenario.outSame, resultSame)
		})
	}
}

func TestIsEquivalentArray(t *testing.T) {
	type args struct {
		expected []byte
		actual   []byte
	}
	tests := []struct {
		name     string
		args     args
		wantSame bool
		wantErr  bool
	}{
		{
			name: "Happy Path - Empty comparison",
			args: args{
				expected: []byte(""),
				actual:   []byte(nil),
			},
			wantSame: true,
			wantErr:  false,
		},
		{
			name: "Happy Path - single match comparison",
			args: args{
				expected: []byte(`[{"a":1,"b":"test"}]`),
				actual:   []byte(`[{"b":"test","a":1}]`),
			},
			wantSame: true,
			wantErr:  false,
		},
		{
			name: "Happy Path - array order matters !!",
			args: args{
				expected: []byte(`[{"a":1,"b":"test"}, {"c":2, "d":"bla"}]`),
				actual:   []byte(`[{"c":2, "d":"bla"}, {"b":"test","a":1}]`),
			},
			wantSame: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, gotSame, err := IsEquivalentArray(tt.args.expected, tt.args.actual)

			require.Equal(t, tt.wantErr, err != nil, "expected error: %t, err: '%s'", tt.wantErr, err)

			assert.Equalf(t, tt.wantSame, gotSame, "IsEquivalentArray(%v, %v)", tt.args.expected, tt.args.actual)
		})
	}
}

func TestIsEquivalentTest(t *testing.T) {
	scenarios := []struct {
		desc        string
		inExpected  string
		inActual    string
		outExpected string
		outActual   string
		outSame     bool
		expectErr   bool
	}{
		{
			desc:        "happy path - empty",
			inExpected:  "{}",
			inActual:    "{}",
			outExpected: "",
			outActual:   "",
			outSame:     true,
			expectErr:   false,
		},
		{
			desc: "happy path - same but different order",
			inExpected: `
{
	"A": "1",
	"B": "2"
}`,
			inActual: `
{
	"B": "2",
	"A": "1"
}`,
			outExpected: "",
			outActual:   "",
			outSame:     true,
			expectErr:   false,
		},
	}

	for _, s := range scenarios {
		scenario := s
		t.Run(scenario.desc, func(t *testing.T) {
			// inputs

			// call object under test
			resultBool := IsEquivalentTest(t, []byte(scenario.inExpected), []byte(scenario.inActual))

			// validation
			assert.True(t, resultBool)
		})
	}
}

func TestIsEquivalentArrayTest(t *testing.T) {
	scenarios := []struct {
		desc        string
		inExpected  string
		inActual    string
		outExpected string
		outActual   string
		outSame     bool
		expectErr   bool
	}{
		{
			desc:        "happy path - empty",
			inExpected:  "[]",
			inActual:    "[]",
			outExpected: "",
			outActual:   "",
			outSame:     true,
			expectErr:   false,
		},
		{
			desc: "happy path - same but different order",
			inExpected: `
[{
	"A": "1",
	"B": "2"
},
{
	"C": "3",
	"D": "4"
}
]`,
			inActual: `
[{
	"B": "2",
	"A": "1"
},
{
	"D": "4",
	"C": "3"
}
]`,
			outExpected: "",
			outActual:   "",
			outSame:     true,
			expectErr:   false,
		},
	}

	for _, s := range scenarios {
		scenario := s
		t.Run(scenario.desc, func(t *testing.T) {
			// inputs

			// call object under test
			resultBool := IsEquivalentArrayTest(t, []byte(scenario.inExpected), []byte(scenario.inActual))

			// validation
			assert.True(t, resultBool)
		})
	}
}

func TestFormatJSON(t *testing.T) {
	in := `{
	"B": "2",
	"C": {
		"FU": 111
	},
	"A": "1"
}`

	expected := `{
	"A": "1",
	"B": "2",
	"C": {
		"FU": 111
	}
}`

	asMap, result, err := FormatJSON([]byte(in))
	require.NoError(t, err)
	require.NotNil(t, asMap)
	require.NotEmpty(t, result)

	assert.Equal(t, expected, string(result))
}

func TestFormatJSON_error(t *testing.T) {
	in := `something invalid`

	_, _, err := FormatJSON([]byte(in))
	require.Error(t, err)
}

func TestFormatArrayJSON(t *testing.T) {
	type args struct {
		src []byte
	}
	tests := []struct {
		name      string
		args      args
		wantAsMap []map[string]interface{}
		expectErr bool
	}{
		{
			name: "Happy path, 1 array, single map value",
			args: args{
				src: []byte("[{ \"a\":true }]"),
			},
			wantAsMap: []map[string]interface{}{
				{
					"a": true,
				},
			},
			expectErr: false,
		},
		{
			name: "Happy path, 2 array, multiple map value",
			args: args{
				src: []byte("[{ \"a\":true , \"b\":\"bla\"},{ \"c\":false , \"d\":\"blabla\"}]"),
			},
			wantAsMap: []map[string]interface{}{
				{
					"b": "bla",
					"a": true,
				},
				{
					"d": "blabla",
					"c": false,
				},
			},
			expectErr: false,
		},
		{
			name: "Happy path, empty source",
			args: args{
				src: []byte(""),
			},
			wantAsMap: []map[string]interface{}{},
			expectErr: false,
		},
		{
			name: "Sad path, invalid JSON",
			args: args{
				src: []byte(" b(-_-)d "),
			},
			wantAsMap: nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAsMap, _, err := FormatArrayJSON(tt.args.src)

			require.Equal(t, tt.expectErr, err != nil, "expected error: %t, err: '%s'", tt.expectErr, err)
			assert.Equalf(t, tt.wantAsMap, gotAsMap, "FormatArrayJSON(%v)", tt.args.src)
		})
	}
}

func TestFormatArrayJSON_error(t *testing.T) {
	in := `something invalid`

	_, _, err := FormatArrayJSON([]byte(in))
	require.Error(t, err)
}
