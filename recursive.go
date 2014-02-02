package main

import "fmt"

func recursiveCopy(p map[string]interface{}, key string,  m interface{}) {
	switch m.(type) {
	case map[string]interface{}:
		p[key] = make(map[string]interface{})
		//newP := p[key]
		//p = p[key]
		for k,v := range m.(map[string]interface{}) {
			recursiveCopy(p[key].(map[string]interface{}), k, v)
		}
	default:
		p[key] = m
	}
}

func main() {
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
