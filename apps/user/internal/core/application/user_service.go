package application

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/subscription-tracker/user/internal/core/domain"
	"golang.org/x/oauth2"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid token")
	ErrSpotifyAuth        = errors.New("spotify authentication failed")
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id string) error
}

type EmailService interface {
	Send(to, subject, message string) error
}

type SpotifyUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"display_name"`
}

type UserService struct {
	repo          UserRepository
	emailService  EmailService
	spotifyConfig oauth2.Config
}

func (s *UserService) GenerateOAuthState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func NewUserService(repo UserRepository, emailService EmailService) *UserService {
	return &UserService{
		repo:         repo,
		emailService: emailService,
		spotifyConfig: oauth2.Config{
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://accounts.spotify.com/authorize",
				TokenURL: "https://accounts.spotify.com/api/token",
			},
			ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("SPOTIFY_REDIRECT_URL"),
			Scopes:       []string{"user-read-email", "user-read-private"},
		},
	}
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}

func (s *UserService) Register(email, name, password string) (*AuthResponse, error) {
	// Check if user already exists
	if _, err := s.repo.FindByEmail(email); err == nil {
		return nil, ErrUserAlreadyExists
	}

	// Create new user
	user, err := domain.NewUser(email, name, password)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: *user}, nil
}

func (s *UserService) Login(email, password string) (*AuthResponse, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := user.ComparePassword(password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: *user}, nil
}

func (s *UserService) GetUserByID(id string) (*domain.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) generateToken(user *domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (s *UserService) UpdateUser(id string, name string, password string) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return ErrUserNotFound
	}
	user.Name = name

	if password != "" {
		if err := user.HashedPassword(password); err != nil {
			return err
		}
	}

	return s.repo.Update(user)
}

func (s *UserService) DeleteUser(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return ErrUserNotFound
	}

	return s.repo.Delete(id)
}

func (s *UserService) InitiatePasswordlessLogin(email string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return ErrUserNotFound
	}

	// Generate login token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
		"type":    "passwordless",
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return err
	}

	loginLink := fmt.Sprintf("http://localhost:3000/login/callback/passwordless?token=%s", tokenString)

	// Send login email
	err = s.emailService.Send(email, "Passwordless Login", fmt.Sprintf("Click the link to login: %s", loginLink))
	if err != nil {
		return fmt.Errorf("failed to send login email: %w", err)
	}

	return nil
}

func (s *UserService) GetSpotifyAuthURL(state string) string {
	return s.spotifyConfig.AuthCodeURL(state)
}

func (s *UserService) HandleSpotifyCallback(code string) (*AuthResponse, error) {
	token, err := s.spotifyConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSpotifyAuth, err)
	}

	client := s.spotifyConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.spotify.com/v1/me")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSpotifyAuth, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSpotifyAuth, err)
	}

	var spotifyUser SpotifyUser
	if err := json.Unmarshal(body, &spotifyUser); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSpotifyAuth, err)
	}

	// Check if user exists
	user, err := s.repo.FindByEmail(spotifyUser.Email)
	if err != nil {
		// Create new user if not exists
		user, err = domain.NewUser(spotifyUser.Email, spotifyUser.Name, "")
		if err != nil {
			return nil, err
		}
		if err := s.repo.Create(user); err != nil {
			return nil, err
		}
	}
	// Update user's Spotify ID
	user.SpotifyCredentials = domain.SpotifyCredentials{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    int(token.ExpiresIn),
	}
	s.repo.Update(user)

	// Generate JWT token
	jwtToken, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: jwtToken, User: *user}, nil
}

func (s *UserService) VerifyLoginToken(tokenString string) (*AuthResponse, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Generate a new regular session token
	newToken, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: newToken, User: *user}, nil
}
