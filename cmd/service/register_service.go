package service

import (
	"errors"
	"fmt"
	"github.com/resend/resend-go/v2"
	"golang.org/x/crypto/bcrypt"
	"log"
	"venecraft-back/cmd/email"
	"venecraft-back/cmd/entity"
	"venecraft-back/cmd/enums"
	"venecraft-back/cmd/repository"
)

type RegisterService interface {
	CreateRegister(register *entity.Register) error
	ApproveRegister(id uint64) (*entity.User, error)
	DenyRegister(id uint64) error
}

type registerService struct {
	registerRepo repository.RegisterRepository
	userRepo     repository.UserRepository
	roleRepo     repository.RoleRepository
	userRoleRepo repository.UserRoleRepository
	emailClient  *email.EmailClient
}

func NewRegisterService(registerRepo repository.RegisterRepository, userRepo repository.UserRepository, roleRepo repository.RoleRepository, userRoleRepo repository.UserRoleRepository) RegisterService {
	return &registerService{
		registerRepo: registerRepo,
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
		emailClient:  email.GetEmailClient(),
	}
}

func (s *registerService) CreateRegister(register *entity.Register) error {
	hashedPassword, err := hashPassword(register.Password)
	if err != nil {
		return err
	}
	register.Password = hashedPassword

	err = s.registerRepo.CreateRegister(register)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			delErr := s.registerRepo.DeleteRegister(register.ID)
			if delErr != nil {
				log.Printf("Error cleaning up registration for user %s: %v", register.Email, delErr)
			} else {
				log.Printf("Cleanup: registration for user %s deleted successfully", register.Email)
			}
		}
	}()

	admins, err := s.userRepo.GetUsersByRole(enums.RoleAdmin)
	if err != nil {
		return fmt.Errorf("failed to fetch admin users: %v", err)
	}

	adminEmails := make([]string, len(admins))
	for i, admin := range admins {
		adminEmails[i] = admin.Email
	}

	err = s.sendUserConfirmationEmail(register.Email, register.Nickname)
	if err != nil {
		log.Printf("Error sending confirmation email to user %s: %v", register.Email, err)
		return fmt.Errorf("failed to send confirmation email to the user")
	}

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

	err = s.sendUserResponseEmail(register.Email, true)
	if err != nil {
		log.Printf("Error sending approval email to user %s: %v", register.Email, err)
	}

	return user, nil
}

func (s *registerService) DenyRegister(id uint64) error {
	register, err := s.registerRepo.GetRegisterByID(id)
	if err != nil {
		return errors.New("registration request not found")
	}

	err = s.registerRepo.DeleteRegister(id)
	if err != nil {
		return err
	}

	err = s.sendUserResponseEmail(register.Email, false)
	if err != nil {
		log.Printf("Error sending denial email to user %s: %v", register.Email, err)
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *registerService) sendUserConfirmationEmail(userEmail string, userName string) error {
	// Use the full relative path for the template
	body, err := email.RenderTemplate("register/user_confirmation.html", map[string]string{"Name": userName})
	if err != nil {
		return fmt.Errorf("failed to render user confirmation email template: %v", err)
	}

	params := &resend.SendEmailRequest{
		From:    "Registration Service <onboarding@jjar.lat>",
		To:      []string{userEmail},
		Html:    body,
		Subject: "Registration Request Created",
	}

	sent, err := s.emailClient.SendEmail(params)
	if err != nil {
		return fmt.Errorf("failed to send user confirmation email: %v", err)
	}

	fmt.Println("User confirmation email sent successfully with ID:", sent.Id)
	return nil
}

func (s *registerService) sendAdminNotificationEmail(adminEmails []string, registerDetails *entity.Register) error {
	// Use the full relative path for the template
	body, err := email.RenderTemplate("register/admin_notification.html", map[string]string{
		"FullName": registerDetails.FullName,
		"Email":    registerDetails.Email,
		"Nickname": registerDetails.Nickname,
	})
	if err != nil {
		return fmt.Errorf("failed to render admin notification email template: %v", err)
	}

	params := &resend.SendEmailRequest{
		From:    "Registration Service <onboarding@jjar.lat>",
		To:      adminEmails,
		Html:    body,
		Subject: "New Registration Request for Review",
	}

	sent, err := s.emailClient.SendEmail(params)
	if err != nil {
		return fmt.Errorf("failed to send admin notification email: %v", err)
	}

	fmt.Println("Admin notification email sent successfully with ID:", sent.Id)
	return nil
}

func (s *registerService) sendUserResponseEmail(userEmail string, accepted bool) error {
	templatePath := "register/user_response.html"
	data := map[string]string{
		"Message": "Your registration request has been approved. Welcome to the platform!",
	}
	if !accepted {
		data["Message"] = "We're sorry, but your registration request has been denied."
	}

	body, err := email.RenderTemplate(templatePath, data)
	if err != nil {
		return fmt.Errorf("failed to render user response email template: %v", err)
	}

	params := &resend.SendEmailRequest{
		From:    "Registration Service <onboarding@jjar.lat>",
		To:      []string{userEmail},
		Html:    body,
		Subject: "Registration Request Status",
	}

	sent, err := s.emailClient.SendEmail(params)
	if err != nil {
		return fmt.Errorf("failed to send user response email: %v", err)
	}

	fmt.Println("User response email sent successfully with ID:", sent.Id)
	return nil
}
