package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Payload struct {
	Action string      `json:"action"`
	Auth   Credentials `json:"auth"`
}

var jwtKey = []byte("microservice_auth")

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				log.Println(e)
				render(w, "login.gohtml")
			}
		}()

		c, err := r.Cookie("userToken")
		if err != nil {
			panic(err)
		}

		tokenStr := c.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			return jwtKey, nil
		})

		if err != nil {
			panic(err)
		}

		// If token is not valid
		if !tkn.Valid {
			panic(err)
		}

		http.Redirect(w, r, "/home", http.StatusFound)
	})

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				log.Println(e)
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}()

		c, err := r.Cookie("userToken")
		if err != nil {
			panic(err)
		}

		tokenStr := c.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			return jwtKey, nil
		})

		if err != nil {
			panic(err)
		}

		// If token is not valid
		if !tkn.Valid {
			panic(err)
		}

		render(w, "test.page.gohtml")
	})

	http.HandleFunc("/get-username", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				log.Println(e)
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}()

		c, err := r.Cookie("userToken")
		if err != nil {
			panic(err)
		}

		tokenStr := c.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			return jwtKey, nil
		})

		if err != nil {
			panic(err)
		}

		// If token is not valid
		if !tkn.Valid {
			panic(err)
		}

		val := struct {
			Email string `json:"email"`
		}{
			claims.Email,
		}

		out, err := json.Marshal(val)
		if err != nil {
			log.Println("ERROR OCCURED: While JSON Marshalling")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusFound)
		_, err = w.Write(out)
		if err != nil {
			log.Println("ERROR OCCURED: While writing response")
		}
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}()
		r.ParseForm()

		email := r.PostFormValue("email")
		pass := r.PostFormValue("password")

		cred := Credentials{
			Email:    email,
			Password: pass,
		}

		pl := Payload{
			Action: "register",
			Auth:   cred,
		}

		jsonData, _ := json.Marshal(pl)
		registerURL := "http://localhost:8080/handle"

		request, err := http.NewRequest("POST", registerURL, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println(err)
			panic(err)
		}

		client := &http.Client{}

		res, err := client.Do(request)
		if err != nil {
			log.Println(err)
			panic(err)
		}
		res.Body.Close()

		c, err := checkJWTAuthorization(res)
		if err != nil {
			panic(err)
		}

		http.SetCookie(w, c)

		http.Redirect(w, r, "/home", http.StatusFound)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if e := recover(); e != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}()

		r.ParseForm()

		email := r.PostFormValue("email")
		pass := r.PostFormValue("password")

		cred := Credentials{
			Email:    email,
			Password: pass,
		}

		pl := Payload{
			Action: "login",
			Auth:   cred,
		}

		jsonData, _ := json.Marshal(pl)
		registerURL := "http://localhost:8080/handle"

		request, err := http.NewRequest("POST", registerURL, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println(err)
			panic(err)
		}

		client := &http.Client{}

		res, err := client.Do(request)
		if err != nil {
			log.Println(err)
			panic(err)
		}
		res.Body.Close()

		c, err := checkJWTAuthorization(res)
		if err != nil {
			panic(err)
		}

		http.SetCookie(w, c)

		http.Redirect(w, r, "/home", http.StatusFound)
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{Name: "userToken", MaxAge: -1, HttpOnly: true, Path: "/"}
		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/", http.StatusFound)
	})

	// This is inside main function
	fmt.Println("Starting front end service on port 80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Panic(err)
	}
}

func render(w http.ResponseWriter, t string) {

	partials := []string{
		"./cmd/web/templates/login.gohtml",
		"./cmd/web/templates/base.layout.gohtml",
		"./cmd/web/templates/header.partial.gohtml",
		"./cmd/web/templates/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("./cmd/web/templates/%s", t))

	templateSlice = append(templateSlice, partials...)

	tmpl, err := template.ParseFiles(templateSlice...)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func checkJWTAuthorization(r *http.Response) (*http.Cookie, error) {
	cook := r.Cookies()
	var c *http.Cookie
	for _, cookie := range cook {
		if cookie.Name == "userToken" {
			c = cookie
			break
		}
	}

	if c == nil {
		return nil, errors.New("cookie couldn't find")
	}

	tokenStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	// If token is not valid
	if !tkn.Valid {
		return nil, errors.New("token is not valid")
	}

	return c, nil
}
