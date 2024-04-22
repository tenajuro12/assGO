package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthCheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/createModule", app.createModuleHandler)
	router.HandlerFunc(http.MethodGet, "/v1/getModule/:id", app.getModuleHandler)
	router.HandlerFunc(http.MethodGet, "/v1/getModule", app.listModuleHandler)
	router.HandlerFunc(http.MethodPut, "/v1/updateModule/:id", app.editModuleHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/deleteModule/:id", app.deleteModuleHandler)

	router.HandlerFunc(http.MethodGet, "/v1/getTeachers", app.listTeacherHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodGet, "/v1/getUser/:id", app.requireActivatedUser(app.requirePermission("users:read", app.getUserHandler)))
	router.HandlerFunc(http.MethodGet, "/v1/getUser", app.requireActivatedUser(app.listUserHandler))
	router.HandlerFunc(http.MethodPut, "/v1/updateUser/:id", app.requireActivatedUser(app.editUserHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/deleteUser/:id", app.requireActivatedUser(app.deleteUserHandler))

	return app.recoverPanic(app.authenticate(router))
}
