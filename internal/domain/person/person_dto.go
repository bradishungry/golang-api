package person

import "github.com/google/uuid"

type PersonDTO struct {
	Name  string
	Email string
	Age   int
	Uuid  uuid.UUID
}
