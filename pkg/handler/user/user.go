package user

import (
	"brief/internal/constant"
	"brief/internal/model"
	"brief/service/user"
	"encoding/json"
	"net/http"

	"brief/utility"

	"github.com/go-chi/chi/v5"
)

// Register - /users - POST
func (base *Controller) Register(w http.ResponseWriter, r *http.Request) {
	req := new(model.User)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	if err := base.Validate.Struct(req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrValidation, utility.ValidationResponse(err, base.Validate), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	if err := user.Register(req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "user created successfully", req)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusCreated)
	w.Write(res)

}

// Login - /users/login - POST
func (base *Controller) Login(w http.ResponseWriter, r *http.Request) {
	req := new(model.UserLogin)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	if err := base.Validate.Struct(req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrValidation, utility.ValidationResponse(err, base.Validate), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	usr, err := user.Login(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "logged in successfully", usr)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

// GetMe - /users - GET
func (base *Controller) GetMe(w http.ResponseWriter, r *http.Request) {
	uId := r.Context().Value("id")

	if uId == nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, "user ID not found", nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	usr, err := user.Get(uId.(string))
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", usr)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// UpdateMe - /users - PATCH
func (base *Controller) UpdateMe(w http.ResponseWriter, r *http.Request) {

	req := new(model.User)
	uId := r.Context().Value("id")

	if uId == nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, "user ID not found", nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	err := user.Update(uId.(string), req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "updated successfully", req)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// VerifyMe - /users/verify - PATCH
func (base *Controller) VerifyMe(w http.ResponseWriter, r *http.Request) {
	uId := r.Context().Value("id")

	if uId == nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, "user ID not found", nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	usr, err := user.Verify(uId.(string))
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "verified user successfully", usr)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// ForgotPassword - /users/forgot-password - POST
func (base *Controller) ForgotPassword(w http.ResponseWriter, r *http.Request) {

	req := new(model.ForgotPassword)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	err := user.ForgotPassword(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "email sent successfully", "")
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// ResetPassword - /users/reset-password - PATCH
func (base *Controller) ResetPassword(w http.ResponseWriter, r *http.Request) {

	req := new(model.ResetPassword)
	uId := r.Context().Value("id")

	if uId == nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, "user ID not found", nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	usr, err := user.ResetPassword(uId.(string), req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "reset password successfully", usr)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// ADMIN ENDPOINTS

// GetAll - /users/get-all - GET
func (base *Controller) GetAll(w http.ResponseWriter, r *http.Request) {

	usrs, err := user.GetAll()
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", usrs)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(res)
}

// GetUserByIdOrEmail - /users/:idOrEmail - GET
func (base *Controller) GetUserByIdOrEmail(w http.ResponseWriter, r *http.Request) {
	idOrEmail := chi.URLParam(r, "idOrEmail")
	usr, err := user.Get(idOrEmail)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", usr)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// LockUser - /users/lock/:idOrEmail - PATCH
func (base *Controller) LockUser(w http.ResponseWriter, r *http.Request) {
	idOrEmail := chi.URLParam(r, "idOrEmail")
	user, err := user.LockUser(idOrEmail)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "user locked successfully", user)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// UnlockUser - /users/unlock/:idOrEmail - PATCH
func (base *Controller) UnlockUser(w http.ResponseWriter, r *http.Request) {
	idOrEmail := chi.URLParam(r, "idOrEmail")
	user, err := user.UnlockUser(idOrEmail)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "user unlocked successfully", user)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
