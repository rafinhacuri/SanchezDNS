package main

import (
	"os"

	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"

	"github.com/rafinhacuri/SanchezDNS/api/controller"
	"github.com/rafinhacuri/SanchezDNS/api/controller/fetch"
	"github.com/rafinhacuri/SanchezDNS/api/controller/insert"
	"github.com/rafinhacuri/SanchezDNS/api/logs"
	"github.com/rafinhacuri/SanchezDNS/api/users"
	"github.com/rafinhacuri/SanchezDNS/api/zonas"
)

type GoRes struct {
	Message string `json:"message"`
}

func main() {
	converter := typescriptify.New()

	converter.Add(GoRes{})
	converter.Add(insert.CreateZoneRequest{})
	converter.Add(logs.LogsResponse{})
	converter.Add(fetch.PdnsZone{})
	converter.Add(controller.SessionRes{})
	converter.Add(fetch.StatisticsResponse{})
	converter.Add(zonas.ZonesResponse{})
	converter.Add(zonas.ZoneFetch{})
	converter.Add(users.User{})

	converter.BackupDir = ""
	converter.CreateInterface = true

	err := os.MkdirAll("../../shared/types/", 0o750)
	if err != nil {
		panic("Failed to create shared/types dir: " + err.Error())
	}

	err = converter.ConvertToFile("../../shared/types/goServer.d.ts")
	if err != nil {
		panic("Failed to generate types: " + err.Error())
	}
}
