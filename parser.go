package parser

import (
	"encoding/json"
	"errors"
	"log"
	"regexp"
	"strings"
)

var (
	ERR_EMPTY_RAW    = errors.New("Empty raw")
	ERR_WRONG_FORMAT = errors.New("Wrong raw format")
)

const (
	_RAW_MATCHED_STRING = iota
	_RAW_KEY
	_RAW_VALUE
)

type (
	Raw struct {
		Fields     map[string]string `json:"fields,omitempty"`
		Expression *Expression       `json:"expression,omitempty"`
	}

	Expression struct {
		LogicOperation map[string][]FieldGroup `json:"logicOperation,omitempty"`
	}

	FieldGroup struct {
		Key map[string][]Field `json:"key"`
	}

	Field struct {
		Value     string `json:"value"`
		Operation string `json:"operation"`
	}
)

const (
	_MATCHED_STRING = iota
	_MS_LOGIC_OPERATION
	_MS_KEY
	_MS_OPERATION
	_MS_VALUE
)

func ParseRawToStruct(raw string) (*Raw, error) {
	if raw == "" {
		return nil, ERR_EMPTY_RAW
	}

	structuredRaw := &Raw{
		Fields: make(map[string]string, 0),
	}

	raw = strings.ReplaceAll(raw, " ", "")
	mss := regexp.MustCompile(`\??([\w\.\%]+)(?:=){1}((?:(?:(?:[^0-9A-Za-z]){2})?[\w\.\%]+(?:(?:[^0-9A-Za-z]){2}[\w\.\%]+)?)*)`).FindAllStringSubmatch(raw, -1)

	for _, ms := range mss {
		if ms[_RAW_KEY] == "expression" {
			structuredRaw.Expression = parseExp(ms[_RAW_VALUE])

			continue
		}

		structuredRaw.Fields[ms[_RAW_KEY]] = ms[_RAW_VALUE]
	}

	return structuredRaw, nil
}

func ParseRawToJSON(raw string) (string, error) {
	jsonData, err := ParseRawToStruct(raw)
	if err != nil {
		return "", err
	}

	return toJSON(jsonData), nil
}

func parseExp(expRaw string) *Expression {
	exp := Expression{map[string][]FieldGroup{}}
	matchedStrings := regexp.MustCompile(`((?:[^0-9A-Za-z]){2})?([\w\.\%]+)(!=|==)([\w\.\%]+)`).FindAllStringSubmatch(expRaw, -1)
	if len(matchedStrings) > 1 {
		matchedStrings[0][_MS_LOGIC_OPERATION] = matchedStrings[1][_MS_LOGIC_OPERATION]
	} else {
		matchedStrings[0][_MS_LOGIC_OPERATION] = "Single"
	}

	for ims := range matchedStrings {
		lo, key := map[string][]FieldGroup{}, map[string][]Field{}

		field := &Field{
			Value:     matchedStrings[ims][_MS_VALUE],
			Operation: matchedStrings[ims][_MS_OPERATION],
		}

		key[matchedStrings[ims][_MS_KEY]] = append(key[matchedStrings[ims][_MS_KEY]], *field)

		fieldGroup := &FieldGroup{
			Key: key,
		}

		lo[matchedStrings[ims][_MS_LOGIC_OPERATION]] = append(lo[matchedStrings[ims][_MS_LOGIC_OPERATION]], *fieldGroup)

		exp.LogicOperation[matchedStrings[ims][_MS_LOGIC_OPERATION]] = append(exp.LogicOperation[matchedStrings[ims][_MS_LOGIC_OPERATION]], lo[matchedStrings[ims][_MS_LOGIC_OPERATION]]...)

	}

	return &exp
}

func toJSON(m any) string {
	js, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return strings.ReplaceAll(string(js), ",", ", ")
}
