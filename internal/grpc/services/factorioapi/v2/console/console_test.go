package consolev2

import (
	"context"
	"testing"

	v2 "github.com/nekomeowww/factorio-rcon-api/v2/apis/factorioapi/v2"
	"github.com/nekomeowww/factorio-rcon-api/v2/internal/libs"
	"github.com/nekomeowww/factorio-rcon-api/v2/internal/rcon/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"
)

func TestCommandEvolution(t *testing.T) {
	t.Parallel()

	logger, err := libs.NewLogger()()
	require.NoError(t, err)
	require.NotNil(t, logger)

	rcon := new(fake.FakeRCON)

	c := NewConsoleService()(NewConsoleServiceParams{
		Logger: logger,
		RCON:   rcon,
	})

	rcon.ExecuteReturns("Nauvis - Evolution factor: 0.9495. (Time 10%) (Pollution 84%) (Spawner kills 5%)\nbpsb-lab-p-Kynazori - Evolution factor: 0.6358. (Time 100%) (Pollution 0%) (Spawner kills 0%)\nbpsb-lab-f-player - Evolution factor: 0.6672. (Time 87%) (Pollution 13%) (Spawner kills 0%)\nbpsb-lab-p-Kunstduenger - Evolution factor: 0.6355. (Time 100%) (Pollution 0%) (Spawner kills 0%)\nbpsb-lab-p-Daemon16 - Evolution factor: 0.6036. (Time 100%) (Pollution 0%) (Spawner kills 0%)\nVulcanus - Evolution factor: 0.5852. (Time 100%) (Pollution 0%) (Spawner kills 0%)\nGleba - Evolution factor: 0.6436. (Time 65%) (Pollution 16%) (Spawner kills 19%)\nFulgora - Evolution factor: 0.5347. (Time 100%) (Pollution 0%) (Spawner kills 0%)\n\n", nil)

	resp, err := c.CommandEvolution(context.Background(), &v2.CommandEvolutionRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Len(t, resp.Evolutions, 8)

	assert.Equal(t, []*v2.Evolution{
		{
			SurfaceName:     "Nauvis",
			EvolutionFactor: 0.9495,
			Time:            10,
			Pollution:       84,
			SpawnerKills:    5,
		},
		{
			SurfaceName:     "bpsb-lab-p-Kynazori",
			EvolutionFactor: 0.6358,
			Time:            100,
			Pollution:       0,
			SpawnerKills:    0,
		},
		{
			SurfaceName:     "bpsb-lab-f-player",
			EvolutionFactor: 0.6672,
			Time:            87,
			Pollution:       13,
			SpawnerKills:    0,
		},
		{
			SurfaceName:     "bpsb-lab-p-Kunstduenger",
			EvolutionFactor: 0.6355,
			Time:            100,
			Pollution:       0,
			SpawnerKills:    0,
		},
		{
			SurfaceName:     "bpsb-lab-p-Daemon16",
			EvolutionFactor: 0.6036,
			Time:            100,
			Pollution:       0,
			SpawnerKills:    0,
		},
		{
			SurfaceName:     "Vulcanus",
			EvolutionFactor: 0.5852,
			Time:            100,
			Pollution:       0,
			SpawnerKills:    0,
		},
		{
			SurfaceName:     "Gleba",
			EvolutionFactor: 0.6436,
			Time:            65,
			Pollution:       16,
			SpawnerKills:    19,
		},
		{
			SurfaceName:     "Fulgora",
			EvolutionFactor: 0.5347,
			Time:            100,
			Pollution:       0,
			SpawnerKills:    0,
		},
	}, resp.Evolutions)

}

func TestCommandEvolutionGet(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		t.Parallel()

		logger, err := libs.NewLogger()()
		require.NoError(t, err)
		require.NotNil(t, logger)

		rcon := new(fake.FakeRCON)

		c := NewConsoleService()(NewConsoleServiceParams{
			Logger: logger,
			RCON:   rcon,
		})

		rcon.ExecuteReturns("\nNauvis - Evolution factor: 0.9495. (Time 10%) (Pollution 84%) (Spawner kills 5%)\n\n", nil)

		resp, err := c.CommandEvolutionGet(context.Background(), &v2.CommandEvolutionGetRequest{
			SurfaceName: "Nauvis",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)

		assert.Equal(t, &v2.Evolution{
			SurfaceName:     "Nauvis",
			EvolutionFactor: 0.9495,
			Time:            10,
			Pollution:       84,
			SpawnerKills:    5,
		}, resp.Evolution)
	})

	t.Run("NotExists", func(t *testing.T) {
		t.Parallel()

		logger, err := libs.NewLogger()()
		require.NoError(t, err)
		require.NotNil(t, logger)

		rcon := new(fake.FakeRCON)

		c := NewConsoleService()(NewConsoleServiceParams{
			Logger: logger,
			RCON:   rcon,
		})

		rcon.ExecuteReturns("Surface \"Nauvis\" does not exist.\n", nil)

		resp, err := c.CommandEvolutionGet(context.Background(), &v2.CommandEvolutionGetRequest{
			SurfaceName: "Nauvis",
		})
		require.Error(t, err)
		require.Nil(t, resp)

		statusErr, ok := status.FromError(err)
		require.True(t, ok)
		require.NotNil(t, statusErr)

		assert.Equal(t, statusErr.Message(), "surface Nauvis does not exist")
	})
}
