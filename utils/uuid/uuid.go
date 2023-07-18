package uuid

import (
	"github.com/google/uuid"
)

func New() string {
	u, _ := uuid.NewRandom()
	return u.String()
}
