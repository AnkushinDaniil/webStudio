package repository_test

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"main.go/internal/entity"
	"main.go/internal/repository"
)

func TestTimeslotItemPostgres_Create(t *testing.T) {
	dataBase, mock, err := sqlmock.Newx()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dataBase.Close()

	rep := repository.NewTimeslotItemPostgres(dataBase)
	query1 := `
		INSERT INTO timeslots_items`
	query2 := `
					INSERT INTO lists_items`

	type input struct {
		listID int
		item   entity.TimeslotItem
	}

	timeNow := time.Now()

	testTable := []struct {
		name         string
		mockBehavior func(args input, ID int)
		input        input
		want         int
		wantErr      bool
	}{
		{
			name: "OK",
			mockBehavior: func(input input, itemID int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(itemID)
				mock.ExpectQuery(query1).
					WithArgs(input.item.Title, input.item.Description, input.item.Start, input.item.End).
					WillReturnRows(rows)

				mock.ExpectExec(query2).
					WithArgs(input.listID, itemID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			input: input{
				listID: 1, item: entity.TimeslotItem{
					ID:          0,
					Title:       "test title",
					Description: "test description",
					Start:       timeNow,
					End:         timeNow,
					Done:        false,
				},
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "Empty Fields",
			mockBehavior: func(input input, itemID int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(itemID).
					RowError(0, errors.New("some error"))
				mock.ExpectQuery(query1).
					WithArgs(input.item.Title, input.item.Description, input.item.Start, input.item.End).
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
			input: input{
				listID: 1, item: entity.TimeslotItem{
					ID:          0,
					Title:       "",
					Description: "test description",
					Start:       timeNow,
					End:         timeNow,
					Done:        false,
				},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "2nd insert error",
			mockBehavior: func(input input, itemID int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(itemID)
				mock.ExpectQuery(query1).
					WithArgs(input.item.Title, input.item.Description, input.item.Start, input.item.End).
					WillReturnRows(rows)

				mock.ExpectExec(query2).
					WithArgs(input.listID, itemID).
					WillReturnError(errors.New("some error"))

				mock.ExpectRollback()
			},
			input: input{
				listID: 1, item: entity.TimeslotItem{
					ID:          0,
					Title:       "test title",
					Description: "test description",
					Start:       timeNow,
					End:         timeNow,
					Done:        false,
				},
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input, testCase.want)

			got, err1 := rep.Create(testCase.input.listID, testCase.input.item)
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

func TestTimeslotItemPostgres_GetAll(t *testing.T) {
	dataBase, mock, err := sqlmock.Newx()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dataBase.Close()

	rep := repository.NewTimeslotItemPostgres(dataBase)
	query := fmt.Sprintf(
		`
			SELECT
			    (.+)
			FROM
			    %s ti
			    INNER JOIN %s li ON (.+)
			    INNER JOIN %s ul ON (.+)
			WHERE (.+)`,
		repository.TimeslotsItemsTable,
		repository.ListsItemsTable,
		repository.UsersListsTable)

	type input struct {
		listID int
		userID int
	}

	timeNow := time.Now()

	testTable := []struct {
		name         string
		mockBehavior func()
		input        input
		want         []entity.TimeslotItem
		wantErr      bool
	}{
		{
			name: "OK",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "beginning", "finish", "done"}).
					AddRow(1, "title1", "description1", timeNow, timeNow, true).
					AddRow(2, "title2", "description2", timeNow, timeNow, false).
					AddRow(3, "title3", "description3", timeNow, timeNow, false)

				mock.ExpectQuery(query).
					WithArgs(1, 1).
					WillReturnRows(rows)
			},
			input: input{
				listID: 1, userID: 1,
			},

			want: []entity.TimeslotItem{
				{
					ID:          1,
					Title:       "title1",
					Description: "description1",
					Start:       timeNow,
					End:         timeNow,
					Done:        true,
				},
				{
					ID:          2,
					Title:       "title2",
					Description: "description2",
					Start:       timeNow,
					End:         timeNow,
					Done:        false,
				},
				{
					ID:          3,
					Title:       "title3",
					Description: "description3",
					Start:       timeNow,
					End:         timeNow,
					Done:        false,
				},
			},
			wantErr: false,
		},
		{
			name: "No Records",
			mockBehavior: func() {
				rows := sqlmock.NewRows(
					[]string{"id", "title", "description", "beginning", "finish", "done"},
				)
				mock.ExpectQuery(query).
					WithArgs(1, 1).
					WillReturnRows(rows)
			},
			input: input{
				listID: 1,
				userID: 1,
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			got, err1 := rep.GetAll(testCase.input.userID, testCase.input.listID)
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

func TestTimeslotItemPostgres_GetById(t *testing.T) {
	dataBase, mock, err := sqlmock.Newx()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dataBase.Close()

	rep := repository.NewTimeslotItemPostgres(dataBase)
	query := fmt.Sprintf(
		`
			SELECT
			    (.+)
			FROM
			    %s ti
			    INNER JOIN %s li ON (.+)
			    INNER JOIN %s ul ON (.+)
			WHERE (.+)`,
		repository.TimeslotsItemsTable,
		repository.ListsItemsTable,
		repository.UsersListsTable,
	)

	type input struct {
		itemID int
		userID int
	}

	timeNow := time.Now()

	testTable := []struct {
		name         string
		mockBehavior func()
		input        input
		want         entity.TimeslotItem
		wantErr      bool
	}{
		{
			name: "OK",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "beginning", "finish", "done"}).
					AddRow(1, "title1", "description1", timeNow, timeNow, true)
				mock.ExpectQuery(query).
					WithArgs(1, 1).
					WillReturnRows(rows)
			},
			input: input{
				itemID: 1,
				userID: 1,
			},
			want: entity.TimeslotItem{
				ID:          1,
				Title:       "title1",
				Description: "description1",
				Start:       timeNow,
				End:         timeNow,
				Done:        true,
			},
			wantErr: false,
		},
		{
			name: "NotFound",
			mockBehavior: func() {
				mock.ExpectQuery(query).
					WithArgs(1, 404).
					WillReturnError(sql.ErrNoRows)
			},
			input: input{
				itemID: 1,
				userID: 404,
			},
			want: entity.TimeslotItem{
				ID:          1,
				Title:       "title1",
				Description: "description1",
				Start:       timeNow,
				End:         timeNow,
				Done:        true,
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			got, err1 := rep.GetByID(testCase.input.userID, testCase.input.itemID)
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

func TestTimeslotItemPostgres_Update(t *testing.T) {
	dataBase, mock, err := sqlmock.Newx()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dataBase.Close()

	rep := repository.NewTimeslotItemPostgres(dataBase)
	query := fmt.Sprintf(
		`
			UPDATE
			    %s ti
			SET
			    (.+)
			FROM
			    %s li,
			    %s ul
			WHERE (.+)`,
		repository.TimeslotsItemsTable,
		repository.ListsItemsTable,
		repository.UsersListsTable,
	)

	type input struct {
		itemID int
		userID int
		update entity.UpdateItemInput
	}

	timeNow := time.Now()
	newTitle := "New title"
	newDescription := "New description"
	newDone := true

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
					WithArgs(newTitle, newDescription, timeNow, timeNow, newDone, 1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: input{
				itemID: 1,
				userID: 1,
				update: entity.UpdateItemInput{
					Title:       &newTitle,
					Description: &newDescription,
					Start:       &timeNow,
					End:         &timeNow,
					Done:        &newDone,
				},
			},

			wantErr: false,
		},
		{
			name: "OK without done",
			mockBehavior: func() {
				mock.ExpectExec(query).
					WithArgs(newTitle, newDescription, timeNow, timeNow, 1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: input{
				itemID: 1,
				userID: 1,
				update: entity.UpdateItemInput{
					Title:       &newTitle,
					Description: &newDescription,
					Start:       &timeNow,
					End:         &timeNow,
					Done:        nil,
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
				itemID: 1,
				userID: 1,
				update: entity.UpdateItemInput{
					Title:       &newTitle,
					Description: &newDescription,
					Start:       nil,
					End:         nil,
					Done:        nil,
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
						    %s ti
						SET
						FROM
						    %s li,
						    %s ul
						WHERE (.+)`,
					repository.TimeslotsItemsTable,
					repository.ListsItemsTable,
					repository.UsersListsTable,
				)).
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: input{
				itemID: 1,
				userID: 1,
				update: entity.UpdateItemInput{
					Title:       nil,
					Description: nil,
					Start:       nil,
					End:         nil,
					Done:        nil,
				},
			},

			wantErr: false,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			err1 := rep.Update(testCase.input.userID, testCase.input.itemID, testCase.input.update)
			if testCase.wantErr {
				require.Error(t, err1)
			} else {
				require.NoError(t, err1)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTimeslotItemPostgres_Delete(t *testing.T) {
	dataBase, mock, err := sqlmock.Newx()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dataBase.Close()

	rep := repository.NewTimeslotItemPostgres(dataBase)
	query := fmt.Sprintf(
		`
			DELETE FROM %s ti USING %s li, %s ul
			WHERE (.+)`,
		repository.TimeslotsItemsTable,
		repository.ListsItemsTable,
		repository.UsersListsTable,
	)

	type input struct {
		itemID int
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
				itemID: 1,
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
				itemID: 404,
				userID: 1,
			},

			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			err1 := rep.Delete(testCase.input.userID, testCase.input.itemID)
			if testCase.wantErr {
				require.Error(t, err1)
			} else {
				require.NoError(t, err1)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
