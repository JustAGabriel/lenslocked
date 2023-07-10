package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/csrf"
	"github.com/justagabriel/lenslocked/models"
	"github.com/justagabriel/lenslocked/util"
	"gorm.io/gorm/logger"
)

type UITemplates struct {
	New            template.Template
	Signin         template.Template
	ForgotPassword template.Template
	CheckYourEmail template.Template
}

type UserController struct {
	templates            UITemplates
	dbService            *models.UserService
	sessionService       *models.SessionService
	passwordResetService *models.PasswordResetService
	emailService         *models.EmailService
}

type templateData struct {
	TemplateBaseData
	Email     string
	CSRFField template.HTML
}

func NewUserController(uiTemplates UITemplates, dbService *models.UserService, sessionService *models.SessionService,
	pwResetService *models.PasswordResetService, emailService *models.EmailService) UserController {
	return UserController{
		templates:            uiTemplates,
		dbService:            dbService,
		sessionService:       sessionService,
		passwordResetService: pwResetService,
		emailService:         emailService,
	}
}

func (uc *UserController) GETSignup(w http.ResponseWriter, r *http.Request) {
	data := templateData{
		CSRFField: csrf.TemplateField(r),
	}

	if err := uc.templates.New.Execute(w, data); err != nil {
		panic(fmt.Errorf("while parsing 'new user' template: %v", err))
	}
}

func (u *UserController) POSTSignup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
	}

	email := r.Form.Get("email")
	pw := r.Form.Get("password")

	pwHash, err := util.GetPasswordHash(pw)
	if err != nil {
		log.Fatal(err)
	}

	newUser := models.User{
		Email:        email,
		PasswordHash: pwHash,
	}
	u.dbService.CreateUser(newUser)

	http.Redirect(w, r, SigninURL, http.StatusFound)
}

func (uc *UserController) GETSignin(w http.ResponseWriter, r *http.Request) {
	data := templateData{
		CSRFField: csrf.TemplateField(r),
		Email:     "",
	}

	if err := uc.templates.Signin.Execute(w, data); err != nil {
		panic(fmt.Errorf("while parsing 'signin' template: %v", err))
	}
}

func (uc *UserController) POSTSignin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
	}

	email := r.Form.Get("email")
	pw := r.Form.Get("password")

	pwHash, err := util.GetPasswordHash(pw)
	if err != nil {
		log.Fatal(err)
	}
	log.Default().Println("received creds: ", email, pwHash)

	usr, err := uc.dbService.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "Username or password did not match", http.StatusUnauthorized)
	}
	log.Default().Printf("found user:\n%+v", usr)

	if err := util.ComparePwAndPwHash(pw, usr.PasswordHash); err != nil {
		http.Error(w, "Username or password did not match", http.StatusUnauthorized)
	}

	s, err := uc.sessionService.GetNewSession(usr.ID)
	if err != nil {
		log.Fatal(err)
	}
	err = util.SetCookie(w, models.SessionCookieName, s.Token)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, UserHomeURL, http.StatusFound)
}

func (uc *UserController) GETSignout(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := util.GetSessionTokenFromCookie(models.SessionCookieName, r)
	if err != nil {
		log.Default().Printf("error while trying to get session: %+v\n", err)
	}

	err2 := uc.sessionService.DeleteSessionByToken(sessionToken)
	if err2 != nil {
		log.Default().Printf("error while trying to delete the session: %+v\n", err2)
	}

	err3 := util.DeleteCookie(w, models.SessionCookieName)
	if err3 != nil {
		log.Default().Printf("error while trying to eelete the session cookie: %+v\n", err3)
	}

	http.Redirect(w, r, WebsiteHomeURL, http.StatusFound)
}

func (uc *UserController) GETForgotPassword(w http.ResponseWriter, r *http.Request) {
	data := templateData{
		Email:     r.FormValue("email"),
		CSRFField: csrf.TemplateField(r),
	}
	err := uc.templates.ForgotPassword.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func (uc UserController) POSTForgotPassword(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	pwReset, err := uc.passwordResetService.GetPasswordReset(email)
	if err != nil {
		// TODO: handle email does not exist
		logger.Default.Error(r.Context(), err.Error())
		http.Error(w, "could not create pw reset", http.StatusInternalServerError)
		return
	}

	urlQuery := url.Values{
		"token": {pwReset.Token},
	}
	// todo: make url configurable
	serverURL := "http://localhost:8000"
	resetURL := serverURL + "/reset-pw?" + urlQuery.Encode()
	uc.emailService.SendForgotPasswordEmail(email, resetURL)

	data := templateData{Email: email}
	uc.templates.CheckYourEmail.Execute(w, data)
}
