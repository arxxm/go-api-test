package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-api-test/pkg/domain"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var req domain.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, http.StatusBadRequest, "")
		return
	}

	if msg := checkUserFields(req); msg != "" {
		responseError(w, http.StatusBadRequest, msg)
		return
	}

	req.CreatedAt = time.Now().UTC()
	id, err := h.service.Users.Create(r.Context(), req)
	if err != nil {
		responseError(w, http.StatusInternalServerError, "an error occurred while creating the user")
		return
	}

	response201(w, id)
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	var req domain.User
	var err error

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, http.StatusBadRequest, "")
		return
	}

	if msg := checkUserFields(req); msg != "" {
		responseError(w, http.StatusBadRequest, msg)
		return
	}

	params := mux.Vars(r)
	idStr := params["id"]

	if idStr == "" {
		responseError(w, http.StatusBadRequest, "id is empty")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		responseError(w, http.StatusBadRequest, "invalid id parameter "+err.Error())
	}

	if err = h.service.Users.Update(r.Context(), int64(id), req); err != nil {
		responseError(w, http.StatusInternalServerError, "an error occurred while updating the user")
		return
	}

	response200(w, nil)
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	idStr := params["id"]

	if idStr == "" {
		responseError(w, http.StatusBadRequest, "id is empty")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		responseError(w, http.StatusBadRequest, "invalid id parameter "+err.Error())
	}

	if err = h.service.Users.Delete(r.Context(), int64(id)); err != nil {
		responseError(w, http.StatusInternalServerError, "an error occurred while deleting the user")
		return
	}

	response200(w, nil)
}

func (h *Handler) getUserByID(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	idStr := params["id"]

	if idStr == "" {
		responseError(w, http.StatusBadRequest, "id is empty")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		responseError(w, http.StatusBadRequest, "invalid id parameter "+err.Error())
	}

	var user domain.User
	if user, err = h.service.Users.GetByID(r.Context(), int64(id)); err != nil {
		responseError(w, http.StatusInternalServerError, "an error occurred while retrieving the user")
		return
	}

	response200(w, user)
}

func (h *Handler) getUsersList(w http.ResponseWriter, r *http.Request) {
	var err error
	var params = &domain.UsersParam{
		Limit:  10,
		Offset: 0,
	}

	limit := int64(0)
	limitStr := r.FormValue("limit")
	offset := int64(50)
	offsetStr := r.FormValue("offset")
	if len(limitStr) != 0 {
		limit, err = strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			responseError(w, http.StatusBadRequest, "invalid limit parameter")
			return
		}
		params.Limit = limit
	}
	if len(offsetStr) != 0 {
		offset, err = strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			responseError(w, http.StatusBadRequest, "invalid offset parameter")
			return
		}
		params.Offset = offset
	}

	users, total, err := h.service.Users.GetList(r.Context(), params)
	if err != nil {
		responseError(w, http.StatusInternalServerError, "an error occurred while retrieving the list of users")
		return
	}

	type response struct {
		Users []domain.User `json:"persons"`
		Total uint64        `json:"total"`
	}

	var res = response{
		Users: users,
		Total: total,
	}

	response200(w, res)
}

func checkUserFields(user domain.User) string {
	if user.Name == "" {
		return `field "Name" is required`
	}
	if user.LastName == "" {
		return `field "LastName" is required`
	}
	if user.Gender == "" {
		return `field "Gender" is required`
	}
	if user.Status == "" {
		return `field "Status" is required`
	}
	return ""
}
