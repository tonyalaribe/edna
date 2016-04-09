package main

import (
	"encoding/json"

	//	"github.com/gorilla/context"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"

	"net/http"
)

// Guardian struct
type Guardian struct {
	Name       string   `json:"name"`
	Phone      string   `json:"phone" bson:"_id,omitempty"`
	Email      string   `json:"email"`
	Occupation string   `json:"occupation"`
	Address    string   `json:"address"`
	Country    string   `json:"country"`
	State      string   `json:"state"`
	City       string   `json:"city"`
	Schools    []string `json:"schools"`
}

//GuardianCollection struct
type GuardianCollection struct {
	Guardians []Guardian `json:"guardians"`
}

//GuardianData acts like GuardianCollection but carries information about a single class
type GuardianData struct {
	Guardian Guardian `json:"guardian"`
}

//GuardianRepo a mongo Collection that could get passed around
type GuardianRepo struct {
	coll *mgo.Collection
}

/* THese are functions that perform the operations on the db. .they are usually,
called by the handlers, in a bid to keep  handlers simple and less bulky.
*/

//Create adds a user to the database
func (r *GuardianRepo) Create(guardian *Guardian, schoolID string) error {

	g := Guardian{}
	err := r.coll.FindId(guardian.Phone).One(&g)
	if err != nil {
		log.Println(err)
	}

	if g.Phone != "" {
		r.coll.UpsertId(guardian.Phone, bson.M{
			"$push": bson.M{
				"schools": schoolID,
			},
		})
	} else {
		err := r.coll.Insert(guardian)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

//Update updates a guardian in the database
func (r *GuardianRepo) Update(guardian *Guardian) error {

	err := r.coll.UpdateId(guardian.Phone, guardian)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//Get gets a class's details from db
func (r *GuardianRepo) Get(slug string) (Guardian, error) {
	var guardian Guardian
	err := r.coll.Find(bson.M{
		"slug": slug,
	}).One(&guardian)

	if err != nil {
		log.Println(err)
		return guardian, err
	}

	return guardian, nil
}

//GetAll gets all user from db
func (r *GuardianRepo) GetAll() ([]Guardian, error) {
	var guardian []Guardian
	err := r.coll.Find(bson.M{}).All(&guardian)

	if err != nil {
		log.Println(err)
		return guardian, err
	}

	return guardian, nil
}

/***************
handlers
***************/

//getStudentsHandler would create a student
func (c *Config) getGuardianHandler(w http.ResponseWriter, r *http.Request) {

	u := GuardianRepo{c.MongoSession.DB(c.MONGODB).C("guardians")}
	guardians, err := u.GetAll()

	err = json.NewEncoder(w).Encode(GuardianCollection{guardians})
	if err != nil {
		log.Println(err)
	}
}

//putGuardianHandler would create a class
func (c *Config) putGuardianHandler(w http.ResponseWriter, r *http.Request) {

	u := GuardianRepo{c.MongoSession.DB(c.MONGODB).C("guardians")}
	guardian := Guardian{}
	err := json.NewDecoder(r.Body).Decode(&guardian)
	if err != nil {
		log.Println(err)
	}
	err = u.Update(&guardian)
	if err != nil {
		log.Println(err)
	}
}
