package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Country struct {
	Name       string
	Capital    string
	Population uint
	Area       uint
	Currency   string
	Languages  []string
	Region     string
	Subregion  string
	Flag       string
}

func (c Country) String() string {
	return fmt.Sprintf(
		"%s\n"+
			"\tCapitale: %s\n"+
			"\tPopulation: %d\n"+
			"\tArea: %d\n"+
			"\tCurrency: %s\n"+
			"\tLanguages: %v\n"+
			"\tRegion: %s\n"+
			"\tSubregion: %s\n"+
			"\tFlag: %s\n",
		c.Name, c.Capital, c.Population, c.Area, c.Currency, c.Languages, c.Region, c.Subregion, c.Flag,
	)
}

func main() {
	json_file, err := os.ReadFile("./example.json")
	if err != nil {
		log.Fatal("Can't read file. Exiting the program.", err)
	}

	var Countries []Country
	err = json.Unmarshal(json_file, &Countries)
	if err != nil {
		log.Fatal("Error.", err)
	}

	first_country := Countries[0]
	fmt.Println(first_country)
}
