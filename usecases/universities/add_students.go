package universities

import (
	"boilerplate/api/authentication/generation"
	"boilerplate/gen"
	"boilerplate/sql"
)

func (s *Service) AddStudents(request gen.AddStudent) (gen.CreationResult, error) {
	query := sql.AddStudents
	var values []interface{}
	for i, student := range request.Students {
		password := generation.GeneratePassword()
		passwordHash, err := generation.HashPassword(password)
		if err != nil {
			return gen.CreationResult{}, err
		}
		values = append(values, request.UniversityId, student.Email)
	}
}
