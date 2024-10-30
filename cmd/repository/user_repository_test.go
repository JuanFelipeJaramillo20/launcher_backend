package repository_test

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"venecraft-back/cmd/entity"
	"venecraft-back/cmd/repository"
)

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&entity.User{}, &entity.Role{})
	return db, err
}

func setupUserRepo() (repository.UserRepository, *gorm.DB) {
	db, _ := setupTestDB()
	repo := repository.NewUserRepository(db)
	return repo, db
}

func createTestUser(db *gorm.DB, nickname, email string) *entity.User {
	user := &entity.User{
		FullName: "Test User",
		Email:    email,
		Nickname: nickname,
		Password: "hashedpassword",
	}
	db.Create(user)
	return user
}

func TestCreateUser(t *testing.T) {
	repo, db := setupUserRepo()
	defer db.Migrator().DropTable(&entity.User{})

	user := &entity.User{
		FullName: "John Doe",
		Email:    "john@example.com",
		Nickname: "johndoe",
		Password: "password123",
	}

	err := repo.CreateUser(user)
	assert.NoError(t, err)

	var fetchedUser entity.User
	db.First(&fetchedUser, user.ID)
	assert.Equal(t, user.Email, fetchedUser.Email)
	assert.Equal(t, user.Nickname, fetchedUser.Nickname)
}

func TestGetAllUsers(t *testing.T) {
	repo, db := setupUserRepo()
	defer db.Migrator().DropTable(&entity.User{})

	createTestUser(db, "testuser1", "test1@example.com")
	createTestUser(db, "testuser2", "test2@example.com")

	users, err := repo.GetAllUsers()
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestGetUserByID(t *testing.T) {
	repo, db := setupUserRepo()
	defer db.Migrator().DropTable(&entity.User{})

	user := createTestUser(db, "testuser", "test@example.com")

	fetchedUser, err := repo.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, fetchedUser.Email)
	assert.Equal(t, user.Nickname, fetchedUser.Nickname)
}

func TestUpdateUser(t *testing.T) {
	repo, db := setupUserRepo()
	defer db.Migrator().DropTable(&entity.User{})

	user := createTestUser(db, "testuser", "test@example.com")
	user.Email = "newemail@example.com"
	err := repo.UpdateUser(user)
	assert.NoError(t, err)

	fetchedUser, _ := repo.GetUserByID(user.ID)
	assert.Equal(t, "newemail@example.com", fetchedUser.Email)
}

func TestDeleteUser(t *testing.T) {
	repo, db := setupUserRepo()
	defer db.Migrator().DropTable(&entity.User{})

	user := createTestUser(db, "testuser", "test@example.com")

	err := repo.DeleteUser(user.ID)
	assert.NoError(t, err)

	_, err = repo.GetUserByID(user.ID)
	assert.Error(t, err)
}

func TestGetUserByEmail(t *testing.T) {
	repo, db := setupUserRepo()
	defer db.Migrator().DropTable(&entity.User{})

	createTestUser(db, "testuser", "test@example.com")

	user, err := repo.GetUserByEmail("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Nickname)
}

func TestGetUserByNickname(t *testing.T) {
	repo, db := setupUserRepo()
	defer db.Migrator().DropTable(&entity.User{})

	createTestUser(db, "testuser", "test@example.com")

	user, err := repo.GetUserByNickname("testuser")
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestGetUsersByRole(t *testing.T) {
	repo, db := setupUserRepo()
	defer db.Migrator().DropTable(&entity.User{})
	defer db.Migrator().DropTable(&entity.Role{})

	// Setup roles and users
	role := &entity.Role{Name: "Admin"}
	db.Create(role)
	user := createTestUser(db, "adminuser", "admin@example.com")
	db.Model(user).Association("Roles").Append(role)

	users, err := repo.GetUsersByRole("Admin")
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "adminuser", users[0].Nickname)
}

func TestGetUserByResetToken(t *testing.T) {
	repo, db := setupUserRepo()
	defer db.Migrator().DropTable(&entity.User{})

	user := createTestUser(db, "testuser", "test@example.com")
	user.RecoverPasswordToken = "testtoken"
	db.Save(user)

	fetchedUser, err := repo.GetUserByResetToken("testtoken")
	assert.NoError(t, err)
	assert.Equal(t, user.Email, fetchedUser.Email)
}

func TestHasRole(t *testing.T) {
	repo, db := setupUserRepo()
	defer db.Migrator().DropTable(&entity.User{})
	defer db.Migrator().DropTable(&entity.Role{})

	role := &entity.Role{Name: "Admin"}
	db.Create(role)
	user := createTestUser(db, "adminuser", "admin@example.com")
	db.Model(user).Association("Roles").Append(role)

	hasRole := repo.HasRole(user.ID, "Admin")
	assert.True(t, hasRole)
}
