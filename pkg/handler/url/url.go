package url

import (
	"brief/internal/constant"
	"brief/internal/model"
	"brief/service/url"
	"brief/utility"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Redirect - /api/v1/{hash} - GET
func (base *Controller) Redirect(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	url, err := url.Redirect(hash)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	http.Redirect(w, r, url.LongURL, http.StatusTemporaryRedirect)
}

// Shorten - /api/v1/url/shorten - POST
func (base *Controller) Shorten(w http.ResponseWriter, r *http.Request) {
	req := new(model.URL)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	// Fetch and set user ID if present
	uInfo := r.Context().Value(struct{}{})
	if uInfo != nil {
		req.UserID = uInfo.(*model.ContextInfo).ID
	}

	if err := base.Validate.Struct(req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrValidation, utility.ValidationResponse(err, base.Validate), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	if err := url.Shorten(req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "successfully created url", req)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

// GetUrls - /api/v1/url - GET
func (base *Controller) GetUrls(w http.ResponseWriter, r *http.Request) {
	uInfo := r.Context().Value(struct{}{}) // fetch user's info from context

	if uInfo == nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, "user ID not found", nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	uId := uInfo.(*model.ContextInfo).ID
	urls, err := url.GetURLs(uId)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", urls)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// Delete - /api/v1/url/{id} - DELETE
func (base *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	urlId := chi.URLParam(r, "id")
	uInfo := r.Context().Value(struct{}{}) // fetch user's info from context

	if uInfo == nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrBinding, "user ID not found", nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	uContextInfo := uInfo.(*model.ContextInfo)
	url, err := url.Delete(uContextInfo, urlId)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successfully deleted url", url)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// ADMIN ENDPOINTS

// GetAll - /api/v1/url/get-all - GET
func (base *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	urls, err := url.GetAll()
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", urls)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(res)
}

// GetUrlsByUserID - /api/v1/url/{user-id} - GET
func (base *Controller) GetUrlsByUserID(w http.ResponseWriter, r *http.Request) {
	uID := chi.URLParam(r, "user-id")
	urls, err := url.GetURLs(uID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed,
			constant.ErrRequest, err.Error(), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", urls)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
