package main

func CopyMap(rawMap map[string]string) map[string]string {
	resultMap := make(map[string]string, len(rawMap))

	for k, v := range rawMap {
		resultMap[k] = v
	}
	return resultMap
}
