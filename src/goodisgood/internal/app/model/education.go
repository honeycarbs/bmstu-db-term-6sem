package model

type EducationPlace struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type EducationProgram struct {
	Field string `json:"field"`
	Level string `json:"level"`
}
