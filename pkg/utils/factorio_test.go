package utils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringListToPlayers(t *testing.T) {
	players, err := StringListToPlayers("  NekoMeow (online)\n  NekoMeow2\n")
	require.NoError(t, err)
	require.Len(t, players, 2)
}

func TestParseWhitelistedPlayers(t *testing.T) {
	type testCase struct {
		input    string
		expected []string
	}

	testCases := []testCase{
		{input: "NekoMeow", expected: []string{"NekoMeow"}},
		{input: "LittleSound and NekoMeow", expected: []string{"LittleSound", "NekoMeow"}},
		{input: "LemonNeko, LittleSound and NekoMeow", expected: []string{"LemonNeko", "LittleSound", "NekoMeow"}},
	}

	for index, tc := range testCases {
		t.Run(strconv.Itoa(index), func(t *testing.T) {
			players := ParseWhitelistedPlayers(tc.input)
			require.Equal(t, tc.expected, players)
		})
	}
}
