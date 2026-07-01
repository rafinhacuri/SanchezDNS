package main

import (
	"os"
	"time"

	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/controller"
	"github.com/rafinhacuri/SanchezDNS/api/controller/fetch"
	"github.com/rafinhacuri/SanchezDNS/api/controller/insert"
	"github.com/rafinhacuri/SanchezDNS/api/logs"
	"github.com/rafinhacuri/SanchezDNS/api/mongo"
	"github.com/rafinhacuri/SanchezDNS/api/users"
	"github.com/rafinhacuri/SanchezDNS/api/zonas"
)

type GoRes struct {
	Message string `json:"message"`
}

func main() {
	converter := typescriptify.New()
	converter.DontExport = true

	converter.ManageType(time.Time{}, typescriptify.TypeOptions{
		TSType: "string",
	})
	converter.ManageType(bson.ObjectID{}, typescriptify.TypeOptions{
		TSType: "string",
	})

	converter.Add(GoRes{})
	converter.Add(insert.CreateZoneRequest{})
	converter.Add(logs.LogsResponse{})
	converter.Add(fetch.PdnsZone{})
	converter.Add(controller.SessionRes{})
	converter.Add(fetch.StatisticsResponse{})
	converter.Add(zonas.ZonesResponse{})
	converter.Add(zonas.ZoneFetch{})
	converter.Add(users.User{})
	converter.Add(mongo.Solicitacao{})
	converter.Add(fetch.SolicitacaoResponse{})
	converter.Add(fetch.CadastroResponse{})

	converter.BackupDir = ""
	converter.CreateInterface = true

  err := os.MkdirAll("../../app/types/", 0o750)
	if err != nil {
		panic("Failed to create app/types dir: " + err.Error())
	}

	err = converter.ConvertToFile("../../app/types/go.d.ts")
	if err != nil {
		panic("Failed to generate types: " + err.Error())
	}
}
