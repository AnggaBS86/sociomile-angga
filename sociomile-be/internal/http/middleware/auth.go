package middleware

import (
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"sociomile-be/internal/http/response"
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	TenantID int64  `json:"tenant_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type AuthMiddleware struct {
	secret []byte
}

const claimsContextKey = "auth_claims"

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{secret: []byte(secret)}
}

func GenerateToken(secret string, userID, tenantID int64, role string, ttlMinutes int) (string, error) {
	now := time.Now()

	claims := &Claims{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(ttlMinutes) * time.Minute)),
			Subject:   strconv.FormatInt(userID, 10),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func (m *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authz := c.Request().Header.Get("Authorization")
		if authz == "" {
			return response.Unauthorized("missing bearer token")
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(authz, "Bearer"))
		if tokenStr == "" {
			return response.Unauthorized("invalid bearer token")
		}

		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return m.secret, nil
		})
		if err != nil || !token.Valid {
			return response.Unauthorized("invalid token")
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			return response.Unauthorized("invalid token claims")
		}

		c.Set(claimsContextKey, claims)

		return next(c)
	}
}

func RequireRoles(roles ...string) echo.MiddlewareFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := GetClaims(c)
			if !ok {
				return response.Unauthorized("missing auth claims")
			}

			if _, exists := allowed[claims.Role]; !exists {
				return response.Forbidden("insufficient role")
			}

			return next(c)
		}
	}
}

func GetClaims(c echo.Context) (*Claims, bool) {
	v := c.Get(claimsContextKey)
	claims, ok := v.(*Claims)

	return claims, ok
}
