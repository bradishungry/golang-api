package api

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/bradishungry/golang-basic-api/internal/db"

	_ "github.com/lib/pq"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const codeLength = 4

type Person struct {
	Name  string
	Email string
	Age   int
	Id    int
	Uuid  string
}

// Took this code from a stackoverflow question, helped greatly for cors issues
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

func queryPersons() ([]Person, error) {
	var dbhandle = db.GetDBHandle()
	rows, err := dbhandle.Query("SELECT * FROM persons")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ppl []Person

	for rows.Next() {
		var ps Person

		err := rows.Scan(&ps.Name, &ps.Email, &ps.Age, &ps.Uuid, &ps.Id)
		if err != nil {
			return nil, err
		}

		ppl = append(ppl, ps)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ppl, nil
}

func insertPersons(ps Person) error {
	var dbhandle = db.GetDBHandle()
	rows, err := dbhandle.Query("INSERT INTO persons (person_name, email, age) values ($1, $2, $3)", ps.Name, ps.Email, ps.Age)
	if err != nil {
		return err
	}
	defer rows.Close()
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func removePerson(pid int) error {
	var dbhandle = db.GetDBHandle()
	rows, err := dbhandle.Query("DELETE FROM persons WHERE person_pk = $1", pid)
	if err != nil {
		return err
	}
	defer rows.Close()
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func selectPersonTest() (Person, error) {
	var dbhandle = db.GetDBHandle()
	rows, err := dbhandle.Query("select * FROM persons order by person_pk desc limit 1")
	if err != nil {
		return Person{}, err
	}
	defer rows.Close()
	var ppl []Person

	for rows.Next() {
		var ps Person

		err := rows.Scan(&ps.Name, &ps.Email, &ps.Age, &ps.Uuid, &ps.Id)
		if err != nil {
			return Person{}, err
		}

		ppl = append(ppl, ps)
	}
	if err = rows.Err(); err != nil {
		return Person{}, err
	}

	return ppl[0], nil
}

func updatePersons(ps Person) error {
	var dbhandle = db.GetDBHandle()
	rows, err := dbhandle.Query(
		"UPDATE persons SET person_name = $1, email = $2, age = $3 WHERE person_pk = $4;",
		ps.Name, ps.Email, ps.Age, ps.Id)
	if err != nil {
		return err
	}
	defer rows.Close()
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func GetCodeEP(w http.ResponseWriter, req *http.Request) {
	b := make([]byte, codeLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	fmt.Fprintf(w, string(b))
}

func GetPeopleEP(w http.ResponseWriter, req *http.Request) {
	persondata, err := queryPersons()
	if err != nil {
		log.Print(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(persondata)
}

func AddPeopleEP(w http.ResponseWriter, req *http.Request) {
	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	}

	//log.Println(formatRequest(req))

	var p Person
	errD := json.NewDecoder(req.Body).Decode(&p)

	if errD != nil {
		http.Error(w, errD.Error(), http.StatusBadRequest)
		log.Println("Error: " + errD.Error())
		return
	}
	err := insertPersons(p)
	if err != nil {
		log.Print(err)
	}
}

func DeletePersonEP(w http.ResponseWriter, req *http.Request) {
	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	}

	//log.Println(formatRequest(req))

	var pid Person
	errD := json.NewDecoder(req.Body).Decode(&pid)

	if errD != nil {
		http.Error(w, errD.Error(), http.StatusBadRequest)
		log.Println("Error: " + errD.Error())
		return
	}
	err := removePerson(pid.Id)
	if err != nil {
		log.Print(err)
	}
}

func UpdatePersonEP(w http.ResponseWriter, req *http.Request) {
	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	}

	//log.Println(formatRequest(req))

	var ps Person
	errD := json.NewDecoder(req.Body).Decode(&ps)

	if errD != nil {
		http.Error(w, errD.Error(), http.StatusBadRequest)
		log.Println("Error: " + errD.Error())
		return
	}
	err := updatePersons(ps)
	if err != nil {
		log.Print(err)
	}
}
