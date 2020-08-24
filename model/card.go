package model

type Card struct {
	ID       int64  `json:id`
	Name     string `json:name`
	Lastname string `json:lastname`
	Surname  string `json:surname`
	Phone    string `json:phone`
	SchoolId int64  `json:schoolId`
}
