package users

type Zone struct {
	Zona    string   `bson:"zona"`
	Leitura []string `bson:"leitura"`
	Escrita []string `bson:"escrita"`
}

type User struct {
	Email     string `json:"email"`
	Permissao string `json:"permissao"`
	Zona      string `json:"zona"`
	ID        string `json:"id"`
}
