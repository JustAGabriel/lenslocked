package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/justagabriel/lenslocked/controllers"
	"github.com/justagabriel/lenslocked/models"
	"github.com/justagabriel/lenslocked/models/migrations"
	"github.com/justagabriel/lenslocked/util"
	"github.com/justagabriel/lenslocked/views"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	baseLayoutFilename = "tailwind"
)

func main() {
	// perform database migration
	db, err := gorm.Open(postgres.Open(models.GetDefaultDBConfig().String()))
	if err != nil {
		panic(err)
	}

	err = models.Migrate(db, migrations.FS)
	if err != nil {
		panic(err)
	}

	// create static UI templates
	homeTemplate := util.Must(views.ParseFS(views.FS, "home", baseLayoutFilename))
	contactTemplate := util.Must(views.ParseFS(views.FS, "contact", baseLayoutFilename))
	FaqTemplate := util.Must(views.ParseFS(views.FS, "faq", baseLayoutFilename))
	signupTemplate := util.Must(views.ParseFS(views.FS, "signup", baseLayoutFilename))
	signinTemplate := util.Must(views.ParseFS(views.FS, "signin", baseLayoutFilename))
	forgotPwTemplate := util.Must(views.ParseFS(views.FS, "forgot-pw", baseLayoutFilename))
	checkYourEmailTemplate := util.Must(views.ParseFS(views.FS, "check-your-email", baseLayoutFilename))
	templates := controllers.UITemplates{
		New:            *signupTemplate,
		Signin:         *signinTemplate,
		ForgotPassword: *forgotPwTemplate,
		CheckYourEmail: *checkYourEmailTemplate,
	}

	// initialize services
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Default().Printf("error while reading '.env': %v\n", err)
	}

	email_host := viper.GetString("mail_host")
	port_str := viper.GetString("mail_port")
	email_port, _ := strconv.Atoi(port_str)
	email_username := viper.GetString("mail_username")
	email_pw := viper.GetString("mail_pw")

	smtpConfig := models.SMTPConfig{
		Host:     email_host,
		Port:     email_port,
		Username: email_username,
		Password: email_pw,
	}

	emailService, err := models.NewEmailService(smtpConfig, models.DefaultSender)
	if err != nil {
		fmt.Println(err)
		return
	}

	userService := models.NewUserService(db)
	pwResetService := models.NewPasswordResetService(&userService, db)
	sessionService := models.NewSessionService(db, &userService)
	userController := controllers.NewUserController(templates, &userService, &sessionService, pwResetService, emailService)

	// register middlewares
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(false),
	)
	r.Use(csrfMw)

	userMiddleware := controllers.NewUserMiddleware(&sessionService)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(userMiddleware.SetUserMiddleware)
		r.Use(userMiddleware.RequireUserMiddleware)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "not implemented :D") // todo: implement user home page
		})
	})

	// register routes
	r.Get(controllers.WebsiteHomeURL, controllers.StaticHandler(homeTemplate))
	r.Get(controllers.WebsiteRootURL, controllers.StaticHandler(homeTemplate))
	r.Get(controllers.WebsiteContactURL, controllers.StaticHandler(contactTemplate))
	r.Get(controllers.WebsiteFaqURL, controllers.FAQ(FaqTemplate))

	r.Get(controllers.SignupURL, userController.GETSignup)
	r.Post(controllers.SignupURL, userController.POSTSignup)

	r.Get(controllers.SigninURL, userController.GETSignin)
	r.Post(controllers.SigninURL, userController.POSTSignin)
	r.Get(controllers.SignoutURL, userController.GETSignout)

	r.Get(controllers.ForgotPasswordURL, userController.GETForgotPassword)
	r.Post(controllers.ForgotPasswordURL, userController.POSTForgotPassword)

	// todo: handle GET/POST "reset-pw" to allow setting a new password via email link

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found, dude!", http.StatusNotFound)
	})

	http.ListenAndServe("localhost:8000", r)
}
