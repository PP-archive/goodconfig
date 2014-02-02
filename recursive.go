package main

import (
	"fmt"
	//"reflect"
	//"os"
)

func recursiveCopy(p map[string]interface{}, key string,  m interface{}) {
	switch m.(type) {
	case map[string]interface{}:
		p[key] = make(map[string]interface{})
		for k,v := range m.(map[string]interface{}) {
			recursiveCopy(p[key].(map[string]interface{}), k, v)
		}
	default:
		p[key] = m
	}
}

type Record struct {
	Value Value
}

type Value interface{}

func (r *Record) Fill(key string, rParent Record) {
	switch rParent.Value.(type) {
	case map[string]Record:
		fmt.Println("hello")
		//os.Exit(1)
		r.Value.(map[string]Record)[key] = Record{Value:make(map[string]Record)}
	for k,v := range rParent.Value.(map[string]Record) {
		nextR := r.Value.(map[string]Record)[key]
		nextR.Fill(k, v)
	}
	default:
		r.Value.(map[string]Record)[key] = Record{Value:rParent.Value}
	}
}

func main() {

	r := Record{Value:make(map[string]Record)}

	r.Value.(map[string]Record)["first"] = Record{Value:make(map[string]Record)}
	r.Value.(map[string]Record)["first"].Value.(map[string]Record)["second"] = Record{Value:make(map[string]interface{})}
	r.Value.(map[string]Record)["first"].Value.(map[string]Record)["second"].Value.(map[string]interface{})["third"] = "hello"

	fmt.Println(r)

	rChild := Record{Value:make(map[string]Record)}

	rChild.Fill("first", r.Value.(map[string]Record)["first"])

	fmt.Println(rChild)
//	fmt.Println(rChild.Value.(map[string]Record)["first"].Value.(map[string]Record)["second"])
	return



	hello := make(map[string]interface{})
	hello["first"] = make(map[string]interface{})
	hello["first"].(map[string]interface{})["second"] = make(map[string]interface{})
	hello["first"].(map[string]interface{})["second"].(map[string]interface{})["last"] = 123

	helloAtom := make(map[string]interface{})
	helloAtom["first"] = 333

	var helloPointer map[string]interface{} = make(map[string]interface{})
	var helloAtomPointer map[string]interface{} = make(map[string]interface{})
	//helloPointer = make(map[string]interface{})

	recursiveCopy(helloPointer, "first", hello["first"])
	recursiveCopy(helloAtomPointer, "first", helloAtom["first"])

	//helloPointer := hello

	helloPointer["first"].(map[string]interface{})["second"].(map[string]interface{})["last"] = 124

	fmt.Println(hello)
	fmt.Println("---")
	fmt.Println(helloPointer)
	fmt.Println("---")

	hello["first"].(map[string]interface{})["second"].(map[string]interface{})["last"] = 125

	fmt.Println(hello)
	fmt.Println("---")
	fmt.Println(helloPointer)

	fmt.Println("---")
	fmt.Println(helloAtomPointer)

}
