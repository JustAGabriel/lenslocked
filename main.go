package main

import (
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

	// register middlewares
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	userService := models.NewUserService(db)
	sessionService := models.NewSessionService(db, &userService)
	userMiddleware := controllers.NewUserMiddleware(&sessionService)
	r.Use(userMiddleware.SetUserMiddleware)

	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(false),
	)
	r.Use(csrfMw)

	// register endpoints
	tpl := util.Must(views.ParseFS(views.FS, "home", baseLayoutFilename))
	r.Get(controllers.WebsiteHomeURL, controllers.StaticHandler(tpl))
	r.Get(controllers.WebsiteRootURL, controllers.StaticHandler(tpl))

	tpl = util.Must(views.ParseFS(views.FS, "contact", baseLayoutFilename))
	r.Get(controllers.WebsiteContactURL, controllers.StaticHandler(tpl))

	tpl = util.Must(views.ParseFS(views.FS, "faq", baseLayoutFilename))
	r.Get(controllers.WebsiteFaqURL, controllers.FAQ(tpl))

	tpl = util.Must(views.ParseFS(views.FS, "signup", baseLayoutFilename))
	tpl2 := util.Must(views.ParseFS(views.FS, "signin", baseLayoutFilename))
	templates := controllers.UITemplates{
		New:    *tpl,
		Signin: *tpl2,
	}

	usersC := controllers.NewUserController(templates, &userService, &sessionService)
	r.Get(controllers.SignupURL, usersC.GETSignup)
	r.Post(controllers.SignupURL, usersC.POSTSignup)

	r.Get(controllers.SigninURL, usersC.GETSignin)
	r.Post(controllers.SigninURL, usersC.POSTSignin)
	r.Get(controllers.SignoutURL, usersC.GETSignout)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found, dude!", http.StatusNotFound)
	})

	http.ListenAndServe("localhost:8000", r)
}
