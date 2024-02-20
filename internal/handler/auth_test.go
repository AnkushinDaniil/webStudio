package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"go.uber.org/mock/gomock"
	"main.go/internal/entity"
	"main.go/internal/service"
	mock_service "main.go/internal/service/mocks"
)

func TestHandler_signUp(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockAuthorization, user entity.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            entity.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"username": "username", "name": "Test Name", "password": "qwerty"}`,
			inputUser: entity.User{
				Username: "username",
				Name:     "Test Name",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user entity.User) {
				r.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{"username": "username"}`,
			inputUser:            entity.User{},
			mockBehavior:         func(r *mock_service.MockAuthorization, user entity.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"username": "username", "name": "Test Name", "password": "qwerty"}`,
			inputUser: entity.User{
				Username: "username",
				Name:     "Test Name",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user entity.User) {
				r.EXPECT().CreateUser(user).Return(0, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Dependencies
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			auth := mock_service.NewMockAuthorization(mockCtrl)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := Handler{services}

			// Init Endpoint
			engine := gin.New()
			engine.POST("/sign-up", handler.signUp)

			// Create Request
			writer := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/sign-up",
				bytes.NewBufferString(testCase.inputBody))

			// Make Request
			engine.ServeHTTP(writer, req)

			// Assert
			assert.Equal(t, writer.Code, testCase.expectedStatusCode)
			assert.Equal(t, writer.Body.String(), testCase.expectedResponseBody)
		})
	}
}
