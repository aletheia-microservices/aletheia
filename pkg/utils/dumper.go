package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"gopkg.in/yaml.v2"

	"analyzer/pkg/logger"
)

const BASE_FOLDER = "output"

// -----------------------
// Yaml Formatting Helpers
// -----------------------

type OrderedProperties []map[string]interface{}
type OrderedPropertyList struct {
	// helper structure for saving properties while ordered
	// if we use a map (instead of a list entry with a single map)
	// then properties are ordered by their keys and not by insertion
	list OrderedProperties
}

func NewOrderedPropertyList() OrderedPropertyList {
	return OrderedPropertyList{}
}
func (lst *OrderedPropertyList) Result() OrderedProperties {
	return lst.list
}
func (lst *OrderedPropertyList) AddOrderedProperty(name string, prop interface{}) {
	fieldsMap := make(map[string]interface{})
	fieldsMap[name] = prop
	lst.list = append(lst.list, fieldsMap)
}

// -------
// Dumpers
// -------

func DumpToDebugFile(data string, appname string, filename string) {
	path := fmt.Sprintf("%s/%s/%s.cpp", BASE_FOLDER, appname, filename)
	// ensure the directory exists
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		logger.Logger.Fatalf("error creating directory %s: %s", dir, err.Error())
	}
	err = os.WriteFile(path, []byte(data), 0644)
	if err != nil {
		logger.Logger.Fatalf("error writing data to %s: %s", path, err.Error())
	}
	logger.Logger.Tracef("[JSON] go file %s", path)
}

func DumpToJSONFile(data interface{}, appname string, filename string) {
	path := fmt.Sprintf("%s/%s/%s.json", BASE_FOLDER, appname, filename)
	// ensure the directory exists
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		logger.Logger.Fatalf("error creating directory %s: %s", dir, err.Error())
	}
	// marshal json data
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		logger.Logger.Fatalf("error marshaling json: %s", err.Error())
	}
	// write to file
	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		logger.Logger.Fatalf("error writing yaml data to %s: %s", path, err.Error())
	}
	logger.Logger.Tracef("[JSON] saved file %s", path)
}

func DeleteFolder(appname string, foldername string) {
	path := fmt.Sprintf("%s/%s/%s", BASE_FOLDER, appname, foldername)
	err := os.RemoveAll(path)
	if err != nil {
		logger.Logger.Fatalf("error deleting folder %s: %s", path, err.Error())
	}
	logger.Logger.Tracef("[UTILS] deleted folder %s", path)
}

func addNewLinesToYaml(input []byte) []byte {
	lines := strings.Split(string(input), "\n")

	var result []string
	for i, line := range lines {
		if i > 0 && line != "" && !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "-") {
			result = append(result, "")
		}
		result = append(result, line)
	}

	return []byte(strings.Join(result, "\n"))
}

func DumpToYamlFile(data interface{}, appname string, filename string) {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		logger.Logger.Fatalf("error marshalling yaml data")
	}

	yamlStr := fixYamlStrings(string(addNewLinesToYaml(yamlData)))
	//yamlStr := string(yamlData)
	path := fmt.Sprintf("%s/%s/%s.yaml", BASE_FOLDER, appname, filename)

	// ensure the directory exists
	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		logger.Logger.Fatalf("error creating directory %s: %s", dir, err.Error())
	}

	// write the YAML data to the file
	err = os.WriteFile(path, []byte(yamlStr), 0644)
	if err != nil {
		logger.Logger.Fatalf("error writing yaml data to %s: %s", path, err.Error())
	}
	logger.Logger.Tracef("[YAML] saved file %s", path)
}

func DumpOutput(data string, appname string, filename string) {
	//yamlStr := string(yamlData)
	path := fmt.Sprintf("%s/%s/%s.out", BASE_FOLDER, appname, filename)

	// ensure the directory exists
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		logger.Logger.Fatalf("[OUT] error creating directory %s: %s", dir, err.Error())
	}

	// write the YAML data to the file
	err = os.WriteFile(path, []byte(data), 0644)
	if err != nil {
		logger.Logger.Fatalf("[OUT] error writing data to %s: %s", path, err.Error())
	}
	logger.Logger.Tracef("[OUT] saved file %s", path)
}

// remove unwanted new lines within items with braces
// all of this still applies to nested braces
// ensure there is always a space between an open brace and next char or a prev char and close brace
// exception goes for when there is just "interface{}"", where we maintain the lack of space
func fixYamlStrings(yamlStr string) string {
	var sb strings.Builder
	inBraces := false      // { }
	inAngles := false      // < >
	inQuotes := false      // " " or ' '
	inParentheses := false // ( )

	for i, r := range yamlStr {
		if !inQuotes {
			switch r {
			// outside quotes: open braces
			case '{':
				inBraces = true
				sb.WriteRune(r)
				/* if i+1 < len(yamlStr) && !unicode.IsSpace(rune(yamlStr[i+1])) && yamlStr[i+1] != '{' && yamlStr[i+1] != '}' {
					sb.WriteRune(' ')
				} */
				continue
			// outside quotes: close braces
			case '}':
				inBraces = false
				/* if sb.Len() > 0 && sb.String()[sb.Len()-1] != ' ' && sb.String()[sb.Len()-1] != '{' {
					sb.WriteRune(' ')
				} */
				sb.WriteRune(r)
				continue
			// outside quotes: open angles
			case '<':
				inAngles = true
				sb.WriteRune(r)
				continue
			// outside quotes: close angles
			case '>':
				inAngles = false
				sb.WriteRune(r)
				continue
			// outside quotes: open parentheses
			case '(':
				inParentheses = true
				sb.WriteRune(r)
				continue
			// outside quotes: close parentheses
			case ')':
				inParentheses = false
				sb.WriteRune(r)
				continue
			// outside quotes: open quotes
			case '"', '\'':
				inQuotes = true
				sb.WriteRune(r)
				continue
			}
		} else {
			switch r {
			// inside quotes: close quotes
			case '"', '\'':
				inQuotes = false
			}
			sb.WriteRune(r)
			continue
		}

		if inBraces || inAngles || inParentheses {
			if r == '\n' {
				// replace newlines inside braces, angles, or parentheses with spaces
				sb.WriteRune(' ')
			} else if unicode.IsSpace(r) {
				// skip multiple spaces inside braces, angles, or parentheses
				if sb.Len() > 0 && sb.String()[sb.Len()-1] != ' ' {
					sb.WriteRune(' ')
				}
			} else {
				sb.WriteRune(r)
			}
		} else {
			if r == '\n' && i+1 < len(yamlStr) && yamlStr[i+1] == ':' {
				// skip newline if next character is colon
				continue
			}
			if i > 0 && len(yamlStr) > i-1 {
				if r == ':' && yamlStr[i-1] == '\n' {
					// add newline after colon
					sb.WriteRune(r)
					sb.WriteRune('\n')
					continue
				}
				if r == '?' && yamlStr[i-1] == '\n' {
					continue
				}
				if i > 1 && len(yamlStr) > i-2 {
					if r == ' ' && yamlStr[i-1] == '?' && yamlStr[i-2] == '\n' {
						continue
					}
				}
				if r == ' ' && yamlStr[i-1] == ':' && yamlStr[i+1] == '-' {
					sb.WriteRune('\t')
					continue
				}
			}
			/* if r == '-' && yamlStr[i-1] == '\n' {
				sb.WriteRune('\t')
				sb.WriteRune(r)
				continue
			} */
			sb.WriteRune(r)
		}
	}

	return sb.String()
}
