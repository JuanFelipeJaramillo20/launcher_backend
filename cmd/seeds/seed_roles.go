package seeds

import (
	"gorm.io/gorm"
	"log"
	"venecraft-back/cmd/entity"
)

func SeedRoles(db *gorm.DB) {
	roles := []entity.Role{
		{Name: "PLAYER"},
		{Name: "ADMIN"},
		{Name: "MODERATOR"},
	}

	for _, role := range roles {
		err := db.Where("name = ?", role.Name).FirstOrCreate(&role).Error
		if err != nil {
			log.Fatalf("Error seeding roles: %v", err)
		}
	}
}
