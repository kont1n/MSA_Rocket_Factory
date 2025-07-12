package v1

import (
	def "github.com/kont1n/MSA_Rocket_Factory/order/internal/client"
	generaredPaymentyV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/payment/v1"
)

var _ def.PaymentClient = (*client)(nil)

type client struct {
	generatedClient generaredPaymentyV1.PaymentServiceClient
}

func NewClient(generatedClient generaredPaymentyV1.PaymentServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
