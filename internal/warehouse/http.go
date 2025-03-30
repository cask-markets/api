package warehouse

import (
	"context"
	flagsService "github.com/flags-gg/go-flags"
	ConfigBuilder "github.com/keloran/go-config"
	"net/http"
)

type System struct {
	Context context.Context
	Config  *ConfigBuilder.Config
}

func NewSystem(cfg *ConfigBuilder.Config) *System {
	return &System{
		Context: context.Background(),
		Config:  cfg,
	}
}

func (s *System) SetContext(ctx context.Context) {
	s.Context = ctx
}

func (s *System) GetWarehouses(w http.ResponseWriter, r *http.Request) {
	flags := flagsService.NewClient(flagsService.WithAuth(flagsService.Auth{
		ProjectID:     s.Config.ProjectProperties["flags_project"].(string),
		AgentID:       s.Config.ProjectProperties["flags_agent"].(string),
		EnvironmentID: s.Config.ProjectProperties["flags_environment"].(string),
	}))

	if flags.Is("warehouses-get").Enabled() {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}
