package statistics

import (
	"encoding/json"
	"strconv"
)

func valor(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}

	var texto string

	err := json.Unmarshal(raw, &texto)
	if err == nil {
		return texto
	}

	var booleano bool

	err = json.Unmarshal(raw, &booleano)
	if err == nil {
		return strconv.FormatBool(booleano)
	}

	var inteiro int64

	err = json.Unmarshal(raw, &inteiro)
	if err == nil {
		return strconv.FormatInt(inteiro, 10)
	}

	var numero float64

	err = json.Unmarshal(raw, &numero)
	if err == nil {
		return strconv.FormatFloat(numero, 'f', -1, 64)
	}

	return string(raw)
}
