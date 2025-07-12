package v1

import (
	"context"
	"net/http"

	orderV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/openapi/order/v1"
)

func (a *api) NewError(ctx context.Context, err error) *orderV1.GenericErrorStatusCode {
	code := orderV1.OptInt{}
	code.SetTo(http.StatusInternalServerError)

	message := orderV1.OptString{}
	message.SetTo(err.Error())

	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    code,
			Message: message,
		},
	}
}
