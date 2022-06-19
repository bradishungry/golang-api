package api

import (
	"log"
	"testing"

	"github.com/bradishungry/golang-basic-api/internal/db"
)

func init() {
	err := db.InitDB("postgres://postgres@localhost?sslmode=disable")
	if err != nil {
		log.Fatalf("test init failed: %s", err)
	}

	_, errEx := db.GetDBHandle().Exec(`CREATE TABLE IF NOT EXISTS public.persons
	(
		person_name character varying(50) COLLATE pg_catalog."default" NOT NULL,
		email character varying(50) COLLATE pg_catalog."default" NOT NULL,
		age integer,
		person_id uuid NOT NULL DEFAULT gen_random_uuid(),
		person_pk SERIAL PRIMARY KEY
		)`)
	if errEx != nil {
		log.Fatalf("test schema creation failed: %s", errEx)
	}
}

func TestAddUser(t *testing.T) {
	person := Person{
		Name:  "Bradley",
		Age:   26,
		Email: "bboswellff6@gmail.com",
	}
	msg := insertPersons(person)
	log.Println(person.Name)
	dbdata, err := selectPersonTest()
	log.Println(dbdata.Name)
	if person.Name != dbdata.Name || err != nil {
		t.Fatalf(`SelectPerson Error: %s, InsertPerson Error: %s`, err, msg)
	}
}

func TestDeleteUser(t *testing.T) {
	person := Person{
		Name:  "UserDel",
		Age:   42,
		Email: "DelUser@gmail.com",
	}
	insertPersons(person)
	dbdata, err := selectPersonTest()
	del := removePerson(dbdata.Id)
	user, errtwo := selectPersonTest()
	if person.Email == user.Email || del != nil || errtwo != nil || err != nil {
		t.Fatalf(`SelectPerson Error: %s, InsertPerson Error: %s`, err, errtwo)
	}
}

func TestUpdateUser(t *testing.T) {
	person := Person{
		Name:  "UserUpdate",
		Age:   87,
		Email: "UpdateUser@gmail.com",
	}
	msg := insertPersons(person)
	dbdata, err := selectPersonTest()
	updatedPerson := Person{
		Name:  "UserUpdate",
		Age:   88,
		Email: "UpdateUser@yahoo.com",
		Id:    dbdata.Id,
	}
	errtwo := updatePersons(updatedPerson)
	user, errtwo := selectPersonTest()
	if user.Email == dbdata.Email || err != nil || errtwo != nil {
		t.Fatalf(`SelectPerson Error: %s, InsertPerson Error: %s`, err, msg)
	}
}
