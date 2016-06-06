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

	}

	return nil
}

func (c *Config) ParseHandler(w http.ResponseWriter, r *http.Request) {
	u := ParseRepo{c.MongoSession.DB(c.MONGODB).C("state")}
	u.AddCsv()
}
