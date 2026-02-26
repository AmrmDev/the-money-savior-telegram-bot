package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateDisplayID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("AM%05d", rand.Intn(90000)+10000)
}