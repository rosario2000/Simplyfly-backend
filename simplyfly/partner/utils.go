package partner

func Contains(arr []string, lookupStr string) bool {
	for _, str := range arr {
		if str == lookupStr {
			return true
		}
	}
	return false
}
