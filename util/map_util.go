package util

func ContainsAllKeys(keys []string, m map[string]interface{}) (notContainedKeys []string) {
	for _, k := range keys {
		if _, exists := m[k]; !exists {
			notContainedKeys = append(notContainedKeys, k)
		}
	}
	return
}
