package main

import (
	"fmt"
	"github.com/PavelPolyakov/goodconfig"
)

func main() {
	Config, err := goodconfig.Parse("./config/application.ini")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Config: ", Config)
}
