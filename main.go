package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/udvarid/don-trade-golang/model"
	collector "github.com/udvarid/don-trade-golang/service"
)

var config = model.Configuration{}

//go:embed resources
var f embed.FS

func main() {
	configFile := flag.String("config", "conf.json", "the Json file contains the configurations")
	environment := flag.String("environment", "local", "where do we run the application, local or internet?")
	flag.Parse()

	configFileInString, _ := f.ReadFile("resources/" + *configFile)
	json.Unmarshal([]byte(configFileInString), &config)

	config.Environment = *environment
	fmt.Println("Collecting data")
	collector.CollectData(&config)
}
