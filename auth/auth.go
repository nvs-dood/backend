package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/EnglederLucas/nvs-dood/authmodels"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

// Middleware decodes the share session cookie and packs the session into context
func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//Get token
			bearer := r.Header.Get("Authorization")
			bearer = strings.Split(bearer, " ")[1]

			//TODO add error handling

			key, err := getKey()
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(bytes.NewReader([]byte(bearer)), jwt.WithVerify(jwa.RS256, key))
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			buf, err := json.MarshalIndent(token, "", "")
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			log.Printf("%s\n", buf)

			iuserEmail, _ := token.Get("email")
			igivenName, _ := token.Get("given_name")
			ifamilyName, _ := token.Get("family_name")
			irole, _ := token.Get("doodrole")
			userEmail := fmt.Sprintf("%v", iuserEmail)
			userName := fmt.Sprintf("%v %v", igivenName, ifamilyName)
			doodRole := fmt.Sprintf("%v", irole)
			isAdmin := false
			if doodRole == "supervisor" {
				isAdmin = true
			}

			user := &authmodels.User{
				ID:    token.Subject(),
				Email: userEmail,
				Name:  &userName,
				Admin: isAdmin,
			}

			ctx := context.WithValue(r.Context(), userCtxKey, user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *authmodels.User {
	raw, _ := ctx.Value(userCtxKey).(*authmodels.User)
	return raw
}

func getKey() (interface{}, error) {

	var jwksURL string = os.Getenv("JWKS_URL")
	if jwksURL == "" {
		jwksURL = `http://localhost:5000/.well-known/openid-configuration/jwks`
	}
	set, err := jwk.Fetch(jwksURL)
	if err != nil {
		log.Printf("failed to parse JWK: %s", err)
		return nil, errors.New("failed to parse JWKs")
	}

	// If you KNOW you have exactly one key, you can just
	// use set.Keys[0]

	keys := set.Keys
	if len(keys) != 1 {
		log.Printf("don't know which signing key to pick: %s", err)
		return nil, errors.New("don't know which signing key to pick")
	}

	var key interface{} //*jwt.Token

	if err := keys[0].Raw(&key); err != nil {
		log.Printf("failed to create public key: %s", err)
		return nil, errors.New("failed to created public key")
	}
	return key, nil
}
