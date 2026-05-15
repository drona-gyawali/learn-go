package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/drona-gyawali/learn-go/internal/storage"
	"github.com/drona-gyawali/learn-go/internal/types"
	"github.com/drona-gyawali/learn-go/internal/utils/response"
	"github.com/go-playground/validator/v10"
)


func New(content string, storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating student...")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return 
		}

		if err  != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return 
		}

		
		if err := validator.New().Struct(student); err != nil{
			validateErrs := err.(validator.ValidationErrors) // data types
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		id , err := storage.CreateStudent(
			student.Name,
			student.Email, 
			student.Age,
		)

		slog.Info("user created sucessfully", slog.String("userId", fmt.Sprint(id)))
		if err != nil {
			slog.Info("user creation occured error", slog.String("userId", fmt.Sprintf(err.Error())))
			response.WriteJson(w, http.StatusInternalServerError, err)
			return 
		}


		response.WriteJson(w, http.StatusCreated, map[string] int64 {"id": id})
	}
}


func GetById(storage  storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Gettting students")

		_id , err  := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("Invalid Id types should be int64")
			response.WriteJson(w, http.StatusBadRequest, err)
			return
		}
		student, err := storage.GetStudentById(_id)
		if err != nil {
			slog.Error("Unable to get student", slog.String("err", fmt.Sprintf(err.Error())))
			response.WriteJson(w, http.StatusBadRequest, err)
			return
		}
		slog.Info("STUDENT DETAILS GET FETCHED")

		response.WriteJson(w, http.StatusOK, student)

	}
}


func GetStudentList(storage  storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting all students")

		student , err := storage.GetStudentList()
		if err != nil {
			slog.Error("Unable to get student List", slog.String("err", fmt.Sprintf(err.Error())))
			response.WriteJson(w, http.StatusBadRequest, err)
			return
		}

		slog.Info("STUDENT LIST GOT FETCHED")

		response.WriteJson(w, http.StatusOK, student)
	}
}


func UpdateStudentView(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("updating student by id")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return 
		}

		if err  != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return 
		}

		
		if err := validator.New().Struct(student); err != nil{
			validateErrs := err.(validator.ValidationErrors) // data types
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

	 

		id, err := storage.UpdateStudent(student) 
		if err != nil {
			slog.Info("user updation occured error", slog.String("userId", fmt.Sprintf(err.Error())))
			response.WriteJson(w, http.StatusInternalServerError, err)
			return 
		}

		slog.Info("user updated sucessfully", slog.String("userId", fmt.Sprint(id)))
		response.WriteJson(w, http.StatusAccepted, student)
	}
}