package main

import (
	"context"
	"log"
	"time"

	"go.uber.org/fx"

	"github.com/nekomeowww/factorio-rcon-api/v2/internal/configs"
	grpcservers "github.com/nekomeowww/factorio-rcon-api/v2/internal/grpc/servers"
	apiserver "github.com/nekomeowww/factorio-rcon-api/v2/internal/grpc/servers/factorioapi/apiserver"
	grpcservices "github.com/nekomeowww/factorio-rcon-api/v2/internal/grpc/services"
	"github.com/nekomeowww/factorio-rcon-api/v2/internal/libs"
	"github.com/nekomeowww/factorio-rcon-api/v2/internal/rcon"
	"github.com/spf13/cobra"
)

var (
	configFilePath string
	envFilePath    string
)

func main() {
	root := &cobra.Command{
		Use: "api-server",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := fx.New(
				fx.Provide(configs.NewConfig("factorio-rcon-api", "api-server", configFilePath, envFilePath)),
				fx.Options(libs.Modules()),
				fx.Options(rcon.Modules()),
				fx.Options(grpcservers.Modules()),
				fx.Options(grpcservices.Modules()),
				fx.Invoke(apiserver.RunGRPCServer()),
				fx.Invoke(apiserver.RunGatewayServer()),
			)

			app.Run()

			stopCtx, stopCtxCancel := context.WithTimeout(context.Background(), time.Minute*5)
			defer stopCtxCancel()

			if err := app.Stop(stopCtx); err != nil {
				return err
			}

			return nil
		},
	}

	root.Flags().StringVarP(&configFilePath, "config", "c", "", "config file path")
	root.Flags().StringVarP(&envFilePath, "env", "e", "", "env file path")

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
