package main

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type State struct {
	Slug  string `bson:"slug"`
	State string `bson:"state"`
}

type Lga struct {
	Lga   string `bson:"lga"`
	State string `bson:"state"`
}

type Country struct {
	Country string `bson:"country"`
}

type ParseRepo struct {
	coll *mgo.Collection
}

func (rp *ParseRepo) AddCountry() error {

	dat, err := ioutil.ReadFile("csv/countries.csv")
	if err != nil {
		log.Fatal(err)
		return err
	}
	r := csv.NewReader(strings.NewReader(string(dat)))
	records, errs := r.ReadAll()
	if errs != nil {
		log.Fatal(errs)
	}
	log.Println(len(records[0]))
	//m := make(map[string]string)
	for i := 1; i < len(records[0]); i++ {
		count := new(Country)
		count.Country = records[0][i]
		rp.coll.Insert(count)
	}

	return nil
}

func (rp *ParseRepo) AddCsv() error {

	dat, err := ioutil.ReadFile("csv/lga.csv")
	if err != nil {
		log.Fatal(err)
		return err
	}
	r := csv.NewReader(strings.NewReader(string(dat)))
	records, errs := r.ReadAll()
	if errs != nil {
		log.Fatal(errs)
	}
	m := make(map[string]string)
	for i := 1; i < len(records); i++ {
		lga := new(Lga)
		lga.Lga = records[i][1]
		lga.State = strings.Replace(records[i][2], " ", "-", -1)
		m[lga.State] = records[i][2]
		rp.coll.Insert(lga)
	}
	for key, value := range m {
		state := new(State)
		state.Slug = key
		state.State = value

	}

	return nil
}

func (rp *ParseRepo) AddCsvS() error {

	dat, err := ioutil.ReadFile("csv/lga.csv")
	if err != nil {
		log.Fatal(err)
		return err
	}
	r := csv.NewReader(strings.NewReader(string(dat)))
	records, errs := r.ReadAll()
	if errs != nil {
		log.Fatal(errs)
	}
	m := make(map[string]string)
	for i := 1; i < len(records); i++ {
		m[strings.Replace(records[i][2], " ", "-", -1)] = records[i][2]
	}
	for key, value := range m {
		state := new(State)
		state.Slug = key
		state.State = value
		rp.coll.Insert(state)
	}

	return nil
}

/*func (rp *ParseRepo) GetStates() (string, error) {
	//  result, err :=
}*/
func (rp *ParseRepo) getStates() ([]State, error) {
	var state []State
	err := rp.coll.Find(bson.M{}).All(&state)
	if err != nil {
		log.Println(err)
	}
	return state, nil
}

func (rp *ParseRepo) getLgas(lga string) ([]Lga, error) {
	var lgas []Lga
	err := rp.coll.Find(bson.M{"state": lga}).All(&lgas)
	if err != nil {
		log.Println(err)
	}
	return lgas, nil
}

func (rp *ParseRepo) getCountries() ([]Country, error) {
	var cs []Country
	err := rp.coll.Find(bson.M{}).All(&cs)
	if err != nil {
		log.Println(err)
	}
	return cs, nil
}

func (c *Config) ParseHandler(w http.ResponseWriter, r *http.Request) {
	u := ParseRepo{c.MongoSession.DB(c.MONGODB).C("lga")}
	u.AddCsv()
	y := ParseRepo{c.MongoSession.DB(c.MONGODB).C("state")}
	y.AddCsvS()
}

func (c *Config) GetCountries(w http.ResponseWriter, r *http.Request) {
	u := ParseRepo{c.MongoSession.DB(c.MONGODB).C("countries")}
	data, err := u.getCountries()
	if err != nil {
		log.Println(err)
	}
	result, _ := json.Marshal(data)
	w.Write(result)

}

func (c *Config) GetStates(w http.ResponseWriter, r *http.Request) {
	u := ParseRepo{c.MongoSession.DB(c.MONGODB).C("state")}
	data, err := u.getStates()
	if err != nil {
		log.Println(err)
	}
	result, _ := json.Marshal(data)
	w.Write(result)

}

func (c *Config) GetLga(w http.ResponseWriter, r *http.Request) {
	u := ParseRepo{c.MongoSession.DB(c.MONGODB).C("lga")}
	data, err := u.getLgas(r.URL.Query().Get("q"))
	if err != nil {
		log.Println(err)
	}
	result, _ := json.Marshal(data)
	w.Write(result)

}

func (c *Config) Test(w http.ResponseWriter, r *http.Request) {
	u := ParseRepo{c.MongoSession.DB(c.MONGODB).C("countries")}
	u.AddCountry()
}
