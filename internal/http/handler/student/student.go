package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ani213/student-api/internal/storage"
	"github.com/ani213/student-api/internal/types"
	"github.com/ani213/student-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func CreateStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logic to create a student
		slog.Info("Creating... student")

		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("request body is empty")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		// Validate the student data
		if err := validator.New().Struct(student); err != nil {
			validationError := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validationError))
			return
		}

		// Create the student in the storage
		lastId, err := storage.CreateStudent(student.Name, student.Age, student.Email)
		if err != nil {
			slog.Error("Failed to create student", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, err.Error())
			return
		}
		slog.Info("Student created successfully", slog.Int64("student_id", lastId))
		response.WriteJson(w, http.StatusCreated, map[string]string{"message": "Student created successfully", "student_id": strconv.Itoa(int(lastId))})
	}

}

func GetStudents(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logic to get all students
		slog.Info("Fetching all students")
		students, err := storage.GetStudents()
		if err != nil {
			slog.Error("Failed to fetch students", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		slog.Info("Students fetched successfully", slog.Int("count", len(students)))
		response.WriteJson(w, http.StatusOK, types.StudentsResponse{
			Students: students,
			Status:   "success",
		})
	}
}

func GetStudentById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logic to get a student by ID
		id := r.PathValue("id")

		slog.Info("Fetching student by ID", slog.String("id", id))
		studentId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("Invalid student ID", slog.String("id", id), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid student ID")))
			return
		}
		student, err := storage.GetStudentById(studentId)
		if err != nil {
			slog.Error("Failed to fetch student", slog.String("id", id), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		slog.Info("Student fetched successfully", slog.Int64("student_id", student.ID))
		response.WriteJson(w, http.StatusOK, types.GetStudentByIdResponse{
			Student: student,
		})
	}
}

func UpdateStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logic to update a student
		id := r.PathValue("id")
		slog.Info("Updating student", slog.String("id", id))
		studentId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("Invalid student ID", slog.String("id", id), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid student ID")))
			return
		}
		body := r.Body
		defer body.Close()
		var student types.Student
		err = json.NewDecoder(body).Decode(&student)
		if err != nil {
			slog.Error("Failed to decode request body", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Validate the student data
		if err := validator.New().Struct(student); err != nil {
			validationError := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validationError))
			return
		}
		// Update the student in the storage
		err = storage.UpdateStudent(studentId, student.Name, student.Age, student.Email)
		if err != nil {
			slog.Error("Failed to update student", slog.String("id", id), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		slog.Info("Student updated successfully", slog.Int64("student_id", studentId))
		response.WriteJson(w, http.StatusOK, map[string]string{"message": "Student updated successfully"})

	}
}
