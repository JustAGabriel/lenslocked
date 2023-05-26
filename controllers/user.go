package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/justagabriel/lenslocked/models"
	"github.com/justagabriel/lenslocked/util"
)

const (
	SessionCookieName = "lenslocked"
)

type UITemplates struct {
	New    template.Template
	Signin template.Template
}

type UserController struct {
	templates      UITemplates
	dbService      *models.UserService
	sessionService *models.SessionService
}

func NewUserController(uiTemplates UITemplates, dbService *models.UserService, sessionService *models.SessionService) UserController {
	return UserController{
		templates:      uiTemplates,
		dbService:      dbService,
		sessionService: sessionService,
	}
}

func (uc *UserController) GetUserBySessionToken(r *http.Request) (models.User, error) {
	sessionToken, err := util.GetSessionTokenFromCookie(SessionCookieName, r)
	if err != nil {
		return models.User{}, fmt.Errorf("error while trying to get session: %+v", err)
	}

	s, err2 := uc.sessionService.GetSessionByToken(sessionToken)
	if err2 != nil {
		return models.User{}, fmt.Errorf("error while trying to get session: %+v", err2)
	}

	u, err3 := uc.dbService.GetUserById(s.UserID)
	if err3 != nil {
		return models.User{}, fmt.Errorf("error while trying to get session: %+v", err3)
	}

	return u, nil
}

func (uc *UserController) GETSignup(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		CSRFField template.HTML
	}

	data := Data{
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
	type SigninData struct {
		CSRFField template.HTML
		Email     string
	}

	data := SigninData{
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
	err = util.SetCookie(w, SessionCookieName, s.Token)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, UserHomeURL, http.StatusFound)
}
