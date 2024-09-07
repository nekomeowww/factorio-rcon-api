package configs

import (
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/nekomeowww/factorio-rcon-api/internal/meta"
	"github.com/nekomeowww/xo"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

type APIServer struct {
	GrpcServerAddr string `json:"grpc_server_addr" yaml:"grpc_server_addr"`
	HttpServerAddr string `json:"http_server_addr" yaml:"http_server_addr"`
}

type Tracing struct {
	OtelCollectorHTTP bool `json:"otel_collector_http" yaml:"otel_collector_http"`
	OtelStdoutEnabled bool `json:"otel_stdout_enabled" yaml:"otel_stdout_enabled"`
}

type Factorio struct {
	RCONHost     string `json:"rcon_host" yaml:"rcon_host"`
	RCONPort     string `json:"rcon_port" yaml:"rcon_port"`
	RCONPassword string `json:"rcon_password" yaml:"rcon_password"`
}

type Config struct {
	meta.Meta `json:"-" yaml:"-"`

	Env       string    `json:"env" yaml:"env"`
	Tracing   Tracing   `json:"tracing" yaml:"tracing"`
	APIServer APIServer `json:"api_server" yaml:"api_server"`
	Factorio  Factorio  `json:"factorio" yaml:"factorio"`
}

func defaultConfig() Config {
	return Config{
		Tracing: Tracing{
			OtelCollectorHTTP: false,
			OtelStdoutEnabled: false,
		},
		APIServer: APIServer{
			GrpcServerAddr: ":24181",
			HttpServerAddr: ":24180",
		},
		Factorio: Factorio{
			RCONHost:     "127.0.0.1",
			RCONPort:     "27015",
			RCONPassword: "123456",
		},
	}
}

func NewConfig(namespace string, app string, configFilePath string, envFilePath string) func() (*Config, error) {
	return func() (*Config, error) {
		configPath := getConfigFilePath(configFilePath)

		lo.Must0(viper.BindEnv("env"))

		lo.Must0(viper.BindEnv("tracing.otel_collector_http"))
		lo.Must0(viper.BindEnv("tracing.otel_stdout_enabled"))

		lo.Must0(viper.BindEnv("api_server.grpc_server_bind"))
		lo.Must0(viper.BindEnv("api_server.http_server_bind"))

		lo.Must0(viper.BindEnv("factorio.rcon_host"))
		lo.Must0(viper.BindEnv("factorio.rcon_port"))
		lo.Must0(viper.BindEnv("factorio.rcon_password"))

		err := loadEnvConfig(envFilePath)
		if err != nil {
			return nil, err
		}

		err = readConfig(configPath)
		if err != nil {
			return nil, err
		}

		config := defaultConfig()

		err = viper.Unmarshal(&config, func(c *mapstructure.DecoderConfig) {
			c.TagName = "yaml"
		})
		if err != nil {
			return nil, err
		}

		xo.PrintJSON(config)

		meta.Env = config.Env
		if meta.Env == "" {
			meta.Env = os.Getenv("ENV")
		}

		config.Meta.Env = config.Env
		config.Meta.App = app
		config.Meta.Namespace = namespace

		return &config, nil
	}
}

func NewTestConfig(envFilePath string) (*Config, error) {
	configPath := tryToMatchConfigPathForUnitTest("")

	if envFilePath != "" {
		err := loadEnvConfig("")
		if err != nil {
			return nil, err
		}
	}

	err := readConfig(configPath)
	if err != nil {
		return nil, err
	}

	config := defaultConfig()
	config.Env = "test"

	err = viper.Unmarshal(&config, func(c *mapstructure.DecoderConfig) {
		c.TagName = "yaml"
	})
	if err != nil {
		return nil, err
	}

	return &config, nil
}
