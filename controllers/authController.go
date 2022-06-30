package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"golang/config"
	"golang/helper"
	"golang/models"
	"golang/services"
	"golang/utils"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
)

type AuthController struct {
	authService services.AuthService
	userService services.UserService
	ctx         context.Context
	db          *gorm.DB
}

func NewAuthController(authService services.AuthService, userService services.UserService, ctx context.Context, db *gorm.DB) AuthController {
	return AuthController{authService, userService, ctx, db}
}

func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var user *models.SignUpInput

	if err := ctx.ShouldBindJSON(&user); err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if _, err := mail.ParseAddress(user.Email); err != nil {
		response := helper.BuildErrorResponse("Email is invalid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// Generate Verification Code
	code := randstr.String(20)
	user.VerificationCode = utils.Encode(code)

	newUser, err := ac.authService.SignUpUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			response := helper.BuildErrorResponse("name or email already exist", err.Error(), helper.EmptyObj{})
			ctx.JSON(http.StatusConflict, response)
			return
		}
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadGateway, response)
		return
	}

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Could not load config", err)
	}

	var firstName = newUser.Name
	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// 👇 Send Email
	emailData := utils.EmailData{
		URL:       config.BaseUrl + "/api/auth/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	err = utils.SendEmail(newUser, &emailData, "verificationCode.html")
	if err != nil {
		response := helper.BuildErrorResponse("There was an error sending email", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadGateway, response)
		return
	}

	message := "We sent an email with a verification code to " + user.Email
	response := helper.BuildResponse("OK", message)
	ctx.JSON(http.StatusCreated, response)
}

func (ac *AuthController) ResendVerification(ctx *gin.Context) {
	var userCredential *models.ResendVerificationInput

	if err := ctx.ShouldBindJSON(&userCredential); err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if _, err := mail.ParseAddress(userCredential.Email); err != nil {
		response := helper.BuildErrorResponse("Email is invalid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	user, err := ac.userService.FindUserByEmail(userCredential.Email)

	if user.Verified {
		response := helper.BuildErrorResponse("Account was verified", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	// Generate Verification Code
	code := randstr.String(20)
	user.VerificationCode = utils.Encode(code)

	if err := ac.db.Save(&user).Error; err != nil {
		response := helper.BuildErrorResponse("failed to process request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadGateway, response)
		return
	}

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Could not load config", err)
	}

	var firstName = user.Name
	// 👇 Send Email
	emailData := utils.EmailData{
		URL:       config.BaseUrl + "/api/auth/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	err = utils.SendEmail(user, &emailData, "verificationCode.html")
	if err != nil {
		response := helper.BuildErrorResponse("there was an error sending email", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadGateway, response)
		return
	}

	message := "We sent an email with a verification code to " + user.Email
	response := helper.BuildResponse("OK", message)
	ctx.JSON(http.StatusCreated, response)
}

func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var credentials *models.SignInInput

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		response := helper.BuildErrorResponse("failed to process request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	user, err := ac.userService.FindUserByEmail(credentials.Email)
	if err != nil {
		response := helper.BuildErrorResponse("there was an error sending email", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if _, err := mail.ParseAddress(credentials.Email); err != nil {
		response := helper.BuildErrorResponse("email is invalid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if !user.Verified {
		response := helper.BuildErrorResponse("you are not verified, please verify your email to login", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	if err := utils.VerifyPassword(user.Password, credentials.Password); err != nil {
		response := helper.BuildErrorResponse("invalid email or Password", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	config, _ := config.LoadConfig()

	// Generate Tokens
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		response := helper.BuildErrorResponse("error create access token", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		response := helper.BuildErrorResponse("error create refresh token", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", config.Domain, false, true)
	ctx.SetCookie("refresh_token", refresh_token, config.RefreshTokenMaxAge*60, "/", config.Domain, false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", config.Domain, false, false)

	result := make(map[string]string)
	result["access_token"] = access_token
	response := helper.BuildResponse("OK", result)
	ctx.JSON(http.StatusOK, response)
}

func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		response := helper.BuildErrorResponse("could not refresh access token", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	config, _ := config.LoadConfig()

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		response := helper.BuildErrorResponse("error validate token", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	user, err := ac.userService.FindUserById(fmt.Sprint(sub))
	if err != nil {
		response := helper.BuildErrorResponse("the user belonging to this token no logger exists", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		response := helper.BuildErrorResponse("error create token", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", config.Domain, false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", config.Domain, false, false)

	response := helper.BuildResponse("OK", access_token)
	ctx.JSON(http.StatusOK, response)
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Could not load config", err)
	}

	ctx.SetCookie("access_token", "", -1, "/", config.Domain, false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", config.Domain, false, true)
	ctx.SetCookie("logged_in", "", -1, "/", config.Domain, false, true)

	response := helper.BuildResponse("OK", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}

func (ac *AuthController) VerifyEmail(ctx *gin.Context) {
	code := ctx.Params.ByName("verificationCode")
	verificationCode := utils.Encode(code)

	var user *models.User

	if err := ac.db.Where("verification_code = ?", verificationCode).First(&user).Error; err != nil {
		response := helper.BuildErrorResponse("error find data", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	user.Verified = true
	user.VerificationCode = ""
	if err := ac.db.Save(&user).Error; err != nil {
		response := helper.BuildErrorResponse("failed to process request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadGateway, response)
		return
	}

	response := helper.BuildResponse("Email verified successfully", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}

func (ac *AuthController) ForgotPassword(ctx *gin.Context) {
	var userCredential *models.ForgotPasswordInput

	if err := ctx.ShouldBindJSON(&userCredential); err != nil {
		response := helper.BuildErrorResponse("failed to process request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if _, err := mail.ParseAddress(userCredential.Email); err != nil {
		response := helper.BuildErrorResponse("Email is invalid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	user, err := ac.userService.FindUserByEmail(userCredential.Email)
	if err != nil {
		response := helper.BuildErrorResponse("failed to process request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadGateway, response)
		return
	}

	if !user.Verified {
		response := helper.BuildErrorResponse("account not verified", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	config, _ := config.LoadConfig()

	// Generate Verification Code
	resetToken := randstr.String(20)
	user.PasswordResetToken = utils.Encode(resetToken)

	if err := ac.db.Save(&user).Error; err != nil {
		response := helper.BuildErrorResponse("failed to process request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadGateway, response)
		return
	}
	var firstName = user.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// 👇 Send Email
	emailData := utils.EmailData{
		URL:       config.BaseUrl + "/api/auth/resetpassword/" + resetToken,
		FirstName: firstName,
		Subject:   "Your password reset token",
	}

	err = utils.SendEmail(user, &emailData, "resetPassword.html")
	if err != nil {
		response := helper.BuildErrorResponse("there was an error sending email", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadGateway, response)
		return
	}
	response := helper.BuildResponse("You will receive a reset email if user with that email exist", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}

func (ac *AuthController) ResendForgotPassword(ctx *gin.Context) {
	var userCredential *models.ForgotPasswordInput

	if err := ctx.ShouldBindJSON(&userCredential); err != nil {
		response := helper.BuildErrorResponse("failed to process request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if _, err := mail.ParseAddress(userCredential.Email); err != nil {
		response := helper.BuildErrorResponse("email is invalid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	user, err := ac.userService.FindUserByEmail(userCredential.Email)

	if !user.Verified {
		response := helper.BuildErrorResponse("account not verified", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	// Generate Verification Code
	resetToken := randstr.String(20)
	user.PasswordResetToken = utils.Encode(resetToken)

	if err := ac.db.Save(&user).Error; err != nil {
		response := helper.BuildErrorResponse("failed to process request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadGateway, response)
		return
	}

	config, _ := config.LoadConfig()

	var firstName = user.Name
	// 👇 Send Email
	emailData := utils.EmailData{
		URL:       config.BaseUrl + "/api/auth/resetpassword/" + resetToken,
		FirstName: firstName,
		Subject:   "Your password reset token",
	}

	err = utils.SendEmail(user, &emailData, "verificationCode.html")
	if err != nil {
		response := helper.BuildErrorResponse("there was an error sending email", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadGateway, response)
		return
	}

	message := "We sent an email with a verification code to " + user.Email
	response := helper.BuildResponse(message, helper.EmptyObj{})
	ctx.JSON(http.StatusCreated, response)
}

func (ac *AuthController) ResetPassword(ctx *gin.Context) {
	var userCredential *models.ResetPasswordInput
	resetToken := ctx.Params.ByName("resetToken")

	if err := ctx.ShouldBindJSON(&userCredential); err != nil {
		response := helper.BuildErrorResponse("failed to process request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	hashedPassword, _ := utils.HashPassword(userCredential.Password)
	passwordResetToken := utils.Encode(resetToken)

	var user *models.User

	if err := ac.db.Where("password_reset_token = ?", passwordResetToken).First(&user).Error; err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadGateway, response)
		return
	}

	user.Password = hashedPassword
	user.PasswordResetToken = ""
	if err := ac.db.Save(&user).Error; err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadGateway, response)
		return
	}

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Could not load config", err)
	}

	ctx.SetCookie("access_token", "", -1, "/", config.Domain, false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", config.Domain, false, true)
	ctx.SetCookie("logged_in", "", -1, "/", config.Domain, false, true)

	response := helper.BuildResponse("Password data updated successfully", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}
