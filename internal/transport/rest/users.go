package rest

import (
	"encoding/json"
	"errors"
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

	if err := checkUserFields(req); err != nil {
		responseError(w, http.StatusBadRequest, err.Error())
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

	if err := checkUserFields(req); err != nil {
		responseError(w, http.StatusBadRequest, err.Error())
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

	response200(w, "ok")
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {

	type request struct {
		ID int64 `json:"id"`
	}

	var req request
	var err error

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, http.StatusBadRequest, "")
		return
	}

	if req.ID == 0 {
		responseError(w, http.StatusBadRequest, "id is empty")
		return
	}

	if err = h.service.Users.Delete(r.Context(), req.ID); err != nil {
		responseError(w, http.StatusInternalServerError, "an error occurred while deleting the user")
		return
	}

	response200(w, "ok")
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

	if err = getUsersFilters(r, params); err != nil {
		responseError(w, http.StatusBadRequest, err.Error())
		return
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

func getUsersFilters(r *http.Request, params *domain.UsersParam) error {
	var err error

	limit := int64(0)
	limitStr := r.FormValue("limit")
	offset := int64(50)
	offsetStr := r.FormValue("offset")
	if len(limitStr) != 0 {
		limit, err = strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			return errors.New(`invalid "limit" parameter`)
		}
		params.Limit = limit
	}
	if len(offsetStr) != 0 {
		offset, err = strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			return errors.New(`invalid "offset"" parameter`)

		}
		params.Offset = offset
	}

	gender := r.FormValue("gender")
	if gender != "" {
		params.Gender = gender
	}

	status := r.FormValue("status")
	if status != "" {
		params.Status = status
	}

	fullName := r.FormValue("full_name")
	if fullName != "" {
		params.FullName = fullName
	}

	orderBy := r.FormValue("order_by")
	if orderBy == "" {
		orderBy = "id"
	}
	params.OrderBy = orderBy
	orderDir := r.FormValue("order_dir")
	if orderDir == "" {
		orderDir = "DESC"
	}
	params.OrderDir = orderDir

	return nil
}

func checkUserFields(user domain.User) error {
	if user.Name == "" {
		return errors.New(`field "Name" is required`)
	}
	if user.LastName == "" {
		return errors.New(`field "Last name" is required`)
	}
	if user.Gender == "" {
		return errors.New(`field "Gender" is required`)
	}
	if user.Status == "" {
		return errors.New(`field "Status" is required`)
	}
	return nil
}
