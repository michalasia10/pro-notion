package infrastructure

import (
	"errors"
	"net/url"
	"strings"
	"time"

	shared "src/internal/modules/shared/domain"
	"src/internal/modules/users/domain"

	"github.com/golang-jwt/jwt/v5"
)

const extensionOAuthAudience = "pro-notion-extension"

type extensionOAuthStateClaims struct {
	RedirectURI string `json:"redirect_uri"`
	jwt.RegisteredClaims
}

type ExtensionOAuthStateJWTService struct {
	secret   []byte
	issuer   string
	clock    shared.Clock
	audience string
}

func NewExtensionOAuthStateJWTService(secret, issuer string, clock shared.Clock) domain.ExtensionOAuthStateService {
	return &ExtensionOAuthStateJWTService{
		secret:   []byte(secret),
		issuer:   issuer,
		clock:    clock,
		audience: extensionOAuthAudience,
	}
}

func (s *ExtensionOAuthStateJWTService) Generate(redirectURI string) (string, error) {
	if err := validateExtensionRedirect(redirectURI); err != nil {
		return "", err
	}

	claims := extensionOAuthStateClaims{
		RedirectURI: redirectURI,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(s.clock.Now().Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(s.clock.Now()),
			Issuer:    s.issuer,
			Audience:  jwt.ClaimStrings{s.audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *ExtensionOAuthStateJWTService) Validate(state string) (string, error) {
	token, err := jwt.ParseWithClaims(state, &extensionOAuthStateClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.secret, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*extensionOAuthStateClaims)
	if !ok || !token.Valid {
		return "", jwt.ErrTokenInvalidClaims
	}
	if s.issuer != "" && claims.Issuer != s.issuer {
		return "", jwt.ErrTokenInvalidClaims
	}
	if !containsAudience(claims.Audience, s.audience) {
		return "", jwt.ErrTokenInvalidClaims
	}
	if err := validateExtensionRedirect(claims.RedirectURI); err != nil {
		return "", err
	}

	return claims.RedirectURI, nil
}

func containsAudience(aud jwt.ClaimStrings, want string) bool {
	for _, a := range aud {
		if a == want {
			return true
		}
	}
	return false
}

func validateExtensionRedirect(redirectURI string) error {
	if redirectURI == "" {
		return errors.New("missing redirect uri")
	}

	parsed, err := url.Parse(redirectURI)
	if err != nil {
		return errors.New("invalid redirect uri")
	}

	if parsed.Scheme != "https" {
		return errors.New("invalid redirect uri scheme")
	}

	if !strings.HasSuffix(parsed.Host, ".chromiumapp.org") {
		return errors.New("invalid redirect uri host")
	}

	return nil
}
