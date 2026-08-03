package users

import "fmt"

func List(zone Zone) []User {
	list := make([]User, 0, len(zone.Leitura)+len(zone.Escrita))

	for index, email := range zone.Leitura {
		list = append(list, User{
			Email:     email,
			Permissao: "leitura",
			Zona:      zone.Zona,
			ID:        fmt.Sprintf("leitura-%d", index),
		})
	}

	for index, email := range zone.Escrita {
		list = append(list, User{
			Email:     email,
			Permissao: "escrita",
			Zona:      zone.Zona,
			ID:        fmt.Sprintf("escrita-%d", index),
		})
	}

	return list
}
