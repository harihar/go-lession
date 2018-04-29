package main

import (
	"testing"
	"fmt"
	"encoding/json"
	"reflect"
)

func TestConvertFileToJSON(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     string
	}{
		{
			"Should combine keys with same parent key",
			"TestXls1.xlsx",
			`{"Countries":{"AL":"Albania","DE":"Germany"},"Spin":{"description":"Your PIN is required","label":"Your PIN"}}`,
		}, {
			"Should create 3 level deep nested JSON",
			"TestXls2.xlsx",
			`{"Countries":{"AL":"Albania","DE":"Germany"},"ProfileForm":{"common":{"conditionText":"You may view and change your personal info in VW portal ","imprintLabel":"Imprint"}},"Spin":{"description":"Your PIN is required","label":"Your PIN"}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertFileToJSON(tt.filePath)
			areEqualJSON, err := AreEqualJSON(got, tt.want)
			if err != nil {
				t.Errorf("Error occured while comparing JSON strings. Error: %v", err)
			}
			if !areEqualJSON {
				t.Errorf("ConvertFileToJSON()\nActual = %v,\nExpected = %v", got, tt.want)
			}
		})
	}
}

func AreEqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}
