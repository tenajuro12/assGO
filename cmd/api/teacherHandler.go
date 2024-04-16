package main

import (
	"fmt"
	"github.com/tenajuro12/assGO/internal/data"
	"github.com/tenajuro12/assGO/internal/validator"
	"net/http"
)

func (app *application) listTeacherHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name    string
		Surname string
		data.Filters
	}
	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Surname = app.readString(qs, "surname", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "name", "surname", "-id", "-name", "-surname"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	teachers_info, err := app.models.GetTeachers(input.Name, input.Surname, input.Filters)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"teachers_info": teachers_info}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	fmt.Fprintf(w, "%+v\n", input)
}
