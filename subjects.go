package main

import (
	"encoding/json"

	"github.com/gorilla/context"

	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"net/http"
)

// Subject struct
type Subject struct {
	ID          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string        `json:"name"`
	Parent      string        `json:"parent"`
	Class       string        `json:"class"`
	Teachers    []string      `json:"teachers"`
	Assessments []Assessment  `json:"assessments"`
}

//SubjectCollection struct
type SubjectCollection struct {
	Subjects []Subject `json:"subjects"`
}

//SubjectData acts like ClassCollection but carries information about a single class
type SubjectData struct {
	Subject Subject `json:"subject"`
}

//SubjectRepo a mongo Collection that could get passed around
type SubjectRepo struct {
	coll *mgo.Collection
}

/* THese are functions that perform the operations on the db. .they are usually,
called by the handlers, in a bid to keep  handlers simple and less bulky.
*/

//Create adds a user to the database
func (r *SubjectRepo) Create(subject *Subject) error {

	err := r.coll.Insert(subject)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//Update adds a user to the database
func (r *SubjectRepo) Update(subject *Subject) error {

	err := r.coll.UpdateId(subject.ID, subject)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//Get gets a class's details from db
func (r *SubjectRepo) Get(id string) (Subject, error) {
	var subject Subject
	err := r.coll.Find(bson.M{
		"_id": bson.ObjectIdHex(id),
	}).One(&subject)

	if err != nil {
		log.Println(err)
		return subject, err
	}
	return subject, nil
}

//GetByName gets a class's details from db
func (r *SubjectRepo) GetByName(name string) (Subject, error) {
	var subject Subject
	err := r.coll.Find(bson.M{
		"name": name,
	}).One(&subject)

	if err != nil {
		log.Println(err)
		return subject, err
	}
	return subject, nil
}

//GetAll gets all user from db
func (r *SubjectRepo) GetAll() ([]Subject, error) {
	var subjects []Subject
	err := r.coll.Find(bson.M{}).All(&subjects)

	if err != nil {
		log.Println(err)
		return subjects, err
	}

	return subjects, nil
}

//GetAllAssignedToTeacher gets all user from db
func (r *SubjectRepo) GetAllAssignedToTeacher(teacher string) ([]Subject, error) {
	var subjects []Subject
	err := r.coll.Find(bson.M{
		"teachers": bson.M{
			"$elemMatch": bson.M{
				"$eq": teacher,
			},
		},
	}).All(&subjects)
	log.Println(subjects)

	if err != nil {
		log.Println(err)
		return subjects, err
	}

	return subjects, nil
}

/***************
handlers
***************/
//createSubjectHandler would create a class
func (c *Config) createSubjectHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)
	u := SubjectRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_subjects")}
	subject := Subject{}
	err := json.NewDecoder(r.Body).Decode(&subject)
	if err != nil {
		log.Println(err)
	}
	err = u.Create(&subject)
	if err != nil {
		log.Println(err)
	}
}

//getSubjectsHandler would create a user/staff
func (c *Config) getSubjectsHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)

	u := SubjectRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_subjects")}
	subjects, err := u.GetAll()

	err = json.NewEncoder(w).Encode(SubjectCollection{subjects})
	if err != nil {
		log.Println(err)
	}
}

//getSubjectHandler would get a subbject
func (c *Config) getSubjectHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)

	name := r.URL.Query().Get("id")

	u := SubjectRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_subjects")}
	subject, err := u.Get(name)

	err = json.NewEncoder(w).Encode(SubjectData{subject})
	if err != nil {
		log.Println(err)
	}
}

//putSubjectHandler would create a class
func (c *Config) putSubjectHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)
	u := SubjectRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_subjects")}
	subject := Subject{}
	err := json.NewDecoder(r.Body).Decode(&subject)
	if err != nil {
		log.Println(err)
	}
	err = u.Update(&subject)
	if err != nil {
		log.Println(err)
	}
}
