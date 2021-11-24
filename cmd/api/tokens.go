package main

import (
	"ecom-api/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

//dummy user for development
var validUser = models.User{
	ID:       10,
	Email:    "admin@admin.com",
	Password: "$2a$12$p8cKatzIyII2V.QbOY9CcOhj23.PLblcu7E.ja42TX3ghWyRRtjjC", //Password: password
}

type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

func (app *application) signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	}

	hashedPassword := validUser.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	}

	var claim jwt.Claims
	claim.Subject = fmt.Sprint(validUser.ID)
	claim.Issued = jwt.NewNumericTime(time.Now())
	claim.NotBefore = jwt.NewNumericTime(time.Now())
	claim.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claim.Issuer = "mydomain.com"
	claim.Audiences = []string{"mydomain.com"}

	jwtBytes, err := claim.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	if err != nil {
		app.errorJSON(w, errors.New("error signing in"))
		return
	}

	app.writeJSON(w, http.StatusOK, string(jwtBytes), "response")
}
