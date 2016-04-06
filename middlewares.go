package main

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

// Middlewares

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				WriteError(w, ErrInternalServer)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

func acceptHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			WriteError(w, ErrNotAcceptable)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func contentTypeHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			WriteError(w, ErrUnsupportedMediaType)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func dbsetter(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		h := strings.Split(r.Host, ".")
		log.Println(h[0])

		session, err := mgo.Dial(MONGOSERVER)
		defer session.Close()
		if err != nil {
			panic(err)
			//log.Println(err)
		}
		session.SetMode(mgo.Monotonic, true)
		col := session.DB(MONGODB).C("schools")
		school := School{}

		if strings.Contains(r.Host, ":8080") || h[0] == "www" {
			err := col.Find(bson.M{
				"_id": "unical",
			}).One(&school)
			if err != nil {
				log.Println(err)
			}
			context.Set(r, "school", school)
			if err != nil {
				log.Println(err)
			}
		} else {
			err = col.Find(bson.M{
				"_id": h[0],
			}).One(&school)
			if err != nil {
				log.Println(err)
			}
			context.Set(r, "school", school)
			if err != nil {
				log.Println(err)
			}
		}
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func bodyHandler(v interface{}) func(http.Handler) http.Handler {
	t := reflect.TypeOf(v)

	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			val := reflect.New(t).Interface()
			err := json.NewDecoder(r.Body).Decode(val)

			if err != nil {
				log.Println(err)
				WriteError(w, ErrBadRequest)
				return
			}

			if next != nil {
				context.Set(r, "body", val)
				next.ServeHTTP(w, r)
			}
		}

		return http.HandlerFunc(fn)
	}

	return m
}

func (ac *Config) frontAuthHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		//for k, v := range r.Header {
		//  log.Println("key:", k, "value:", v)
		//}

		var tokenValue string

		// check if we have a cookie with out tokenName
		headerToken := r.Header.Get("X-AUTH-TOKEN")
		//log.Println(headerToken)

		if headerToken != "" {
			tokenValue = headerToken
		} else {

			tokenCookie, err := r.Cookie(ac.Token)
			if err != nil {
				log.Println(err)
			}
			//log.Println(ac.token)
			//log.Println(tokenCookie)

			switch {
			case err == http.ErrNoCookie:
				WriteError(w, ErrNoAuth)
				return

			case err != nil:
				//w.WriteHeader(http.StatusInternalServerError)
				//fmt.Fprintln(w, "Error while Parsing cookie!")
				log.Printf("Cookie parse error: %v\n", err)
				//next.ServeHTTP(w, r)
				WriteError(w, ErrInternalServer)
				return
			}

			tokenValue = tokenCookie.Value

		}
		// validate the token
		token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (interface{}, error) {
			// since we only use the one private key to sign the tokens, we also only use its public counter part to verify
			return ac.Public, nil
		})

		// branch out into the possible error from signing
		switch err.(type) {

		case nil: // no error

			if !token.Valid { // but may still be invalid
				WriteError(w, ErrBadToken)
			}

			context.Set(r, "User", token.Claims["User"])
			next.ServeHTTP(w, r)

		case *jwt.ValidationError: // something was wrong during the validation
			vErr := err.(*jwt.ValidationError)

			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				WriteError(w, ErrBadToken)
				return

			default:
				WriteError(w, ErrBadToken)
				log.Printf("ValidationError error: %+v\n", vErr.Errors)
				return

			}

		default: // something else went wrong
			WriteError(w, ErrBadToken)
			return
		}

	}
	return http.HandlerFunc(fn)

}
