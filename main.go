package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/justagabriel/lenslocked/controllers"
	"github.com/justagabriel/lenslocked/models"
	"github.com/justagabriel/lenslocked/util"
	"github.com/justagabriel/lenslocked/views"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	baseLayoutFilename = "tailwind"
)

func main() {
	db, err := gorm.Open(postgres.Open(models.GetDefaultDBConfig().String()))
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(false),
	)
	r.Use(csrfMw)

	tpl := util.Must(views.ParseFS(views.FS, "home", baseLayoutFilename))
	r.Get("/home", controllers.StaticHandler(tpl))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = util.Must(views.ParseFS(views.FS, "contact", baseLayoutFilename))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = util.Must(views.ParseFS(views.FS, "faq", baseLayoutFilename))
	r.Get("/faq", controllers.FAQ(tpl))

	tpl = util.Must(views.ParseFS(views.FS, "signup", baseLayoutFilename))
	tpl2 := util.Must(views.ParseFS(views.FS, "signin", baseLayoutFilename))
	templates := controllers.UITemplates{
		New:    *tpl,
		Signin: *tpl2,
	}

	userService := models.NewUserService(db)
	usersC := controllers.NewUserController(templates, &userService)
	r.Get("/signup", usersC.GETSignup)
	r.Post("/signup", usersC.POSTSignup)

	r.Get("/signin", usersC.GETSignin)
	r.Post("/signin", usersC.POSTSignin)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found, dude!", http.StatusNotFound)
	})

	http.ListenAndServe("localhost:8000", r)
}
