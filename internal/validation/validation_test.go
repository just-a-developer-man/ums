package validation

import (
	"reflect"
	"strings"
	"testing"
	"ums/internal/dto"
	"ums/internal/mocks/validatormocks"

	govalidator "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestNewValidator(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Create validator",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewValidator()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewValidator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestValidator_ValidateUID(t *testing.T) {
	testValidator, err := NewValidator()
	if err != nil {
		t.Fatalf("Failed to create validator: %s", err.Error())
	}
	type fields struct {
		validate *Validator
	}
	type args struct {
		uid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// Корректные UUID
		{
			name: "Valid UUID Version 5",
			fields: fields{
				validate: testValidator,
			},
			args: args{
				uid: "6ba7b810-9dad-51d1-80b4-00c04fd430c8",
			},
			wantErr: false,
		},

		// invalid UUID and corner cases
		{
			name: "Valid UUID Version 1",
			fields: fields{
				validate: testValidator,
			},
			args: args{
				uid: "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			},
			wantErr: true,
		},
		{
			name: "Valid UUID Version 2",
			fields: fields{
				validate: testValidator,
			},
			args: args{
				uid: "6ba7b811-9dad-21d1-80b4-00c04fd430c8",
			},
			wantErr: true,
		},
		{
			name: "Valid UUID Version 3",
			fields: fields{
				validate: testValidator,
			},
			args: args{
				uid: "6ba7b810-9dad-31d1-80b4-00c04fd430c8",
			},
			wantErr: true,
		},
		{
			name: "Valid UUID Version 4",
			fields: fields{
				validate: testValidator,
			},
			args: args{
				uid: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			},
			wantErr: true,
		},

		{
			name: "Invalid UUID - Too Short",
			fields: fields{
				validate: testValidator,
			},
			args: args{
				uid: "6ba7b810-9dad-11d1-80b4",
			},
			wantErr: true,
		},
		{
			name: "Invalid UUID - Too Long",
			fields: fields{
				validate: testValidator,
			},
			args: args{
				uid: "6ba7b810-9dad-11d1-80b4-00c04fd430c8-extra",
			},
			wantErr: true,
		},
		{
			name: "Invalid UUID - Missing Dashes",
			fields: fields{
				validate: testValidator,
			},
			args: args{
				uid: "6ba7b8109dad11d180b400c04fd430c8",
			},
			wantErr: true,
		},
		{
			name: "Invalid UUID - Invalid Characters",
			fields: fields{
				validate: testValidator,
			},
			args: args{
				uid: "6ba7b810-9dad-11d1-80b4-00c04fd430cg",
			},
			wantErr: true,
		},
		{
			name: "Invalid UUID - Empty String",
			fields: fields{
				validate: testValidator,
			},
			args: args{
				uid: "",
			},
			wantErr: true,
		},
		{
			name: "Invalid UUID - Null Value",
			fields: fields{
				validate: testValidator,
			},
			args: args{
				uid: "\x00\x00\x00\x00-\x00\x00-\x00\x00-\x00\x00-\x00\x00\x00\x00\x00\x00",
			},
			wantErr: true,
		},
		{
			name: "Invalid UUID - Incorrect Variant",
			fields: fields{
				validate: testValidator,
			},
			args: args{
				uid: "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			},
			wantErr: true,
		},
		{
			name: "Invalid UUID - Incorrect Version",
			fields: fields{
				validate: testValidator,
			},
			args: args{
				uid: "6ba7b810-9dad-61d1-80b4-00c04fd430c8",
			},
			wantErr: true,
		},
		{
			name: "Invalid UUID - Random String",
			fields: fields{
				validate: testValidator,
			},
			args: args{
				uid: "random-string-not-a-uuid",
			},
			wantErr: true,
		},
		{
			name: "Invalid UUID - Only Dashes",
			fields: fields{
				validate: testValidator,
			},
			args: args{
				uid: "----------------",
			},
			wantErr: true,
		},
		{
			name: "Invalid UUID - Mixed Case",
			fields: fields{
				validate: testValidator,
			},
			args: args{
				uid: "6BA7B810-9DAD-11D1-80b4-00C04fD430C8",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.fields.validate
			if err := v.ValidateUID(tt.args.uid); (err != nil) != tt.wantErr {
				t.Errorf("Validator.ValidateUID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validatePassword(t *testing.T) {
	type args struct {
		fl govalidator.FieldLevel
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// Valid passwords
		{
			name: "Valid password with all requirements",
			args: args{
				fl: func() govalidator.FieldLevel {
					mock := validatormocks.NewMockFieldLevel(t)
					mock.EXPECT().Field().Return(reflect.ValueOf("A1b2C3!"))
					return mock
				}(),
			},
			want: true,
		},
		{
			name: "Valid password with minimum length",
			args: args{
				fl: func() govalidator.FieldLevel {
					mock := validatormocks.NewMockFieldLevel(t)
					mock.EXPECT().Field().Return(reflect.ValueOf("A1b!"))
					return mock
				}(),
			},
			want: false, // Updated: Minimum length is now 6
		},
		{
			name: "Valid password with maximum length",
			args: args{
				fl: func() govalidator.FieldLevel {
					mock := validatormocks.NewMockFieldLevel(t)
					mock.EXPECT().Field().Return(reflect.ValueOf(strings.Join([]string{strings.Repeat("A", 60), "A1b!"}, "")))
					return mock
				}(),
			},
			want: true,
		},

		// Invalid passwords
		{
			name: "Password too short",
			args: args{
				fl: func() govalidator.FieldLevel {
					mock := validatormocks.NewMockFieldLevel(t)
					mock.EXPECT().Field().Return(reflect.ValueOf("A1b!"))
					return mock
				}(),
			},
			want: false,
		},
		{
			name: "Password too long",
			args: args{
				fl: func() govalidator.FieldLevel {
					mock := validatormocks.NewMockFieldLevel(t)
					mock.EXPECT().Field().Return(reflect.ValueOf(strings.Repeat("A", 65)))
					return mock
				}(),
			},
			want: false,
		},
		{
			name: "Password with spaces",
			args: args{
				fl: func() govalidator.FieldLevel {
					mock := validatormocks.NewMockFieldLevel(t)
					mock.EXPECT().Field().Return(reflect.ValueOf("A1 b2C3!"))
					return mock
				}(),
			},
			want: false,
		},
		{
			name: "Password without special characters",
			args: args{
				fl: func() govalidator.FieldLevel {
					mock := validatormocks.NewMockFieldLevel(t)
					mock.EXPECT().Field().Return(reflect.ValueOf("A1b2C3"))
					return mock
				}(),
			},
			want: false,
		},
		{
			name: "Password without uppercase letters",
			args: args{
				fl: func() govalidator.FieldLevel {
					mock := validatormocks.NewMockFieldLevel(t)
					mock.EXPECT().Field().Return(reflect.ValueOf("a1b2c3!"))
					return mock
				}(),
			},
			want: false,
		},
		{
			name: "Password without lowercase letters",
			args: args{
				fl: func() govalidator.FieldLevel {
					mock := validatormocks.NewMockFieldLevel(t)
					mock.EXPECT().Field().Return(reflect.ValueOf("A1B2C3!"))
					return mock
				}(),
			},
			want: false,
		},
		{
			name: "Password without digits",
			args: args{
				fl: func() govalidator.FieldLevel {
					mock := validatormocks.NewMockFieldLevel(t)
					mock.EXPECT().Field().Return(reflect.ValueOf("Abcdef!"))
					return mock
				}(),
			},
			want: false,
		},
		{
			name: "Password with non-printable or non-ASCII characters",
			args: args{
				fl: func() govalidator.FieldLevel {
					mock := validatormocks.NewMockFieldLevel(t)
					mock.EXPECT().Field().Return(reflect.ValueOf("A1b2C3!\x00"))
					return mock
				}(),
			},
			want: false,
		},

		// Edge cases
		{
			name: "Empty password",
			args: args{
				fl: func() govalidator.FieldLevel {
					mock := validatormocks.NewMockFieldLevel(t)
					mock.EXPECT().Field().Return(reflect.ValueOf(""))
					return mock
				}(),
			},
			want: false,
		},
		{
			name: "Password with only special characters",
			args: args{
				fl: func() govalidator.FieldLevel {
					mock := validatormocks.NewMockFieldLevel(t)
					mock.EXPECT().Field().Return(reflect.ValueOf("!@#$%^&*"))
					return mock
				}(),
			},
			want: false,
		},
		{
			name: "Password with only digits",
			args: args{
				fl: func() govalidator.FieldLevel {
					mock := validatormocks.NewMockFieldLevel(t)
					mock.EXPECT().Field().Return(reflect.ValueOf("12345678"))
					return mock
				}(),
			},
			want: false,
		},
		{
			name: "Password with only uppercase letters",
			args: args{
				fl: func() govalidator.FieldLevel {
					mock := validatormocks.NewMockFieldLevel(t)
					mock.EXPECT().Field().Return(reflect.ValueOf("ABCDEFGH"))
					return mock
				}(),
			},
			want: false,
		},
		{
			name: "Password with only lowercase letters",
			args: args{
				fl: func() govalidator.FieldLevel {
					mock := validatormocks.NewMockFieldLevel(t)
					mock.EXPECT().Field().Return(reflect.ValueOf("abcdefgh"))
					return mock
				}(),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, validatePassword(tt.args.fl))
		})
	}
}

func TestValidator_ValidateStruct(t *testing.T) {
	validate, err := NewValidator()
	if err != nil {
		t.Fatalf("Cannot create validator: %s", err.Error())
	}
	type fields struct {
		validate *Validator
	}
	type args struct {
		data *dto.CreateUserWithRoleRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// Положительный случай: все поля заполнены корректно
		{
			name: "Valid data",
			fields: fields{
				validate: validate,
			},
			args: args{
				data: &dto.CreateUserWithRoleRequest{
					Name:     "valid-user-name",
					Email:    "test@example.com",
					Password: "StrongPass1!",
					Role:     "user",
				},
			},
			wantErr: false,
		},
		// Отрицательные случаи для Name
		{
			name: "Empty Name",
			fields: fields{
				validate: validate,
			},
			args: args{
				data: &dto.CreateUserWithRoleRequest{
					Name:     "",
					Email:    "test@example.com",
					Password: "StrongPass1!",
					Role:     "user",
				},
			},
			wantErr: true,
		},
		{
			name: "Big Name",
			fields: fields{
				validate: validate,
			},
			args: args{
				data: &dto.CreateUserWithRoleRequest{
					Name:     strings.Repeat("A", 65),
					Email:    "test@example.com",
					Password: "StrongPass1!",
					Role:     "user",
				},
			},
			wantErr: true,
		},
		{
			name: "Small Name",
			fields: fields{
				validate: validate,
			},
			args: args{
				data: &dto.CreateUserWithRoleRequest{
					Name:     "a",
					Email:    "test@example.com",
					Password: "StrongPass1!",
					Role:     "user",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid Name format",
			fields: fields{
				validate: validate,
			},
			args: args{
				data: &dto.CreateUserWithRoleRequest{
					Name:     "invalid name!",
					Email:    "test@example.com",
					Password: "StrongPass1!",
					Role:     "user",
				},
			},
			wantErr: true,
		},
		// Отрицательные случаи для Email
		{
			name: "Empty Email",
			fields: fields{
				validate: validate,
			},
			args: args{
				data: &dto.CreateUserWithRoleRequest{
					Name:     "valid-user-name",
					Email:    "",
					Password: "StrongPass1!",
					Role:     "user",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid Email format",
			fields: fields{
				validate: validate,
			},
			args: args{
				data: &dto.CreateUserWithRoleRequest{
					Name:     "valid-user-name",
					Email:    "invalid-email",
					Password: "StrongPass1!",
					Role:     "user",
				},
			},
			wantErr: true,
		},
		// Отрицательные случаи для Password
		{
			name: "Empty Password",
			fields: fields{
				validate: validate,
			},
			args: args{
				data: &dto.CreateUserWithRoleRequest{
					Name:     "valid-user-name",
					Email:    "test@example.com",
					Password: "",
					Role:     "user",
				},
			},
			wantErr: true,
		},
		{
			name: "Password too short",
			fields: fields{
				validate: validate,
			},
			args: args{
				data: &dto.CreateUserWithRoleRequest{
					Name:     "valid-user-name",
					Email:    "test@example.com",
					Password: "12345", // Меньше 6 символов
					Role:     "user",
				},
			},
			wantErr: true,
		},
		{
			name: "Password too long",
			fields: fields{
				validate: validate,
			},
			args: args{
				data: &dto.CreateUserWithRoleRequest{
					Name:     "valid-user-name",
					Email:    "test@example.com",
					Password: strings.Repeat("a", 65), // Больше 64 символов
					Role:     "user",
				},
			},
			wantErr: true,
		},
		// Отрицательные случаи для Role
		{
			name: "Empty Role",
			fields: fields{
				validate: validate,
			},
			args: args{
				data: &dto.CreateUserWithRoleRequest{
					Name:     "valid-user-name",
					Email:    "test@example.com",
					Password: "StrongPass1!",
					Role:     "",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid Role",
			fields: fields{
				validate: validate,
			},
			args: args{
				data: &dto.CreateUserWithRoleRequest{
					Name:     "valid-user-name",
					Email:    "test@example.com",
					Password: "StrongPass1!",
					Role:     "guest", // Недопустимая роль
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.fields.validate
			if err := v.ValidateStruct(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Validator.ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateUserName(t *testing.T) {
	type args struct {
		fl govalidator.FieldLevel
	}
	// Mock FieldLevel implementation
	mockFieldLevel := func(username string) govalidator.FieldLevel {
		mock := validatormocks.NewMockFieldLevel(t)
		mock.EXPECT().Field().Return(reflect.ValueOf(username))
		return mock
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid username with minimum length",
			args: args{
				fl: mockFieldLevel("abc"),
			},
			want: true,
		},
		{
			name: "Valid username with maximum length",
			args: args{
				fl: mockFieldLevel(strings.Repeat("a", 64)),
			},
			want: true,
		},
		{
			name: "Username too short",
			args: args{
				fl: mockFieldLevel("ab"),
			},
			want: false,
		},
		{
			name: "Username too long",
			args: args{
				fl: mockFieldLevel(strings.Repeat("a", 65)),
			},
			want: false,
		},
		{
			name: "Username with invalid characters",
			args: args{
				fl: mockFieldLevel("user@name"),
			},
			want: false,
		},
		{
			name: "Empty username",
			args: args{
				fl: mockFieldLevel(""),
			},
			want: false,
		},
		{
			name: "Username with special characters allowed by regex",
			args: args{
				fl: mockFieldLevel("user-name_123"),
			},
			want: true,
		},
		{
			name: "Username with only numbers",
			args: args{
				fl: mockFieldLevel("12345"),
			},
			want: true,
		},
		{
			name: "Username starting with a number",
			args: args{
				fl: mockFieldLevel("1username"),
			},
			want: true,
		},
		{
			name: "Username with leading space",
			args: args{
				fl: mockFieldLevel(" username"),
			},
			want: false,
		},
		{
			name: "Username with trailing space",
			args: args{
				fl: mockFieldLevel("username "),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateUserName(tt.args.fl); got != tt.want {
				t.Errorf("validateUserName() = %v, want %v", got, tt.want)
			}
		})
	}
}
