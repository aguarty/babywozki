package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"go.uber.org/zap"
)

type authData struct {
	User     string `json:"username"`
	Password string `json:"password"`
}

type token struct {
	Token string `json:"token"`
}

type MyCustomClaims struct {
	User string `json:"user"`
	Exp  int64  `json:"exp"`
	jwt.StandardClaims
}

//Authenticator - custom middleware check auth
func (a *application) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//token := TokenFromCookie(r)
		token, claims, err := jwtauth.FromContext(r.Context()) //claims
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		if token == nil || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		user := claims["user"].(string)
		if user != "john.doe" {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		t := time.Now().UTC()
		exp := int64(claims["exp"].(float64))
		if time.Unix(exp, 0).Unix() < t.Unix() {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

//Verifier - custom middleware vaerify auth
func Verifier(ja *jwtauth.JWTAuth, salt string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return jwtauth.Verify(ja, TokenFromCookie(salt))(next)
	}
}

//TokenFromCookie - return token from header
func TokenFromCookie(salt string) func(r *http.Request) string {
	return func(r *http.Request) string {
		if cookie, err := r.Cookie(serviceName); err == nil {
			c := string(decrypt([]byte(cookie.Value), salt))
			return c
		}
		return ""
	}
}

//TokenFromHeader - return token from header
func TokenFromHeader(r *http.Request) string {
	// Get token from authorization header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) >= 6 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

// loginHandler -
func (a *application) loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := authData{}
		resp := ApiResp{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			a.logger.Error("login", zap.Error(err))
			resp.Err = "Invalid input"
			sendResponse(a.logger, w, http.StatusBadRequest, resp)
			return
		}

		if req.Password == "foobar" && req.User == "john.doe" {
			a.tokenAuth = jwtauth.New("HS256", []byte(a.cfg.Secure.Secret), nil)
			_, tokenString, _ := a.tokenAuth.Encode(&MyCustomClaims{
				User: req.User,
				Exp:  time.Now().Add(time.Hour * 5).Unix(),
			})

			tokenString = string(encrypt([]byte(tokenString), a.cfg.Secure.Salt))

			cookie := &http.Cookie{
				Name:     serviceName,
				Value:    tokenString,
				Path:     "/",
				Expires:  time.Now().UTC().Add(time.Hour * 5),
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)

			sendResponse(a.logger, w, http.StatusOK, token{Token: tokenString})
		} else {
			resp.Err = "bad login"
			sendResponse(a.logger, w, http.StatusUnauthorized, resp)
		}

	}
}
