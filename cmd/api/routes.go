package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthCheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/createModule", app.createModuleHandler)
	router.HandlerFunc(http.MethodGet, "/v1/getModule/:id", app.getModuleHandler)
	router.HandlerFunc(http.MethodGet, "/v1/getModule", app.listModuleHandler)
	router.HandlerFunc(http.MethodPut, "/v1/updateModule/:id", app.editModuleHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/deleteModule/:id", app.deleteModuleHandler)

	router.HandlerFunc(http.MethodGet, "/v1/getTeachers", app.listTeacherHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	return router
}
