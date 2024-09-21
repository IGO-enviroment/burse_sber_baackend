package sql

const GetStudentByEmail = `
	select s.id, s.email, s.password_digest, s.is_admin, s.is_student, s.is_company, s.is_university from public.users s
	where s.email = '%s'
`
