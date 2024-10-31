package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringListToPlayers(t *testing.T) {
	players, err := StringListToPlayers("  NekoMeow (online)\n  NekoMeow2\n")
	require.NoError(t, err)
	require.Len(t, players, 2)
}
