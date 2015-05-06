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
