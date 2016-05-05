package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"

	//	"github.com/gorilla/context"

	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Guardian struct
type Guardian struct {
	Name       string `json:"name"`
	Phone      string `json:"phone" bson:"_id,omitempty"`
	Email      string `json:"email"`
	Occupation string `json:"occupation"`
	Address    string `json:"address"`
	Country    string `json:"country"`
	State      string `json:"state"`
	Pin2       string `json:"Pin2"`
	True       string
	Pin        []byte   `json:"-" bson:"pin"`
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
func randSe(n int) string {
	letters := []rune("abcdefghij123klmnop890qrstuvwxyz4567")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

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

func (r *GuardianRepo) RequestPin(rr string) error {
	user := Guardian{}
	pin := randSe(6)
	err := r.coll.Find(bson.M{
		"_id": rr,
	}).One(&user)
	if err != nil {
		log.Println(err)
		return err
	}

	if user.Phone != "" {
		user.Pin, _ = bcrypt.GenerateFromPassword([]byte(pin), Cost)
		r.coll.Update(bson.M{
			"_id": user.Phone,
		}, bson.M{
			"$set": bson.M{
				"pin": user.Pin,
			},
		})
		req, err := http.Get("http://api.clickatell.com/http/sendmsg?" + "api_id=3596410&user=digitalforte&password=digitalforte9!!&to=" + user.Phone + "&text=" + url.QueryEscape("Welcome to Edna Iparent Mobile App use the Pin below to sign in "+pin))
		if err != nil {
			fmt.Printf("%s", err)
			return err
		}
		defer req.Body.Close()
		contents, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Printf("%s", err)
			return err
		}
		fmt.Printf("%s\n", string(contents))

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

func (r *GuardianRepo) AuthGuardian(rr string, x string) (Guardian, error) {
	var guardian Guardian
	err := r.coll.Find(bson.M{
		"_id": rr,
	}).One(&guardian)
	if err != nil {
		log.Println(err)
		return guardian, err
	}
	err = bcrypt.CompareHashAndPassword(guardian.Pin, []byte(x))
	if err != nil {
		return guardian, err
	}
	guardian.True = "true"
	return guardian, nil
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

func (c *Config) VerifyGuardian(w http.ResponseWriter, r *http.Request) {
	tmp := r.URL.Query().Get("no")
	log.Println(tmp)
	u := GuardianRepo{c.MongoSession.DB(c.MONGODB).C("guardians")}
	err := u.RequestPin(tmp)
	if err != nil {
		log.Println(err)
	}
}

func (c *Config) AuthGuardianHandler(w http.ResponseWriter, r *http.Request) {
	tmp := r.URL.Query().Get("no")
	u := GuardianRepo{c.MongoSession.DB(c.MONGODB).C("guardians")}
	//guardian := Guardian{}
	//buf := new(bytes.Buffer)
	//buf.ReadFrom(r.Body)
	//s := buf.String()
	//log.Println(s)

	/*	err := json.NewDecoder(r.Body).Decode(&guardian)
		if err != nil {
			log.Println(err)
		}
	*/

	//log.Println(guardian.Pin2)
	log.Println(r.FormValue("Pin2"))
	guardian, err := u.AuthGuardian(tmp, r.FormValue("Pin2"))
	if err != nil {
		log.Println(err)
	}

	res, err := json.Marshal(guardian)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

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
