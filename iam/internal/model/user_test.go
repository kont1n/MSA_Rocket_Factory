package model

import (
	"testing"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Валидный пароль",
			password: "StrongP@ss123!",
			wantErr:  false,
		},
		{
			name:     "Слишком короткий пароль",
			password: "Test1!",
			wantErr:  true,
		},
		{
			name:     "Нет заглавных букв",
			password: "strongp@ss123!",
			wantErr:  true,
		},
		{
			name:     "Нет строчных букв",
			password: "STRONGP@SS123!",
			wantErr:  true,
		},
		{
			name:     "Нет цифр",
			password: "StrongP@ssword!",
			wantErr:  true,
		},
		{
			name:     "Нет специальных символов",
			password: "StrongPassword123",
			wantErr:  true,
		},
		{
			name:     "Слишком длинный пароль",
			password: string(make([]byte, 130)),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "Валидный email",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "Email с поддоменом",
			email:   "user@mail.example.com",
			wantErr: false,
		},
		{
			name:    "Невалидный email без @",
			email:   "testexample.com",
			wantErr: true,
		},
		{
			name:    "Невалидный email без домена",
			email:   "test@",
			wantErr: true,
		},
		{
			name:    "Пустой email",
			email:   "",
			wantErr: true,
		},
		{
			name:    "Слишком длинный email",
			email:   string(make([]byte, 260)) + "@example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateLogin(t *testing.T) {
	tests := []struct {
		name    string
		login   string
		wantErr bool
	}{
		{
			name:    "Валидный логин",
			login:   "user123",
			wantErr: false,
		},
		{
			name:    "Логин с подчеркиванием",
			login:   "user_name",
			wantErr: false,
		},
		{
			name:    "Слишком короткий логин",
			login:   "ab",
			wantErr: true,
		},
		{
			name:    "Логин со спецсимволами",
			login:   "user@name",
			wantErr: true,
		},
		{
			name:    "Логин с пробелами",
			login:   "user name",
			wantErr: true,
		},
		{
			name:    "Пустой логин",
			login:   "",
			wantErr: true,
		},
		{
			name:    "Слишком длинный логин",
			login:   string(make([]byte, 55)),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLogin(tt.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRegistrationInfo_Validate(t *testing.T) {
	tests := []struct {
		name    string
		info    *UserRegistrationInfo
		wantErr bool
	}{
		{
			name: "Валидные данные",
			info: &UserRegistrationInfo{
				Login:    "testuser",
				Email:    "test@example.com",
				Password: "StrongP@ss123!",
			},
			wantErr: false,
		},
		{
			name: "Невалидный логин",
			info: &UserRegistrationInfo{
				Login:    "te", // слишком короткий
				Email:    "test@example.com",
				Password: "StrongP@ss123!",
			},
			wantErr: true,
		},
		{
			name: "Невалидный email",
			info: &UserRegistrationInfo{
				Login:    "testuser",
				Email:    "invalid-email",
				Password: "StrongP@ss123!",
			},
			wantErr: true,
		},
		{
			name: "Невалидный пароль",
			info: &UserRegistrationInfo{
				Login:    "testuser",
				Email:    "test@example.com",
				Password: "weak", // слишком слабый
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.info.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRegistrationInfo.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
