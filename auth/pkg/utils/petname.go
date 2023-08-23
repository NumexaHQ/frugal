package utils

import (
	"math/rand"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
)

func GenerateOrganizationName() string {
	rand.Seed(time.Now().UnixNano())
	return petname.Generate(3, "-")
}

func GenerateProjectName() string {
	rand.Seed(time.Now().UnixNano())
	return petname.Generate(2, "-")
}
