package users

type Zone struct {
	Zona    string   `bson:"zona"`
	Leitura []string `bson:"leitura"`
	Escrita []string `bson:"escrita"`
}
