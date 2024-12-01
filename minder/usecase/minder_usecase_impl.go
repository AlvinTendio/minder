package usecase

import (
	"context"
	"log"
	"net/http"

	"github.com/AlvinTendio/minder/common"
	minder_model "github.com/AlvinTendio/minder/minder/model"
	"github.com/AlvinTendio/minder/minder/repository"
)

type minderUsecaseImpl struct {
	MinderRepo repository.MinderRepository
}

func NewMinderUsecaseImpl(minderRepo repository.MinderRepository) MinderUsecase {
	return &minderUsecaseImpl{
		MinderRepo: minderRepo,
	}
}

func (u *minderUsecaseImpl) Register(ctx context.Context, req *minder_model.RegisterReq) (res *common.HTTPResponse, err error) {
	data, err := u.MinderRepo.Register(ctx, req)

	if err != nil || data == 0 {
		log.Println(ctx, "Error ", err)
		res = &common.HTTPResponse{
			HTTPStatus:      http.StatusInternalServerError,
			ResponseCode:    common.StatusInternalServerErrorResponseCode,
			ResponseMessage: common.StatusInternalServerErrorResponseMessage,
		}
		return
	}

	res = &common.HTTPResponse{
		HTTPStatus:      http.StatusOK,
		ResponseCode:    common.StatusOKResponseCode,
		ResponseMessage: common.StatusOKResponseMessage,
	}

	return
}

func (u *minderUsecaseImpl) Login(ctx context.Context, req *minder_model.LoginReq) (res *common.HTTPResponse, err error) {
	data, err := u.MinderRepo.Login(ctx, req)

	if err != nil || data == nil {
		log.Println(ctx, "Error ", err)
		res = &common.HTTPResponse{
			HTTPStatus:      http.StatusInternalServerError,
			ResponseCode:    common.StatusInternalServerErrorResponseCode,
			ResponseMessage: common.StatusInternalServerErrorResponseMessage,
		}
		return
	}

	res = &common.HTTPResponse{
		HTTPStatus:      http.StatusOK,
		ResponseCode:    common.StatusOKResponseCode,
		ResponseMessage: common.StatusOKResponseMessage,
		Data:            data,
	}

	return
}
func (u *minderUsecaseImpl) UpgradeAccount(ctx context.Context, id uint64) (res *common.HTTPResponse, err error) {
	data, err := u.MinderRepo.UpgradeAccount(ctx, id)

	if err != nil || data == 0 {
		log.Println(ctx, "Error ", err)
		res = &common.HTTPResponse{
			HTTPStatus:      http.StatusInternalServerError,
			ResponseCode:    common.StatusInternalServerErrorResponseCode,
			ResponseMessage: common.StatusInternalServerErrorResponseMessage,
		}
		return
	}

	res = &common.HTTPResponse{
		HTTPStatus:      http.StatusOK,
		ResponseCode:    common.StatusOKResponseCode,
		ResponseMessage: common.StatusOKResponseMessage,
	}

	return
}
func (u *minderUsecaseImpl) GetTargetUser(ctx context.Context, id uint64) (res *common.HTTPResponse, err error) {
	upgradeStatus, err := u.MinderRepo.GetUserUpgradeStatus(ctx, id)

	if err != nil {
		log.Println(ctx, "Error ", err)
		res = &common.HTTPResponse{
			HTTPStatus:      http.StatusInternalServerError,
			ResponseCode:    common.StatusInternalServerErrorResponseCode,
			ResponseMessage: common.StatusInternalServerErrorResponseMessage,
		}
		return
	}

	if !upgradeStatus {
		viewCount, err := u.MinderRepo.GetUserViewCount(ctx, id)

		if err != nil || viewCount == nil || *viewCount >= 10 {
			log.Println(ctx, "Error ", err)
			res = &common.HTTPResponse{
				HTTPStatus:      http.StatusInternalServerError,
				ResponseCode:    common.StatusInternalServerErrorResponseCode,
				ResponseMessage: common.StatusInternalServerErrorResponseMessage,
			}
			return res, err
		}
	}

	data, err := u.MinderRepo.GetTargetUser(ctx, id)

	if err != nil || data == nil {
		log.Println(ctx, "Error ", err)
		res = &common.HTTPResponse{
			HTTPStatus:      http.StatusInternalServerError,
			ResponseCode:    common.StatusInternalServerErrorResponseCode,
			ResponseMessage: common.StatusInternalServerErrorResponseMessage,
		}
		return
	}

	total, err := u.MinderRepo.InsertSwipe(ctx, id, data.UserId)

	if err != nil || total == 0 {
		log.Println(ctx, "Error ", err)
		res = &common.HTTPResponse{
			HTTPStatus:      http.StatusInternalServerError,
			ResponseCode:    common.StatusInternalServerErrorResponseCode,
			ResponseMessage: common.StatusInternalServerErrorResponseMessage,
		}
		return
	}

	res = &common.HTTPResponse{
		HTTPStatus:      http.StatusOK,
		ResponseCode:    common.StatusOKResponseCode,
		ResponseMessage: common.StatusOKResponseMessage,
		Data:            data,
	}
	return
}
func (u *minderUsecaseImpl) Swipe(ctx context.Context, req *minder_model.SwipeReq) (res *common.HTTPResponse, err error) {
	data, err := u.MinderRepo.UpdateSwipe(ctx, req)

	if err != nil || data == 0 {
		log.Println(ctx, "Error ", err)
		res = &common.HTTPResponse{
			HTTPStatus:      http.StatusInternalServerError,
			ResponseCode:    common.StatusInternalServerErrorResponseCode,
			ResponseMessage: common.StatusInternalServerErrorResponseMessage,
		}
		return
	}

	res = &common.HTTPResponse{
		HTTPStatus:      http.StatusOK,
		ResponseCode:    common.StatusOKResponseCode,
		ResponseMessage: common.StatusOKResponseMessage,
	}

	return
}
