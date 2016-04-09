package main

import (
	"encoding/json"

	"github.com/gorilla/context"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"

	"net/http"
)

// Student struct
type Student struct {
	ID       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string        `json:"name"`
	Parent   string        `json:"parent"`
	Teachers []string      `json:"teachers"`
}

//StudentCollection struct
type StudentCollection struct {
	Students []Student `json:"students"`
}

//StudentData acts like StudentCollection but carries information about a single class
type StudentData struct {
	Student Student `json:"student"`
}

//StudentRepo a mongo Collection that could get passed around
type StudentRepo struct {
	coll *mgo.Collection
}

/* THese are functions that perform the operations on the db. .they are usually,
called by the handlers, in a bid to keep  handlers simple and less bulky.
*/

//Create adds a user to the database
func (r *StudentRepo) Create(student *Student) error {

	err := r.coll.Insert(student)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//Update updates a student in the database
func (r *StudentRepo) Update(student *Student) error {

	err := r.coll.UpdateId(student.ID, student)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//Get gets a class's details from db
func (r *StudentRepo) Get(slug string) (Student, error) {
	var student Student
	err := r.coll.Find(bson.M{
		"slug": slug,
	}).One(&student)

	if err != nil {
		log.Println(err)
		return student, err
	}

	return student, nil
}

//GetAll gets all user from db
func (r *StudentRepo) GetAll() ([]Student, error) {
	var student []Student
	err := r.coll.Find(bson.M{}).All(&student)

	if err != nil {
		log.Println(err)
		return student, err
	}

	return student, nil
}

/***************
handlers
***************/
//createStudentHandler would create a class
func (c *Config) createStudentHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)
	u := StudentRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_students")}
	student := Student{}
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		log.Println(err)
	}
	err = u.Create(&student)
	if err != nil {
		log.Println(err)
	}
}

//getStudentsHandler would create a student
func (c *Config) getStudentsHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)

	u := StudentRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_students")}
	students, err := u.GetAll()

	err = json.NewEncoder(w).Encode(StudentCollection{students})
	if err != nil {
		log.Println(err)
	}
}

//putStudentHandler would create a class
func (c *Config) putStudentHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)
	u := StudentRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_students")}
	student := Student{}
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		log.Println(err)
	}
	err = u.Update(&student)
	if err != nil {
		log.Println(err)
	}
}
