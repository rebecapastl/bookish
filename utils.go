package main

import "strconv"

func SanitizeAuthorName(name *string) *string {
	if name == nil {
		anounymous := "anonymous"
		name = &anounymous
	}
	return name
}

func SanitizeIdNumber(id *string) (*int, error) {
	intId, err := strconv.Atoi(*id)
	if err != nil {
		return nil, err
	}
	return &intId, nil
}
