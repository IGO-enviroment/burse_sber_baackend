package entity

type Student struct {
	Id             int
	PasswordDigest string
	Email          string
	IsAdmin        bool
	IsStudent      bool
	IsCompany      bool
	IsUniversity   bool
}
