package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	//"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// User struct holds information about each users skills, aids in marshalling to json and storing on the database
type User struct {
	ID          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:",omitempty"`
	Type        string        `json:"type" bson:",omitempty"`
	Image       string        `json:"image" bson:",omitempty"`
	UpdateImage string        `json:"updateimage" bson:",omitempty"`
	Phone       string        `json:"phone" bson:",omitempty"`
	Email       string        `json:"email" bson:",omitempty"`
	P           string        `json:"password" bson:",omitempty"`
	Password    []byte        `json:"-" bson:",omitempty"`
	Roles       []string      `json:"roles" bson:",omitempty"`
}

//UserCollection holds a slice of User structs within a Data key, to conform with the json api schema spec
type UserCollection struct {
	Users []User `json:"users"`
}

//UserData acts like SchoolData but carries information about a single school
type UserData struct {
	User User `json:"user"`
}

//UserRepo a mongo Collection that could get passed around
type UserRepo struct {
	coll *mgo.Collection
}

/* THese are functions that perform the operations on the db. .they are usually,
called by the handlers, in a bid to keep  handlers simple and less bulky.
*/

//CreateAdmin adds an admin to the database
func (r *UserRepo) CreateAdmin(school *School) error {
	user := User{}
	user.Name = school.AdminName
	user.Image = "/img/avatar.png"
	user.Email = school.AdminEmail
	user.Password = school.AdminPassword
	user.Phone = school.AdminPhone
	user.Roles = []string{"admin"}

	err := r.coll.Insert(user)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//Create adds a user to the database
func (r *UserRepo) Create(user *User) error {
	phash, err := bcrypt.GenerateFromPassword([]byte(user.P), Cost)
	if err != nil {
		log.Println(err)
		return err
	}
	user.Password = phash
	user.P = ""
	user.Roles = []string{"staff"}

	if user.Type == "Teaching Staff" {
		user.Roles = append(user.Roles, "teacher")
	}

	avatars := []string{"/img/avatars/avatar.png", "/img/avatars/avatar2.jpg", "/img/avatars/avatar3.jpg", "/img/avatars/avatar4.png", "/img/avatars/avatar5.png"}
	user.Image = avatars[rand.Intn(len(avatars))]

	err = r.coll.Insert(user)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//Update adds a user to the database
func (r *UserRepo) Update(user *User) error {
	//log.Println(user)
	err := r.coll.UpdateId(user.ID, bson.M{
		"$set": bson.M{
			"name":  user.Name,
			"email": user.Email,
			"phone": user.Phone,
			"image": user.Image,
			"type":  user.Type,
			"roles": user.Roles,
		},
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//Get gets a users details from db
func (r *UserRepo) Get(username string) (User, error) {
	var user User
	err := r.coll.Find(bson.M{
		"$or": []bson.M{
			bson.M{"email": username},
			bson.M{"phone": username},
		},
	}).One(&user)

	if err != nil {
		log.Println(err)
		return user, err
	}

	return user, nil
}

//GetAll gets all user from db
func (r *UserRepo) GetAll() ([]User, error) {
	var users []User
	err := r.coll.Find(bson.M{}).All(&users)

	if err != nil {
		log.Println(err)
		return users, err
	}

	return users, nil
}

/***************
handlers
***************/

//LoginPost post would taake care of authentication, and return json
func (c *Config) LoginPost(w http.ResponseWriter, r *http.Request) {

	x := struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Remember bool   `json:"remember"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&x)
	if err != nil {
		log.Println(err)
	}

	school := context.Get(r, "school").(School)

	u := UserRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_users")}

	user, err := u.Get(x.Username)
	if err != nil {
		log.Println(err)
		//http.Redirect(w, r, "/login.html", http.StatusNotFound)
		WriteError(w, ErrNotFound)
		return
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(x.Password))
	if err != nil {
		log.Println(err)
		WriteError(w, ErrWrongPassword)
		//http.Redirect(w, r, "/login.html", http.StatusNotFound)
		return
	}

	// create a signer for rsa 256
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	// set our claims
	t.Claims["AccessToken"] = user.Roles
	t.Claims["User"] = user
	t.Claims["UserID"] = user.ID.Hex()

	//log.Println(user.ID.Hex())

	//log.Println("the tokened user is", user)
	// set the expire time
	// see http://tools.ietf.org/html/draft-ietf-oauth-json-web-token-20#section-4.1.4
	if x.Remember {
		t.Claims["exp"] = time.Now().Add(time.Hour * 3000).Unix()
	} else {
		t.Claims["exp"] = time.Now().Add(time.Hour * 10).Unix()
	}
	tokenString, err := t.SignedString(c.Private)

	if err != nil {
		WriteError(w, ErrInternalServer)
		log.Printf("Token Signing error: %v\n", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:       c.Token,
		Value:      tokenString,
		Path:       "/",
		RawExpires: "0",
	})
	//w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := struct {
		User    User   `json:"user"`
		Message string `json:"message"`
		Token   string `json:"token"`
	}{
		User:    user,
		Message: "Token succesfully generated",
		Token:   tokenString,
	}
	//log.Println(response)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
	}

}

//getMeHandler would get logged in users details
func (c *Config) getMeHandler(w http.ResponseWriter, r *http.Request) {
	user, err := userget(r)

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Println(err)
	}

}

//createUserHandler would create a user/staff
func (c *Config) createUserHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)
	u := UserRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_users")}
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
	}
	err = u.Create(&user)
	if err != nil {
		log.Println(err)
	}

}

//updateUserHandler would update a user/staff
func (c *Config) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)
	u := UserRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_users")}
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
	}

	if user.UpdateImage != "" {
		//log.Println(user.UpdateImage)

		bucket := c.S3Bucket

		byt, err := base64.StdEncoding.DecodeString(strings.Split(user.UpdateImage, "base64,")[1])
		if err != nil {
			log.Println(err)
		}

		meta := strings.Split(user.UpdateImage, "base64,")[0]
		newmeta := strings.Replace(strings.Replace(meta, "data:", "", -1), ";", "", -1)
		//imagename := randSeq(30)

		err = bucket.Put(school.ID+"/"+user.ID.Hex(), byt, newmeta, s3.PublicReadWrite)
		if err != nil {
			log.Println(err)
		}

		//log.Println(bucket.URL(school.ID + "/" + user.ID.Hex()))

		user.Image = bucket.URL(school.ID + "/" + user.ID.Hex())
	}
	err = u.Update(&user)
	if err != nil {
		log.Println(err)
	}
}

//getUsersHandler would create a user/staff
func (c *Config) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)
	//log.Println(school)
	u := UserRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_users")}
	users, err := u.GetAll()

	err = json.NewEncoder(w).Encode(UserCollection{users})
	if err != nil {
		log.Println(err)
	}
}
