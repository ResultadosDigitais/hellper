package commands

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestGetStringInt64(test *testing.T) {
	scenarios := []struct {
		name         string
		input        string
		output       int64
		expectError  bool
		errorMessage string
	}{
		{
			name:   "Convert string to int",
			input:  "0",
			output: 0,
		},
		{
			name:         "Returns error when try to convert",
			input:        "1.2",
			expectError:  true,
			errorMessage: "strconv.ParseInt: parsing \"1.2\": invalid syntax",
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				result, err := getStringInt64(scenario.input)
				if scenario.expectError {
					if err == nil {
						t.Fatal("an error was expected, but not occurred")
					}
					assert.EqualError(t, err, scenario.errorMessage)
				} else {
					if err != nil {
						t.Fatal(
							"an error occurred, but was not expected\n",
							"error: ",
							err,
						)
					}
					assert.Equal(t, scenario.output, result, "The return should be the same")
				}

			},
		)
	}
}

func TestGetSeverityLevelText(test *testing.T) {
	scenarios := []struct {
		name   string
		input  int64
		output string
	}{
		{
			name:   "Test case 0",
			input:  0,
			output: "SEV0 - All hands on deck",
		},
		{
			name:   "Test case 1",
			input:  1,
			output: "SEV1 - Critical impact to many users",
		},
		{
			name:   "Test case 2",
			input:  2,
			output: "SEV2 - Minor issue that impacts ability to use product",
		},
		{
			name:   "Test case 3",
			input:  3,
			output: "SEV3 - Minor issue not impacting ability to use product",
		},
		{
			name:   "Test case default",
			input:  123456,
			output: "",
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				result := getSeverityLevelText(scenario.input)
				require.Equal(t, scenario.output, result, "The return should be the same")
			},
		)
	}
}

func TestConvertTimestamp(test *testing.T) {
	scenarios := []struct {
		name         string
		input        string
		output       time.Time
		expectError  bool
		errorMessage string
	}{
		{
			name:         "When the string is empty",
			input:        "",
			expectError:  true,
			errorMessage: "Empty Timestamp",
		},
		{
			name:   "Convert string to timestamp",
			input:  "1512085950.000216",
			output: time.Unix(1512085950, 216),
		},
		{
			name:         "An error is returned when seconds is invalid",
			input:        "1512085950.abcde",
			expectError:  true,
			errorMessage: "strconv.ParseInt: parsing \"abcde\": invalid syntax",
		},
		{
			name:         "An error is returned when minutes is invalid",
			input:        "abcde.000216",
			expectError:  true,
			errorMessage: "strconv.ParseInt: parsing \"abcde\": invalid syntax",
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				result, err := convertTimestamp(scenario.input)
				if scenario.expectError {
					if err == nil {
						t.Fatal("an error was expected, but not occurred")
					}
					assert.EqualError(t, err, scenario.errorMessage)
				} else {
					if err != nil {
						t.Fatal(
							"an error occurred, but was not expected\n",
							"error: ",
							err,
						)
					}
					assert.Equal(t, scenario.output, result, "The return should be the same")
				}

			},
		)
	}
}
