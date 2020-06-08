package auth

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/example/scalars/model"
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

			bearer := r.Header.Get("Authorization")
			///print(bearer)
			//if bearer == nil{
			//	user := &model.User{ID: "3", Name: "Fritz"}
			//}
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
