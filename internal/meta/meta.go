package meta

var (
	Version    = "1.0.0"
	LastCommit = "abcdefgh"
	Env        = "dev"
)

type Meta struct {
	Namespace  string `json:"namespace" yaml:"namespace"`
	App        string `json:"app" yaml:"app"`
	Version    string `json:"version" yaml:"version"`
	LastCommit string `json:"last_commit" yaml:"last_commit"`
	Env        string `json:"env" yaml:"env"`
}

func NewMeta(namespace, app string) *Meta {
	return &Meta{
		Namespace:  namespace,
		App:        app,
		Version:    Version,
		LastCommit: LastCommit,
		Env:        Env,
	}
}
