package entity

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func TestNewUser(t *testing.T) {
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    User
		wantErr bool
	}{
		{
			name: "Valid email and password",
			args: args{
				email:    "test@example.com",
				password: "password",
			},
			want: User{
				Email:       "test@example.com",
				Password:    "password",
				CreatedAt:   time.Time{},
				UpdatedAt:   time.Time{},
				SuspendedAt: sql.NullTime{Valid: false},
			},
			wantErr: false,
		},
		{
			name: "Invalid email",
			args: args{
				email:    "not-an-email",
				password: "password",
			},
			want:    User{},
			wantErr: true,
		},
		{
			name: "Email too long",
			args: args{
				email:    "loremipsumhasbeentheindustrysstandarddummytexteversincethe1500swhenanunknownprintertookagalleyoftypeandscrambledittomake@example.com",
				password: "password",
			},
			want:    User{},
			wantErr: true,
		},
		{
			name: "Empty password",
			args: args{
				email:    "test@example.com",
				password: "",
			},
			want:    User{},
			wantErr: true,
		},
		{
			name: "Password too short",
			args: args{
				email:    "test@example.com",
				password: "123",
			},
			want:    User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_OnSave(t *testing.T) {
	timeNow := time.Now().UTC()
	passwordPlain := "password"
	passwordHashed, _ := bcrypt.GenerateFromPassword([]byte(passwordPlain), bcrypt.DefaultCost)

	type args struct {
		Password  string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	tests := []struct {
		name    string
		fields  User
		want    args
		wantErr bool
	}{
		{
			name: "Update user",
			fields: User{
				ID:          1,
				Email:       "test@example.com",
				Password:    string(passwordHashed),
				CreatedAt:   time.Time{},
				UpdatedAt:   time.Time{},
				SuspendedAt: sql.NullTime{Valid: false},
			},
			want: args{
				Password:  passwordPlain,
				CreatedAt: timeNow,
				UpdatedAt: timeNow,
			},
			wantErr: false,
		},
		{
			name: "Create user",
			fields: User{
				ID:          0,
				Email:       "test@example.com",
				Password:    "password",
				CreatedAt:   time.Time{},
				UpdatedAt:   time.Time{},
				SuspendedAt: sql.NullTime{Valid: false},
			},
			want: args{
				Password:  "password",
				CreatedAt: timeNow,
				UpdatedAt: timeNow,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := tt.fields
			err := user.OnSave()
			if err != nil {
				t.Errorf("User.OnSave() error = %v", err)
			}

			if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tt.want.Password)); err != nil {
				t.Errorf("User.OnSave() Password = %v, want %v", user.Password, tt.want.Password)
			}
			if reflect.DeepEqual(user.UpdatedAt, tt.want.UpdatedAt) {
				t.Errorf("User.OnSave() UpdatedAt = %v, want %v", user.UpdatedAt, tt.want.UpdatedAt)
			}
		})
	}
}

func TestUser_CheckPassword(t *testing.T) {
	type fields struct {
		ID          int
		Email       string
		Password    string
		CreatedAt   time.Time
		UpdatedAt   time.Time
		SuspendedAt sql.NullTime
	}
	type args struct {
		plainPassword string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Correct password",
			fields: fields{
				ID:          0,
				Email:       "test@example.com",
				Password:    "password",
				CreatedAt:   time.Time{},
				UpdatedAt:   time.Time{},
				SuspendedAt: sql.NullTime{Valid: false},
			},
			args: args{
				plainPassword: "password",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Incorrect password",
			fields: fields{
				ID:          0,
				Email:       "test@example.com",
				Password:    "password",
				CreatedAt:   time.Time{},
				UpdatedAt:   time.Time{},
				SuspendedAt: sql.NullTime{Valid: false},
			},
			args: args{
				plainPassword: "incorrect",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:          tt.fields.ID,
				Email:       tt.fields.Email,
				Password:    tt.fields.Password,
				CreatedAt:   tt.fields.CreatedAt,
				UpdatedAt:   tt.fields.UpdatedAt,
				SuspendedAt: tt.fields.SuspendedAt,
			}
			_ = u.OnSave()
			got, err := u.CheckPassword(tt.args.plainPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckPassword() got = %v, want %v", got, tt.want)
			}
		})
	}
}
