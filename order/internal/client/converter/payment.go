package converter

import (
	generaredPaymentV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/payment/v1"
)

func ToProtoPaymentMethod(paymentMethod string) generaredPaymentV1.PaymentMethod {
	intPaymentMethod, ok := generaredPaymentV1.PaymentMethod_value["PAYMENT_METHOD_"+paymentMethod]
	if !ok {
		intPaymentMethod = int32(generaredPaymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED)
	}
	return generaredPaymentV1.PaymentMethod(intPaymentMethod)
}
