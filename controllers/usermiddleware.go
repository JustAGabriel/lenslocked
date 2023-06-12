package controllers

import (
	"context"
	"net/http"

	"github.com/justagabriel/lenslocked/models"
	"gorm.io/gorm/logger"
)

type UserMiddleware struct {
	sessionService *models.SessionService
}

type ctxKey string

const (
	key ctxKey = "user"
)

func NewUserMiddleware(sessionService *models.SessionService) *UserMiddleware {
	return &UserMiddleware{
		sessionService: sessionService,
	}
}

func withUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, key, user)
}

func getUser(context context.Context) *models.User {
	val := context.Value(key)
	user, ok := val.(*models.User)
	if !ok {
		logger.Default.Warn(context, "could not get user from request context")
		return nil
	}

	return user
}

func (um *UserMiddleware) RequireUserMiddleware(handler http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		u := getUser(r.Context())
		if u != nil {
			http.Redirect(w, r, "/signin", http.StatusUnauthorized)
		}
		handler.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}

func (um *UserMiddleware) SetUserMiddleware(handler http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		user, err := um.sessionService.GetUserFromRequest(r)
		if err != nil {
			handler.ServeHTTP(w, r)
			return
		}

		ctx := withUser(r.Context(), &user)
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}
