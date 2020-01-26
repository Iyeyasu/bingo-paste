package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type configParser struct {
	configMap map[interface{}]interface{}
}

func newConfigParser(filename string) *configParser {
	log.Infof("Reading configuration file '%s'", filename)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logError(err.Error())
		return nil
	}

	parser := new(configParser)
	parser.configMap = make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(data), parser.configMap)
	if err != nil {
		logError(err.Error())
		return nil
	}

	return parser
}

func (parser *configParser) hasOption(name string) bool {
	_, err := parser.getOption(name)
	return err == nil
}

func (parser *configParser) getOption(name string) (interface{}, error) {
	path := strings.Split(name, ".")
	m := parser.configMap
	for i := 0; i < len(path); i++ {
		val, ok := m[path[i]]
		if !ok {
			break
		}

		if i == len(path)-1 {
			return val, nil
		}

		childMap, ok := val.(map[interface{}]interface{})
		if !ok {
			break
		}

		m = childMap
	}

	return nil, fmt.Errorf("option '%s' not found", name)
}

func (parser *configParser) getBool(name string, defaultVal bool) bool {
	log.Tracef("Parsing bool field '%s' (default: %t)", name, defaultVal)

	opt, err := parser.getOption(name)
	if err != nil {
		log.Tracef("Using default value for field %s", name)
		return defaultVal
	}

	val, ok := opt.(bool)
	if !ok {
		logOptionError(name, fmt.Sprintf("expected bool, got %T", opt))
		return false
	}

	log.Tracef("Parsed field '%s: %t'", name, val)
	return val
}

func (parser *configParser) getInt(name string, defaultVal int) int {
	log.Tracef("Parsing int field '%s' (default: %d)", name, defaultVal)

	opt, err := parser.getOption(name)
	if err != nil {
		log.Tracef("Using default value for field %s", name)
		return defaultVal
	}

	val, ok := opt.(int)
	if !ok {
		logOptionError(name, fmt.Sprintf("expected int, got %T", opt))
		return 0
	}

	log.Tracef("Parsed field '%s: %d'", name, val)
	return val
}

func (parser *configParser) getStr(name string, defaultVal string) string {
	log.Tracef("Parsing string field '%s' (default: %s)", name, defaultVal)

	opt, err := parser.getOption(name)
	if err != nil {
		log.Tracef("Using default value for field %s", name)
		return defaultVal
	}

	val, ok := opt.(string)
	if !ok {
		logOptionError(name, fmt.Sprintf("expected string, got %T", opt))
		return ""
	}

	log.Tracef("Parsed field '%s: %s'", name, val)
	return val
}

func (parser *configParser) getIntArray(name string, defaultVal []int) []int {
	log.Tracef("Parsing int array field '%s' (default: %v)", name, defaultVal)

	opt, err := parser.getOption(name)
	if err != nil {
		log.Tracef("Using default value for field %s", name)
		return defaultVal
	}

	arr, ok := opt.([]interface{})
	if !ok {
		logOptionError(name, fmt.Sprintf("expected []interface{}, got %T", opt))
		return nil
	}

	val := make([]int, len(arr))
	for i := 0; i < len(arr); i++ {
		val[i], ok = arr[i].(int)
		if !ok {
			logOptionError(fmt.Sprintf("%s[%d]", name, i), fmt.Sprintf("expected int, got %T", arr[i]))
			return nil
		}
	}

	log.Tracef("Parsed field '%s: %v'", name, val)
	return val
}

func (parser *configParser) getStrArray(name string, defaultVal []string) []string {
	log.Tracef("Parsing string array field '%s' (default: %v)", name, defaultVal)

	opt, err := parser.getOption(name)
	if err != nil {
		log.Tracef("Using default value for field %s", name)
		return defaultVal
	}

	arr, ok := opt.([]interface{})
	if !ok {
		logOptionError(name, fmt.Sprintf("expected []interface{}, got %T", opt))
		return nil
	}

	val := make([]string, len(arr))
	for i := 0; i < len(arr); i++ {
		val[i], ok = arr[i].(string)
		if !ok {
			logOptionError(fmt.Sprintf("%s[%d]", name, i), fmt.Sprintf("expected string, got %T", arr[i]))
			return nil
		}
	}

	log.Tracef("Parsed field '%s: %v'", name, val)
	return val
}

func logError(msg string) {
	log.Fatalf("Failed to parse configuration file: %s", msg)
}

func logOptionError(name string, msg string) {
	log.Fatalf("Failed to parse configuration field '%s': %s", name, msg)
}
