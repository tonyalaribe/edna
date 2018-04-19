package main

import (
	"encoding/json"

	"github.com/gorilla/context"

	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"net/http"
)

// Session struct
type Session struct {
	ID    string `json:"id" bson:"_id"`
	Start string `json:"start"`
	End   string `json:"end"`
}

//SessionRepo a mongo Collection that could get passed around
type SessionRepo struct {
	coll *mgo.Collection
}

/* THese are functions that perform the operations on the db. .they are usually,
called by the handlers, in a bid to keep  handlers simple and less bulky.
*/

//Create adds a user to the database
func (r *SessionRepo) Create(session *Session) error {

	err := r.coll.Insert(session)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//Session gets all user from db
func (r *SessionRepo) GetAll() ([]Session, error) {
	sessions := []Session{}
	err := r.coll.Find(bson.M{}).All(&sessions)
	if err != nil {
		log.Println(err)
		return sessions, err
	}

	return sessions, nil
}

/***************
handlers
***************/
//createSessionHandler would create a class
func (c *Config) createSessionHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)
	u := SessionRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_sessions")}
	session := Session{}
	err := json.NewDecoder(r.Body).Decode(&session)
	if err != nil {
		log.Println(err)
	}
	err = u.Create(&session)
	if err != nil {
		log.Println(err)
	}
}

//getUsersHandler would create a user/staff
func (c *Config) getSessionHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)

	u := SessionRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_sessions")}
	sessions, err := u.GetAll()
	err = json.NewEncoder(w).Encode(sessions)
	if err != nil {
		log.Println(err)
	}
}
