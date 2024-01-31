package main

import (
	"amikom-pedia-api/app"
	"amikom-pedia-api/helper"
	"amikom-pedia-api/module/image/image_repository"
	"amikom-pedia-api/module/image/image_service"
	"amikom-pedia-api/module/login/login_controller"
	"amikom-pedia-api/module/login/login_repository"
	"amikom-pedia-api/module/login/login_service"
	"amikom-pedia-api/module/otp/otp_controller"
	"amikom-pedia-api/module/otp/otp_repository"
	"amikom-pedia-api/module/otp/otp_service"
	"amikom-pedia-api/module/register/register_controller"
	"amikom-pedia-api/module/register/register_repository"
	"amikom-pedia-api/module/register/register_service"
	"amikom-pedia-api/module/user/user_controller"
	"amikom-pedia-api/module/user/user_repository"
	"amikom-pedia-api/module/user/user_service"
	"amikom-pedia-api/utils"
	"amikom-pedia-api/utils/aws"
	"amikom-pedia-api/utils/mail"
	"amikom-pedia-api/utils/token"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
)

func main() {

	config, err := utils.LoadConfig(".")
	helper.PanicIfError(err)
	//jwt init
	tokenMaker, err := token.NewJWTMaker(config.TokenSymetricKey)
	helper.PanicIfError(err)

	gmailSender := mail.NewGmailSender(config.EmailName, config.EmailSender, config.EmailPassword)
	db := app.NewDB(config.DBDriver, config.DBSource)
	validate := utils.CustomValidator()

	awsSession, err := aws.NewSessionAWSS3(config)
	helper.PanicIfError(err)

	sesS3 := aws.NewAwsS3(awsSession, config.AWSS3Bucket)

	//REPOSITORY
	userRepository := user_repository.NewUserRepository()
	registerRepository := register_repository.NewRegisterRepository()
	otpRepository := otp_repository.NewOtpRepository()
	loginRepository := login_repository.NewLoginRepository()
	imageRepository := image_repository.NewImageRepository()

	//SERVICE
	userService := user_service.NewUserService(userRepository, otpRepository, gmailSender, db, validate)
	registerService := register_service.NewRegisterService(registerRepository, otpRepository, gmailSender, db, validate)
	loginService := login_service.NewLoginService(tokenMaker, loginRepository, db, validate)
	otpService := otp_service.NewOtpService(otpRepository, registerRepository, userRepository, gmailSender, db, validate, tokenMaker)
	imageService := image_service.NewImageService(imageRepository, db, sesS3)

	//CONTROLLER
	userController := user_controller.NewUserController(userService, imageService)
	registerController := register_controller.NewRegisterController(registerService)
	loginController := login_controller.NewLoginController(loginService)
	otpController := otp_controller.NewOtpController(otpService)

	router := app.NewRouter(tokenMaker, userController, registerController, otpController, loginController)

	server := http.Server{
		Addr:    config.ServerAddress,
		Handler: router,
	}

	err = server.ListenAndServe()
	helper.PanicIfError(err)
	fmt.Println("Starting web server at port 3000")
}
