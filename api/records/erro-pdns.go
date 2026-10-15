package records

import (
	"encoding/json"
	"errors"
	"strconv"
)

type PdnsError struct {
	Status int
	Corpo  string
}

func (e *PdnsError) Error() string {
	return "pdns: " + strconv.Itoa(e.Status) + " " + e.Corpo
}

func RecusaPdns(err error) (string, bool) {
	var e *PdnsError
	if !errors.As(err, &e) || (e.Status != 400 && e.Status != 422) {
		return "", false
	}

	var corpo struct {
		Error string `json:"error"`
	}

	if json.Unmarshal([]byte(e.Corpo), &corpo) == nil && corpo.Error != "" {
		return corpo.Error, true
	}

	return e.Corpo, true
}
