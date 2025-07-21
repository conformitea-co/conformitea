package infrastructure

import (
	"fmt"

	domainOrganization "conformitea/domain/organization"
	domainTeam "conformitea/domain/team"
	domainUser "conformitea/domain/user"
	"conformitea/infrastructure/config"
	"conformitea/infrastructure/database"
	"conformitea/infrastructure/gateway/hydra"
	"conformitea/infrastructure/gateway/microsoft"
	"conformitea/infrastructure/logger"
	"conformitea/infrastructure/persistence/organization"
	"conformitea/infrastructure/persistence/team"
	"conformitea/infrastructure/persistence/user"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Persistence struct {
	user         domainUser.UserRepository
	team         domainTeam.TeamRepository
	organization domainOrganization.OrganizationRepository
}

type Container struct {
	config          config.Config
	logger          *zap.Logger
	database        *gorm.DB
	hydraClient     *hydra.HydraClient
	microsoftClient *microsoft.OAuthClient
	persistence     Persistence
}

var container *Container

func Initialize(lc config.LoggerConfig, dc config.DatabaseConfig, hc config.HydraConfig, oc config.OAuthConfig) (*Container, error) {
	l, err := logger.Initialize(lc)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	db, err := database.Initialize(dc, l)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	h, err := hydra.Initialize(hc)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Hydra client: %w", err)
	}

	ms, err := microsoft.Initialize(oc.Microsoft)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Microsoft OAuth client: %w", err)
	}

	container = &Container{
		config: config.Config{
			LoggerConfig:   lc,
			DatabaseConfig: dc,
			HydraConfig:    hc,
			OAuthConfig:    oc,
		},
		logger:          l,
		database:        db,
		hydraClient:     h,
		microsoftClient: ms,
		persistence: Persistence{
			user:         &user.UserRepository{},
			team:         &team.TeamRepository{},
			organization: &organization.OrganizationRepository{},
		},
	}

	return container, nil
}

func (c *Container) GetLogger() *zap.Logger {
	return c.logger
}

func (c *Container) GetDatabase() *gorm.DB {
	return c.database
}

func (c *Container) GetHydraClient() *hydra.HydraClient {
	return c.hydraClient
}

func (c *Container) GetMicrosoftClient() *microsoft.OAuthClient {
	return c.microsoftClient
}

func (c *Container) GetPersistence() Persistence {
	return c.persistence
}

func (p *Persistence) GetUserRepository() domainUser.UserRepository {
	return p.user
}

func (p *Persistence) GetTeamRepository() domainTeam.TeamRepository {
	return p.team
}

func (p *Persistence) GetOrganizationRepository() domainOrganization.OrganizationRepository {
	return p.organization
}
