package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/AlvinTendio/minder/common"
	common_http "github.com/AlvinTendio/minder/delivery/http"
	minder_model "github.com/AlvinTendio/minder/minder/model"
	"github.com/AlvinTendio/minder/minder/usecase"
	validatorfmt "github.com/AlvinTendio/minder/validator-fmt"
	validator "github.com/go-playground/validator/v10"
)

type MinderHandler struct {
	MinderUsecase usecase.MinderUsecase
}

func NewMinderHandler(minderUsecase usecase.MinderUsecase) {
	h := &MinderHandler{
		MinderUsecase: minderUsecase,
	}

	common_http.Route(http.MethodPost, "/register", h.Register, "Register")
	common_http.Route(http.MethodPost, "/login", h.Login, "Login")
	common_http.Route(http.MethodPut, "/upgrade-account/([0-9]+)", h.UpgradeAccount, "Upgrade Account")
	common_http.Route(http.MethodGet, "/get-target-user/([0-9]+)", h.GetTargetUser, "GetTargetUser")
	common_http.Route(http.MethodPut, "/swipe", h.Swipe, "Swipe")
}

func (h *MinderHandler) Register(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	registerReq := &minder_model.RegisterReq{}
	err := json.NewDecoder(req.Body).Decode(registerReq)
	if err != nil {
		log.Println("Error in POST parameters : ", err)
	}

	validate := validator.New()
	validate.RegisterValidation("dateformat", validatorfmt.DateFormatValidator)
	err = validate.Struct(registerReq)
	if err != nil {
		resp := &common.HTTPResponse{
			HTTPStatus:      http.StatusBadRequest,
			ResponseCode:    common.StatusBadRequestErrorResponseCode,
			ResponseMessage: common.StatusBadRequestErrorResponseMessage,
		}
		common_http.ResponseWrite(req, rw, resp, resp.HTTPStatus)
		return
	}

	result, err := h.MinderUsecase.Register(ctx, registerReq)
	if err != nil {
		log.Println(ctx, "[delivery:http:handler] : Exception Register", err)
		common_http.ResponseWrite(req, rw, result, http.StatusInternalServerError)
		return
	}
	common_http.ResponseWrite(req, rw, result, result.HTTPStatus)
}

func (h *MinderHandler) Login(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	loginReq := &minder_model.LoginReq{}
	err := json.NewDecoder(req.Body).Decode(loginReq)
	if err != nil {
		log.Println("Error in POST parameters : ", err)
	}

	validate := validator.New()
	err = validate.Struct(loginReq)
	if err != nil {
		resp := &common.HTTPResponse{
			HTTPStatus:      http.StatusBadRequest,
			ResponseCode:    common.StatusBadRequestErrorResponseCode,
			ResponseMessage: common.StatusBadRequestErrorResponseMessage,
		}
		common_http.ResponseWrite(req, rw, resp, resp.HTTPStatus)
		return
	}

	result, err := h.MinderUsecase.Login(ctx, loginReq)
	if err != nil {
		log.Println(ctx, "[delivery:http:handler] : Exception Login", err)
		common_http.ResponseWrite(req, rw, result, http.StatusInternalServerError)
		return
	}
	common_http.ResponseWrite(req, rw, result, result.HTTPStatus)
}

func (h *MinderHandler) UpgradeAccount(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id := getParamUint64(rw, req)

	if id == 0 {
		resp := &common.HTTPResponse{
			HTTPStatus:      http.StatusBadRequest,
			ResponseCode:    common.StatusBadRequestErrorResponseCode,
			ResponseMessage: common.StatusBadRequestErrorResponseMessage,
		}
		common_http.ResponseWrite(req, rw, resp, resp.HTTPStatus)
		return
	}

	result, err := h.MinderUsecase.UpgradeAccount(ctx, id)
	if err != nil {
		log.Println(ctx, "[delivery:http:handler] : Exception Upgrade Account", err)
		common_http.ResponseWrite(req, rw, result, http.StatusInternalServerError)
		return
	}
	common_http.ResponseWrite(req, rw, result, result.HTTPStatus)
}

func (h *MinderHandler) GetTargetUser(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id := getParamUint64(rw, req)

	if id == 0 {
		resp := &common.HTTPResponse{
			HTTPStatus:      http.StatusBadRequest,
			ResponseCode:    common.StatusBadRequestErrorResponseCode,
			ResponseMessage: common.StatusBadRequestErrorResponseMessage,
		}
		common_http.ResponseWrite(req, rw, resp, resp.HTTPStatus)
		return
	}

	result, err := h.MinderUsecase.GetTargetUser(ctx, id)
	if err != nil {
		log.Println(ctx, "[delivery:http:handler] : Exception Get Target User", err)
		common_http.ResponseWrite(req, rw, result, http.StatusInternalServerError)
		return
	}
	common_http.ResponseWrite(req, rw, result, result.HTTPStatus)
}

func (h *MinderHandler) Swipe(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	swipeReq := &minder_model.SwipeReq{}
	err := json.NewDecoder(req.Body).Decode(swipeReq)
	if err != nil {
		log.Println("Error in POST parameters : ", err)
	}

	validate := validator.New()
	err = validate.Struct(swipeReq)
	if err != nil {
		resp := &common.HTTPResponse{
			HTTPStatus:      http.StatusBadRequest,
			ResponseCode:    common.StatusBadRequestErrorResponseCode,
			ResponseMessage: common.StatusBadRequestErrorResponseMessage,
		}
		common_http.ResponseWrite(req, rw, resp, resp.HTTPStatus)
		return
	}

	result, err := h.MinderUsecase.Swipe(ctx, swipeReq)
	if err != nil {
		log.Println(ctx, "[delivery:http:handler] : Exception Register", err)
		common_http.ResponseWrite(req, rw, result, http.StatusInternalServerError)
		return
	}
	common_http.ResponseWrite(req, rw, result, result.HTTPStatus)
}

func getParamUint64(rw http.ResponseWriter, req *http.Request) uint64 {
	id, err := strconv.ParseUint(common_http.Param(req, 0), 0, 64)
	if err != nil {
		log.Println("err -> ", err)
		return 0
	}
	return id
}
