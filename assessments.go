package main

import (
	"encoding/json"

	"github.com/gorilla/context"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"

	"net/http"
)

type Assessment struct {
	Name       string `json:"name"`
	Upperlimit string `json:"upperlimit"`
	Percentage string `json:"percentage"`
}

type StudentAssessments struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	StudentID   string        `json:"studentid"`
	Name        string        `json:"name"`
	Subject     string        `json:"subject"`
	Class       string        `json:"class"`
	Assessments []struct {
		ID    bson.ObjectId `json:"name" bson:"_id"`
		Name  string        `json:"name"`
		Score string        `json:"score"`
	}
}

//StudentAssessmentCollection struct
type StudentAssessmentCollection struct {
	Subjects []StudentAssessments `json:"subjects"`
}

//StudentAssessmentData acts like ClassCollection but carries information about a single class
type StudentAssessmentData struct {
	Subject Subject `json:"subject"`
}

//StudentAssessmentRepo a mongo Collection that could get passed around
type StudentAssessmentRepo struct {
	coll *mgo.Collection
}

func (r *SubjectRepo) CreateAssessment(subjectid string, assessment Assessment) error {
	err := r.coll.Update(bson.M{
		"_id": bson.ObjectIdHex(subjectid),
	},
		bson.M{
			"$push": bson.M{
				"assessments": assessment,
			},
		})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *StudentAssessmentRepo) UpsertStudentAssessmentData(s StudentAssessments) error {
	_, err := r.coll.Upsert(bson.M{
		"studentid": s.StudentID,
		"class":     s.Class,
		"subject":   s.Subject,
	}, bson.M{
		"$push": bson.M{
			"$addToSet": bson.M{
				"assessments": bson.M{
					"$each": s.Assessments,
				},
			},
		},
		"$setOnInsert": s,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *StudentAssessmentRepo) GetAssessments() ([]StudentAssessments, error) {
	result := []StudentAssessments{}
	err := r.coll.Find(bson.M{}).All(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}

//func RecordAssessment(studentid, Assessment)

//newAssessmentHandler would create an assessment
func (c *Config) newAssessmentHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)
	u := SubjectRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_subject")}
	assessment := Assessment{}
	err := json.NewDecoder(r.Body).Decode(&assessment)
	if err != nil {
		log.Println(err)
	}

	subjectid := r.URL.Query().Get("id")
	log.Println(subjectid)
	err = u.CreateAssessment(subjectid, assessment)
	if err != nil {
		log.Println(err)
	}

}

func (c *Config) getStudentsAndAssessmentsHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)
	StudentAssessmentRepo := StudentAssessmentRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_assessments")}
	StudentRepo := StudentRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_students")}
	assessments, err := StudentAssessmentRepo.GetAssessments()
	if err != nil {
		log.Println(err)
	}

	class := r.URL.Query().Get("class")
	//subject := r.URL.Query().Get("subject")
	students, err := StudentRepo.GetAllStudentsInParentClass(class)
	if err != nil {
		log.Println(err)
	}

	returnStudents := []StudentAssessments{}

	for _, student := range students {
		log.Println(student)
		s := StudentAssessments{}

		s.Name = student.FirstName + " " + student.LastName
		log.Println(s.Name)
		s.StudentID = student.ID.Hex()

		//as := assessments.Assessments

		for i, a := range assessments {
			log.Println(a.ID)
			log.Println(student.ID)
			if a.StudentID == student.ID.Hex() {
				s.Assessments = append(s.Assessments, a.Assessments...)
				assessments[i] = assessments[len(assessments)-1]
				assessments = assessments[:len(assessments)-1]
			}

		}

		returnStudents = append(returnStudents, s)
		//log.Println(s)
	}

	log.Println(returnStudents)
	//log.Println(students)

	err = json.NewEncoder(w).Encode(returnStudents)
	if err != nil {
		log.Println(err)
	}

}

func (c *Config) addStudentAssessmentsHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)
	sRepo := StudentAssessmentRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_assessments")}

	log.Println(sRepo)

	s := StudentAssessments{}
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		log.Println(err)
	}

	err = sRepo.UpsertStudentAssessmentData(s)
	if err != nil {
		log.Println(err)
	}

}
