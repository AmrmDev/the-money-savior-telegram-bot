package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateExpenseID() string {
	max := big.NewInt(999999)
	n, _ := rand.Int(rand.Reader, max)
	return fmt.Sprintf("AM%06d", n.Int64())
}
