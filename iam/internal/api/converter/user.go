package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	iamV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/iam/v1"
)

// ToProtoUser конвертирует доменную модель User в protobuf
func ToProtoUser(user *model.User) *iamV1.User {
	if user == nil {
		return nil
	}

	return &iamV1.User{
		Uuid:      user.UUID.String(),
		Info:      ToProtoUserInfo(user),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

// ToProtoUserInfo конвертирует информацию о пользователе в protobuf
func ToProtoUserInfo(user *model.User) *iamV1.UserInfo {
	if user == nil {
		return nil
	}

	return &iamV1.UserInfo{
		Login:               user.Login,
		Email:               user.Email,
		NotificationMethods: ToProtoNotificationMethods(user.NotificationMethods),
	}
}

// ToProtoNotificationMethods конвертирует массив способов уведомления в protobuf
func ToProtoNotificationMethods(methods []model.NotificationMethod) []*iamV1.NotificationMethod {
	if methods == nil {
		return nil
	}

	protoMethods := make([]*iamV1.NotificationMethod, len(methods))
	for i, method := range methods {
		protoMethods[i] = ToProtoNotificationMethod(&method)
	}

	return protoMethods
}

// ToProtoNotificationMethod конвертирует способ уведомления в protobuf
func ToProtoNotificationMethod(method *model.NotificationMethod) *iamV1.NotificationMethod {
	if method == nil {
		return nil
	}

	return &iamV1.NotificationMethod{
		ProviderName: method.ProviderName,
		Target:       method.Target,
	}
}

// ToModelNotificationMethods конвертирует массив способов уведомления из protobuf в доменную модель
func ToModelNotificationMethods(protoMethods []*iamV1.NotificationMethod) []model.NotificationMethod {
	if protoMethods == nil {
		return nil
	}

	methods := make([]model.NotificationMethod, len(protoMethods))
	for i, protoMethod := range protoMethods {
		methods[i] = ToModelNotificationMethod(protoMethod)
	}

	return methods
}

// ToModelNotificationMethod конвертирует способ уведомления из protobuf в доменную модель
func ToModelNotificationMethod(protoMethod *iamV1.NotificationMethod) model.NotificationMethod {
	return model.NotificationMethod{
		ProviderName: protoMethod.ProviderName,
		Target:       protoMethod.Target,
	}
}

// ToModelUserRegistrationInfo конвертирует данные регистрации из protobuf в доменную модель
func ToModelUserRegistrationInfo(req *iamV1.RegisterRequest) *model.UserRegistrationInfo {
	if req == nil || req.Info == nil || req.Info.UserInfo == nil {
		return nil
	}

	return &model.UserRegistrationInfo{
		Login:               req.Info.UserInfo.Login,
		Email:               req.Info.UserInfo.Email,
		Password:            req.Info.Password,
		NotificationMethods: ToModelNotificationMethods(req.Info.UserInfo.NotificationMethods),
	}
}
