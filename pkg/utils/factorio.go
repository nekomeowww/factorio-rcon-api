package utils

import (
	"strconv"
	"strings"
	"time"

	v1 "github.com/nekomeowww/factorio-rcon-api/v2/apis/factorioapi/v1"
	v2 "github.com/nekomeowww/factorio-rcon-api/v2/apis/factorioapi/v2"
	"github.com/nekomeowww/factorio-rcon-api/v2/pkg/apierrors"
	"github.com/samber/lo"
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
		if len(parts) > 2 {
			return nil, apierrors.NewErrBadRequest().WithDetailf("failed to parse admins: %s due to parts not equals 2", line).AsStatus()
		}

		player := &v1.Player{
			Username: parts[0],
		}
		if len(parts) == 2 {
			player.Online = parts[1] == "(online)"
		}

		players = append(players, player)
	}

	return players, nil
}

func PrefixedStringCommaSeparatedListToPlayers(list string, prefix string) ([]*v1.Player, error) {
	// output:
	// SomePrefix: NekoMeow, NekoMeow2\n
	withoutPrefix := strings.TrimPrefix(list, prefix+": ")
	split := strings.Split(withoutPrefix, ",")
	split = lo.Map(split, func(item string, _ int) string { return strings.TrimSpace(item) })

	return lo.Map(split, func(item string, _ int) *v1.Player {
		return &v1.Player{Username: item}
	}), nil
}

func MapV1PlayerToV2Player(v1Player *v1.Player) *v2.Player {
	return &v2.Player{
		Username: v1Player.Username,
		Online:   v1Player.Online,
	}
}

func MapV1PlayersToV2Players(v1Players []*v1.Player) []*v2.Player {
	return lo.Map(v1Players, func(item *v1.Player, _ int) *v2.Player { return MapV1PlayerToV2Player(item) })
}

func ParseDuration(input string) (time.Duration, error) {
	// Split the input string into parts
	parts := strings.Fields(input)

	// Initialize the total duration
	var totalDuration time.Duration

	// Iterate over the parts and parse the time values
	for i := 0; i < len(parts); i++ {
		switch parts[i] {
		case "days":
			if i > 0 {
				days, err := strconv.ParseInt(parts[i-1], 10, 64)
				if err != nil {
					return 0, err
				}

				totalDuration += time.Duration(days) * 24 * time.Hour
			}
		case "hours":
			if i > 0 {
				hours, err := time.ParseDuration(parts[i-1] + "h")
				if err != nil {
					return 0, err
				}

				totalDuration += hours
			}
		case "minutes":
			if i > 0 {
				minutes, err := time.ParseDuration(parts[i-1] + "m")
				if err != nil {
					return 0, err
				}

				totalDuration += minutes
			}
		case "seconds":
			if i > 0 {
				seconds, err := time.ParseDuration(parts[i-1] + "s")
				if err != nil {
					return 0, err
				}

				totalDuration += seconds
			}
		}
	}

	return totalDuration, nil
}

func ParseWhitelistedPlayers(input string) []string {
	// Replace " and " with a comma to unify the delimiters
	input = strings.ReplaceAll(input, " and ", ", ")

	// Split the players by commas and trim spaces
	players := strings.Split(input, ", ")
	for i := range players {
		players[i] = strings.TrimSpace(players[i])
	}

	return players
}
