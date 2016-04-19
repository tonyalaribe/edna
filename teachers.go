package main

import (
	"encoding/json"
	"github.com/gorilla/context"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

//TeacherAssignmentCollection struct
type TeacherAssignmentCollection struct {
	Subjects []Subject `json:"subjects"`
}

//GetAllTeachers gets all user from db
func (r *UserRepo) GetAllTeachers() ([]User, error) {
	var users []User
	err := r.coll.Find(bson.M{
		"roles": bson.M{
			"$elemMatch": bson.M{
				"$eq": "teacher",
			},
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

	u := UserRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_users")}
	users, err := u.GetAllTeachers()

	err = json.NewEncoder(w).Encode(UserCollection{users})
	if err != nil {
		log.Println(err)
	}
}

func (c *Config) getTeacherAssignmentsHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)

	user, err := userget(r)
	if err != nil {
		log.Println(err)
	}

	me := UserRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_users")}
	me_data, err := me.Get(user.Phone)
	if err != nil {
		log.Println(err)
	}

	u := SubjectRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_subjects")}

	subjects, err := u.GetAllAssignedToTeacher(me_data.ID.Hex())
	if err != nil {
		log.Println(err)
	}

	err = json.NewEncoder(w).Encode(TeacherAssignmentCollection{subjects})
	if err != nil {
		log.Println(err)
	}

}
