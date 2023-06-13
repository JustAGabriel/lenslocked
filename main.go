package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/justagabriel/lenslocked/controllers"
	"github.com/justagabriel/lenslocked/models"
	"github.com/justagabriel/lenslocked/models/migrations"
	"github.com/justagabriel/lenslocked/util"
	"github.com/justagabriel/lenslocked/views"
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
	templates := controllers.UITemplates{
		New:    *signupTemplate,
		Signin: *signinTemplate,
	}

	// register middlewares
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	userService := models.NewUserService(db)
	sessionService := models.NewSessionService(db, &userService)
	userMiddleware := controllers.NewUserMiddleware(&sessionService)
	userController := controllers.NewUserController(templates, &userService, &sessionService)

	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(false),
	)
	r.Use(csrfMw)

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

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found, dude!", http.StatusNotFound)
	})

	// todo:
	// login -> why 'require auth' meachnism not working? (db error + general logic should be tested)
	// general, debugging: master debug configs in vs code

	http.ListenAndServe("localhost:8000", r)
}
