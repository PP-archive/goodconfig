// package goodconfig provides the basic *.ini files parsing options
// and supports the sections inheritance
package goodconfig

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Config structure describes the overall configuration object structure
type Config struct {
	// it's possible to set the default section using the SetDefaultSection function
	defaultSection string
	// the map of sections
	Sections       map[string]Section
	// the array which stores the errors
	errors         []error
}

// passing the default section key to this function, you will set the defaultSection for the certain config object
// after setting the defaultSection, you are able to call Get method directly for the config object
func (c *Config) SetDefaultSection(key string) error {
	if _, ok := c.Sections[key]; ok {
		c.defaultSection = key
		return nil
	} else {
		return errors.New("Can't set the default section, because the section '%s' doesn't exist")
	}
}

// add error to the errors array
func (c *Config) AddError(e error) {
	c.errors = append(c.errors, e)
}

// get all the available errors and clear the errors array
func (c *Config) GetErrors() []error {
	errors := c.errors
	c.errors = make([]error, 0)
	return errors
}

// get the record which is stored by the key
// if you're using the Get function on Config object, it's expected that the defaultSection is already set
func (c *Config) Get(key string) *Record {
	if len(c.defaultSection) == 0 {
		c.AddError(errors.New("The default section is not defined"))
		return &Record{}
	}

	return c.Section(c.defaultSection).Get(key)
}

type Section struct {
	config *Config
	Value map[string]Record
}

type Record struct {
	config *Config
	Value Value
}

type Value interface{}

func (c *Config) Section(key string) *Section {
	if _, ok := c.Sections[key]; ok {
		s := c.Sections[key]
		return &s
	} else {
		c.AddError(errors.New(fmt.Sprintf("There is no Section for key '%s'", key)))
		return &Section{}
	}
}

func (s *Section) Get(key string) *Record {
	if _, ok := s.Value[key]; ok {
		r := s.Value[key]
		return &r
	} else {
		s.config.AddError(errors.New(fmt.Sprintf("There is no Record for key '%s'", key)))
		return &Record{}
	}
}

func (r *Record) Get(key string) *Record {
	switch r.Value.(type) {
	case map[string]Record:
		if _, ok := r.Value.(map[string]Record)[key]; ok {
			retR := r.Value.(map[string]Record)[key]
			return &retR
		} else {
			r.config.AddError(errors.New(fmt.Sprintf("There is no Record for key '%s'", key)))
			return &Record{}
		}
	default:
		r.config.AddError(errors.New(fmt.Sprintf("The Record value is incorrect", key)))
		return &Record{}
	}
}

func (r Record) ToString() string {
	stringValue := fmt.Sprintf("%s", r.Value)
	return stringValue
}

func (r Record) ToInt() int {
	intValue, _ := strconv.Atoi(r.Value.(string))
	return intValue
}

func (r Record) ToFloat() float64 {
	floatValue, _ := strconv.ParseFloat(r.Value.(string),64)
	return floatValue
}

func (r Record) ToBool() bool {
	boolValue, _ := strconv.ParseBool(r.Value.(string))

	return boolValue
}

func (r Record) ToMap() map[string]Record {
	mapValue := r.Value.(map[string]Record)

	return mapValue
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Parse(filename string) (error) {

	fmt.Println("c address: ", &c)

	var err error

	// check if file exists
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("File %s doesn't exist", filename))
	}

	var Content []byte

	Content, err = ioutil.ReadFile(filename)

	// split by sections
	var Lines []string = strings.Split(string(Content), "\n")

	c.Sections = make(map[string]Section)

	var activeSection string

	for i, line := range Lines {

		line = strings.TrimSpace((line))

		switch {
		case len(line) == 0:
			// skip, empty line
		case line[:1] == ";":
			// skip, because that's a comment
		case line[:1] == "[" && line[len(line)-1:] == "]":
			// it's a section
			sectionName := line[1 : len(line)-1]
			sectionParts := strings.Split(sectionName, ":")

			if len(sectionParts) > 2 {
				return errors.New(fmt.Sprintf("Only one parent section is allowed, line: %s (%s)", i, line))
			}

			activeSection = strings.TrimSpace(sectionParts[0])

			c.Sections[activeSection] = Section{config: c, Value: make(map[string]Record)}

			if len(sectionParts) == 2 {
				// then we should first copy the values of the parent section
				parentSection := strings.TrimSpace(sectionParts[1])

				if len(c.Sections[parentSection].Value) == 0 {
					return errors.New(fmt.Sprintf("The parent section was not declared yet, line: %s (%s)", i, line))
				}

				c.Sections[activeSection] = c.Sections[parentSection]
			}
		default:
			// seems this is a simple config value
			parts := strings.Split(line, "=")

			if len(parts) < 2 {
				return errors.New(fmt.Sprintf("Incorrect format, line: %s (%s)", i, line))
			}

			// defining the key and value
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(strings.Join(parts[1:], "="))

			keyParts := strings.Split(key, ".")

			if len(keyParts) > 1 {
				var record *Record = &Record{}

				for i, part := range keyParts {

					switch {
					case i == 0:
						// inserting to the section
						if _, ok := c.Sections[activeSection].Value[part]; ok {
							*record = c.Sections[activeSection].Value[part]
						} else {
							record = &Record{config: c, Value: make(map[string]Record)}
							c.Sections[activeSection].Value[part] = *record
						}
					case (i + 1) < len(keyParts):
						// if the branch was not yet created
						if _, ok := record.Value.(map[string]Record)[part]; ok {
							*record = record.Value.(map[string]Record)[part]
						} else {
							// if the subsection does not exist
							tmp := Record{config: c, Value: make(map[string]Record)}
							record.Value.(map[string]Record)[part] = tmp
							record = &tmp
						}
					case (i + 1) == len(keyParts):
						// if that was the last key
						record.Value.(map[string]Record)[part] = Record{config: c, Value: value}
					}
				}
			} else {
				// if it's not a composite key
				c.Sections[activeSection].Value[key] = Record{config: c, Value: value}
			}
		}

	}

	return nil
}
