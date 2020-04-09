package log

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testOut struct {
	name     string
	output   string
	expected Out
}

func TestOut(test *testing.T) {
	scenarios := []testOut{
		{
			name:     "Creates default STDOUT Out",
			output:   "",
			expected: STDOUT,
		},
		{
			name:     "Creates a stdout Out",
			output:   "stdout",
			expected: STDOUT,
		},
		{
			name:     "Creates a STDOUT Out",
			output:   "STDOUT",
			expected: STDOUT,
		},
		{
			name:     "Creates a file Out",
			output:   "pathtoafile",
			expected: Out("pathtoafile"),
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				var (
					out Out
					err = out.Set(scenario.output)
				)
				assert.Nil(t, err, "Out.Set error")
				assert.Exactly(t, scenario.expected, out, "out instance")
				assert.NotZero(t, out.String(), "string result value")
			},
		)
	}
}

type testLevel struct {
	name     string
	level    string
	expected Level
}

func TestLevel(test *testing.T) {
	scenarios := []testLevel{
		{
			name:     "Creates default DEBUG Level",
			level:    "",
			expected: DEBUG,
		},
		{
			name:     "Creates an invalid Level",
			level:    "invalid",
			expected: DEBUG,
		},
		{
			name:     "Creates a debug Level",
			level:    "debug",
			expected: DEBUG,
		},
		{
			name:     "Creates a DEBUG Level",
			level:    "DEBUG",
			expected: DEBUG,
		},
		{
			name:     "Creates a info Level",
			level:    "info",
			expected: INFO,
		},
		{
			name:     "Creates a INFO Level",
			level:    "INFO",
			expected: INFO,
		},
		{
			name:     "Creates a error Level",
			level:    "error",
			expected: ERROR,
		},
		{
			name:     "Creates a ERROR Level",
			level:    "ERROR",
			expected: ERROR,
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				var (
					level Level
					err   = level.Set(scenario.level)
				)
				assert.Nil(t, err, "Level.Set error")
				assert.Exactly(t, scenario.expected, level, "level instance")
				assert.NotZero(t, level.String(), "string result value")
			},
		)
	}
}
