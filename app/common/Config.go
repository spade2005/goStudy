package common

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Config struct {
	data map[string]map[string]string // Maps sections to options to values.
}

var (
	ConfigObj      *Config
	DefaultSection = "default" // Default section name (must be lower-case).

	// Maximum allowed depth when recursively substituing variable names.
	DepthValues = 200

	// Strings accepted as bool.
	BoolStrings = map[string]bool{
		"0":     false,
		"1":     true,
		"f":     false,
		"false": false,
		"n":     false,
		"no":    false,
		"off":   false,
		"on":    true,
		"t":     true,
		"true":  true,
		"y":     true,
		"yes":   true,
	}

	varRegExp = regexp.MustCompile(`%\(([a-zA-Z0-9_.\-]+)\)s`)
)

func init() {
	path, _ := os.Getwd()
	ConfigObj, _ = ReadConfigFile(path + "/app/config/main.ini")
}

func (c *Config) AddSection(section string) bool {

	section = strings.ToLower(section)

	if _, ok := c.data[section]; ok {
		return false
	}
	c.data[section] = make(map[string]string)

	return true
}

func (c *Config) RemoveSection(section string) bool {

	section = strings.ToLower(section)

	switch _, ok := c.data[section]; {
	case !ok:
		return false
	case section == DefaultSection:
		return false // default section cannot be removed
	default:
		for o, _ := range c.data[section] {
			delete(c.data[section], o)
		}
		delete(c.data, section)
	}

	return true
}

func (c *Config) AddOption(section string, option string, value string) bool {

	c.AddSection(section) // make sure section exists

	section = strings.ToLower(section)
	option = strings.ToLower(option)

	_, ok := c.data[section][option]
	c.data[section][option] = value

	return !ok
}

func (c *Config) RemoveOption(section string, option string) bool {

	section = strings.ToLower(section)
	option = strings.ToLower(option)

	if _, ok := c.data[section]; !ok {
		return false
	}

	_, ok := c.data[section][option]
	delete(c.data[section], option)

	return ok
}

func NewConfigFile() *Config {

	c := new(Config)
	c.data = make(map[string]map[string]string)

	c.AddSection(DefaultSection) // default section always exists

	return c
}

func ReadConfigFile(fname string) (*Config, error) {

	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	c := NewConfigFile()
	if err := c.read(bufio.NewReader(file)); err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) read(buf *bufio.Reader) error {

	var section, option string

	for {
		l, err := buf.ReadString('\n') // parse line-by-line
		if err == io.EOF {
			if len(l) == 0 {
				break
			}
		} else if err != nil {
			return err
		}

		l = strings.TrimSpace(l)
		// switch written for readability (not performance)
		switch {
		case len(l) == 0: // empty line
			continue

		case l[0] == '#': // comment
			continue

		case l[0] == ';': // comment
			continue

		case len(l) >= 3 && strings.ToLower(l[0:3]) == "rem":
			// comment (for windows users)
			continue

		case l[0] == '[' && l[len(l)-1] == ']': // new section
			option = "" // reset multi-line value
			section = strings.TrimSpace(l[1 : len(l)-1])
			c.AddSection(section)

		case section == "": // not new section and no section defined so far
			return errors.New("Section not found: must start with section")

		default: // other alternatives
			i := firstIndex(l, []byte{'=', ':'})
			switch {
			case i > 0: // option and value
				i := firstIndex(l, []byte{'=', ':'})
				option = strings.TrimSpace(l[0:i])
				value := strings.TrimSpace(stripComments(l[i+1:]))
				c.AddOption(section, option, value)

			case section != "" && option != "":
				// continuation of multi-line value
				prev, _ := c.GetRawString(section, option)
				value := strings.TrimSpace(stripComments(l))
				c.AddOption(section, option, prev+"\n"+value)

			default:
				return errors.New(fmt.Sprintf("Could not parse line: %s", l))
			}
		}
	}

	return nil
}

func (c *Config) GetRawString(section string, option string) (string, error) {

	section = strings.ToLower(section)
	option = strings.ToLower(option)

	if _, ok := c.data[section]; ok {

		if value, ok := c.data[section][option]; ok {
			return value, nil
		}

		return "", errors.New(fmt.Sprintf("Option not found: %s", option))
	}

	return "", errors.New(fmt.Sprintf("Section not found: %s", section))
}

func stripComments(l string) string {

	// comments are preceded by space or TAB
	for _, c := range []string{" ;", "\t;", " #", "\t#"} {
		if i := strings.Index(l, c); i != -1 {
			l = l[0:i]
		}
	}

	return l
}

func firstIndex(s string, delim []byte) int {

	for i := 0; i < len(s); i++ {
		for j := 0; j < len(delim); j++ {
			if s[i] == delim[j] {
				return i
			}
		}
	}

	return -1
}

func (c *Config) GetString(section string, option string) (string, error) {

	value, err := c.GetRawString(section, option)
	if err != nil {
		return "", err
	}

	section = strings.ToLower(section)

	var i int

	for i = 0; i < DepthValues; i++ { // keep a sane depth

		vr := varRegExp.FindStringSubmatchIndex(value)
		if len(vr) == 0 {
			break
		}

		noption := value[vr[2]:vr[3]]
		noption = strings.ToLower(noption)

		// search variable in default section
		nvalue, _ := c.data[DefaultSection][noption]
		if _, ok := c.data[section][noption]; ok {
			nvalue = c.data[section][noption]
		}

		if nvalue == "" {
			return "", errors.New(fmt.Sprintf("Option not found: %s", noption))
		}

		// substitute by new value and take off leading '%(' and trailing ')s'
		value = value[0:vr[2]-2] + nvalue + value[vr[3]+2:]
	}

	if i == DepthValues {
		return "",
			errors.New(
				fmt.Sprintf(
					"Possible cycle while unfolding variables: max depth of %d reached",
					strconv.Itoa(DepthValues),
				),
			)
	}

	return value, nil
}
