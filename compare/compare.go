package compare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIsEquivalent is a convenience method for calling IsEquivalent
func IsEquivalentTest(t *testing.T, expected, actual []byte) bool {
	expectedFormatted, resultFormatted, resultSame, resultErr := IsEquivalent(expected, actual)
	if !assert.NoError(t, resultErr) {
		fmt.Printf("Actual:\n%s\n", string(actual))
		fmt.Print("\n-----------------------------------\n")
		fmt.Printf("Expected:\n%s\n", string(expected))

		return false
	}

	if !assert.True(t, resultSame) {
		assert.Equal(t, expectedFormatted, resultFormatted)
		return false
	}

	return true
}

// IsEquivalentArrayTest is a convenience method for calling IsEquivalentArray
func IsEquivalentArrayTest(t *testing.T, expected, actual []byte) bool {
	expectedFormatted, resultFormatted, resultSame, resultErr := IsEquivalentArray(expected, actual)
	require.NoError(t, resultErr)

	if !assert.True(t, resultSame) {
		assert.Equal(t, expectedFormatted, resultFormatted)
		return false
	}

	return true
}

// IsEquivalent compares if supplied []byte are equivalent from a logical JSON perspective.
// Returns strings with the first difference (use with assert.Equal for handy viewing), a boolean to indicate the success and an error.
//
//nolint:nakedret
func IsEquivalent(expected, actual []byte) (expectedVersion, actualVersion string, same bool, err error) {
	expectedMap, expectedFormatted, err := FormatJSON(expected)
	if err != nil {
		err = fmt.Errorf("error while formatting expected: %w", err)

		return
	}

	actualMap, actualFormatted, err := FormatJSON(actual)
	if err != nil {
		err = fmt.Errorf("error while formatting actual: %w", err)

		return
	}

	if bytes.Equal(expectedFormatted, actualFormatted) {
		same = true
		return
	}

	expectedVersion, actualVersion = compareMaps(expectedMap, actualMap)

	if len(expectedVersion)+len(actualVersion) > 0 {
		return
	}

	// default to "everything"
	expectedVersion = string(expectedFormatted)
	actualVersion = string(actualFormatted)

	return
}

// IsEquivalentArray compares if supplied []byte are equivalent from a logical JSON perspective when marshaled into Array.
// Returns strings with the first difference (use with assert.Equal for handy viewing), a boolean to indicate the success and an error.
func IsEquivalentArray(expected, actual []byte) (expectedVersion, actualVersion string, same bool, err error) {
	expectedMap, expectedFormatted, err := FormatArrayJSON(expected)
	if err != nil {
		err = fmt.Errorf("error while formatting array expected: %w", err)

		return
	}

	actualMap, actualFormatted, err := FormatArrayJSON(actual)
	if err != nil {
		err = fmt.Errorf("error while formatting array actual: %w", err)

		return
	}

	if bytes.Equal(expectedFormatted, actualFormatted) {
		same = true
		return
	}

	if len(expectedMap) != len(actualMap) {
		same = false
		return
	}

	for i, currentExpectedMap := range expectedMap {
		expectedVersion, actualVersion = compareMaps(currentExpectedMap, actualMap[i])
	}

	if len(expectedVersion)+len(actualVersion) > 0 {
		return
	}

	// default to "everything"
	expectedVersion = string(expectedFormatted)
	actualVersion = string(actualFormatted)

	return expectedVersion, actualVersion, same, err
}

func compareMaps(expectedMap, actualMap map[string]interface{}) (expectedVersion, actualVersion string) {
	for key, expectedValue := range expectedMap {
		actualValue := actualMap[key]

		switch expectedValue.(type) { //nolint:gocritic
		case map[string]interface{}:
			expectedConverted, _ := expectedValue.(map[string]interface{})

			actualConverted, ok := actualValue.(map[string]interface{})
			if !ok {
				expectedItemFormatted, _ := json.MarshalIndent(expectedMap, "", "\t")
				actualItemFormatted, _ := json.MarshalIndent(actualMap, "", "\t")

				return string(expectedItemFormatted), string(actualItemFormatted)
			}

			expectedVersion, actualVersion = compareMaps(expectedConverted, actualConverted)

			if len(expectedVersion)+len(actualVersion) > 0 {
				return
			}
		}

		expectedItemFormatted, _ := json.MarshalIndent(expectedValue, "", "\t")
		actualItemFormatted, _ := json.MarshalIndent(actualValue, "", "\t")

		if !bytes.Equal(expectedItemFormatted, actualItemFormatted) {
			expectedOut, _ := json.MarshalIndent(expectedMap, "", "\t")
			actualOut, _ := json.MarshalIndent(actualMap, "", "\t")

			expectedVersion = string(expectedOut)
			actualVersion = string(actualOut)

			return expectedVersion, actualVersion
		}
	}

	return expectedVersion, actualVersion
}

func FormatJSON(src []byte) (asMap map[string]interface{}, formatted []byte, err error) {
	asMap = map[string]interface{}{}

	err = json.Unmarshal(src, &asMap)
	if err != nil {
		return nil, nil, err
	}

	formatted, err = json.MarshalIndent(asMap, "", "\t")
	if err != nil {
		return nil, nil, err
	}

	return asMap, formatted, nil
}

func FormatArrayJSON(src []byte) (asMap []map[string]interface{}, formatted []byte, err error) {
	asMap = []map[string]interface{}{}

	if len(src) == 0 {
		return asMap, src, nil
	}

	err = json.Unmarshal(src, &asMap)
	if err != nil {
		return nil, nil, err
	}

	formatted, err = json.MarshalIndent(asMap, "", "\t")
	if err != nil {
		return nil, nil, err
	}

	return asMap, formatted, nil
}
