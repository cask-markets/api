package main

import (
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	"github.com/cask-warehouse/api/internal"
	ConfigBuilder "github.com/keloran/go-config"
)

var (
	BuildVersion = "0.0.1"
	BuildHash    = "unknown"
	ServiceName  = "service"
)

type ProjectConfig struct{}

func (pc ProjectConfig) Build(cfg *ConfigBuilder.Config) error {
	type FlagsService struct {
		ProjectID     string `env:"FLAGS_PROJECT_ID" envDefault:"flags-gg"`
		AgentID       string `env:"FLAGS_AGENT_ID" envDefault:"orchestrator"`
		EnvironmentID string `env:"FLAGS_ENVIRONMENT_ID" envDefault:"orchestrator"`
	}

	type PC struct {
		StripeSecret string `env:"STRIPE_SECRET" envDefault:"stripe_secret"`
		RailwayPort  string `env:"PORT" envDefault:"3000"`
		OnRailway    bool   `env:"ON_RAILWAY" envDefault:"false"`
		Flags        FlagsService
	}
	p := PC{}

	if err := env.Parse(&p); err != nil {
		return logs.Errorf("Failed to parse services: %v", err)
	}
	if cfg.ProjectProperties == nil {
		cfg.ProjectProperties = make(map[string]interface{})
	}
	cfg.ProjectProperties["stripeKey"] = p.StripeSecret
	cfg.ProjectProperties["railway_port"] = p.RailwayPort
	cfg.ProjectProperties["on_railway"] = p.OnRailway

	cfg.ProjectProperties["flags_agent"] = p.Flags.AgentID
	cfg.ProjectProperties["flags_environment"] = p.Flags.EnvironmentID
	cfg.ProjectProperties["flags_project"] = p.Flags.ProjectID

	return nil
}

func main() {
	logs.Logf("Starting %s version %s (build %s)", ServiceName, BuildVersion, BuildHash)
	c := ConfigBuilder.NewConfigNoVault()

	err := c.Build(
		ConfigBuilder.Local,
		ConfigBuilder.Mongo,
		ConfigBuilder.Bugfixes,
		ConfigBuilder.Clerk,
		ConfigBuilder.Resend,
		ConfigBuilder.WithProjectConfigurator(ProjectConfig{}))
	if err != nil {
		logs.Fatalf("Failed to build config: %v", err)
	}

	if err := internal.New(c).Start(); err != nil {
		logs.Fatalf("Failed to start service: %v", err)
	}
}
