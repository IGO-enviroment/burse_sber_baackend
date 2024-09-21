package students

import (
	"boilerplate/api/authentication/generation"
	"boilerplate/entity"
	"boilerplate/gen"
	"boilerplate/jwt"
	"boilerplate/sql"
	"errors"
	"fmt"
	"time"
)

func (s Service) Authenticate(authRequest gen.Login) (gen.TokenReponse, error) {
	if authRequest.Email == "" || authRequest.Password == "" {
		return gen.TokenReponse{}, errors.New("empty password or email")
	}

	user, err := s.getStudentByEmail(authRequest.Email)
	if err != nil {
		return gen.TokenReponse{}, err
	}

	if !generation.CheckPasswordHash(authRequest.Password, user.PasswordDigest) {
		return gen.TokenReponse{}, errors.New("incorrect credentials")
	}

	claims := generation.AccessTokenClaims{
		UserId:            user.Id,
		Email:             user.Email,
		IsStudent:         user.IsStudent,
		IsAdmin:           user.IsAdmin,
		IsOrganization:    user.IsCompany,
		IsUniversity:      user.IsUniversity,
		CreationTimestamp: time.Now().Unix(),
		TTL:               s.settings.AccessTokenTTL,
	}
	accessToken := jwt.GetToken(claims, s.settings.JwtSecret)

	return gen.TokenReponse{
		AccessToken: accessToken,
		ExpiresIn:   s.settings.AccessTokenTTL,
	}, nil
}

func (s Service) getStudentByEmail(email string) (entity.Student, error) {
	query := fmt.Sprintf(sql.GetStudentByEmail, email)
	rows, err := s.pg.Query(query)
	if err != nil {
		return entity.Student{}, err
	}

	defer rows.Close()

	var student entity.Student
	for rows.Next() {
		err := rows.Scan(
			&student.Id,
			&student.Email,
			&student.PasswordDigest,
			&student.IsAdmin,
			&student.IsStudent,
			&student.IsCompany,
			&student.IsUniversity,
		)
		if err != nil {
			return entity.Student{}, err
		}
	}

	return student, nil
}
