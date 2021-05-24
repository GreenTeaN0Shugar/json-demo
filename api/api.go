package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	IsMale bool   `json:"isMale"`
}

var currentID = 1
var users = []User{}

func ServeHTTP() {

	fmt.Println("Server is about to start...")

	http.HandleFunc("/users/", UserHandler)
	http.ListenAndServe(":8080", nil)
}

func UserHandler(w http.ResponseWriter, r *http.Request) {

	id := -1

	if r.URL.Path != "/users/" {
		idRaw := r.URL.Path[len("/users/"):]

		var err error
		id, err = strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, "cannot parse ID body: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:

		if id == -1 {
			GetUsers(w, r)
			return
		}

		GetUser(id, w, r)
		return

	case http.MethodPost:
		CreateUser(w, r)
		return

	case http.MethodDelete:
		DeleteUser(id, w, r)
		return

	case http.MethodPut:
		UpdateUser(id, w, r)
		return
	}
}

func GetUsers(writter http.ResponseWriter, request *http.Request) {

	bodyBytes, err := json.Marshal(users)
	if err != nil {
		http.Error(writter, "cannot read body: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writter.Write(bodyBytes)
}

func GetUser(id int, writter http.ResponseWriter, request *http.Request) {
	i := getIndexByID(id)
	if i == -1 {
		writter.WriteHeader(http.StatusNotFound)
		return
	}
	u := users[i]

	bodyBytes, err := json.Marshal(u)
	if err != nil {
		http.Error(writter, "cannot read body: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writter.Write(bodyBytes)
}

func CreateUser(writter http.ResponseWriter, request *http.Request) {

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writter, "cannot read body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	u := &User{}
	err = json.Unmarshal(body, u)
	if err != nil {
		http.Error(writter, "cannot unmarshal body: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(string(body))
	fmt.Printf("%#v\n", u)

	u.ID = currentID
	users = append(users, *u)
	currentID++
	writter.WriteHeader(http.StatusCreated)
}

func DeleteUser(id int, writter http.ResponseWriter, request *http.Request) {
	i := getIndexByID(id)
	if i == -1 {
		writter.WriteHeader(http.StatusNotFound)
		return
	}

	users = append(users[:i], users[i+1:]...)
	writter.WriteHeader(http.StatusOK)
}

func UpdateUser(id int, writter http.ResponseWriter, request *http.Request) {
	i := getIndexByID(id)
	if i == -1 {
		writter.WriteHeader(http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writter, "cannot read body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userNew := &User{}
	err = json.Unmarshal(body, userNew)
	if err != nil {
		http.Error(writter, "cannot unmarshal body: "+err.Error(), http.StatusBadRequest)
		return
	}

	users[i].Age = userNew.Age
	users[i].Name = userNew.Name
	users[i].IsMale = userNew.IsMale

	fmt.Println(string(body))
	fmt.Printf("%#v\n", userNew)

	writter.WriteHeader(http.StatusOK)
}

func getIndexByID(userID int) int {
	for i, u := range users {
		if u.ID == userID {
			return i
		}
	}

	return -1
}
