package libs

import "go.uber.org/fx"

func Modules() fx.Option {
	return fx.Options(
		fx.Provide(NewLogger()),
		fx.Provide(NewOtel()),
	)
}
