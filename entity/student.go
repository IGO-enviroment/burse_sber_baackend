package entity

import "time"

type Student struct {
	Id             int
	PasswordDigest string
	Email          string
	IsAdmin        bool
	IsStudent      bool
	IsCompany      bool
	IsUniversity   bool
	UniversityId   int
	FirstName      string
	LastName       string
	MiddleName     string
	PhoneNumber    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
