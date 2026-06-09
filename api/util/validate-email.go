package util

import "regexp"

func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}

	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)

	return re.MatchString(email)
}
