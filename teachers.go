package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"gopkg.in/mgo.v2/bson"
)

//TeacherAssignmentCollection struct
type TeacherAssignmentCollection struct {
	Subjects []Subject `json:"subjects"`
	Classes  []Class   `json:"classes"`
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

func (c *Config) getClassesAssignedToTeacherHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)

	user, err := userget(r)
	if err != nil {
		log.Println(err)
	}

	me := UserRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_users")}
	meData, err := me.Get(user.Phone)
	if err != nil {
		log.Println(err)
	}

	sRepo := SubjectRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_subjects")}

	subjects, err := sRepo.GetAllAssignedToTeacher(meData.ID.Hex())
	if err != nil {
		log.Println(err)
	}

	cRepo := ClassRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_classes")}

	returnSubjects := []Subject{}

	for _, s := range subjects {
		childClasses, err := cRepo.GetAllChildClasses(s.Parent)
		if err != nil {
			log.Println(err)
		}

		for _, cc := range childClasses {
			subj := Subject{}
			subj.ID = s.ID
			subj.Name = s.Name
			subj.Parent = s.Parent
			subj.Class = cc.Name
			subj.Teachers = s.Teachers
			subj.Assessments = s.Assessments

			returnSubjects = append(returnSubjects, subj)
		}

	}

	classes, err := cRepo.GetClassesAssignedToTeacher(meData.ID.Hex())

	if err != nil {
		log.Println(err)
	}

	err = json.NewEncoder(w).Encode(TeacherAssignmentCollection{returnSubjects, classes})
	if err != nil {
		log.Println(err)
	}

}
