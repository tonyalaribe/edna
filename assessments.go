package main

import (
	"encoding/json"

	"github.com/gorilla/context"

	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"net/http"
)

//Assessment is
type Assessment struct {
	Name       string `json:"name"`
	Upperlimit string `json:"upperlimit"`
	Percentage string `json:"percentage"`
}

//StudentAssessments is
type StudentAssessments struct {
	ID          bson.ObjectId       `json:"id,omitempty" bson:"_id,omitempty"`
	StudentID   string              `json:"studentid"`
	Name        string              `json:"name"` //Student Name
	Subject     string              `json:"subject"`
	SubjectInfo Subject             `json:"subjectinfo"`
	Class       string              `json:"class"`
	Assessments []StudentAssessment `json:"assessments" bson:",omitempty"`
}

//StudentAssessment is
type StudentAssessment struct {
	ID    bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string        `json:"name"`
	Score int           `json:"score"`
}

//SingleStudentAssessment for  retrieving sudent assesment from api, and to aid json marshalling
type SingleStudentAssessment struct {
	ID             bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	StudentID      string        `json:"studentid"`
	Name           string        `json:"name"`
	Subject        string        `json:"subject"`
	Class          string        `json:"class"`
	AssessmentName string        `json:"assessmentname"`
	Score          int           `json:"score"`
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

//CreateAssessment is for creating assessments on a subject
func (r *SubjectRepo) CreateAssessment(subjectid string, assessment Assessment) error {
	//subjectidhex := bson.ObjectIdHex(subjectid)
	//log.Println(subjectidhex)
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

func (r *StudentAssessmentRepo) UpsertStudentAssessmentData(s SingleStudentAssessment) error {

	log.Println(s)

	assessment := StudentAssessment{}
	assessment.Name = s.AssessmentName
	assessment.Score = s.Score

	/*
		_, err := r.coll.Upsert(bson.M{
			"studentid": s.StudentID,
			"class":     s.Class,
			"subject":   s.Subject,
		}, bson.M{
			"$addToSet": bson.M{
				"assessments": bson.M{
					"name":  assessment.Name,
					"score": assessment.Score,
				},
			},

			"$setOnInsert": bson.M{
				"studentid": s.StudentID,
				"class":     s.Class,
				"subject":   s.Subject,
				"name":      s.Name,
			},
		})

	*/

	err := r.coll.Update(bson.M{
		"studentid":        s.StudentID,
		"class":            s.Class,
		"subject":          s.Subject,
		"assessments.name": assessment.Name,
	}, bson.M{
		"$set": bson.M{
			"assessments.$.score": assessment.Score,
		},
	})

	if err != nil {
		log.Println(err)

		_, err := r.coll.Upsert(bson.M{
			"studentid": s.StudentID,
			"class":     s.Class,
			"subject":   s.Subject,
		}, bson.M{
			"$push": bson.M{
				"assessments": assessment,
			},
			"$setOnInsert": bson.M{
				"studentid": s.StudentID,
				"class":     s.Class,
				"subject":   s.Subject,
			},
		},
		)

		if err != nil {
			log.Println(err)
			return err
		}

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

func (r *StudentAssessmentRepo) GetAssessmentsOfAStudent(student string) ([]StudentAssessments, error) {
	result := []StudentAssessments{}
	err := r.coll.Find(bson.M{
		"studentid": student,
	}).All(&result)
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
	u := SubjectRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_subjects")}
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
		//log.Println(student)
		s := StudentAssessments{}

		s.Name = student.FirstName + " " + student.LastName
		//log.Println(s.Name)
		s.StudentID = student.ID.Hex()

		//as := assessments.Assessments

		for _, a := range assessments {
			//log.Println(a.ID)
			//log.Println(student.ID)
			if a.StudentID == student.ID.Hex() {
				s.Assessments = append(s.Assessments, a.Assessments...)

				//assessments[i] = assessments[len(assessments)-1]
				//assessments = assessments[:len(assessments)-1]
			}

		}

		if len(s.Assessments) < 1 {
			s.Assessments = []StudentAssessment{}
		}
		//log.Println(s.Assessments)
		returnStudents = append(returnStudents, s)
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

	s := SingleStudentAssessment{}
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		log.Println(err)
	}

	err = sRepo.UpsertStudentAssessmentData(s)
	if err != nil {
		log.Println(err)
	}

}
func (c *Config) GetAssessmentsOfAStudentHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)
	sRepo := StudentAssessmentRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_assessments")}
	subjectRepo := SubjectRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_subjects")}

	studentID := r.URL.Query().Get("id")

	assessments, err := sRepo.GetAssessmentsOfAStudent(studentID)
	if err != nil {
		log.Println(err)
	}

	for i := range assessments {
		subject, err := subjectRepo.GetByName(assessments[i].Subject)
		if err != nil {
			log.Println(err)
		}
		assessments[i].SubjectInfo = subject
	}

	err = json.NewEncoder(w).Encode(assessments)
	if err != nil {
		log.Println(err)
	}
}

//GetAssessmentsOfAStudentMobile get assesment of a student for mobile
func (c *Config) GetAssessmentsOfAStudentMobile(w http.ResponseWriter, r *http.Request) {

	sRepo := StudentAssessmentRepo{c.MongoSession.DB(c.MONGODB).C(r.URL.Query().Get("school") + "_assessments")}
	subjectRepo := SubjectRepo{c.MongoSession.DB(c.MONGODB).C(r.URL.Query().Get("school") + "_subjects")}

	studentID := r.URL.Query().Get("id")

	assessments, err := sRepo.GetAssessmentsOfAStudent(studentID)
	if err != nil {
		log.Println(err)
	}

	for i := range assessments {
		subject, err := subjectRepo.GetByName(assessments[i].Subject)
		if err != nil {
			log.Println(err)
		}
		assessments[i].SubjectInfo = subject
	}

	err = json.NewEncoder(w).Encode(assessments)
	if err != nil {
		log.Println(err)
	}
}
