package tools

func Has(v []int, s int) bool {
	for _, a := range v {
		if a == s {
			return true
		}
	}
	return false
}
