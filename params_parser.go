package resque_status

import "encoding/json"

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
