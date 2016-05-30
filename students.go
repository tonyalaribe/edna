package main

import (
	"encoding/json"

	"github.com/gorilla/context"

	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"net/http"
)

// Student struct
type Student struct {
	ID                          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	SignUpDate                  string        `json:"signupdate,omitempty" bson:",omitempty"`
	Class                       string        `json:"class, omitempty" bson:",omitempty"`
	FirstName                   string        `json:"firstname,omitempty" bson:",omitempty"`
	MiddleName                  string        `json:"middlename,omitempty" bson:",omitempty"`
	LastName                    string        `json:"lastname,omitempty" bson:",omitempty"`
	DateOfBirth                 string        `json:"dateofbirth,omitempty" bson:",omitempty"`
	Gender                      string        `json:"gender,omitempty" bson:",omitempty"`
	BloodGroup                  string        `json:"bloodgroup,omitempty" bson:",omitempty"`
	Nationality                 string        `json:"nationality,omitempty" bson:",omitempty"`
	StateOfOrigin               string        `json:"stateoforigin,omitempty" bson:",omitempty"`
	State                       string        `json:"state,omitempty" bson:",omitempty"`
	PermanentAddress            string        `json:"permanentaddress,omitempty" bson:",omitempty"`
	Country                     string        `json:"country,omitempty" bson:",omitempty"`
	City                        string        `json:"city,omitempty" bson:",omitempty"`
	Phone                       string        `json:"phone,omitempty" bson:",omitempty"`
	Email                       string        `json:"email,omitempty" bson:",omitempty"`
	PreviousSchoolName          string        `json:"previous_schoolname,omitempty" bson:",omitempty"`
	PreviousSchoolAddress       string        `json:"previous_schooladdress,omitempty" bson:",omitempty"`
	PreviousSchoolQualification string        `json:"previous_schoolqualification,omitempty" bson:",omitempty"`
	GuardianRelationship        string        `json:"guardian_relationship,omitempty" bson:",omitempty"`
	GuardianName                string        `json:"guardian_name,omitempty" bson:",omitempty"`
	GuardianMobile              string        `json:"guardian_mobile,omitempty" bson:",omitempty"`
	GuardianEmail               string        `json:"guardian_email,omitempty" bson:",omitempty"`
	GuardianOccupation          string        `json:"guardian_occupation,omitempty" bson:",omitempty"`
	GuardianAddress             string        `json:"guardian_address,omitempty" bson:",omitempty"`
	GuardianCountry             string        `json:"guardian_country,omitempty" bson:",omitempty"`
	GuardianState               string        `json:"guardian_state,omitempty" bson:",omitempty"`
	GuardianCity                string        `json:"guardian_city,omitempty" bson:",omitempty"`
	Guardian2Relationship       string        `json:"guardian2_relationship,omitempty" bson:",omitempty"`
	Guardian2Name               string        `json:"guardian2_name,omitempty" bson:",omitempty"`
	Guardian2Mobile             string        `json:"guardian2_mobile,omitempty" bson:",omitempty"`
	Guardian2Email              string        `json:"guardian2_email,omitempty" bson:",omitempty"`
	Guardian2Occupation         string        `json:"guardian2_occupation,omitempty" bson:",omitempty"`
	Guardian2Address            string        `json:"guardian2_address,omitempty" bson:",omitempty"`
	Guardian2Country            string        `json:"guardian2_country,omitempty" bson:",omitempty"`
	Guardian2State              string        `json:"guardian2_state,omitempty" bson:",omitempty"`
	Guardian2City               string        `json:"guardian2_city,omitempty" bson:",omitempty"`
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

//Get gets a student's details from db
func (r *StudentRepo) Get(id string) (Student, error) {
	var student Student
	err := r.coll.Find(bson.M{
		"_id": bson.ObjectIdHex(id),
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

//GetAllStudentsInClass gets all user from db
func (r *StudentRepo) GetAllStudentsInClass(class string) ([]Student, error) {
	var student []Student
	err := r.coll.Find(bson.M{
		"class": class,
	}).All(&student)

	if err != nil {
		log.Println(err)
		return student, err
	}

	return student, nil
}

//GetAllStudentsInParentClass gets all students in all clases under a parent
func (r *StudentRepo) GetAllStudentsInParentClass(class string) ([]Student, error) {
	students := []Student{}
	students, err := r.GetAllStudentsInClass(class)
	if err != nil {
		log.Println(err)
		return students, err
	}

	classes := []Class{}
	err = r.coll.Find(bson.M{
		"parent": class,
	}).All(&classes)
	if err != nil {
		log.Println(err)
		return students, err
	}

	for _, c := range classes {
		s, err := r.GetAllStudentsInClass(c.Name)
		if err != nil {
			log.Println(err)
			return students, err
		}

		students = append(students, s...)

	}

	return students, nil
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
	tmp := []string{}
	tmp = append(tmp, school.ID)

	guardian1 := Guardian{
		Name:       student.GuardianName,
		Phone:      student.GuardianMobile,
		Email:      student.GuardianEmail,
		Occupation: student.GuardianOccupation,
		Address:    student.GuardianAddress,
		Country:    student.GuardianCountry,
		State:      student.GuardianState,
		City:       student.GuardianCity,
		Schools:    tmp,
	}

	guardian2 := Guardian{
		Name:       student.Guardian2Name,
		Phone:      student.Guardian2Mobile,
		Email:      student.Guardian2Email,
		Occupation: student.Guardian2Occupation,
		Address:    student.Guardian2Address,
		Country:    student.Guardian2Country,
		State:      student.Guardian2State,
		City:       student.Guardian2City,
	}
	g := GuardianRepo{c.MongoSession.DB(c.MONGODB).C("guardians")}
	if guardian1.Phone != "" {
		err = g.Create(&guardian1, school.ID)
		if err != nil {
			log.Println(err)
		}
	}
	if guardian2.Phone != "" {
		err = g.Create(&guardian2, school.ID)
		if err != nil {
			log.Println(err)
		}
	}

}

//getStudentsHandler would return students
func (c *Config) getStudentsHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)

	u := StudentRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_students")}
	students, err := u.GetAll()

	err = json.NewEncoder(w).Encode(StudentCollection{students})
	if err != nil {
		log.Println(err)
	}
}

//getStudentHandler will return a student
func (c *Config) getStudentHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)
	id := r.URL.Query().Get("id")

	u := StudentRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_students")}
	student, err := u.Get(id)

	err = json.NewEncoder(w).Encode(StudentData{student})
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

//getStudentsHandler would create a student
func (c *Config) getStudentsInClassHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)

	class := r.URL.Query().Get("class")
	u := StudentRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_students")}
	students, err := u.GetAllStudentsInParentClass(class)
	if err != nil {
		log.Println(err)
	}
	err = json.NewEncoder(w).Encode(StudentCollection{students})
	if err != nil {
		log.Println(err)
	}
}

//getClassDataHandler would create a student
func (c *Config) getClassDataHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)

	class := r.URL.Query().Get("class")
	u := StudentRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_students")}
	students, err := u.GetAllStudentsInParentClass(class)
	if err != nil {
		log.Println(err)
	}

	result := struct {
		Count int    `json:"count"`
		Name  string `json:"name"`
	}{
		Count: len(students),
		Name:  class,
	}
	err = json.NewEncoder(w).Encode(&result)
	if err != nil {
		log.Println(err)
	}
}
