package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"sso/interal/lib/jwt"
	"sso/interal/storage"
	"sso/interal/domain/models"
	
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)
type Auth struct {
	log *slog.Logger
	usrSaver  UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte)(uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string)(models.User, error)
	IsAdmin(ctx context.Context, userID int64)(bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int)(models.App, error)
}

// New creates a new instance of Auth with the given logger and user provider and app provider functions and returns a new instance of Auth
func New (log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log: log,
		usrSaver: userSaver,
		usrProvider: userProvider,
		appProvider: appProvider,
		tokenTTL: tokenTTL,
	}
}

//Login checks if the user is logged in and returns an error if it is not logged
func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("Attempting to authenticate user")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound){
			a.log.Warn("User not found", err)

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("Password does not match the hash of the user", err)
		
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully authenticated user")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to create token for user", err)
		
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

//RegisterNewUser register the user with the app provider and returns an error if the user is already registered
func (a *Auth) RegisterNewUser(ctx context.Context, email string, pass string)(int64, error){
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("Registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate pass hash for new user", err)
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", err)
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User is registered running hash environment")

	return id, nil
}

//IsAdmin checks if the user is an administrator and returns an error if it is not an administrator
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "Auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("userID", userID),
	)

	log.Info("Checking if user is an administrator for userID")
	
	

}
