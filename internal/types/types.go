package types

type Student struct {
	ID    int64  `json:"id"`
	Name  string `json:"name" validate:"required"`
	Age   int    `json:"age" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type StudentsResponse struct {
	Students []Student `json:"students"`
	Status   string    `json:"status"`
}

type GetStudentByIdResponse struct {
	Student Student `json:"student"`
}
