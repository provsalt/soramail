package random

import (
	"math/rand"
	"time"
)

// AddressProvider is an interface that provides email addresses
type AddressProvider interface {
	// RandomEmail is the function that provides randomization of email addresses.
	// The input parameter of this function specifies the domain name while the output should be an email address
	// that has the same domain as the input
	RandomEmail(string) string
}

// DefaultRandomizer is the default randomizer used by soramail.
type DefaultRandomizer struct {
	// Length is the length of the expected length of the email address before the @ part of the
	// email address.
	Length uint
}

func (dr DefaultRandomizer) RandomEmail(domain string) string {
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.New(rand.NewSource(time.Now().UnixNano()))

	localPart := make([]byte, dr.Length)
	for i := range localPart {
		localPart[i] = charset[rand.Intn(len(charset))]
	}

	return string(localPart) + "@" + domain
}
