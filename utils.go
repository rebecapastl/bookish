package main

func SanitizeAuthorName(name *string) *string {
	if name == nil {
		anounymous := "anonymous"
		name = &anounymous
	}
	return name
}