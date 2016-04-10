package main

import (
	"log"
	"net/http"
	"os"

	//	"github.com/gorilla/context"
	//	"github.com/justinas/alice" //A middleware chaining library.

	"github.com/gorilla/context"
	"github.com/justinas/alice"
	"github.com/rs/cors"
)

const (
	//Cost is the, well, cost of the bcrypt encryption used for storing user
	//passwords in the database. It determines the amount of processing power to
	// be used while hashing and saalting the password. The higher, the cost,
	//the more secure the password hash, and also the more cpu cycles used for
	//password related processes like comparing hasshes during authentication
	//or even hashing a new password.
	Cost int = 5
)

type Conf struct {
	AuthToken string
}

var (
	Token Conf
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	Token = Conf{
		AuthToken: "iBamebYJnnpXUCjP6tmeYpUkMw3HA57pdSs7c1qc.H92jTVMAg112xmPQnnwAH",
	}
}

func main() {
	//REDISADDR, REDISPW, MONGOSERVER, MONGODB, Public, Private, RootURL, AWSBucket := checks()

	//config := generateConfig()

	commonHandlers := alice.New(context.ClearHandler, loggingHandler, recoverHandler)
	router := NewRouter()

	//router.Post("/api/v0.1/auth", commonHandlers.ThenFunc(appC.authHandler))

	router.HandleMethodNotAllowed = false
	router.NotFound = http.FileServer(http.Dir("./static")).ServeHTTP
	//api routes for iparent
	router.Get("/api/child", commonHandlers.ThenFunc(ChildHandler))
	router.Get("/api/board", commonHandlers.ThenFunc(BoardHandler))
	router.Post("/send", commonHandlers.ThenFunc(RegParent))
	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Println("No Global port has been defined, using default port :8080")

		PORT = "8080"

	}

	handler := cors.New(cors.Options{
		//		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedOrigins: []string{"*"},

		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-Auth-Token", "*"},
		Debug:            false,
	}).Handler(router)
	log.Println("serving ")

	//open.Run("http://localhost:" + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, handler))
}
