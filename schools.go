package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"

	"github.com/gorilla/context"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"net"
	"net/http"
	//"net/smtp"
	"net/url"
)

// School struct holds information about each users skills, aids in marshalling to json and storing on the database
type School struct {
	ID              string `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string `json:"name"`
	Address         string `json:"address"`
	Domain          string `json:"domain"`
	AdminPhone      string `json:"adminphone"`
	AdminName       string `json:"adminname"`
	AdminEmail      string `json:"adminemail"`
	Password        string `json:"password"`
	AdminPassword   []byte `json:"adminpassword"`
	VerificationKey string `json:"verificationKey"`
	Verified        bool   `json:"verified"`
}

type recacptchaResponse struct {
	Success     bool   `json:"success"`
	ChallengeTs string `json:"challenge_ts"` // timestamp of the challenge load (ISO format yyyy-MM-dd'T'HH:mm:ssZZ)
	Hostname    string `json:"hostname"`
}

//SchoolCollection holds a slice of School structs within a Data key, to conform with the json api schema spec
type SchoolCollection struct {
	Schools []School `json:"schools"`
}

//SchoolData acts like SchoolData but carries information about a single school
type SchoolData struct {
	School School `json:"school"`
}

//SchoolRepo a mongo Collection that could get passed around
type SchoolRepo struct {
	coll *mgo.Collection
}

/* THese are functions that perform the operations on the db. .they are usually,
called by the handlers, in a bid to keep  handlers simple and less bulky.
*/

//Create adds a skill to the database, based on it's owner
func (r *SchoolRepo) Create(school *School) error {
	id := school.ID

	phash, err := bcrypt.GenerateFromPassword([]byte(school.Password), Cost)
	if err != nil {
		log.Println(err)
		return err
	}
	school.AdminPassword = phash
	school.Password = ""
	school.VerificationKey = randSeq(20)
	_, err = r.coll.UpsertId(id, school)
	if err != nil {
		return err
	}

	return nil
}

//Update updates a school in the database
func (r *SchoolRepo) Update(school *School) error {

	err := r.coll.UpdateId(school.ID, bson.M{
		"$set": bson.M{
			"name":       school.Name,
			"address":    school.Address,
			"adminname":  school.AdminName,
			"adminemail": school.AdminEmail,
			"adminphone": school.AdminPhone,
		},
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//Verify Completes the verification Process from the link and mark school as verified
func (r *SchoolRepo) Verify(adminEmail string, schoolID string, verificationKey string, rootURL string) (School, error) {
	school := School{}
	err := r.coll.Find(bson.M{
		"adminemail": adminEmail,
		"_id":        schoolID,
	}).One(&school)
	if err != nil {
		return school, err
	}

	log.Println(school.VerificationKey + " vs " + verificationKey)

	if school.VerificationKey == verificationKey {
		var domain string
		if rootURL != "localhost:8080" {
			//host, _, _ := net.SplitHostPort(rootURL)
			domain = school.ID + "." + rootURL
		} else {
			domain = rootURL + "/dashboard"
		}
		log.Println(domain)
		r.coll.Update(bson.M{
			"adminemail": adminEmail,
			"_id":        schoolID,
		}, bson.M{
			"$set": bson.M{
				"verified": true,
				"domain":   domain,
			},
		})
		school.Domain = domain
		return school, nil
	}
	return school, errors.New("wrong verification key")
}

/***********************************************
HANDLERS
***********************************************/

//NewSchool is a handler for registering new schoolls
func (c *Config) NewSchool(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	school := School{}

	school.ID = r.FormValue("school_id")
	school.Name = r.FormValue("school_name")
	school.AdminName = r.FormValue("admin_name")
	school.Password = r.FormValue("password")
	school.AdminPhone = r.FormValue("admin_phone")
	school.AdminEmail = r.FormValue("admin_email")

	log.Println(school)

	//recaptcha verfication

	recaptcha := r.FormValue("g-recaptcha-response")

	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{"secret": {"6Lf5yBsTAAAAANVM9JnJ8u8mFCg9t4clPSCvY65Z"}, "response": {recaptcha}})
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/registrationerror.html", http.StatusFound)
		return
	}

	rResponse := recacptchaResponse{}
	err = json.NewDecoder(resp.Body).Decode(&rResponse)
	if err != nil {
		log.Println(err)
	}

	if rResponse.Success == false {
		http.Redirect(w, r, "/registrationerror.html", http.StatusInternalServerError)
		return
	}

	//recaptcha verification ends here

	x := SchoolRepo{c.MongoSession.DB(c.MONGODB).C("schools")}

	err = x.Create(&school)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/registrationerror.html", http.StatusInternalServerError)
		return
	}
	verificationURL := c.RootURL + "/verify?key=" + school.VerificationKey + "&email=" + school.AdminEmail + "&id=" + school.ID
	log.Println(verificationURL)

	client := &http.Client{}

	// ...

	str := `{"to":{"` + school.AdminEmail + `":"` + school.AdminName + `"}, "from":["anthonyalaribe@gmail.com","Edna - School Management System"], "subject":"Edna: Verify your Account", "html":"You created a School named <strong>` + school.Name + `</strong>. Please click the verification link below,  to verify your account.<br/> <a href='` + verificationURL + `'>` + verificationURL + `</a>"}`

	mesg := bytes.NewReader([]byte(str))

	req, err := http.NewRequest("POST", "https://api.sendinblue.com/v2.0/email", mesg)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/registrationerror.html", http.StatusInternalServerError)
		return
	}

	req.Header.Add("api-key", "2BsIqZ9XWMp6YKUk")
	_, err = client.Do(req)

	if err != nil {

		log.Println(err)
		http.Redirect(w, r, "/registrationerror.html", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/success.html", http.StatusFound)
}

//ValidateReg validates registration
func (r *SchoolRepo) ValidateReg(schoolID string) (string, error) {
	var htm string
	school := School{}
	err := r.coll.Find(bson.M{
		"_id": schoolID,
	}).One(&school)
	if err != nil {
		htm = "<p class='text-success'>School Id is Available</p>"
		return htm, err
	}
	htm = "<p class='text-danger'>This School Id already exists</p>"
	return htm, err
}

//CheckEmail checcks email
func (r *SchoolRepo) CheckEmail(email string) (string, error) {
	var htm string
	school := School{}
	err := r.coll.Find(bson.M{
		"adminemail": email,
	}).One(&school)
	if err != nil {
		htm = ""
		return htm, err
	}
	htm = "<p class='text-danger'>This Email is already associated with another account  " + school.ID + "+</p>"
	return htm, err
}

//CheckEmailHandler checks email
func (c *Config) CheckEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("id")
	x := SchoolRepo{c.MongoSession.DB(c.MONGODB).C("schools")}
	msg, err := x.CheckEmail(email)
	if err != nil {
		log.Println(err)

	}
	w.Write([]byte(msg))
}

//ValidateRegHandler validates registration
func (c *Config) ValidateRegHandler(w http.ResponseWriter, r *http.Request) {
	schoolID := r.URL.Query().Get("id")
	x := SchoolRepo{c.MongoSession.DB(c.MONGODB).C("schools")}
	school, err := x.ValidateReg(schoolID)
	if err != nil {
		log.Println(err)

	}
	w.Write([]byte(school))
}

//VerifySchool handles verifying users who follow a link sent to their email upon registration
func (c *Config) VerifySchool(w http.ResponseWriter, r *http.Request) {
	verificationKey := r.URL.Query().Get("key")
	adminEmail := r.URL.Query().Get("email")
	schoolID := r.URL.Query().Get("id")
	x := SchoolRepo{c.MongoSession.DB(c.MONGODB).C("schools")}
	school, err := x.Verify(adminEmail, schoolID, verificationKey, c.RootURL)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/verificationerror.html", http.StatusNotAcceptable)
		return
	}

	u := UserRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_users")}
	err = u.CreateAdmin(&school)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/verificationerror.html", http.StatusNotAcceptable)
		return
	}

	log.Println(school.Domain)
	http.Redirect(w, r, "http://"+school.Domain, http.StatusFound)
}

//UpdateSchoolHandler updates the school details
func (c *Config) UpdateSchoolHandler(w http.ResponseWriter, r *http.Request) {

	school := School{}
	json.NewDecoder(r.Body).Decode(&school)
	x := SchoolRepo{c.MongoSession.DB(c.MONGODB).C("schools")}
	err := x.Update(&school)
	if err != nil {
		log.Println(err)
	}
}

//GetSchoolHandler updates the school details
func (c *Config) GetSchoolHandler(w http.ResponseWriter, r *http.Request) {

	school := context.Get(r, "school").(School)
	err := json.NewEncoder(w).Encode(school)
	if err != nil {
		log.Println(err)
	}

}
