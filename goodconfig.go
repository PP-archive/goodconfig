package goodconfig

import (
	"fmt"
	"os"
	"errors"
	"io/ioutil"
	//"regexp"
	"strings"
	//"reflect"
	"reflect"
)

type Config map[string]Section
type Section map[string]Record
type Record interface{}

func Parse(filename string) (*Config, error) {
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

			fmt.Println("section parts: ", sectionParts)

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
				var record map[string]Record

				record = make(map[string]Record)
				config[activeSection] = record

				for i, part := range keyParts {
					if((i+1) < len(keyParts)) {
						tmp := make(map[string]Record)

						record[part] = tmp
						record = tmp

					} else {
						record[part] = value
					}
				}

				fmt.Println(keyParts[0], config[activeSection][keyParts[0]])
			} else {
				config[activeSection][key] = value
			}
		}

	}


	fmt.Printf("ggg %v", config)

	return &Config{}, nil
}
