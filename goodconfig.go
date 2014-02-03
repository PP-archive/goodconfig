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
	defaultSection string
	// the map of sections
	Sections       map[string]Section
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

// check if there are yet unseen errors
func (c *Config) HasErrors() bool {
	return len(c.errors) > 0
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

// struct which describes the section inside the *.ini file
type Section struct {
	config *Config
	Value map[string]Record
}

// struct which describes the record,
// Record Value could be a map of other records
type Record struct {
	config *Config
	Value Value
}

// struct which describes the Value,
// Value is the final entity in the config hierarchy
type Value interface{}

// returns the Section by key
func (c *Config) Section(key string) *Section {
	if _, ok := c.Sections[key]; ok {
		s := c.Sections[key]
		return &s
	} else {
		c.AddError(errors.New(fmt.Sprintf("There is no Section for key '%s'", key)))
		return &Section{}
	}
}

// returns the Record of the Section by key
func (s *Section) Get(key string) *Record {
	if _, ok := s.Value[key]; ok {
		r := s.Value[key]
		return &r
	} else {
		s.config.AddError(errors.New(fmt.Sprintf("There is no Record for key '%s'", key)))
		return &Record{}
	}
}

// fill the section with the other section, so we receive the copy of the maps, not the pointers
func (s *Section) Fill(sParent Section) {
	for k, v := range sParent.Value {
		// defining the type of the child Value element
		switch v.Value.(type) {
		case map[string]Record:
			s.Value[k] = Record{config: sParent.config, Value:make(map[string]Record)}

			nextR := s.Value[k]
			nextR.Fill(v)
		default:
			s.Value[k] = v
		}

	}
}


// fill the record with the other record, so we receive the copy of the maps, not the pointers
func (r *Record) Fill(rParent Record) {
	switch rParent.Value.(type) {
	case map[string]Record:
	for k, v := range rParent.Value.(map[string]Record) {
		// defining the type of the child Value element
		switch v.Value.(type) {
		case map[string]Record:
			r.Value.(map[string]Record)[k] = Record{config: rParent.config, Value:make(map[string]Record)}

			nextR := r.Value.(map[string]Record)[k]
			nextR.Fill(v)
		default:
			r.Value.(map[string]Record)[k] = v
		}
	}
	default:
	r.Value = rParent.Value
	}
}


// returns the Record by key (which is located inside the other Record)
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

// returns the Value casted to string
func (r Record) ToString() string {
	stringValue := fmt.Sprintf("%s", r.Value)
	return stringValue
}

// returns the Value casted to int
func (r Record) ToInt() int {
	intValue, _ := strconv.Atoi(r.Value.(string))
	return intValue
}

// returns the Value casted to float64
func (r Record) ToFloat() float64 {
	floatValue, _ := strconv.ParseFloat(r.Value.(string),64)
	return floatValue
}

// returns the Value casted to bool
func (r Record) ToBool() bool {
	boolValue, _ := strconv.ParseBool(r.Value.(string))

	return boolValue
}

// returns the Value casted to map[string]Record
func (r Record) ToMap() map[string]Record {
	mapValue := r.Value.(map[string]Record)

	return mapValue
}

// creates the new Config object
func NewConfig() *Config {
	return &Config{}
}

// parses the ini formatted file (filename)
// builds the hierarchy of Sections and Records
// after the Parse, you are able to use the Get methods
func (c *Config) Parse(filename string) (error) {
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

				// filling the child section
				activeSectionVariable := c.Sections[activeSection]
				parentSectionVariable := c.Sections[parentSection]


				// actually the filling process
				activeSectionVariable.Fill(parentSectionVariable)

//				fmt.Println("active: ", activeSectionVariable)
//				fmt.Println("---")
//				fmt.Println("parent: ", parentSectionVariable)
//				os.Exit(1)
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
						// whether the branch was not yet created
						if _, ok := record.Value.(map[string]Record)[part]; ok {
							*record = record.Value.(map[string]Record)[part]
						} else {
							// if the subsection does not exist
							var newRecord *Record = &Record{config: c, Value: make(map[string]Record)}
							record.Value.(map[string]Record)[part] = *newRecord
							record = newRecord
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
