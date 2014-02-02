package main

import (
	"fmt"
	"github.com/PavelPolyakov/goodconfig"
	"reflect"
)

func main() {
	Config := goodconfig.NewConfig()
	err := Config.Parse("./config/application.ini")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Key: ")
	fmt.Println("full query:")
	fmt.Println(Config.Section("production").Get("parent").Get("child").Get("subchild").ToInt())

	Config.SetDefaultSection("production")
	fmt.Println("production:")
	fmt.Println(Config.Get("parent").Get("child").Get("subchild").ToInt())


	fmt.Println(Config.Get("float").ToFloat(), reflect.TypeOf(Config.Get("float").ToFloat()))

	for key, _ := range Config.Get("parent").ToMap() {
		fmt.Println("Key :" , key)
	}

	fmt.Println("bool value:")
	v := Config.Get("parent").Get("child").Get("helloBool").ToBool()
	fmt.Println(v, reflect.TypeOf(v))

	tmp := Config.Get("parent").Get("chi44ld")

	p := &Config
	fmt.Println("Config address: ", p)


	errs := Config.GetErrors()
	fmt.Println("Tmp: ", tmp)
	fmt.Println("Errors: ", errs)
	fmt.Println("Errors: ", len(errs))
	fmt.Println("Errors: ", len(Config.GetErrors()))
}
