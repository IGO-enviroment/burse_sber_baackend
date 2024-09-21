package sql

const AddStudents = `
	insert into public.users (university_id, email, password_digest, first_name, last_name, middle_name, phone_number, is_admin, is_student, is_company, is_university, created_at, updated_at) values 
`
