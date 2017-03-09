package main

import (
	"reflect"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

var singleKeyInput = `
key: value
`

var multipleKeyInput = `
audit-api:
  api:
    imageTag: latest
`
var listInput = `
audit-api:
  - api:
      imageTag: latest
  - api:
      imageTag: something
`
var primitiveListInput = `audit-api:
- api1
- api2
`

var booleanValueInput = `
key: true
`

func readYaml(t *testing.T, input string) (output map[interface{}]interface{}) {
	err := yaml.Unmarshal([]byte(input), &output)
	if err != nil {
		t.Fatalf("expected nil, got %s", err)
	}

	return
}

func testWriteYaml(t *testing.T, input interface{}) string {
	res, err := yaml.Marshal(input)
	if err != nil {
		t.Fatalf("expected nil, got %s", err)
	}

	return string(res)
}

func TestUpdateYamlSingleKeyNotFound(t *testing.T) {
	yaml := readYaml(t, singleKeyInput)

	_, err := updateYaml(yaml, "test", "123")
	if err != ErrKeyNotFound {
		t.Errorf("expected %s, got %s", ErrKeyNotFound, err)
	}
}

func TestUpdateYamlSingleKeySuccess(t *testing.T) {
	yaml := readYaml(t, singleKeyInput)

	output, err := updateYaml(yaml, "key", 123)
	if err != nil {
		t.Fatalf("expected nil, got %s", err)
	}

	res := testWriteYaml(t, output)
	expectedYaml := `key: 123
`
	if res != expectedYaml {
		t.Fatalf("expected:\n%#v\n,got:\n%#v", expectedYaml, res)
	}

}

func TestUpdateYamlMultipleKeyNotFound(t *testing.T) {
	yaml := readYaml(t, multipleKeyInput)

	_, err := updateYaml(yaml, "card-api.api.latest", "123")
	if err != ErrKeyNotFound {
		t.Errorf("expected %s, got %s", ErrKeyNotFound, err)
	}
}

func TestUpdateYamlMultipleKeyInvalidObject(t *testing.T) {
	var input = `
audit-api:
  api: latest
`
	yaml := readYaml(t, input)

	_, err := updateYaml(yaml, "audit-api.api.imageTag", 123)
	if err != ErrInvalidKey {
		t.Errorf("expected %s, got %v", ErrInvalidObject, err)
	}

}

func TestUpdateYamlMultipleKey(t *testing.T) {
	yaml := readYaml(t, multipleKeyInput)

	output, err := updateYaml(yaml, "audit-api.api.imageTag", 123)
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}

	res := testWriteYaml(t, output)
	expectedYaml := `audit-api:
  api:
    imageTag: 123
`
	if res != expectedYaml {
		t.Fatalf("expected:\n%#v\n,got:\n%#v", expectedYaml, res)
	}
}

func TestUpdateYamlList(t *testing.T) {
	yaml := readYaml(t, listInput)

	output, err := updateYaml(yaml, "audit-api.1.api.imageTag", 123)
	if err != nil {
		t.Fatalf("expected nil, got %s", err)
	}

	res := testWriteYaml(t, output)
	expectedYaml := `audit-api:
- api:
    imageTag: latest
- api:
    imageTag: 123
`
	if res != expectedYaml {
		t.Fatalf("expected:\n%#v\n,got:\n%#v", expectedYaml, res)
	}
}

func TestUpdateYamlListWithPrimitiveElements(t *testing.T) {
	yaml := readYaml(t, primitiveListInput)

	output, err := updateYaml(yaml, "audit-api.1", 123)
	if err != nil {
		t.Fatalf("expected nil, got %s", err)
	}

	res := testWriteYaml(t, output)
	expectedYaml := `audit-api:
- api1
- 123
`
	if res != expectedYaml {
		t.Fatalf("expected:\n%#v\n,got:\n%#v", expectedYaml, res)
	}
}

func TestUpdateYamlListWithOOIndexPrimitiveElements(t *testing.T) {
	yaml := readYaml(t, primitiveListInput)

	_, err := updateYaml(yaml, "audit-api.3", 123)
	if err != ErrKeyNotFound {
		t.Fatalf("expected %s, got %s", ErrKeyNotFound, err)
	}

}

func TestUpdateYamlListOutOfIndex(t *testing.T) {
	yaml := readYaml(t, primitiveListInput)

	_, err := updateYaml(yaml, "audit-api.3", 123)
	if err != ErrKeyNotFound {
		t.Fatalf("expected %s, got %s", ErrKeyNotFound, err)
	}
}

func TestUpdateYamlListInvalidIndex(t *testing.T) {
	yaml := readYaml(t, primitiveListInput)

	_, err := updateYaml(yaml, "audit-api.key", 123)
	if err != ErrInvalidIndex {
		t.Fatalf("expected %s, got %s", ErrInvalidIndex, err)
	}

}

func TestUpdateYamlListPrimitiveElementInvalidKey(t *testing.T) {
	yaml := readYaml(t, primitiveListInput)

	_, err := updateYaml(yaml, "audit-api.1.test", 123)
	if err != ErrInvalidKey {
		t.Fatalf("expected %s, got %s", ErrInvalidKey, err)
	}

}

func TestUpdateYamlBooleanValueInput(t *testing.T) {
	yaml := readYaml(t, booleanValueInput)

	output, err := updateYaml(yaml, "key", false)
	if err != nil {
		t.Fatalf("expected nil, got %s", err)
	}

	res := testWriteYaml(t, output)
	expectedYaml := `key: false
`
	if res != expectedYaml {
		t.Fatalf("expected:\n%#v\n,got:\n%#v", expectedYaml, res)
	}
}

func TestParseInput(t *testing.T) {
	keyPairs = "kubernetes.password=123,user-api.api.imageTag=1q2w3e"

	output := parseInput()
	expected := []string{"kubernetes.password", "123", "user-api.api.imageTag", "1q2w3e"}
	if !reflect.DeepEqual(expected, output) {
		t.Fatalf("expected %+v, got %+v", expected, output)
	}
}

func TestParseInputWithSpaces(t *testing.T) {
	keyPairs = "    kubernetes.password=123,    user-api.api.imageTag=1q2w3e   "

	output := parseInput()
	expected := []string{"kubernetes.password", "123", "user-api.api.imageTag", "1q2w3e"}
	if !reflect.DeepEqual(expected, output) {
		t.Fatalf("expected %+v, got %+v", expected, output)
	}
}
