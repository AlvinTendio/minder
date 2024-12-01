package usecase

import (
	"context"

	"github.com/AlvinTendio/minder/common"
	minder_model "github.com/AlvinTendio/minder/minder/model"
)

type MinderUsecase interface {
	Register(ctx context.Context, req *minder_model.RegisterReq) (res *common.HTTPResponse, err error)
	Login(ctx context.Context, req *minder_model.LoginReq) (res *common.HTTPResponse, err error)
	UpgradeAccount(ctx context.Context, id uint64) (res *common.HTTPResponse, err error)
	GetTargetUser(ctx context.Context, id uint64) (res *common.HTTPResponse, err error)
	Swipe(ctx context.Context, req *minder_model.SwipeReq) (res *common.HTTPResponse, err error)
}
