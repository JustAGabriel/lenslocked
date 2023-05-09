package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/justagabriel/lenslocked/models"
	"github.com/justagabriel/lenslocked/util"
)

type UITemplates struct {
	New    template.Template
	Signin template.Template
}

type UserController struct {
	templates UITemplates
	dbService *models.UserService
}

func NewUserController(uiTemplates UITemplates, dbService *models.UserService) UserController {
	return UserController{
		templates: uiTemplates,
		dbService: dbService,
	}
}

func (uc *UserController) New(w http.ResponseWriter, r *http.Request) {
	if err := uc.templates.New.Execute(w, nil); err != nil {
		panic(fmt.Errorf("while parsing 'new user' template: %v", err))
	}
}

func (u *UserController) Create(w http.ResponseWriter, r *http.Request) {
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

	fmt.Fprintf(w, "email: %s, pw: %s", email, pw)
}

func (uc *UserController) GETSignin(w http.ResponseWriter, r *http.Request) {
	if err := uc.templates.Signin.Execute(w, nil); err != nil {
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

	// todo: get user by email
	log.Default().Println("received creds: ", email, pwHash)
}
