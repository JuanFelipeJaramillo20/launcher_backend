package seeds

import (
	"gorm.io/gorm"
	"log"
	"os"
	"venecraft-back/cmd/entity"
)

func SeedUsers(db *gorm.DB) {
	var adminRole entity.Role
	err := db.Where("name = ?", "ADMIN").First(&adminRole).Error
	if err != nil {
		log.Fatalf("Error finding ADMIN role: %v", err)
	}

	adminUser := entity.User{
		FullName: os.Getenv("ADMIN_FULL_NAME"),
		Email:    os.Getenv("ADMIN_EMAIL"),
		Nickname: os.Getenv("ADMIN_NICKNAME"),
		Password: os.Getenv("ADMIN_PASSWORD"),
	}

	err = db.Where("email = ?", adminUser.Email).FirstOrCreate(&adminUser).Error
	if err != nil {
		log.Fatalf("Error seeding admin user: %v", err)
	}
	log.Printf("Admin user %s seeded successfully", adminUser.Email)

	adminUserRole := entity.UserRole{
		UserID: adminUser.ID,
		RoleID: adminRole.ID,
	}

	err = db.Where("user_id = ? AND role_id = ?", adminUser.ID, adminRole.ID).FirstOrCreate(&adminUserRole).Error
	if err != nil {
		log.Fatalf("Error assigning ADMIN role to admin user: %v", err)
	}
	log.Println("ADMIN role assigned to admin user successfully")
}
