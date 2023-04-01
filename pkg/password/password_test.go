package password

import (
	"blogapi/pkg/utils"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    func(t *testing.T, hashedPassword string)
		wantErr bool
	}{
		{
			name: "Password Hashed",
			args: args{
				email: utils.RandomEmail(),
				password: utils.RandomString(6),
			},
			want: func(t *testing.T, hashedPassword string) {
				require.NotEmpty(t, hashedPassword)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashPassword(tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.want(t, got)
		})
	}
}

func TestCheckPassword(t *testing.T) {

	email := utils.RandomEmail()
	password := utils.RandomString(6)
	hashedPassword, _ := HashPassword(email, password)
	wrongHashedPassoword, _:= HashPassword(utils.RandomEmail(), password)

	type args struct {
		email          string
		password       string
		hashedPassword string
	}
	tests := []struct {
		name    string
		args    args
		want    func(t *testing.T, err error)
		wantErr bool
	}{
		{
			name: "Password Match",
			args: args{
				email: email,
				password: password,
				hashedPassword: hashedPassword,
			},
			want: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "Password Don't Match",
			args: args{
				email: email,
				password: password,
				hashedPassword: wrongHashedPassoword,
			},
			want: func(t *testing.T, err error) {
				require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {			
			err := CheckPassword(tt.args.email, tt.args.password, tt.args.hashedPassword); 
			tt.want(t, err)
		})
	}
}
