package main

import (
	"github.com/tenajuro12/assGO/internal/data"
	"net/http"
)

func (app *application) createModuleHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Modile_name     string `json:"module_name"`
		Module_duration int32  `json:"module_duration"`
		ExamType        string `json:"exam_type"`
	}

	err := app.readJSON(w, r, &input)
	module := &data.Models{}
}
