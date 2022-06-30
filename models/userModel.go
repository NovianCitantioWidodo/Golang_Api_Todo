package models

import (
	"time"
)

type User struct {
	ID                 int       `json:"id,omitempty" bson:"_id,omitempty"`
	Name               string    `json:"name,omitempty" bson:"name,omitempty"`
	Email              string    `json:"email,omitempty" bson:"email,omitempty"`
	Password           string    `json:"password" bson:"password" binding:"required,min=8"`
	VerificationCode   string    `json:"verificationCode,omitempty" bson:"verificationCode"`
	PasswordResetToken string    `json:"passwordResetToken,omitempty" bson:"passwordResetToken,omitempty"`
	Role               string    `json:"role,omitempty" bson:"role,omitempty"`
	Verified           bool      `json:"verified" bson:"verified"`
	CreatedAt          time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" bson:"updated_at"`
}

func (User) TableName() string {
	return "user"
}

// 👈 SignUpInput struct
type SignUpInput struct {
	Name             string    `json:"name" bson:"name" binding:"required"`
	Email            string    `json:"email" bson:"email" binding:"required"`
	Password         string    `json:"password" bson:"password" binding:"required,min=8"`
	PasswordConfirm  string    `json:"passwordConfirm" bson:"passwordConfirm,omitempty" binding:"required"`
	VerificationCode string    `json:"verificationCode,omitempty" bson:"verificationCode,omitempty"`
	Role             string    `json:"role" bson:"role"`
	Verified         bool      `json:"verified" bson:"verified"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
}

// 👈 SignInInput struct
type SignInInput struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

// 👈 UserResponse struct
type UserResponse struct {
	ID        int       `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string    `json:"name,omitempty" bson:"name,omitempty"`
	Email     string    `json:"email,omitempty" bson:"email,omitempty"`
	Role      string    `json:"role,omitempty" bson:"role,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// 👈 ForgotPasswordInput struct
type ResendVerificationInput struct {
	Email string `json:"email" binding:"required"`
}

// 👈 ForgotPasswordInput struct
type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required"`
}

// 👈 ResetPasswordInput struct
type ResetPasswordInput struct {
	Password string `json:"password" binding:"required"`
}

type UserEdit struct {
	Name     string `json:"name" bson:"name" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required,min=8"`
}

func FilteredResponse(user *User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
