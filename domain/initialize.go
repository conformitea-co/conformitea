package domain

import (
	"conformitea/domain/organization"
	"conformitea/domain/team"
	"conformitea/domain/user"
)

type Container struct {
	user         *user.UserService
	team         *team.TeamService
	organization *organization.OrganizationService
}

func Initialize(ur user.UserRepository, tr team.TeamRepository, or organization.OrganizationRepository) (*Container, error) {
	us := user.Initialize(ur)
	ts := team.Initialize(tr)
	os := organization.Initialize(or)

	return &Container{
		user:         us,
		team:         ts,
		organization: os,
	}, nil
}

func (c *Container) GetUserService() *user.UserService {
	return c.user
}

func (c *Container) GetTeamService() *team.TeamService {
	return c.team
}

func (c *Container) GetOrganizationService() *organization.OrganizationService {
	return c.organization
}
