package rpc

import (
	"context"
	"fmt"
	proto "github.com/FireFly4ik/Lavka-auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	api proto.AuthServiceClient
}

func NewAuth(addr string) (*AuthClient, error) {
	const op = "auth.New"

	cc, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, fmt.Errorf("%s: dial: %w", op, err)
	}

	return &AuthClient{api: proto.NewAuthServiceClient(cc)}, nil
}

func (c *AuthClient) RefreshRequest(ctx context.Context, token string) (string, string, string, error) {
	const op = "auth.RefreshRequest"

	resp, err := c.api.Refresh(ctx, &proto.RefreshRequest{
		Token: token,
	})
	if err != nil {
		return "", "", "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.RefreshToken, resp.AccessToken, resp.Message, nil
}

func (c *AuthClient) LoginRequest(ctx context.Context, login string, password string, userAgent string, ip string) (string, string, string, error) {
	const op = "auth.LoginRequest"

	resp, err := c.api.Login(ctx, &proto.LoginRequest{
		Login:     login,
		Password:  password,
		UserAgent: userAgent,
		IpAddress: ip,
	})
	if err != nil {
		return "", "", "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.AccessToken, resp.RefreshToken, resp.Message, nil
}

func (c *AuthClient) SSOLoginRequest(ctx context.Context, ssoToken string) (string, string, string, error) {
	const op = "auth.SSOLoginRequest"

	resp, err := c.api.SSOLogin(ctx, &proto.SSOLoginRequest{
		SsoToken: ssoToken,
	})
	if err != nil {
		return "", "", "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.AccessToken, resp.RefreshToken, resp.Message, nil
}

func (c *AuthClient) RegisterRequest(ctx context.Context, username string, email string, password string, redirectUrl string) (string, error) {
	const op = "auth.RegisterRequest"

	resp, err := c.api.Register(ctx, &proto.RegisterRequest{
		Username:    username,
		Email:       email,
		Password:    password,
		RedirectUrl: redirectUrl,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Message, nil
}

func (c *AuthClient) ResetPasswordRequest(ctx context.Context, verificationCode string, password string) (string, error) {
	const op = "auth.ResetPasswordRequest"

	resp, err := c.api.ResetPassword(ctx, &proto.ResetPasswordRequest{
		ResetToken: verificationCode,
		Password:   password,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Message, nil
}

func (c *AuthClient) ConfirmEmailRequest(ctx context.Context, verificationCode string) (string, error) {
	const op = "auth.ConfirmEmailRequest"

	resp, err := c.api.ConfirmEmail(ctx, &proto.ConfirmEmailRequest{
		VerificationToken: verificationCode,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Message, nil
}

func (c *AuthClient) ForgotPasswordRequest(ctx context.Context, email string, redirectUrl string) (string, error) {
	const op = "auth.ForgotPasswordRequest"

	resp, err := c.api.ForgotPassword(ctx, &proto.ForgotPasswordRequest{
		Email:       email,
		RedirectUrl: redirectUrl,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Message, nil
}

func (c *AuthClient) LogoutRequest(ctx context.Context, token string) (string, error) {
	const op = "auth.LogoutRequest"

	resp, err := c.api.Logout(ctx, &proto.LogoutRequest{
		RefreshToken: token,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Message, nil
}

func (c *AuthClient) ClientInfoRequest(ctx context.Context, token string) (string, string, *string, *string, string, error) {
	const op = "auth.ClientInfoRequest"

	resp, err := c.api.ClientInfo(ctx, &proto.ClientInfoRequest{
		Token: token,
	})
	if err != nil {
		return "", "", nil, nil, "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Email, resp.Username, &resp.Phone, &resp.Avatar, resp.Message, nil
}
