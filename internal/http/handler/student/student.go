package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ani213/student-api/internal/types"
	"github.com/ani213/student-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func CreateStudent() http.HandlerFunc {
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

		response.WriteJson(w, http.StatusCreated, map[string]string{"message": "Student created successfully", "student_id": strconv.Itoa(student.ID)})
	}

}
