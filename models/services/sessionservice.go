package models

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/justagabriel/lenslocked/models"
	"github.com/justagabriel/lenslocked/util"
	"gorm.io/gorm"
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
	db.AutoMigrate(&models.Session{})
	return SessionService{
		db:          db,
		userService: userService,
	}
}

func (ss *SessionService) GetNewSession(userId uint) (*models.Session, error) {
	// implement deletion of previous sessions associated with the given user id
	// usecase: user cookie was deleted -> auth will fail, new user session created -- conflict in db since
	// userid already is meantioned in other sesssion entry
	userExistsQueryResult := ss.db.Model(&models.Session{}).First(&models.Session{UserID: userId})
	hasUserExistingSession := userExistsQueryResult.Error == nil && userExistsQueryResult.RowsAffected == 1
	if hasUserExistingSession {
		userSessionDeletionResult := ss.db.Unscoped().Model(&models.Session{}).Where(&models.Session{UserID: userId}).Delete(&models.Session{})
		if userSessionDeletionResult.Error != nil {
			log.Default().Printf("error while deleting session: %s", userSessionDeletionResult.Error)
		}
	}

	token := util.GetSessionToken()
	hashedToken := hashToken(token)
	s := &models.Session{
		UserID: userId,
		Token:  hashedToken,
	}

	creation_transaction := ss.db.Model(&models.Session{}).Create(s)
	if creation_transaction.Error != nil {
		return nil, fmt.Errorf("error while creating new session: %s", creation_transaction.Error)
	}

	s.Token = token
	return s, nil
}

func (ss *SessionService) GetSessionByToken(unhashedToken string) (models.Session, error) {
	hashedToken := hashToken(unhashedToken)
	session := models.Session{
		Token: hashedToken,
	}

	res := ss.db.Model(&models.Session{}).First(&session)
	if res.Error != nil {
		return models.Session{}, res.Error
	}

	return session, nil
}

func (ss *SessionService) DeleteSessionByToken(token string) error {
	result := ss.db.Delete(&models.Session{Token: token})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (ss *SessionService) GetUserFromRequest(r *http.Request) (models.User, error) {
	sessionToken, err := util.GetSessionTokenFromCookie(SessionCookieName, r)
	if err != nil {
		return models.User{}, fmt.Errorf("error while trying to get session: %+v", err)
	}

	s, err2 := ss.GetSessionByToken(sessionToken)
	if err2 != nil {
		return models.User{}, fmt.Errorf("error while trying to get session: %+v", err2)
	}

	u, err3 := ss.userService.GetUserById(s.UserID)
	if err3 != nil {
		return models.User{}, fmt.Errorf("error while trying to get session: %+v", err3)
	}

	return u, nil
}
