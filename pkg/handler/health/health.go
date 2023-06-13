package health

import (
	"encoding/json"
	"fmt"
	"net/http"

	"brief/internal/model"
	"brief/service/ping"
	"brief/utility"

	"github.com/go-playground/validator/v10"
)

type Controller struct {
	Validate *validator.Validate
	Logger   *utility.Logger
}

func (base *Controller) Post(w http.ResponseWriter, r *http.Request) {
	var req model.Ping

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	if err := base.Validate.Struct(&req); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validate), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	if !ping.ReturnTrue() {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "ping failed", fmt.Errorf("ping failed"), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	base.Logger.Info("ping successfull")

	rd := utility.BuildSuccessResponse(http.StatusOK, "ping successfull", req.Message)
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(res)
}

func (base *Controller) Get(w http.ResponseWriter, r *http.Request) {
	if !ping.ReturnTrue() {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "ping failed", fmt.Errorf("ping failed"), nil)
		res, _ := json.Marshal(rd)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	base.Logger.Info("ping successfull")
	rd := utility.BuildSuccessResponse(http.StatusOK, "ping successfull", map[string]interface{}{"user": "user object"})
	res, _ := json.Marshal(rd)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(res)
	return
}
