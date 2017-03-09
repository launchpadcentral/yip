package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	ErrKeyNotFound   = errors.New("key not found")
	ErrInvalidObject = errors.New("invalid object")
	ErrInvalidIndex  = errors.New("invalid index")
	ErrInvalidKey    = errors.New("invalid key")

	keyPairs   string
	outputFile string
	inputFile  string
)

func parseFlags() {
	flag.StringVar(&keyPairs, "key-pairs", "", "key-pairs input needs to be in key=pair format. i.e. api.tag=123,srv.image=234")
	flag.StringVar(&inputFile, "f", "", "input file name")
	flag.StringVar(&outputFile, "o", "", "output file name")
	flag.Parse()

	if keyPairs == "" {
		log.Fatalf("-key-pairs flag is not set")
	}

	if inputFile == "" {
		log.Fatalf("-f is not set")
	}
}

func main() {

	updates := validateInput()
	t := mustParseYaml()

	var err error
	for i := 0; i < len(updates); i += 2 {
		var newValue interface{}
		newValue = updates[i+1]
		if intVal, err := strconv.Atoi(updates[i+1]); err == nil {
			newValue = intVal
		}
		t, err = updateYaml(t, updates[i], newValue)
		if err != nil {
			log.Fatalf("Error with key %s: %s", updates[i], err)
			return
		}
	}

	output(t)
}

func validateInput() []string {
	parseFlags()
	updates := parseInput()

	if len(updates)%2 != 0 {
		log.Fatalf("-key-pairs input needs to be in key=pair format. i.e. api.tag=123,srv.image=234")
	}

	return updates
}

func mustParseYaml() interface{} {
	data := mustReadYaml()

	var t interface{}
	if err := yaml.Unmarshal([]byte(data), &t); err != nil {
		log.Fatalf("Could not unmarshal yaml file: %s", err)
	}

	return t
}

func output(t interface{}) {
	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("Could not marshal yaml file: %s", err)
	}

	if outputFile == "" {
		fmt.Printf("%+v\n", string(d))
		return
	}

	writeYaml(d)
}

func mustReadYaml() []byte {
	dat, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Could not read input file: %s", err)
	}

	return dat
}

func writeYaml(output []byte) {
	if err := ioutil.WriteFile(outputFile, output, 0644); err != nil {
		log.Fatalf("Could not write the output to the file: %s", err)
	}
}

func parseInput() []string {
	pairs := strings.Split(keyPairs, ",")
	input := []string{}
	for i := range pairs {
		input = append(input, strings.Split(strings.Trim(pairs[i], " "), "=")...)
	}

	return input
}

func updateYaml(input interface{}, key string, value interface{}) (interface{}, error) {
	keys := strings.Split(key, ".")

	currentKey := keys[0]

	nextKey := strings.Join(keys[1:], ".")
	switch input.(type) {
	case map[interface{}]interface{}:
		subMap := input.(map[interface{}]interface{})

		val, err := updateYaml(subMap[currentKey], nextKey, value)
		if err != nil {
			return nil, err
		}
		subMap[currentKey] = val
		return subMap, nil
	case []interface{}:
		subArray := input.([]interface{})
		index, err := strconv.Atoi(currentKey)
		if err != nil {
			return nil, ErrInvalidIndex
		}

		if index >= len(subArray) {
			return nil, ErrKeyNotFound
		}

		val, err := updateYaml(subArray[index], nextKey, value)
		if err != nil {
			return nil, err
		}
		subArray[index] = val

		return subArray, nil
	case string, int, bool:
		if len(keys) == 0 || keys[0] == "" {
			return value, nil
		}
		return nil, ErrInvalidKey
	}

	return nil, ErrKeyNotFound
}
