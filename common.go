package main

func reverseCopy(s []string) []string {
	len := len(s)
	r := make([]string, len)

	for i := range s {
		r[len-(i+1)] = s[i]
	}

	return r
}
