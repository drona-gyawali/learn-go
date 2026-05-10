package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/drona-gyawali/learn-go/internal/types"
	"github.com/drona-gyawali/learn-go/internal/utils/response"
	"github.com/go-playground/validator/v10"
)


func New(content string) http.HandlerFunc {
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
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string] string {"sucess": "ok"})
	}
}