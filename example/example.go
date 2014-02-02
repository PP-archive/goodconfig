package main

import (
	"fmt"
	"github.com/PavelPolyakov/goodconfig"
	"os"
	"reflect"
)

func main() {
	Config := goodconfig.NewConfig()
	err := Config.Parse("./config/application.ini")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Config.SetDefaultSection("development")

	intValue := Config.Get("intValue").ToInt()
	fmt.Println("intValue:", intValue, reflect.TypeOf(intValue));

	boolValue := Config.Get("boolValue").ToBool()
	fmt.Println("boolValue:", boolValue, reflect.TypeOf(boolValue));

	stringValue := Config.Get("stringValue").ToString()
	fmt.Println("stringValue:", stringValue, reflect.TypeOf(stringValue));

	// incorrect call of the child key
	incorrectChild := Config.Section("production").Get("parent").Get("chiled").ToString() // typo

	if(Config.HasErrors()) {
		fmt.Println(Config.GetErrors())
		fmt.Println("Child: ", incorrectChild)
	}

	// correct call of the child key
	correctChild := Config.Section("production").Get("parent").Get("child").Get("subchild").ToString()
	fmt.Println("Subchild:", correctChild)

	fmt.Println("Config: ", Config)

}
