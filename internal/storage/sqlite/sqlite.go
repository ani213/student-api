package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/ani213/student-api/internal/config"
	"github.com/ani213/student-api/internal/types"
	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

type Sqlite struct {
	Db *sql.DB
}

func (s *Sqlite) CreateStudent(name string, age int, email string) (int64, error) {
	// This function should implement the logic to create a student in the SQLite database
	// For now, it returns a dummy value
	stmt, err := s.Db.Prepare("INSERT INTO students (name, age, email) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(name, age, email)

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return lastId, nil

}

func (s *Sqlite) GetStudents() ([]types.Student, error) {

	stmt, err := s.Db.Prepare("SELECT id, name, age, email FROM students")
	if err != nil {
		return []types.Student{}, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return []types.Student{}, err
	}
	defer rows.Close()
	var students []types.Student
	for rows.Next() {
		var student types.Student
		err = rows.Scan(&student.ID, &student.Name, &student.Age, &student.Email)
		if err != nil {
			return []types.Student{}, err
		}
		students = append(students, student)
	}
	if err = rows.Err(); err != nil {
		return []types.Student{}, err
	}
	if len(students) == 0 {
		return []types.Student{}, nil
	}
	return students, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, age, email FROM students WHERE id = ?")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()
	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.ID, &student.Name, &student.Age, &student.Email)
	if err == sql.ErrNoRows {
		return types.Student{}, fmt.Errorf("no student found with ID %d", id)
	}
	if err != nil {
		return types.Student{}, err
	}
	return student, nil
}
func (s *Sqlite) UpdateStudent(id int64, name string, age int, email string) error {
	stmt, err := s.Db.Prepare("UPDATE students SET name = ?, age = ?, email = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(name, age, email, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no student found with ID %d", id)
	}
	return nil
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students ( 
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	age INTEGER NOT NULL,
	email TEXT NOT NULL UNIQUE
	)`)
	if err != nil {
		return nil, err
	}
	return &Sqlite{
		Db: db,
	}, nil
}
