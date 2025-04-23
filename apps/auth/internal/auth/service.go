package auth

import (
	"context"
	"errors"
	database "github.com/FireFly4ik/Lavka-auth/internal/db"
	"github.com/FireFly4ik/Lavka-auth/internal/models"
	"github.com/FireFly4ik/Lavka-auth/pkg/utils"
	pb "github.com/FireFly4ik/Lavka-auth/proto"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

func (s *AuthServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.RegisterResponse{Message: "Could not create user"}, err
		}
	}

	if user.ID != uuid.Nil && user.IsEmailVerified {
		return &pb.RegisterResponse{Message: "User already registered"}, nil
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return &pb.RegisterResponse{Message: "Could not hash password"}, err
	}

	newUser := models.User{
		Email:    req.Email,
		Password: hashedPassword,
		Username: req.Username,
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return &pb.RegisterResponse{Message: "Could not start transaction"}, tx.Error
	}

	if user.ID == uuid.Nil {
		if err := tx.Create(&newUser).Error; err != nil {
			tx.Rollback()
			return &pb.RegisterResponse{Message: "Could not create user"}, err
		}
		user = newUser
	} else {
		if err := tx.Model(&user).Updates(newUser).Error; err != nil {
			tx.Rollback()
			return &pb.RegisterResponse{Message: "Could not update user"}, err
		}
	}

	token := uuid.New().String()
	verification := models.Verification{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Minute * 15),
	}

	if err := tx.Create(&verification).Error; err != nil {
		tx.Rollback()
		return &pb.RegisterResponse{Message: "Could not create verification"}, err
	}

	if err := tx.Commit().Error; err != nil {
		return &pb.RegisterResponse{Message: "Could not commit transaction"}, err
	}

	type EmailData struct{ URL string }
	data := EmailData{
		URL: req.RedirectUrl + "?token=" + token,
	}

	if err := s.EmailService.SendEmail(user.Email, "Verify email", "verify_email", data); err != nil {
		return &pb.RegisterResponse{Message: "Could not send email"}, err
	}

	return &pb.RegisterResponse{Message: "User registered, please verify your email"}, nil
}

func (s *AuthServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user models.User
	if err := database.DB.Where("email = ? AND is_email_verified = ?", req.Login, true).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	session := models.Session{
		UserID:    user.ID,
		Agent:     req.UserAgent,
		IP:        req.IpAddress,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	if err := database.DB.Create(&session).Error; err != nil {
		return nil, errors.New("could not create session")
	}

	accessToken, err := s.JwtService.GenerateAccessToken(time.Hour*4, utils.JwtCustomClaim{
		UserID: user.ID.String(),
		Role:   string(user.Role),
	})
	if err != nil {
		return nil, errors.New("could not generate access token")
	}

	refreshToken, err := s.JwtService.GenerateRefreshToken(time.Hour*24*7, utils.JwtCustomClaim{
		UserID: user.ID.String(),
		Role:   string(user.Role),
	}, session.ID.String())
	if err != nil {
		return nil, errors.New("could not generate refresh token")
	}

	return &pb.LoginResponse{
		Message:      "User logged in successfully",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthServiceServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	claims, err := s.JwtService.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if err := database.DB.Model(&models.Session{}).Where("id = ?", claims.ID).
		Updates(map[string]interface{}{"is_active": false, "deleted_at": time.Now()}).Error; err != nil {
		return nil, errors.New("could not logout user")
	}

	return &pb.LogoutResponse{Message: "User logged out successfully"}, nil
}

func (s *AuthServiceServer) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	claims, err := s.JwtService.ValidateToken(req.Token)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	var session models.Session
	if err := database.DB.Where("id = ? AND is_active = ?", claims.ID, true).
		Preload("User").First(&session).Error; err != nil {
		return nil, errors.New("invalid or expired session")
	}

	accessToken, err := s.JwtService.GenerateAccessToken(time.Hour*4, utils.JwtCustomClaim{
		UserID: claims.UserID,
		Role:   claims.Role,
	})
	if err != nil {
		return nil, errors.New("could not generate access token")
	}

	refreshToken, err := s.JwtService.GenerateRefreshToken(time.Hour*24*7, utils.JwtCustomClaim{
		UserID: claims.UserID,
		Role:   claims.Role,
	}, session.ID.String())
	if err != nil {
		return nil, errors.New("could not generate refresh token")
	}

	return &pb.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "Token refreshed successfully",
	}, nil
}

func (s *AuthServiceServer) ConfirmEmail(ctx context.Context, req *pb.ConfirmEmailRequest) (*pb.ConfirmEmailResponse, error) {
	var verification models.Verification

	tx := database.DB.Begin()
	if err := tx.Where("token = ?", req.VerificationToken).First(&verification).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&models.User{}).Where("id = ?", verification.UserID).Update("is_email_verified", true).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Delete(&verification).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Where("user_id = ?", verification.UserID).Delete(&models.Verification{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &pb.ConfirmEmailResponse{Message: "Email verified successfully."}, nil
}

type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
	Logger       *zap.Logger
	JwtService   utils.JWTService
	EmailService utils.EmailService
}

func (s *AuthServiceServer) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return &pb.ForgotPasswordResponse{Message: "If your email is registered, you'll receive a password reset link"}, nil
	}

	token := uuid.New().String()
	reset := models.Reset{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := database.DB.Create(&reset).Error; err != nil {
		return nil, errors.New("could not create password reset")
	}

	type ResetData struct{ URL string }
	data := ResetData{
		URL: req.RedirectUrl + "?token=" + token,
	}

	if err := s.EmailService.SendEmail(user.Email, "Reset password", "reset_password", data); err != nil {
		return nil, errors.New("could not send reset email")
	}

	return &pb.ForgotPasswordResponse{Message: "If your email is registered, you'll receive a password reset link"}, nil
}

func (s *AuthServiceServer) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("could not start transaction")
	}

	var reset models.Reset
	if err := tx.Where("token = ?", req.ResetToken).Preload("User").First(&reset).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("invalid or expired reset token")
	}

	if reset.ExpiresAt.Before(time.Now()) {
		tx.Rollback()
		return nil, errors.New("reset token has expired")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		tx.Rollback()
		return &pb.ResetPasswordResponse{Message: "Could not hash password"}, err
	}

	if err := tx.Model(&reset.User).Update("password", hashedPassword).Error; err != nil {
		tx.Rollback()
		return &pb.ResetPasswordResponse{Message: "Could not update password"}, err
	}

	if err := tx.Delete(&reset).Error; err != nil {
		tx.Rollback()
		return &pb.ResetPasswordResponse{Message: "Could not delete reset token"}, err
	}

	if err := tx.Commit().Error; err != nil {
		return &pb.ResetPasswordResponse{Message: "Could not commit transaction"}, err
	}

	return &pb.ResetPasswordResponse{Message: "Password changed successfully"}, nil
}

func (s *AuthServiceServer) ClientInfo(ctx context.Context, req *pb.ClientInfoRequest) (*pb.ClientInfoResponse, error) {
	claims, err := s.JwtService.ValidateToken(req.Token)

	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	var user models.User

	if err := database.DB.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		return nil, errors.New("invalid token")
	}

	var phone string
	var avatar string

	if user.Phone != nil {
		phone = *user.Phone
	} else {
		phone = ""
	}

	if user.Avatar != nil {
		avatar = *user.Avatar
	} else {
		avatar = ""
	}

	return &pb.ClientInfoResponse{
		Email:    user.Email,
		Username: user.Username,
		Phone:    phone,
		Avatar:   avatar,
		Message:  "",
	}, nil
}
