package main

import (
	"encoding/json"

	"github.com/gorilla/context"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"

	"net/http"
)

// Class struct
type Class struct {
	Name     string   `json:"name" bson:"_id"`
	Parent   string   `json:"parent"`
	Teachers []string `json:"teachers"`
}

//ClassCollection struct
type ClassCollection struct {
	Classes []Class `json:"classes"`
}

//ClassData acts like ClassCollection but carries information about a single class
type ClassData struct {
	Class Class `json:"class"`
}

//ClassRepo a mongo Collection that could get passed around
type ClassRepo struct {
	coll *mgo.Collection
}

/* THese are functions that perform the operations on the db. .they are usually,
called by the handlers, in a bid to keep  handlers simple and less bulky.
*/

//Create adds a user to the database
func (r *ClassRepo) Create(class *Class) error {

	err := r.coll.Insert(class)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//Update adds a user to the database
func (r *ClassRepo) Update(class *Class) error {

	err := r.coll.UpdateId(class.Name, class)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//Get gets a class's details from db
func (r *ClassRepo) Get(slug string) (Class, error) {
	var class Class
	err := r.coll.Find(bson.M{
		"slug": slug,
	}).One(&class)

	if err != nil {
		log.Println(err)
		return class, err
	}

	return class, nil
}

//GetAll gets all user from db
func (r *ClassRepo) GetAll() ([]Class, error) {
	classes := []Class{}
	err := r.coll.Find(bson.M{}).All(&classes)
	if err != nil {
		log.Println(err)
		return classes, err
	}

	return classes, nil
}

//GetAllChildClasses gets all user from db
func (r *ClassRepo) GetAllChildClasses(parent string) ([]Class, error) {
	classes := []Class{}
	err := r.coll.Find(bson.M{
		"parent": parent,
	}).All(&classes)
	if err != nil {
		log.Println(err)
		return classes, err
	}

	return classes, nil
}

//GetClassesAssignedToTeacher getss all classess assigned to a teacher
func (r *ClassRepo) GetClassesAssignedToTeacher(teacher string) ([]Class, error) {
	classes := []Class{}
	err := r.coll.Find(bson.M{
		"teachers": bson.M{
			"$elemMatch": bson.M{
				"$eq": teacher,
			},
		},
	}).All(&classes)

	if err != nil {
		log.Println(err)
		return classes, err
	}

	return classes, nil
}

/***************
handlers
***************/
//createClassHandler would create a class
func (c *Config) createClassHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)
	u := ClassRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_classes")}
	class := Class{}
	err := json.NewDecoder(r.Body).Decode(&class)
	if err != nil {
		log.Println(err)
	}
	err = u.Create(&class)
	if err != nil {
		log.Println(err)
	}
}

//getUsersHandler would create a user/staff
func (c *Config) getClassesHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)

	u := ClassRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_classes")}
	classes, err := u.GetAll()
	err = json.NewEncoder(w).Encode(ClassCollection{classes})
	if err != nil {
		log.Println(err)
	}
}

//putClassHandler would create a class
func (c *Config) putClassHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)
	u := ClassRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_classes")}
	class := Class{}
	err := json.NewDecoder(r.Body).Decode(&class)
	if err != nil {
		log.Println(err)
	}
	err = u.Update(&class)
	if err != nil {
		log.Println(err)
	}
}
