package hexid

import (
	"crypto/rand"
	"fmt"
	"regexp"
)

var (
	pattern = `^[0-9a-f]{24}$`
)

func Generate() (string, error) {
	b := make([]byte, 12)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func Validate(id string) bool {
	matched, _ := regexp.MatchString(pattern, id)
	return matched
}
