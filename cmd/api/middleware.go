package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/geekilx/restaurantAPI/internal/models"
	"github.com/geekilx/restaurantAPI/internal/validator"
)

func (app *application) panicRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {

			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")

				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)

	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if app.cfg.Limiter.Enabled {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			currentSecond := time.Now().Unix()

			key := fmt.Sprintf("rateLimit:%s:%d", ip, currentSecond)

			ctx := r.Context()
			pipe := app.redis.TxPipeline()

			incr := pipe.Incr(ctx, key)
			pipe.Expire(ctx, key, 5*time.Second)

			_, err = pipe.Exec(ctx)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			if incr.Val() > int64(app.cfg.Limiter.Rps) {
				app.rateLimitExceededResponse(w, r)
				return
			}

		}
		next.ServeHTTP(w, r)

	})

}

func (app *application) authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			r = app.setUserContext(w, r, models.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if models.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.failedValidationResponse(w, r, v)
			return
		}

		var userID int64
		var user *models.User

		// --- REDIS LOGIC START ---

		idStr, err := app.redis.Get(r.Context(), "token:"+token).Result()
		if err == nil {
			userID, _ = strconv.ParseInt(idStr, 10, 64)

			userJSON, err := app.redis.Get(r.Context(), "user:"+idStr).Result()
			if err == nil {
				err = json.Unmarshal([]byte(userJSON), &user)
				if err == nil {
					r = app.setUserContext(w, r, user)
					next.ServeHTTP(w, r)
					return
				}
			}

			// if we get to this point, it means that we could not found the token specified in the request
			// in the redis cache, so we fallback to query the database
		} else {
			userID, err = app.models.Tokens.GetByToken(token)
			if err != nil {
				switch {
				case errors.Is(err, models.ErrRecordNotFound):
					app.invalidAuthenticationTokenResponse(w, r)
				default:
					app.serverErrorResponse(w, r, err)
				}
				return
			}
		}

		// if we get to this point, it means that we found the token in the redis cache but we could not find the user struct in the database
		// so we fallback to query the database
		user, err = app.models.Users.GetUser(userID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// create redis cache for user
		if userBytes, err := json.Marshal(user); err == nil {
			app.redis.Set(r.Context(), "user:"+strconv.FormatInt(user.ID, 10), userBytes, 24*time.Hour)
		}
		// --- REDIS LOGIC END ---
		r = app.setUserContext(w, r, user)

		next.ServeHTTP(w, r)
	})

}

func (app *application) requiredAutheicatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := app.getUserContext(r)
		if models.IsAnonymous(user) {
			app.authorizationRequierd(w, r)
			return
		}

		next.ServeHTTP(w, r)

	})
}

func (app *application) requiredActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := app.getUserContext(r)
		if !user.IsActive {
			app.requiredActivationResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})

	return app.requiredAutheicatedUser(fn)
}

func (app *application) requirePermissions(code string, next http.HandlerFunc) http.HandlerFunc {

	fn := func(w http.ResponseWriter, r *http.Request) {

		user := app.getUserContext(r)
		if models.IsAnonymous(user) {
			app.notPermittedResponse(w, r)
			return
		}

		permissions, err := app.models.Permissions.GetForAllUser(user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		if !permissions.Include(code) {
			app.notPermittedResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}

	return app.requiredActivatedUser(fn)

}
