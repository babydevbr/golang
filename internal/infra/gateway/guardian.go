// Package gateway ...
package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/babydevbr/golang/internal/domain"
	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/shaj13/go-guardian/v2/auth/strategies/basic"
	"github.com/shaj13/go-guardian/v2/auth/strategies/jwt"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
	"github.com/shaj13/libcache"

	// guard libcache.
	_ "github.com/shaj13/libcache/fifo"
)

type userService interface {
	GetByEmail(email string) (*domain.User, error)
}

type passwordService interface {
	Encrypt(raw string) string
	Compare(hash, s string) error
}

// Guardian middleware.
type Guardian struct {
	whitelist       map[string]string
	keeper          jwt.SecretsKeeper
	strategy        union.Union
	cache           libcache.Cache
	userService     userService
	passwordService passwordService
}

// NewGuardian Guardian factory.
func NewGuardian(keeper jwt.SecretsKeeper, u userService, p passwordService) *Guardian {
	duration := 5
	cache := libcache.FIFO.New(0)
	cache.SetTTL(time.Minute * time.Duration(duration))
	cache.RegisterOnExpired(func(key, _ interface{}) {
		cache.Peek(key)
	})

	guard := &Guardian{
		keeper:          keeper,
		cache:           cache,
		whitelist:       map[string]string{"/users/signup": "POST"},
		userService:     u,
		passwordService: p,
	}

	basicStrategy := basic.NewCached(guard.ValidateEmailAndPassword(), cache)
	jwtStrategy := jwt.New(cache, keeper)
	strategy := union.New(jwtStrategy, basicStrategy)

	guard.strategy = strategy

	return guard
}

// ValidateEmailAndPassword Compare Email and Password.
func (am *Guardian) ValidateEmailAndPassword() basic.AuthenticateFunc {
	return func(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
		e, err := am.userService.GetByEmail(userName)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		err = am.passwordService.Compare(e.Password, password)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		return auth.NewDefaultUser(e.Name, e.ID, []string{e.Role}, nil), nil
	}
}

// GenerateToken return a JWT token.
func (am *Guardian) GenerateToken(w http.ResponseWriter, r *http.Request) {
	u := auth.User(r)
	exp := jwt.SetExpDuration(time.Hour)

	token, err := jwt.IssueAccessToken(u, am.keeper, exp)
	if err != nil {
		log.Printf("Error %s\n", err.Error())

		_, _ = w.Write([]byte(err.Error()))

		return
	}

	body := fmt.Sprintf(`{"name":"%s","token":"%s"}`, u.GetUserName(), token)
	_, _ = w.Write([]byte(body))
}

// AuthMiddleware Authentication middleware.
func (am *Guardian) AuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if _, ok := am.whitelist[r.URL.Path]; ok {
		next.ServeHTTP(w, r)

		return
	}

	log.Println("Executing Auth Middleware")

	_, user, err := am.strategy.AuthenticateRequest(r)
	if err != nil {
		log.Println(err)

		code := http.StatusUnauthorized
		http.Error(w, http.StatusText(code), code)

		return
	}

	log.Printf("User %s Authenticated\n", user.GetUserName())
	r = auth.RequestWithUser(user, r)

	next(w, r)
}

// GetUserID return User ID.
func (am *Guardian) GetUserID(r *http.Request) string {
	u := auth.User(r)

	return u.GetID()
}
