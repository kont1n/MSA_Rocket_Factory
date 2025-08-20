package app

import (
	"context"

	iamV1API "github.com/kont1n/MSA_Rocket_Factory/iam/internal/api/iam/v1"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/service"
	iamService "github.com/kont1n/MSA_Rocket_Factory/iam/internal/service/iam"
	iamV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/iam/v1"
)

type diContainer struct {
	iamAPIv1   iamV1.IamServiceServer
	iamService service.IamService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) IamV1API(ctx context.Context) iamV1.IamServiceServer {
	if d.iamAPIv1 == nil {
		d.iamAPIv1 = iamV1.NewAPI(d.IamService(ctx))
	}
	return d.iamAPIv1
}

func (d *diContainer) IamService(ctx context.Context) service.IamService {
	if d.iamService == nil {
		d.iamService = iamService.NewService()
	}
	return d.iamService
}
