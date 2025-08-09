package commands

import (
	"conformitea/app/auth"
	cmd "conformitea/cmd/config"
	"conformitea/domain"
	"conformitea/infrastructure"
	"conformitea/server"
	serverConfig "conformitea/server/config"
	"conformitea/server/types"

	"github.com/spf13/cobra"
)

func ServeCmd(config cmd.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "serve",
		Short:   "Start the ConformiTea server",
		Aliases: []string{"s", "start"},
		RunE: func(cmd *cobra.Command, args []string) error {
			server, err := initializeServer(config)
			if err != nil {
				return err
			}

			if err := server.Start(); err != nil {
				return err
			}

			return nil
		},
	}
}

func initializeServer(c cmd.Config) (types.Server, error) {
	ic, err := initializeInfrastructure(c)
	if err != nil {
		return nil, err
	}

	dc, err := initializeDomain(ic.GetPersistence())
	if err != nil {
		return nil, err
	}

	auth := initializeApp(dc, ic)

	sc := serverConfig.Config{
		General:    c.GeneralConfig,
		HTTPServer: c.HTTPServerConfig,
		Redis:      c.RedisConfig,
	}

	return server.Initialize(sc, ic.GetLogger(), auth)
}

func initializeApp(dc *domain.Container, ic *infrastructure.Container) *auth.Auth {
	auth := auth.Initialize(
		ic.GetDatabase(),
		dc.GetUserService(),
		ic.GetMicrosoftClient(),
		ic.GetHydraClient(),
	)

	return auth
}

func initializeDomain(p infrastructure.Persistence) (*domain.Container, error) {
	container, err := domain.Initialize(p.GetUserRepository(), p.GetTeamRepository(), p.GetOrganizationRepository())
	if err != nil {
		return nil, err
	}

	return container, nil
}

func initializeInfrastructure(c cmd.Config) (*infrastructure.Container, error) {
	container, err := infrastructure.Initialize(c.LoggerConfig, c.DatabaseConfig, c.HydraConfig, c.OAuthConfig)
	if err != nil {
		return nil, err
	}

	return container, nil
}
