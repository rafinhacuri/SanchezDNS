package users

import (
	"strconv"
	"strings"
)

func ParseID(id string) (string, int, bool) {
	permissao, position, found := strings.Cut(strings.TrimSpace(id), "-")
	if !found || (permissao != "leitura" && permissao != "escrita") {
		return "", 0, false
	}

	index, err := strconv.Atoi(position)
	if err != nil || index < 0 {
		return "", 0, false
	}

	return permissao, index, true
}
