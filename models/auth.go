package models

import "time"

type UserRegisterRequest struct {
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	UserID       int64                      `json:"userID"`
	RoleName     string                     `json:"roleName"`
	AccessToken  string                     `json:"accessToken"`
	RefreshToken string                     `json:"refreshToken"`
	Permission   []UserRolePermissionModels `json:"permission"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type FindUserRoleResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	RoleID   int    `json:"roleID"`
	RoleName string `json:"roleName"`
}

type CurrentUserModels struct {
	ID    int64  `json:"id"`
	Role  string `json:"role"`
	Email string `json:"email"`
}

type UserGenerateOTPRequest struct {
	Phone string `json:"phone"`
}

type OTPModels struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"userID"`
	OtpHash   string    `json:"otpHash"`
	CreatedAt string    `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	UsedAt    time.Time `json:"usedAt"`
	IsUsed    bool      `json:"isUsed"`
}

type UserValidateOtpRequest struct {
	UserID  int64  `json:"userID"`
	OtpHash string `json:"otpHash"`
}

type UserFindUserByEmailRequest struct {
	Email string `json:"email"`
}
