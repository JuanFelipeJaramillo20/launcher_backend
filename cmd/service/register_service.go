package service

import (
	"errors"
	"fmt"
	"github.com/resend/resend-go/v2"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"venecraft-back/cmd/entity"
	"venecraft-back/cmd/enums"
	"venecraft-back/cmd/repository"
)

type RegisterService interface {
	CreateRegister(register *entity.Register) error
	ApproveRegister(id uint64) (*entity.User, error)
}

type registerService struct {
	registerRepo repository.RegisterRepository
	userRepo     repository.UserRepository
	roleRepo     repository.RoleRepository
	userRoleRepo repository.UserRoleRepository
	emailClient  *resend.Client
}

func NewRegisterService(registerRepo repository.RegisterRepository, userRepo repository.UserRepository, roleRepo repository.RoleRepository, userRoleRepo repository.UserRoleRepository) RegisterService {
	apiKey := os.Getenv("RESEND_API_KEY")
	emailClient := resend.NewClient(apiKey)

	return &registerService{
		registerRepo: registerRepo,
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
		emailClient:  emailClient,
	}
}

func (s *registerService) CreateRegister(register *entity.Register) error {
	// Hash the password before saving
	hashedPassword, err := hashPassword(register.Password)
	if err != nil {
		return err
	}
	register.Password = hashedPassword

	// Attempt to create the registration record
	err = s.registerRepo.CreateRegister(register)
	if err != nil {
		return err
	}

	// Prepare cleanup in case of a later failure by deferring the deletion of the register record
	defer func() {
		if err != nil { // Only executes if there's an error
			delErr := s.registerRepo.DeleteRegister(register.ID)
			if delErr != nil {
				log.Printf("Error cleaning up registration for user %s: %v", register.Email, delErr)
			} else {
				log.Printf("Cleanup: registration for user %s deleted successfully", register.Email)
			}
		}
	}()

	// Fetch all admins to notify them of the new registration
	admins, err := s.userRepo.GetUsersByRole(enums.RoleAdmin)
	if err != nil {
		return fmt.Errorf("failed to fetch admin users: %v", err)
	}

	// Collect all admin email addresses
	adminEmails := make([]string, len(admins))
	for i, admin := range admins {
		adminEmails[i] = admin.Email
	}

	// Attempt to send a confirmation email to the user
	err = s.sendUserConfirmationEmail(register.Email)
	if err != nil {
		log.Printf("Error sending confirmation email to user %s: %v", register.Email, err)
		return fmt.Errorf("failed to send confirmation email to the user")
	}

	// Attempt to notify admins about the new registration request
	err = s.sendAdminNotificationEmail(adminEmails, register)
	if err != nil {
		log.Printf("Error sending notification email to admins about registration by user %s: %v", register.Email, err)
		return fmt.Errorf("failed to send registration notification email to admins")
	}

	log.Println("Registration emails sent successfully.")
	return nil
}

func (s *registerService) ApproveRegister(id uint64) (*entity.User, error) {
	register, err := s.registerRepo.GetRegisterByID(id)
	if err != nil {
		return nil, errors.New("registration request not found")
	}

	user := &entity.User{
		FullName: register.FullName,
		Email:    register.Email,
		Nickname: register.Nickname,
		Password: register.Password,
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	playerRole, err := s.roleRepo.GetRoleByName(enums.RolePlayer)
	if err != nil {
		return nil, errors.New("failed to assign role: PLAYER role not found")
	}

	userRole := &entity.UserRole{
		UserID: user.ID,
		RoleID: playerRole.ID,
	}
	err = s.userRoleRepo.AssignRole(userRole)
	if err != nil {
		return nil, err
	}

	err = s.registerRepo.DeleteRegister(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *registerService) sendUserConfirmationEmail(userEmail string) error {
	params := &resend.SendEmailRequest{
		From:    "Registration Service <onboarding@jjar.lat>",
		To:      []string{userEmail},
		Html:    "<strong>Your registration request has been created successfully. Please wait for an admin to review your request. You will be notified by email once your request is processed.</strong>",
		Subject: "Registration Request Created",
	}

	sent, err := s.emailClient.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send user confirmation email: %v", err)
	}

	fmt.Println("User confirmation email sent successfully with ID:", sent.Id)
	return nil
}

func (s *registerService) sendAdminNotificationEmail(adminEmails []string, registerDetails *entity.Register) error {
	emailContent := fmt.Sprintf(
		"<strong>A new registration request has been received.</strong><br><br>"+
			"<b>Full Name:</b> %s<br><b>Email:</b> %s<br><b>Nickname:</b> %s<br>",
		registerDetails.FullName,
		registerDetails.Email,
		registerDetails.Nickname,
	)

	params := &resend.SendEmailRequest{
		From:    "Registration Service <onboarding@jjar.lat>",
		To:      adminEmails,
		Html:    emailContent,
		Subject: "New Registration Request for Review",
	}

	sent, err := s.emailClient.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send admin notification email: %v", err)
	}

	fmt.Println("Admin notification email sent successfully with ID:", sent.Id)
	return nil
}
