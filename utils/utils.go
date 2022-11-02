package utils

import (
	"fmt"
	//
	"github.com/gruntwork-io/terratest/modules/terraform"
	//
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

type MSI = map[string]interface{}

const horizontalLine = "**********************"

var errorQueue = make([]error, 0)

func WatchForError(e error) {
	if e != nil {
		errorQueue = append(errorQueue, e)
	}
}

func DumpErrors(t *testing.T) {
	if len(errorQueue) > 0 {
		for len(errorQueue) > 0 {
			t.Log(errorQueue[0])
			errorQueue = errorQueue[1:]
		}
		t.FailNow()
	}
}

func StringByOS(stringMap MSI) string {
	ret, ok := stringMap[runtime.GOOS]
	if !ok {
		WatchForError(fmt.Errorf("missing string for operating system: %s\nin map: %v", runtime.GOOS, stringMap))
		return "os_string_error_occurred"
	}
	return ret.(string)
}

func FileSha1(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error in filesha1: %s", err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func CheckPlanExpectations(t *testing.T, expectedResources []MSI, plan *terraform.PlanStruct) {
	for _, resource := range expectedResources {
		if resource["address"] != nil {
			// Process single address
			if resource["fields_override"] == nil {
				// not counted
				t.Log(horizontalLine)
				t.Logf("Checking plan for address: %s\n\n", resource["address"].(string))
				t.Log(horizontalLine)

				if plan.ResourcePlannedValuesMap[resource["address"].(string)] == nil {
					WatchForError(fmt.Errorf("x_x Required address not found in plan: %s", resource["address"].(string)))
				} else {
					planResource := plan.ResourcePlannedValuesMap[resource["address"].(string)]

					if resource["fields"] != nil {
						for fieldKey, fieldValue := range resource["fields"].(MSI) {
							WatchForError(CheckValuesEqualE(t, fieldKey, fieldValue, planResource.AttributeValues, resource["address"].(string)))
						}
					}
				}
			} else {
				// Counted string address - address["0"], address["1"]..., address["n"]
				// Number of checks is based off size of fields_override
				for _, override := range resource["fields_override"].([]MSI) {
					found := false
					for i := range resource["fields_override"].([]MSI) {
						address := fmt.Sprintf(`%s["%d"]`, resource["address"].(string), i)

						t.Log(horizontalLine)
						t.Logf("^-^ Checking plan for address: %s\n\n", address)
						t.Log(horizontalLine)

						if plan.ResourcePlannedValuesMap[address] == nil {
							WatchForError(fmt.Errorf("x_x Required address not found in plan: %s", address))
						} else {
							planResource := plan.ResourcePlannedValuesMap[address]

							if resource["fields"] != nil {
								for fieldKey, fieldValue := range resource["fields"].(MSI) {
									err := CheckValuesEqualE(t, fieldKey, fieldValue, planResource.AttributeValues, address)
									if err != nil {
										break
									}
								}
							}
							for fieldKey, fieldValue := range override {
								err := CheckValuesEqualE(t, fieldKey, fieldValue, planResource.AttributeValues, address)
								if err != nil {
									break
								}
							}
							found = true
						}
					}
					if !found {
						WatchForError(fmt.Errorf("x_x Did not find matching resource for counted field override. \nAddress: %s\nfields: %v\nfields_override: %v",
							resource["address"].(string), resource["fields"].(MSI), override))
					}
				}
			}
		} else if resource["addresses"] != nil {
			// Process multiple addresses
			for i, address := range resource["addresses"].([]string) {
				t.Log(horizontalLine)
				t.Logf("^-^ Checking plan for address: %s\n\n", address)
				t.Log(horizontalLine)

				if plan.ResourcePlannedValuesMap[address] == nil {
					WatchForError(fmt.Errorf("x_x Required address not found in plan: %s", address))
				} else {
					planResource := plan.ResourcePlannedValuesMap[address]
					if resource["fields"] != nil {
						for fieldKey, fieldValue := range resource["fields"].(MSI) {
							WatchForError(CheckValuesEqualE(t, fieldKey, fieldValue, planResource.AttributeValues, address))
						}
					}
					if resource["fields_override"] != nil {
						for fieldKey, fieldValue := range resource["fields_override"].([]MSI)[i] {
							WatchForError(CheckValuesEqualE(t, fieldKey, fieldValue, planResource.AttributeValues, address))
						}
					}
				}
			}
		}
	}
}

func CheckValuesE(t *testing.T, expectedKey string, expectedValue interface{}, attributeValues MSI, allowPartialMatch bool, context string) error {
	if attributeValues[expectedKey] == nil {
		if expectedValue == nil {
			return nil
		} else {
			return fmt.Errorf("x_x Error context %s: expected key %s not found. Could it be nested improperly?", context, expectedKey)
		}
	}
	switch v := expectedValue.(type) {
	case string:
		// Defer an inline function to catch panics from type assumptions
		defer func(expectedKey string) {
			if err := recover(); err != nil {
				WatchForError(fmt.Errorf("x_x Error context %s: Incorrect type assumption. %s is not type string\n%s", context, expectedKey, err))
			}
		}(expectedKey)
		// Try the type assumption
		foundValue := attributeValues[expectedKey].(string)

		// Comparing string could be equal or regex based
		// Change the log to reflect this
		compareString := "=="
		if allowPartialMatch {
			compareString = "~="
		}

		t.Logf("Checking %s: '%s' %s '%s'\n\n\n", expectedKey, expectedValue.(string), compareString, foundValue)

		if allowPartialMatch {
			// if !string.Contains(foundValue, expectedValue.(string)) {
			matched, err := regexp.MatchString(expectedValue.(string), foundValue)
			t.Logf("regex matched: %v, err: %v, expected: %v, found %v\n\n\n", matched, err, expectedValue, foundValue)
			if err != nil {
				return err
			} else if !matched {
				return fmt.Errorf("x_x Error context %s: found value %s did not contain regexp expected %s", context, foundValue, expectedValue.(string))
			}
		} else if expectedValue.(string) != foundValue {
			return fmt.Errorf("x_x Error context %s: found value %s did not match expected %s", context, foundValue, expectedValue.(string))
		}
	case int:
		// Defer an inline function to catch panics from type assumptions
		defer func(expectedKey string) {
			if err := recover(); err != nil {
				WatchForError(fmt.Errorf("x_x Error context %s: Incorrect type assumption. %s is not type integer\n%s", context, expectedKey, err))
			}
		}(expectedKey)
		// try the type assumption
		foundValue := int(attributeValues[expectedKey].(float64))
		t.Logf("Checking %s: '%d' == '%d'\n\n\n", expectedKey, expectedValue.(int), foundValue)
		if expectedValue.(int) != foundValue {
			return fmt.Errorf("x_x Error context %s: found value %s did not match expected %d", context, foundValue, expectedValue.(int))
		}
	case float64:
		// Defer an inline function to catch points from type assumptions
		defer func(expectedKey string) {
			if err := recover(); err != nil {
				WatchForError(fmt.Errorf("x_x Error context %s: Incorrect type assumption. %s is not type float64\n%s", context, expectedKey, err))
			}
		}(expectedKey)
		// try the type assumption
		foundValue := attributeValues[expectedKey].(float64)
		t.Logf("Checking %s: '%v' == '%v'\n\n\n", expectedKey, expectedValue.(float64), foundValue)
		if expectedValue.(float64) != foundValue {
			return fmt.Errorf("x_x Error context %s: found value %s did not match expected %d", context, foundValue, expectedValue.(float64))
		}
	case MSI:
		t.Logf("Checking map %v == %v\n\n", expectedKey, attributeValues[expectedKey])
		// Defer an inline function to catch points from type assumptions
		defer func(expectedKey string) {
			if err := recover(); err != nil {
				WatchForError(fmt.Errorf("x_x Error context %s: Incorrect type assumption. %s is not type map[string]interface\n%s", context, expectedKey, err))
			}
		}(expectedKey)
		// try the type assumption
		foundValue := attributeValues[expectedKey].(MSI)
		for fieldKey, fieldValue := range expectedValue.(MSI) {
			var err error
			if allowPartialMatch {
				err = CheckValuesContainsE(t, fieldKey, fieldValue, foundValue, fmt.Sprintf("%s -> %s", context, fieldKey))
			} else {
				err = CheckValuesEqualE(t, fieldKey, fieldValue, foundValue, fmt.Sprintf("%s -> %s", context, fieldKey))
			}
			if err != nil {
				return err
			}
		}
	case []interface{}:
		t.Logf("Checking slice %v == %v\n\n", expectedValue, attributeValues[expectedKey])
		// Defer an inline function to catch points from type assumptions
		defer func(expectedKey string) {
			if err := recover(); err != nil {
				WatchForError(fmt.Errorf("x_x Error context %s: Incorrect type assumption. %s is not type []interface\n%s", context, expectedKey, err))
			}
		}(expectedKey)
		// try the type assumption
		foundValues := attributeValues[expectedKey].([]interface{})

		for _, eValue := range expectedValue.([]interface{}) {
			t.Logf("Checking existence of %s: %v in %v\n\n\n", expectedKey, eValue, foundValues)
			found := false
			for _, fValue := range foundValues {
				if allowPartialMatch {
					switch eValue.(type) {
					case string:
						//
						matched, err := regexp.MatchString(eValue.(string), fValue.(string))
						if err != nil {
							return err
						} else if matched {
							found = true
						}
					default:
						if eValue == fValue {
							found = true
						}
					}
				} else if eValue == fValue {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("x_x Error context %s: found \n%v\n did not contain expected value: %v", context, foundValues, eValue)
			}
		}
	case []MSI:
		t.Logf("Checking slice of maps %v == %v\n\n", expectedValue, attributeValues[expectedKey])
		// Defer an inline function to catch points from type assumptions
		defer func(expectedKey string) {
			if err := recover(); err != nil {
				WatchForError(fmt.Errorf("x_x Error context %s: Incorrect type assumption. %s is not type map[string]interface\n%s", context, expectedKey, err))
			}
		}(expectedKey)
		// try the type assumption
		foundValues := attributeValues[expectedKey].([]interface{})
		var err error

		for _, eMap := range expectedValue.([]MSI) {
			t.Logf("Checking existence of %s: %v in %v\n\n\n", expectedKey, eMap, foundValues)
			var found bool
			for _, fValue := range foundValues {
				found = true
				t.Logf("   Checking: %v in %v\n\n", eMap, fValue)
				for eK, eV := range eMap {
					t.Logf("   Checking: %v:%v in %v\n\n", eK, eV, fValue)
					if allowPartialMatch {
						err = CheckValuesContainsE(t, eK, eV, fValue.(map[string]interface{}), fmt.Sprintf("%s -> %s", context, eK))
					} else {
						err = CheckValuesEqualE(t, eK, eV, fValue.(map[string]interface{}), fmt.Sprintf("%s -> %s", context, eK))
					}
					if err != nil {
						t.Log("not found\n\n")
						found = false
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				return fmt.Errorf("x_x Error context %s: found \n%v\n did not contain expected value: %v\n\nError last message: %s\n\n\n", context, foundValues, eMap, err)
			}
		}
	case bool:
		t.Logf("Checking bool %v == %v\n\n", expectedValue, attributeValues[expectedKey])
		// Defer an inline function to catch points from type assumptions
		defer func(expectedKey string) {
			if err := recover(); err != nil {
				WatchForError(fmt.Errorf("x_x Error context %s: Incorrect type assumption. %s is not type bool\n%s", context, expectedKey, err))
			}
		}(expectedKey)
		// try the type assumption
		foundValue := attributeValues[expectedKey].(bool)
		t.Logf("   Checking: %s: '%t' == '%t'\n\n", expectedKey, expectedValue.(bool), foundValue)
		if expectedValue.(bool) != foundValue {
			return fmt.Errorf("x_x Error context %s: found value %T did not contain expected %T\n\n\n", context, foundValue, expectedValue.(bool))
		}
	case func(*testing.T, string) error:
		// Defer an inline function to catch points from type assumptions
		defer func(expectedKey string) {
			if err := recover(); err != nil {
				WatchForError(fmt.Errorf("x_x Error context %s: Incorrect type assumption. %s is not type string\n%s", context, expectedKey, err))
			}
		}(expectedKey)
		// try the type assumption
		foundValue := attributeValues[expectedKey].(string)

		return expectedValue.(func(*testing.T, string) error)(t, foundValue)

	default:
		return fmt.Errorf("x_x Error context %s: found value %T did not contain expected %v\n\n\n", context, v, expectedValue)
	}
	return nil
}

func CheckValuesEqualE(t *testing.T, expectedKey string, expectedValue interface{}, attributeValues MSI, context string) error {
	return CheckValuesE(t, expectedKey, expectedValue, attributeValues, false, context)
}

func CheckValuesEqual(t *testing.T, expectedKey string, expectedValue interface{}, attributeValues MSI, context string) {
	err := CheckValuesE(t, expectedKey, expectedValue, attributeValues, false, context)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func CheckValuesContainsE(t *testing.T, expectedKey string, expectedValue interface{}, attributeValues MSI, context string) error {
	return CheckValuesE(t, expectedKey, expectedValue, attributeValues, true, context)
}

func CheckValuesContains(t *testing.T, expectedKey string, expectedValue interface{}, attributeValues MSI, context string) {
	err := CheckValuesContainsE(t, expectedKey, expectedValue, attributeValues, context)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func Contains(expected string) func(*testing.T, string) error {
	return func(t *testing.T, found string) error {
		t.Logf("   Checking: %s work in %s\n\n", expected, found)
		if !strings.Contains(found, expected) {
			return fmt.Errorf("found string: %s did not contain: \"%s\"", found, expected)
		}
		t.Logf("found\n\n")

		return nil
	}
}

func ContainsAll(expected []string) func(*testing.T, string) error {
	return func(t *testing.T, found string) error {
		for _, s := range expected {
			t.Logf("   Checking: %s in %s\n\n", s, found)
			if !strings.Contains(found, s) {
				return fmt.Errorf("found string: %s did not contain: \"%s\"", found, s)
			}
			t.Logf("found\n\n")
		}

		return nil
	}
}

func Regex(expected string) func(*testing.T, string) error {
	return func(t *testing.T, found string) error {
		matched, err := regexp.MatchString(expected, found)
		t.Logf("   regex matched: %v, err: %v, expected: %v, found: %v\n\n", matched, err, expected, found)
		if err != nil {
			return err
		} else if !matched {
			return fmt.Errorf("x_x Error: found value %s did not contain regexp expected %s", found, expected)
		}

		return nil
	}
}
