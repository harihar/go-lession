package main

import (
	"github.com/tealeg/xlsx"
	"log"
	"strings"
	"encoding/json"
	"os"
)

func main() {
	jsonString := ConvertFileToJSON("/translator/UC08 - Profile Registration - Translation - Master.xlsx")
	println(jsonString)
	file1, err := os.Create("translations.json")
	if err != nil {
		log.Fatalf("Error creating file. %v", err)
	}
	_, err = file1.WriteString(jsonString)
	if err != nil {
		log.Fatalf("Error writing to the file. %v", err)
	}
}

func ConvertFileToJSON(filePath string) string {
	currentDir, _ := os.Getwd()
	file, err := xlsx.OpenFile(currentDir + "/" + filePath)
	if err != nil {
		log.Panicf("Error while opening excel file - %s. Error - %v", filePath, err)
	}
	translationJSONMap := map[string]interface{}{}
	for rowIndex, row := range file.Sheets[0].Rows {
		if rowIndex == 0 {
			continue
		}
		translationKey := row.Cells[0].Value
		if translationKey == "" {
			break
		}
		translationValue := row.Cells[3].Value

		convertToNestedMap(translationKey, translationValue, translationJSONMap)
	}
	bytes, err := json.MarshalIndent(translationJSONMap, "", "  ")
	return string(bytes)
}

func convertToNestedMap(translationKey, translationValue string, translationJSONMap map[string]interface{}) {
	keyParts := strings.Split(translationKey, ".")
	currentKeyValueMap := translationJSONMap
	// For each key part traverse the map and check if the key already exists
	// If not, insert the key
	// Finally, after several nesting, set the value
	for i, keyPart := range keyParts {
		// End of all parent key traversal
		if i == len(keyParts)-1 {
			break
		}
		// If the key is not found in the map, insert it and set an empty string-to-interface map as value
		if val, ok := currentKeyValueMap[keyPart]; !ok {
			newKeyValueMap := map[string]interface{}{}
			currentKeyValueMap[keyPart] = newKeyValueMap
			currentKeyValueMap = newKeyValueMap
		} else {
			// Key already exists, copy it to the current key variable. Will inspect for further nesting if required.
			currentKeyValueMap = val.(map[string]interface{})
		}
	}
	// Set the value of the key read from the excel sheet
	currentKeyValueMap[keyParts[len(keyParts)-1]] = translationValue
}
