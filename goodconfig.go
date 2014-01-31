package goodconfig

import (
	"fmt"
	"os"
	"errors"
	"io/ioutil"
	"strings"
)

type Config map[string]Section
type Section map[string]Record
type Record struct {
	Value Value
}

type Value interface{}

//func (r Record) String() string {
//	return fmt.Sprintf("%s", r)
//}

func Parse(filename string) (*Config, error) {
//	var r Record
//	r.Value = "hello 22"
//	fmt.Println(r.Value)
//
//	return &Config{}, nil
//
	var err error

	// check if file exists
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("File %s doesn't exist", filename))
	}

	var Content []byte

	Content, err = ioutil.ReadFile(filename)

	// split by sections
	var Lines []string = strings.Split(string(Content), "\n")

	config := make(Config)

	var activeSection string

	for i, line := range Lines {

		switch {
		case len(line) == 0:
			// skip, empty line
		case line[:1] == ";":
			// skip, because that's a comment
		case line[:1] == "[" && line[len(line)-1:] == "]":
			// it's a section
			sectionName := line[1:len(line)-1]
			sectionParts := strings.Split(sectionName, ":")

			if(len(sectionParts) > 2) {
				return nil, errors.New(fmt.Sprintf("Only one parent section is allowed, line: %s (%s)", i, line))
			}

			activeSection = sectionParts[0]

			config[activeSection] = make(Section)

			if(len(sectionParts) == 2) {
				// then we should first copy the values of the parent section
				parentSection:= strings.TrimSpace(sectionParts[1])

				if(config[parentSection] == nil) {
					return nil, errors.New(fmt.Sprintf("The parent section was not declared yet, line: %s (%s)", i, line))
				}

				config[activeSection] = config[parentSection]
			}
		default:
			// seems this is a simple config value
			parts := strings.Split(line, "=")

			if(len(parts) < 2){
				return nil, errors.New(fmt.Sprintf("Incorrect format, line: %s (%s)", i, line))
			}

			// defining the key and value
			key := parts[0]
			value := strings.Join(parts[1:], "=")

			keyParts := strings.Split(key, ".")

			if(len(keyParts) > 1) {
				var record *Record

				for i, part := range keyParts {
					switch {
					case i == 0:
						fmt.Println("hello", i)
						record = &Record{Value:make(map[string]Record)}

						config[activeSection][part] = *record
					case (i+1) < len(keyParts):
						fmt.Println("hello", i)
						tmpMap := record.Value.(map[string]Record)

						// if the branch was not yet created
						if _, ok := tmpMap[part]; ok {
							fmt.Println("Exist", part)
							r := tmpMap[part]
							record = &r
						} else {
							fmt.Println("Doesn't exist", part)
							tmp := Record{Value:make(map[string]Record)}
							tmpMap[part]= tmp
							record = &tmp
						}
						fmt.Printf("%v\n", config)
					case (i+1) == len(keyParts):
						fmt.Println("hello", i)
//
//						record.Value = value
						//fmt.Printf("Here I am %v\n", record.Value)
						finalMap := record.Value.(map[string]Record)
						finalMap[part] = Record{Value:value}

					}
				}
			} else {
				config[activeSection][key] = Record{Value:value}
			}
		}

	}

	fmt.Printf("%v", config)

	return &config, nil
}
