package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	iamV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/iam/v1"
)

// ToProtoSession конвертирует доменную модель Session в protobuf
func ToProtoSession(session *model.Session) *iamV1.Session {
	if session == nil {
		return nil
	}

	return &iamV1.Session{
		Uuid:      session.UUID.String(),
		CreatedAt: timestamppb.New(session.CreatedAt),
		UpdatedAt: timestamppb.New(session.UpdatedAt),
		ExpiresAt: timestamppb.New(session.ExpiresAt),
	}
}
