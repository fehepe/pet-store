package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const (
	UserContextKey     = contextKey("user")
	UserTypeContextKey = contextKey("userType")
)

type UserType string

const (
	UserTypeMerchant UserType = "merchant"
	UserTypeCustomer UserType = "customer"
)

type User struct {
	Username string
	Type     UserType
}

// Hardcoded users for demo purposes - in production, these would be in a database
var users = map[string]struct {
	PasswordHash string
	Type         UserType
}{
	"merchant1": {
		PasswordHash: "$2a$10$YourHashedPasswordHere", // Password: merchant123
		Type:         UserTypeMerchant,
	},
	"customer1": {
		PasswordHash: "$2a$10$YourHashedPasswordHere", // Password: customer123
		Type:         UserTypeCustomer,
	},
	"customer2": {
		PasswordHash: "$2a$10$YourHashedPasswordHere", // Password: customer123
		Type:         UserTypeCustomer,
	},
}

func init() {
	users["merchant1"] = struct {
		PasswordHash string
		Type         UserType
	}{
		PasswordHash: hashPassword("merchant123"),
		Type:         UserTypeMerchant,
	}
	users["customer1"] = struct {
		PasswordHash string
		Type         UserType
	}{
		PasswordHash: hashPassword("customer123"),
		Type:         UserTypeCustomer,
	}
	users["customer2"] = struct {
		PasswordHash string
		Type         UserType
	}{
		PasswordHash: hashPassword("customer123"),
		Type:         UserTypeCustomer,
	}
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			return
		}

		const prefix = "Basic "
		if !strings.HasPrefix(authHeader, prefix) {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(authHeader[len(prefix):])
		if err != nil {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		username, password := parts[0], parts[1]

		user, ok := users[username]
		if !ok {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, username)
		ctx = context.WithValue(ctx, UserTypeContextKey, user.Type)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUser(ctx context.Context) (string, error) {
	user, ok := ctx.Value(UserContextKey).(string)
	if !ok {
		return "", errors.New("user not found in context")
	}
	return user, nil
}

func GetUserType(ctx context.Context) (UserType, error) {
	userType, ok := ctx.Value(UserTypeContextKey).(UserType)
	if !ok {
		return "", errors.New("user type not found in context")
	}
	return userType, nil
}

func RequireMerchant(ctx context.Context) error {
	userType, err := GetUserType(ctx)
	if err != nil {
		return err
	}
	if userType != UserTypeMerchant {
		return errors.New("merchant access required")
	}
	return nil
}

func RequireCustomer(ctx context.Context) error {
	userType, err := GetUserType(ctx)
	if err != nil {
		return err
	}
	if userType != UserTypeCustomer {
		return errors.New("customer access required")
	}
	return nil
}

// ConditionalAuthMiddleware allows certain public queries without authentication
func ConditionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the request body to check for public queries
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(strings.NewReader(string(body)))

		// Parse the GraphQL request
		var gqlRequest struct {
			Query string `json:"query"`
		}
		if err := json.Unmarshal(body, &gqlRequest); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		queryLower := strings.ToLower(gqlRequest.Query)
		isMutation := strings.Contains(queryLower, "mutation")
		isPublicQuery := strings.Contains(queryLower, "liststores") ||
			strings.Contains(queryLower, "availablepets")
		if !isMutation && isPublicQuery {
			next.ServeHTTP(w, r)
			return
		}

		// For all other queries, require authentication
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			return
		}

		const prefix = "Basic "
		if !strings.HasPrefix(authHeader, prefix) {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(authHeader[len(prefix):])
		if err != nil {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		credentials := string(decoded)
		parts := strings.SplitN(credentials, ":", 2)
		if len(parts) != 2 {
			http.Error(w, "Invalid credentials format", http.StatusUnauthorized)
			return
		}

		username, password := parts[0], parts[1]
		user, exists := users[username]
		if !exists {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), UserContextKey, username)
		ctx = context.WithValue(ctx, UserTypeContextKey, user.Type)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
