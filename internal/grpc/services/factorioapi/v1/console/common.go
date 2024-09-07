package consolev1

import (
	"strings"

	v1 "github.com/nekomeowww/factorio-rcon-api/apis/factorioapi/v1"
	"github.com/nekomeowww/factorio-rcon-api/pkg/apierrors"
)

func StringListToPlayers(list string) ([]*v1.Player, error) {
	// output:
	//   NekoMeow (online)\n
	//   NekoMeow2 (offline)\n
	split := strings.Split(strings.TrimSuffix(list, "\n"), "\n")
	players := make([]*v1.Player, 0, len(split))
	for _, line := range split {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			return nil, apierrors.NewErrBadRequest().WithDetailf("failed to parse admins: %s due to parts not equals 2", line).AsStatus()
		}

		players = append(players, &v1.Player{
			Username: parts[0],
			Online:   parts[1] == "(online)",
		})
	}

	return players, nil
}
