package main

import (
	"errors"
	"fmt"
	"github.com/tenajuro12/assGO/internal/data"
	"github.com/tenajuro12/assGO/internal/validator"
	"net/http"
)

func (app *application) createModuleHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Module_name     string `json:"module_name"`
		Module_duration int32  `json:"module_duration"`
		ExamType        string `json:"exam_type"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	module := &data.ModuleInfo{
		ModuleName:     input.Module_name,
		ModuleDuration: input.Module_duration,
		ExamType:       input.ExamType,
	}
	//v := validator.New()
	/*if data.ValidateModule(v, module); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}*/
	err = app.models.ModuleInfo.Insert(module)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", module.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"module info": module}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) getModuleHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	module, err := app.models.ModuleInfo.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"module": module}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) editModuleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	module, err := app.models.ModuleInfo.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		Module_name     string `json:"module_name"`
		Module_duration int32  `json:"module_duration"`
		ExamType        string `json:"exam_type"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	module.ModuleName = input.Module_name
	module.ModuleDuration = input.Module_duration
	module.ExamType = input.ExamType

	err = app.models.ModuleInfo.Update(module)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"module": module}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteModuleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.ModuleInfo.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) listModuleHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		ModuleName string
		ExamType   string
		data.Filters
	}
	v := validator.New()

	qs := r.URL.Query()

	input.ModuleName = app.readString(qs, "module_name", "")
	input.ExamType = app.readString(qs, "exam_type", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "module_name", "module_duration", "-id", "-module_name", "-module_duration"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	module_info, err := app.models.ModuleInfo.GetAll(input.ModuleName, input.ExamType, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"module_info": module_info}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	fmt.Fprintf(w, "%+v\n", input)
}
