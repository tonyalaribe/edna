package main

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"log"
	"os"
)

var (
	MONGOSERVER string
	MONGODB     string
)

//Config struct is important, as all handlers will  be methods implementing
//this struct.  THis technique will make it easy to pass common resources like a
// database connection, between connections
type Config struct {
	MONGOSERVER string
	MONGODB     string

	MongoSession *mgo.Session

	Public   []byte
	Private  []byte
	RootURL  string
	S3Bucket *s3.Bucket
	Token    string
}

func generateConfig() (config Config) {

	config.MONGOSERVER = os.Getenv("MONGOLAB_URI")

	if config.MONGOSERVER == "" {
		log.Println("No mongo server address set, resulting to default address")
		config.MONGOSERVER = "localhost:27017"
	}
	MONGOSERVER = config.MONGOSERVER
	log.Println("MONGOSERVER is ", config.MONGOSERVER)

	config.MONGODB = os.Getenv("MONGODB")
	if config.MONGODB == "" {
		log.Println("No Mongo database name set, resulting to default")
		config.MONGODB = "edna"
	}
	MONGODB = config.MONGODB
	log.Println("MONGODB is ", config.MONGODB)

	session, err := mgo.Dial(config.MONGOSERVER)
	if err != nil {
		panic(err)
		//log.Println(err)
	}
	session.SetMode(mgo.Monotonic, true)

	config.MongoSession = session

	AWSBucket := os.Getenv("AWSBucket")
	if AWSBucket == "" {
		log.Println("No AWSBucket set, resulting to default")
		AWSBucket = "edna"
	}
	log.Println("AWS Bucket is ", AWSBucket)

	auth, err := aws.EnvAuth()
	if err != nil {
		//panic(err)
		log.Println("no aws ish")
	}
	s := s3.New(auth, aws.USWest2)
	s3bucket := s.Bucket(AWSBucket)

	config.S3Bucket = s3bucket

	config.Public, err = ioutil.ReadFile("app.rsa.pub")
	if err != nil {
		log.Fatal("Error reading public key")
		return
	}

	config.Private, err = ioutil.ReadFile("app.rsa")
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}

	config.RootURL = os.Getenv("RootURL")
	if config.RootURL == "" {
		config.RootURL = "http://localhost:8080"
	}

	config.Token = "AccessToken"

	return
}
