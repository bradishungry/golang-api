package postgres

import "github.com/google/uuid"

type Person struct {
	Name  string
	Email string
	Age   int
	Id    int
	Uuid  uuid.UUID
}
