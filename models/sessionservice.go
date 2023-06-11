package models

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/justagabriel/lenslocked/util"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	SessionTokenLength = 32
	SessionCookieName  = "lenslocked"
)

type SessionService struct {
	db          *gorm.DB
	userService *UserService
}

func hashToken(token string) string {
	hashedTokenBytes := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hashedTokenBytes[:])
}

func NewSessionService(db *gorm.DB, userService *UserService) SessionService {
	db.AutoMigrate(&Session{})
	return SessionService{
		db:          db,
		userService: userService,
	}
}

func (ss *SessionService) GetNewSession(userId uint) (*Session, error) {
	token := util.GetSessionToken()
	hashedToken := hashToken(token)
	s := &Session{
		UserID: userId,
		Token:  hashedToken,
	}

	_ = ss.db.Model(&Session{}).Create(s)

	return s, nil
}

func (ss *SessionService) GetSessionByToken(unhashedToken string) (Session, error) {
	hashedToken := hashToken(unhashedToken)
	session := Session{
		Token: hashedToken,
	}

	res := ss.db.Where(&session).First(&session)
	if res.Error != nil {
		return Session{}, res.Error
	}

	return session, nil
}

func (ss *SessionService) DeleteSessionByToken(token string) error {
	result := ss.db.Delete(&Session{Token: token})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (ss *SessionService) GetUserFromRequest(r *http.Request) (User, error) {
	sessionToken, err := util.GetSessionTokenFromCookie(SessionCookieName, r)
	if err != nil {
		return User{}, fmt.Errorf("error while trying to get session: %+v", err)
	}

	s, err2 := ss.GetSessionByToken(sessionToken)
	if err2 != nil {
		return User{}, fmt.Errorf("error while trying to get session: %+v", err2)
	}

	u, err3 := ss.userService.GetUserById(s.UserID)
	if err3 != nil {
		return User{}, fmt.Errorf("error while trying to get session: %+v", err3)
	}

	return u, nil
}

func (ss *SessionService) SetUserMiddleware(handler http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		user, err := ss.GetUserFromRequest(r)
		if err != nil {
			handler.ServeHTTP(w, r)
			return
		}

		ctx := WithUser(r.Context(), &user)
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}

type ctxKey string

const (
	key ctxKey = "user"
)

func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, key, user)
}

func GetUser(context context.Context) *User {
	val := context.Value(key)
	user, ok := val.(*User)
	if !ok {
		logger.Default.Warn(context, "could not get user from request context")
		return nil
	}

	return user
}
