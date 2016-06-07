package main

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2"
)

type State struct {
	Slug  string `bson:"slug"`
	State string `bson:"state"`
}

type Lga struct {
	Lga   string `bson:"lga"`
	State string `bson:"state"`
}
type ParseRepo struct {
	coll *mgo.Collection
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

func (c *Config) ParseHandler(w http.ResponseWriter, r *http.Request) {
	u := ParseRepo{c.MongoSession.DB(c.MONGODB).C("lga")}
	u.AddCsv()
	y := ParseRepo{c.MongoSession.DB(c.MONGODB).C("state")}
	y.AddCsvS()
}
