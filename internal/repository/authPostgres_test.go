package repository_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"main.go/internal/entity"
	"main.go/internal/repository"
)

func TestAuthPostgres_CreateUser(t *testing.T) {
	dataBase, mock, err := sqlmock.Newx()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dataBase.Close()

	rep := repository.NewAuthPostgres(dataBase)
	query := fmt.Sprintf(
		`
			INSERT INTO %s`,
		repository.UsersTable,
	)

	type input struct {
		user entity.User
	}

	testTable := []struct {
		name         string
		mockBehavior func(args input, ID int)
		input        input
		want         int
		wantErr      bool
	}{
		{
			name: "OK",
			mockBehavior: func(input input, userID int) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(userID)
				mock.ExpectQuery(query).
					WithArgs(input.user.Name, input.user.Username, input.user.Password).
					WillReturnRows(rows)
			},
			input: input{
				user: entity.User{
					ID:       0,
					Name:     "name",
					Username: "username",
					Password: "password",
				},
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "Empty Fields",
			mockBehavior: func(input input, userID int) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(userID).
					RowError(0, errors.New("some error"))
				mock.ExpectQuery(query).
					WithArgs(input.user.Name, input.user.Username, input.user.Password).
					WillReturnRows(rows)
			},
			input: input{
				user: entity.User{
					ID:       0,
					Name:     "",
					Username: "username",
					Password: "password",
				},
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input, testCase.want)

			got, err1 := rep.CreateUser(testCase.input.user)
			if testCase.wantErr {
				require.Error(t, err1)
			} else {
				require.NoError(t, err1)
				require.Equal(t, testCase.want, got)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthPostgres_GetUser(t *testing.T) {
	dataBase, mock, err := sqlmock.Newx()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dataBase.Close()

	rep := repository.NewAuthPostgres(dataBase)
	query := fmt.Sprintf(`
		SELECT
		    (.+)
		FROM
		    %s
		WHERE (.+)`, repository.UsersTable)
	user := entity.User{
		ID:       1,
		Name:     "name1",
		Username: "username1",
		Password: "passwordHash1",
	}

	type input struct {
		username     string
		passwordHash string
	}

	testTable := []struct {
		name         string
		mockBehavior func(args input)
		input        input
		want         entity.User
		wantErr      bool
	}{
		{
			name: "OK",
			mockBehavior: func(input input) {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "password_hash"}).
					AddRow(user.ID, user.Name, user.Username, user.Password)
				mock.ExpectQuery(query).
					WithArgs(input.username, input.passwordHash).
					WillReturnRows(rows)
			},
			input: input{
				username:     user.Username,
				passwordHash: user.Password,
			},
			want:    user,
			wantErr: false,
		},
		{
			name: "Wrong username",
			mockBehavior: func(input input) {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "password_hash"}).
					AddRow(user.ID, user.Name, user.Username, user.Password).
					RowError(0, errors.New("some error"))
				mock.ExpectQuery(query).
					WithArgs(input.username, input.passwordHash).
					WillReturnRows(rows)
			},
			input: input{
				username:     "wrong_username",
				passwordHash: user.Password,
			},
			want: entity.User{
				ID:       0,
				Name:     "",
				Username: "",
				Password: "",
			},
			wantErr: true,
		},
		{
			name: "Wrong password",
			mockBehavior: func(input input) {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "password_hash"}).
					AddRow(user.ID, user.Name, user.Username, user.Password).
					RowError(0, errors.New("some error"))
				mock.ExpectQuery(query).
					WithArgs(input.username, input.passwordHash).
					WillReturnRows(rows)
			},
			input: input{
				username:     user.Username,
				passwordHash: "wrong_password",
			},
			want: entity.User{
				ID:       0,
				Name:     "",
				Username: "",
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input)

			got, err1 := rep.GetUser(testCase.input.username, testCase.input.passwordHash)
			if testCase.wantErr {
				require.Error(t, err1)
			} else {
				require.NoError(t, err1)
				require.Equal(t, testCase.want, got)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
