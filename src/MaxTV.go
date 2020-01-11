package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func GetScreenStat(t time.Time) {

	fmt.Println("Timer >>>")

	resp, err := http.Get(ScreenStatUrl)
	if err != nil {
		log.Fatalln(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Println(string(body))

	type incomingJson struct {
		total    int
		entities []IncomingBuilding
	}

	var response incomingJson
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return
	}

	var Attempt StatGetAttempt
	Attempt.Hash = GetHash()
	Db.Create(&Attempt)

	for _, record := range response.entities {
		var building Building
		building.maxtv_id = record.id
		building.name = record.name
		building.address = record.address
		building.attempt = Attempt.Hash
		Db.Create(&building)

	}
}
