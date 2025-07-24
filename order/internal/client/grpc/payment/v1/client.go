package v1

import (
	generaredPaymentV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/payment/v1"
)

type paymentClient struct {
	generatedClient generaredPaymentV1.PaymentServiceClient
}

func NewClient(generatedClient generaredPaymentV1.PaymentServiceClient) *paymentClient {
	return &paymentClient{
		generatedClient: generatedClient,
	}
}
