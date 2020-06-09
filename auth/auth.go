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

	"github.com/99designs/gqlgen/example/scalars/model"
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
func Middleware( /*db *sql.DB*/ ) func(http.Handler) http.Handler {
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

			// TODO transfer into Usermodel
			user := &model.User{
				ID: 3,
				Address: model.Address{ID: 3, Location: &model.Point{
					X: 815, Y: 815,
				}},
				Name: bearer,
			}

			//Get UID and then User

			ctx := context.WithValue(r.Context(), userCtxKey, user)

			r = r.WithContext(ctx)
			//w.Write([]byte("aösdkfaposdj"))
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *model.User {
	raw, _ := ctx.Value(userCtxKey).(*model.User)
	return raw
}

// import (
// 	"bytes"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"log"

// 	"github.com/lestrrat-go/jwx/jwa"
// 	"github.com/lestrrat-go/jwx/jwk"
// 	"github.com/lestrrat-go/jwx/jwt"
// )

//const bearer = `eyJhbGciOiJSUzI1NiIsImtpZCI6IjctZElZWHZUTUJzZmdENVlNcGQ1VEEiLCJ0eXAiOiJhdCtqd3QifQ.eyJuYmYiOjE1OTE3MjI3MTEsImV4cCI6MTU5MTcyNjMxMSwiaXNzIjoiaHR0cHM6Ly9sb2NhbGhvc3Q6NTAwMCIsImNsaWVudF9pZCI6InBvc3RtYW4tY2xpZW50Iiwic3ViIjoiMGZiZGIzOGUtNTkzNC00M2M2LWI5MGUtNWJjZWRmNzQyNTJlIiwiYXV0aF90aW1lIjoxNTkxNjgzMjYzLCJpZHAiOiJsb2NhbCIsIndlYnNpdGUiOiJodHRwOi8vYm9iLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjoidHJ1ZSIsImVtYWlsIjoiQm9iU21pdGhAZW1haWwuY29tIiwiZmFtaWx5X25hbWUiOiJTbWl0aCIsImdpdmVuX25hbWUiOiJCb2IiLCJuYW1lIjoiQm9iIFNtaXRoIiwicHJlZmVycmVkX3VzZXJuYW1lIjoiYm9iIiwiZG9vZHJvbGUiOiJzdXBlcnZpc29yIiwic2NvcGUiOlsib3BlbmlkIiwicHJvZmlsZSIsImVtYWlsIiwiZG9vZCJdLCJhbXIiOlsicHdkIl19.O2yhapqoC6_UjxwOEBW7lVU2Hj2AVCix8t27gBZ7DVJ4Xr5-P1JgYW-pz-Shb9Tc1vHcPsx-3xXp3f3XUplqacvyu2dx9bTmET4XMdP98u27veor1VZmiCyjf7Fyvphh1jtHZ9cJSivFatFtnkyFOpk-UU402-EfkLk4gzc11W4lAOnfgP_hWjHvVJhlf9yEH4Pb7LREwGib_ms2wXqDDz71FAnbyO26MaYwrqNxOaseCMgom5wVnvr2_6wHskwtjaV1JlEUOr5WjwpICL9pE9TRG3eMgrc1lcrH5jrIRBFCaG9wezMqZ_FQ_wctAuHCUW09yxq2cBsG7SLPJGoZug`

const jwksURL = `https://localhost:5000/.well-known/openid-configuration/jwks`

func main() {

}

func getKey( /*bearer *jwt.Token*/ ) (interface{}, error) {
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
