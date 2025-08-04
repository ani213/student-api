package storage

import "github.com/ani213/student-api/internal/types"

type Storage interface {
	CreateStudent(name string, age int, email string) (int64, error)
	GetStudents() ([]types.Student, error)
	GetStudentById(id int64) (types.Student, error)
	UpdateStudent(id int64, name string, age int, email string) error
}
