package rest

import (
	"github.com/gorilla/mux"
	"go-api-test/internal/service"
	"net/http"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()

	api := router.PathPrefix("/api/v1").Subrouter()

	{
		api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			response200(w, "ok")
		}).Methods("GET")

		users := api.PathPrefix("/users").Subrouter()
		{
			users.HandleFunc("/", h.getUsersList).Methods("GET")
			users.HandleFunc("/", h.createUser).Methods("POST")
			users.HandleFunc("/", h.updateUser).Methods("PUT")
			users.HandleFunc("/{id}", h.deleteUser).Methods("DELETE")
			users.HandleFunc("/{id}", h.getUserByID).Methods("GET")

		}
	}
	//{
	//	api.GET("/health", func(context *gin.Context) {
	//		response200(context, nil)
	//	})
	//
	//	persons := api.Group("/persons")
	//	{
	//		//persons.GET("/", h.getPersonsList)
	//		//persons.PATCH("/:id", h.updatePerson)
	//		//persons.POST("/", h.createPerson)
	//		//persons.DELETE("/:id", h.deletePerson)
	//	}
	//}

	return router
}
