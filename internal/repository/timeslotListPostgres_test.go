package repository_test

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"main.go/internal/entity"
	"main.go/internal/repository"
)

func TestTimeslotsListPostgres_Create(t *testing.T) {
	dataBase, mock, err := sqlmock.Newx()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer dataBase.Close()

	rep := repository.NewTimeslotListPostgres(dataBase)
	query1 := fmt.Sprintf(
		`
			INSERT INTO %s `,
		repository.TimeslotListsTable,
	)
	query2 := fmt.Sprintf(
		`
			INSERT INTO %s `,
		repository.UsersListsTable,
	)

	type input struct {
		userID int
		list   entity.TimeslotsList
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
			mockBehavior: func(input input, listID int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(listID)
				mock.ExpectQuery(query1).
					WithArgs(input.list.Title, input.list.Description).
					WillReturnRows(rows)

				mock.ExpectExec(query2).
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			input: input{
				userID: 1,
				list: entity.TimeslotsList{
					ID:          0,
					Title:       "test title",
					Description: "test description",
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Empty Fields",
			mockBehavior: func(input input, listID int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(listID).
					RowError(0, errors.New("some error"))
				mock.ExpectQuery(query1).
					WithArgs(input.list.Title, input.list.Description).
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
			input: input{
				userID: 1,
				list: entity.TimeslotsList{
					ID:          0,
					Title:       "",
					Description: "test description",
				},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "2nd insert error",
			mockBehavior: func(input input, listID int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(listID)
				mock.ExpectQuery(query1).
					WithArgs(input.list.Title, input.list.Description).
					WillReturnRows(rows)

				mock.ExpectExec(query2).
					WithArgs(input.userID, listID).
					WillReturnError(errors.New("some error"))

				mock.ExpectRollback()
			},
			input: input{
				userID: 1, list: entity.TimeslotsList{
					ID:          0,
					Title:       "test title",
					Description: "test description",
				},
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input, testCase.want)

			got, err1 := rep.Create(testCase.input.userID, testCase.input.list)
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

func TestTimeslotsListPostgres_GetAll(t *testing.T) {
	dataBase, mock, err := sqlmock.Newx()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dataBase.Close()

	rep := repository.NewTimeslotListPostgres(dataBase)
	query := fmt.Sprintf(
		`
			SELECT
			    (.+)
			FROM
			    %s tl
			    INNER JOIN %s ul ON (.+)
			WHERE
			    (.+)`,
		repository.TimeslotListsTable,
		repository.UsersListsTable)

	type input struct {
		listID int
		userID int
	}

	testTable := []struct {
		name         string
		mockBehavior func()
		input        input
		want         []entity.TimeslotsList
		wantErr      bool
	}{
		{
			name: "OK",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow(1, "title1", "description1").
					AddRow(2, "title2", "description2").
					AddRow(3, "title3", "description3")

				mock.ExpectQuery(query).
					WithArgs(1).
					WillReturnRows(rows)
			},
			input: input{
				listID: 1, userID: 1,
			},

			want: []entity.TimeslotsList{
				{
					ID:          1,
					Title:       "title1",
					Description: "description1",
				},
				{
					ID:          2,
					Title:       "title2",
					Description: "description2",
				},
				{
					ID:          3,
					Title:       "title3",
					Description: "description3",
				},
			},
			wantErr: false,
		},
		{
			name: "No Records",
			mockBehavior: func() {
				rows := sqlmock.NewRows(
					[]string{"id", "title", "description"},
				)
				mock.ExpectQuery(query).
					WithArgs(1).
					WillReturnRows(rows)
			},
			input: input{
				listID: 1, userID: 1,
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			got, err1 := rep.GetAll(testCase.input.userID)
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

func TestTimeslotsListPostgres_GetById(t *testing.T) {
	dataBase, mock, err := sqlmock.Newx()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dataBase.Close()

	rep := repository.NewTimeslotListPostgres(dataBase)
	query := fmt.Sprintf(
		`
			SELECT
			    (.+)
			FROM
			    %s tl
			    INNER JOIN %s ul ON (.+)
			WHERE
			    (.+)`,
		repository.TimeslotListsTable,
		repository.UsersListsTable,
	)

	type input struct {
		listID int
		userID int
	}

	testTable := []struct {
		name         string
		mockBehavior func()
		input        input
		want         entity.TimeslotsList
		wantErr      bool
	}{
		{
			name: "OK",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow(1, "title1", "description1")
				mock.ExpectQuery(query).
					WithArgs(1, 1).
					WillReturnRows(rows)
			},
			input: input{
				listID: 1,
				userID: 1,
			},
			want: entity.TimeslotsList{
				ID:          1,
				Title:       "title1",
				Description: "description1",
			},
			wantErr: false,
		},
		{
			name: "NotFound",
			mockBehavior: func() {
				mock.ExpectQuery(query).
					WithArgs(404, 1).
					WillReturnError(sql.ErrNoRows)
			},
			input: input{
				listID: 1,
				userID: 404,
			},
			want: entity.TimeslotsList{
				ID:          0,
				Title:       "",
				Description: "",
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			got, err1 := rep.GetByID(testCase.input.userID, testCase.input.listID)
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

func TestTimeslotsListPostgres_Update(t *testing.T) {
	dataBase, mock, err := sqlmock.Newx()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dataBase.Close()

	rep := repository.NewTimeslotListPostgres(dataBase)
	query := fmt.Sprintf(
		`
			UPDATE
			    %s tl
			SET
			    (.+)
			FROM
			    %s ul
			WHERE
			    (.+)`,
		repository.TimeslotListsTable,
		repository.UsersListsTable,
	)

	type input struct {
		listID int
		userID int
		update entity.UpdateListInput
	}

	newTitle := "New title"
	newDescription := "New description"

	type mockBehavior func()

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		input        input
		wantErr      bool
	}{
		{
			name: "OK",
			mockBehavior: func() {
				mock.ExpectExec(query).
					WithArgs(newTitle, newDescription, 1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: input{
				listID: 1,
				userID: 1,
				update: entity.UpdateListInput{
					Title:       &newTitle,
					Description: &newDescription,
				},
			},

			wantErr: false,
		},
		{
			name: "OK without done",
			mockBehavior: func() {
				mock.ExpectExec(query).
					WithArgs(newTitle, newDescription, 1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: input{
				listID: 1,
				userID: 1,
				update: entity.UpdateListInput{
					Title:       &newTitle,
					Description: &newDescription,
				},
			},

			wantErr: false,
		},
		{
			name: "OK without done and time",
			mockBehavior: func() {
				mock.ExpectExec(query).
					WithArgs(newTitle, newDescription, 1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: input{
				listID: 1,
				userID: 1,
				update: entity.UpdateListInput{
					Title:       &newTitle,
					Description: &newDescription,
				},
			},

			wantErr: false,
		},
		{
			name: "OK empty",
			mockBehavior: func() {
				mock.ExpectExec(fmt.Sprintf(
					`
			UPDATE
			    %s tl
			SET
			FROM
			    %s ul
			WHERE
			    (.+)`,
					repository.TimeslotListsTable,
					repository.UsersListsTable,
				)).
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: input{
				listID: 1,
				userID: 1,
				update: entity.UpdateListInput{
					Title:       nil,
					Description: nil,
				},
			},

			wantErr: false,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			err1 := rep.Update(testCase.input.userID, testCase.input.listID, testCase.input.update)
			if testCase.wantErr {
				require.Error(t, err1)
			} else {
				require.NoError(t, err1)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTimeslotsListPostgres_Delete(t *testing.T) {
	dataBase, mock, err := sqlmock.Newx()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dataBase.Close()

	rep := repository.NewTimeslotListPostgres(dataBase)
	query := fmt.Sprintf(
		`
			DELETE FROM %s tl USING %s ul
			WHERE (.+)`,
		repository.TimeslotListsTable,
		repository.UsersListsTable,
	)

	type input struct {
		listID int
		userID int
	}

	testTable := []struct {
		name         string
		mockBehavior func()
		input        input
		wantErr      bool
	}{
		{
			name: "OK",
			mockBehavior: func() {
				mock.ExpectExec(query).
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: input{
				listID: 1,
				userID: 1,
			},

			wantErr: false,
		},
		{
			name: "Not found",
			mockBehavior: func() {
				mock.ExpectExec(query).
					WithArgs(1, 404).
					WillReturnError(sql.ErrNoRows)
			},
			input: input{
				listID: 404,
				userID: 1,
			},

			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			err1 := rep.Delete(testCase.input.userID, testCase.input.listID)
			if testCase.wantErr {
				require.Error(t, err1)
			} else {
				require.NoError(t, err1)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
