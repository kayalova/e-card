package model

type StudentCard struct {
	ID          int64  `json:id`
	Name        string `json:name`
	Lastname    string `json:lastname`
	Surname     string `json:surname`
	PhoneNumber string `json:phoneNumber`
	SchoolName  string `json:schoolName`
}
