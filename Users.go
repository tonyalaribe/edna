package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// User struct holds information about each users skills, aids in marshalling to json and storing on the database
type User struct {
	ID       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string        `json:"name"`
	Type     string        `json:"type"`
	Image    string        `json:"image"`
	Phone    string        `json:"phone"`
	Email    string        `json:"email"`
	P        string        `json:"password"`
	Password []byte        `json:"-"`
	Roles    []string      `json:"roles"`
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

	avatars := []string{"/img/avatar.png", "/img/avatar2.jpg", "/img/avatar5.png"}
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

	err := r.coll.UpdateId(user.ID, user)
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

//updateUserHandler would create a user/staff
func (c *Config) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)
	u := UserRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_users")}
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
	}
	err = u.Update(&user)
	if err != nil {
		log.Println(err)
	}
}

//getUsersHandler would create a user/staff
func (c *Config) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	school := context.Get(r, "school").(School)
	log.Println(school)
	u := UserRepo{c.MongoSession.DB(c.MONGODB).C(school.ID + "_users")}
	users, err := u.GetAll()

	err = json.NewEncoder(w).Encode(UserCollection{users})
	if err != nil {
		log.Println(err)
	}
}
