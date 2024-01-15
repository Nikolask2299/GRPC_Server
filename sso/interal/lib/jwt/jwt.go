package jwt

import (
	"sso/interal/domain/models"

	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	clams := token.Claims.(jwt.MapClaims)

	clams["uid"] = user.ID
	clams["email"] = user.Email
	clams["exp"] = time.Now().Add(duration).Unix()
	clams["app_id"] = app.ID

	secr := []byte(app.Secret)
	tokenSting, err := token.SignedString(secr)
	if err != nil {
		return "", err
	}

	return tokenSting, nil
}