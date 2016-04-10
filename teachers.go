package main

import (
	"encoding/json"
	"github.com/gorilla/context"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

//GetAllTeachers gets all user from db
func (r *UserRepo) GetAllTeachers() ([]User, error) {
	var users []User
	err := r.coll.Find(bson.M{
		"$elemMatch": bson.M{
			"$eq": "",
		},
	}).All(&users)

	if err != nil {
		log.Println(err)
		return users, err
	}

	return users, nil
}

/***************
handlers
***************/

//getTeachers Handler would get list of teachers
func (c *Config) getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)
	log.Println(school)
	u := UserRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_users")}
	users, err := u.GetAllTeachers()

	err = json.NewEncoder(w).Encode(UserCollection{users})
	if err != nil {
		log.Println(err)
	}
}
