package ssoservice

import (
	"context"
	"log/slog"

	ssomodels "github.com/DimTur/lp_api_gateway/internal/clients/sso/models.go"
	"github.com/go-playground/validator/v10"
)

type AuthServiceProvider interface {
	RegisterUser(ctx context.Context, newUser *ssomodels.RegisterUser) (*ssomodels.RegisterResp, error)
	LoginUser(ctx context.Context, logUser *ssomodels.LogIn) (*ssomodels.LogInResp, error)
	LogInViaTg(ctx context.Context, email *ssomodels.LogInViaTg) (*ssomodels.LogInViaTgResp, error)
	CheckOTPAndLogIn(ctx context.Context, otp *ssomodels.CheckOTPAndLogIn) (*ssomodels.CheckOTPAndLogInResp, error)
	UpdateUserInfo(ctx context.Context, newInfo *ssomodels.UpdateUserInfo) (*ssomodels.UpdateUserInfoResp, error)
	AuthCheck(ctx context.Context, authCheck *ssomodels.AuthCheck) (*ssomodels.AuthCheckResp, error)
}

type LgServiceProvider interface {
	CreateLearningGroup(ctx context.Context, newLGroup *ssomodels.CreateLearningGroup) (*ssomodels.CreateLearningGroupResp, error)
	GetLearningGroupByID(ctx context.Context, lgID *ssomodels.GetLgByID) (*ssomodels.GetLgByIDResp, error)
	UpdateLearningGroup(ctx context.Context, updFields *ssomodels.UpdateLearningGroup) (*ssomodels.UpdateLearningGroupResp, error)
	DeleteLearningGroup(ctx context.Context, lgID *ssomodels.DelLgByID) (*ssomodels.DelLgByIDResp, error)
	GetLearningGroups(ctx context.Context, uID *ssomodels.GetLGroups) (*ssomodels.GetLGroupsResp, error)
	UserIsLearnerIn(ctx context.Context, user *ssomodels.UserIsLearnerIn) ([]string, error)
}

type SsoService struct {
	Log          *slog.Logger
	Validator    *validator.Validate
	AuthProvider AuthServiceProvider
	LgProvider   LgServiceProvider
}

func New(
	log *slog.Logger,
	validator *validator.Validate,
	authProvider AuthServiceProvider,
	lgProvider LgServiceProvider,
) *SsoService {
	return &SsoService{
		Log:          log,
		Validator:    validator,
		AuthProvider: authProvider,
		LgProvider:   lgProvider,
	}
}
