package main

import (
	"github.com/ironpark/genpig"
)

func main() {
	searchDirs := []string{}
	fileNames := []string{}
	// Default Values
	config := Config{}
	// Json config file load
	genpig.LoadJsonConfig(searchDirs, fileNames, &config)

}
