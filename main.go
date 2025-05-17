package main

import (
	"log"

	"github.com/thomas-osgood/rawdog/server"
)

func main() {
	var err error
	var ts *server.TeamServer

	ts, err = server.NewTeamServer()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Fatal(ts.Start())
}
