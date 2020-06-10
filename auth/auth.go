package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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

			key, _ := getKey()

			token, _ := jwt.Parse(bytes.NewReader([]byte(bearer)), jwt.WithVerify(jwa.RS256, key))

			buf, _ := json.MarshalIndent(token, "", "")
			//TODO convert buffer to string and read data

			fmt.Printf("%s\n", buf)
			print(key)

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

var jwksURL string = `https://localhost:5000/.well-known/openid-configuration/jwks`

func getKey() (interface{}, error) {
	set, err := jwk.Fetch(jwksURL)
	if err != nil {
		log.Printf("failed to parse JWK: ${}%s", err)
		return nil, errors.New("failed to parse JWKs")
	}

	// If you KNOW you have exactly one key, you can just
	// use set.Keys[0]
	keys := set.LookupKeyID("7-dIYXvTMBsfgD5YMpd5TA")
	if len(keys) == 0 {
		log.Printf("failed to lookup key: %s", err)
		return nil, errors.New("failed to lookup keys")
	}

	var key interface{} //*jwt.Token

	if err := keys[0].Raw(&key); err != nil {
		log.Printf("failed to create public key: %s", err)
		return nil, errors.New("failed to created public key")
	}
	return key, nil
}
