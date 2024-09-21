package universities

import (
	"boilerplate/api/authentication/generation"
	"boilerplate/gen"
	"boilerplate/sql"
	"strconv"
	"time"
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
		createdAt := time.Now()
		values = append(values, request.UniversityId, student.Email, passwordHash, student.FirstName, student.LastName, student.MiddleName, student.PhoneNumber, false, true, false, false, createdAt, createdAt)
		numFields := 13
		n := i * numFields
		query += `(`
		for j := 0; j < numFields; j++ {
			query += `$` + strconv.Itoa(n+j+1) + `,`
		}
		query = query[:len(query)-1] + `),`
	}
	query = query[:len(query)-1]
	query += " RETURNING id"
	rows, err := s.pg.Exec(query, values...)
	if err != nil {
		return gen.CreationResult{}, err
	}
	_ = rows
	return gen.CreationResult{}, nil
}
