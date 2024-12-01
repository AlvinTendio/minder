package repository

import (
	"context"

	minder_model "github.com/AlvinTendio/minder/minder/model"
)

type MinderRepository interface {
	Register(ctx context.Context, req *minder_model.RegisterReq) (data int64, err error)
	Login(ctx context.Context, req *minder_model.LoginReq) (data *minder_model.UserData, err error)
	UpgradeAccount(ctx context.Context, id uint64) (data int64, err error)
	GetUserUpgradeStatus(ctx context.Context, id uint64) (data bool, err error)
	GetUserViewCount(ctx context.Context, id uint64) (total *int64, err error)
	GetTargetUser(ctx context.Context, id uint64) (data *minder_model.TargetUserData, err error)
	InsertSwipe(ctx context.Context, id uint64, targetId int64) (data int64, err error)
	UpdateSwipe(ctx context.Context, req *minder_model.SwipeReq) (data int64, err error)
}
