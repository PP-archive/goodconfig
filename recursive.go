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

func (r *Record) Fill(rParent Record) {
	switch rParent.Value.(type) {
	case map[string]Record:
		for k, v := range rParent.Value.(map[string]Record) {
				// defining the type of the child Value element
				switch v.Value.(type) {
					case map[string]Record:
					r.Value.(map[string]Record)[k] = Record{Value:make(map[string]Record)}
				case map[string]interface{}:
					r.Value.(map[string]Record)[k] = Record{Value:make(map[string]interface{})}
					default:
				}

				nextR := r.Value.(map[string]Record)[k]
				nextR.Fill(v)
		}
	case map[string]interface{}:
		for k, v := range rParent.Value.(map[string]interface{}) {
			r.Value.(map[string]interface{})[k] = v
		}
	default:
	}
}

func main() {

	r := Record{Value:make(map[string]Record)}

	r.Value.(map[string]Record)["first"] = Record{Value:make(map[string]Record)}
	r.Value.(map[string]Record)["first"].Value.(map[string]Record)["second"] = Record{Value:make(map[string]interface{})}
	r.Value.(map[string]Record)["first"].Value.(map[string]Record)["second"].Value.(map[string]interface{})["third"] = "hello"

	fmt.Println(r)

	rChild := Record{Value:make(map[string]Record)}

	fmt.Println("before: ", rChild)
	rChild.Fill(r)

	fmt.Println("after: ", rChild)
	rChild.Value.(map[string]Record)["first"].Value.(map[string]Record)["second"].Value.(map[string]interface{})["third"] = 333
	fmt.Println(rChild.Value.(map[string]Record)["first"].Value.(map[string]Record)["second"].Value.(map[string]interface{})["third"])
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
