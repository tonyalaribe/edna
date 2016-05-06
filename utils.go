package main

import (
	//	"log"
	"math/rand"

	"net/http"

	"github.com/gorilla/context"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/mgo.v2/bson"
)

func randSeq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]

	}
	return string(b)

}

func userget(r *http.Request) (User, error) {
	id := context.Get(r, "UserID")
	u := context.Get(r, "User")

	var user User
	err := mapstructure.Decode(u, &user)
	user.ID = bson.ObjectIdHex(id.(string))

	if err != nil {
		return user, err
	}
	return user, nil

}
