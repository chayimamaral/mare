package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TenantClaims struct {
	ID          string `json:"id"`
	Active      bool   `json:"active"`
	Nome        string `json:"nome"`
	Contato     string `json:"contato"`
	Plano       string `json:"plano"`
	SchemaName  string `json:"schema_name,omitempty"`
	IsVecMaster bool   `json:"is_vec_master,omitempty"`
}

type Claims struct {
	UserID           string       `json:"-"`
	Nome             string       `json:"nome"`
	Email            string       `json:"email"`
	Tenant           TenantClaims `json:"tenant"`
	Role             string       `json:"role"`
	RepresentativeID string       `json:"representative_id,omitempty"`
	// Sem omitempty: lista vazia deve ir como [] no JWT para o middleware não confundir com token legado.
	FeatureSlugs     []string     `json:"feature_slugs"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secret []byte

	expiresIn time.Duration
}

func NewTokenService(secret string) *TokenService {
	return &TokenService{
		secret:    []byte(secret),
		expiresIn: 30 * 24 * time.Hour,
	}
}

func (s *TokenService) Generate(claims Claims) (string, error) {
	now := time.Now()
	claims.RegisteredClaims = jwt.RegisteredClaims{
		Subject:   claims.UserID,
		ExpiresAt: jwt.NewNumericDate(now.Add(s.expiresIn)),
		IssuedAt:  jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *TokenService) Parse(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
