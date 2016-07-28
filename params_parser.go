package resque_status

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func ParseInt(params map[string]interface{}, name string) (int, error) {
	jsonNumber := params[name]

	intValue, err := jsonNumber.(json.Number).Int64()
	if err != nil {
		return 0, err
	}
	return int(intValue), nil
}

func ParseIntArray(params map[string]interface{}, name string) ([]int, error) {
	jsonNumbers := params[name].([]interface{})
	intNumbers := []int{}

	for _, jsonNumber := range jsonNumbers {
		intValue, err := jsonNumber.(json.Number).Int64()
		if err != nil {
			return nil, err
		}
		intNumbers = append(intNumbers, int(intValue))

	}
	return intNumbers, nil
}

func ParseJsonParam(params map[string]interface{}, name string) ([]map[string]string, error) {
	parsedJson := []map[string]string{}
	parsedInterfaces := []map[string]interface{}{}
	err := json.Unmarshal([]byte(params[name].(string)), &parsedInterfaces)
	if err != nil {
		return nil, err
	}
	for _, productData := range parsedInterfaces {
		productDataMap := map[string]string{}
		for productAttr, productAttrVal := range productData {
			switch v := productAttrVal.(type) {
			case int:
				productDataMap[productAttr] = strconv.FormatInt(int64(v), 10)
			case float64:
				productDataMap[productAttr] = strconv.FormatFloat(v, 'f', -1, 64)
			case string:
				productDataMap[productAttr] = v
			case bool:
				if v {
					productDataMap[productAttr] = "true"
				} else {
					productDataMap[productAttr] = "false"
				}
			default:
				productDataMap[productAttr] = fmt.Sprintf("%v", v)
			}
		}
		parsedJson = append(parsedJson, productDataMap)
	}
	return parsedJson, nil
}
