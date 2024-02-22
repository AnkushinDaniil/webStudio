package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"main.go/internal/entity"
)

func TestTimeslotItemPostgres_Create(t *testing.T) {
	dataBase, mock, err := sqlmock.Newx()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dataBase.Close()

	rep := NewTimeslotItemPostgres(dataBase)

	type input struct {
		listID int
		item   entity.TimeslotItem
	}

	type mockBehavior func(args input, ID int)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		input        input
		want         int
		wantErr      bool
	}{
		{
			name: "OK",
			mockBehavior: func(input input, itemID int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(itemID)
				mock.ExpectQuery(`INSERT INTO timeslots_items`).
					WithArgs(input.item.Title, input.item.Description, input.item.Start, input.item.End).
					WillReturnRows(rows)

				mock.ExpectExec(`INSERT INTO lists_items`).
					WithArgs(input.listID, itemID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			input: input{
				listID: 1, item: entity.TimeslotItem{
					Title:       "test title",
					Description: "test description",
					Start:       time.Now(),
					End:         time.Now(),
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
				mock.ExpectQuery(`INSERT INTO timeslots_items`).
					WithArgs(input.item.Title, input.item.Description, input.item.Start, input.item.End).
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
			input: input{
				listID: 1, item: entity.TimeslotItem{
					Title:       "",
					Description: "test description",
					Start:       time.Now(),
					End:         time.Now(),
				},
			},
			wantErr: true,
		},
		{
			name: "2nd insert error",
			mockBehavior: func(input input, itemID int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(itemID)
				mock.ExpectQuery(`INSERT INTO timeslots_items`).
					WithArgs(input.item.Title, input.item.Description, input.item.Start, input.item.End).
					WillReturnRows(rows)

				mock.ExpectExec(`INSERT INTO lists_items`).
					WithArgs(input.listID, itemID).
					WillReturnError(errors.New("some error"))

				mock.ExpectRollback()
			},
			input: input{
				listID: 1, item: entity.TimeslotItem{
					Title:       "test title",
					Description: "test description",
					Start:       time.Now(),
					End:         time.Now(),
				},
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.input, testCase.want)

			got, err := rep.Create(testCase.input.listID, testCase.input.item)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
