package main

import (
	"log"
	"rawdog/server"
)

func main() {
	var err error
	var ts *server.TeamServer

	ts, err = server.NewTeamServer()
	if err != nil {
		log.Fatal(err.Error())
	}

	ts.Start()
}
