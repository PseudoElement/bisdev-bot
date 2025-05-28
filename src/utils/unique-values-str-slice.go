package utils

func FilterUnique(s []string) []string {
	unique := make([]string, 0, len(s))
	m := make(map[string]int8)
	for _, name := range s {
		_, ok := m[name]
		if !ok {
			m[name] = 0
			unique = append(unique, name)
		}
	}

	return unique
}
