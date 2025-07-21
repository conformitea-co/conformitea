package organization

import (
	domain "conformitea/domain/organization"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationRepository struct{}

func (o *OrganizationRepository) GetOrganizationByID(DB *gorm.DB, id uuid.UUID) (domain.Organization, error) {
	var organization Organization

	if err := DB.Where("id = ?", id).First(&organization).Error; err != nil {
		return domain.Organization{}, err
	}

	return domain.Organization{
		ID:   organization.ID,
		Name: organization.Name,
	}, nil
}
