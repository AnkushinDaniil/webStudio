package rest //nolint:testpackage // need to use handler.userIdentity.

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"go.uber.org/mock/gomock"
	mock_service "main.go/internal/controller/rest/mocks"
	"main.go/internal/service"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorizationService, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorizationService, token string) {
				s.EXPECT().ParseToken(token).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
		},
		{
			name:        "Invalid header name",
			headerName:  "",
			headerValue: "",
			token:       "",
			mockBehavior: func(_ *mock_service.MockAuthorizationService, _ string) {
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"empty authorization header"}`,
		},
		{
			name:                 "Empty header name",
			headerName:           "",
			headerValue:          "Beerer token",
			token:                "token",
			mockBehavior:         func(_ *mock_service.MockAuthorizationService, _ string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"empty authorization header"}`,
		},
		{
			name:                 "Invalid header name",
			headerName:           "Authorization",
			headerValue:          "Beerer token",
			token:                "token",
			mockBehavior:         func(_ *mock_service.MockAuthorizationService, _ string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid authorization header"}`,
		},
		{
			name:                 "Empty Token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "token",
			mockBehavior:         func(_ *mock_service.MockAuthorizationService, _ string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"token is empty"}`,
		},
		{
			name:                 "Empty Token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "token",
			mockBehavior:         func(_ *mock_service.MockAuthorizationService, _ string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"token is empty"}`,
		},
		{
			name:        "Parse Error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorizationService, token string) {
				s.EXPECT().ParseToken(token).Return(0, errors.New("invalid token"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid token"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init deps
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			auth := mock_service.NewMockAuthorizationService(ctrl)
			testCase.mockBehavior(auth, testCase.token)

			services := &service.Service{
				Authorization: auth,
				TimeslotList:  nil,
				TimeslotItem:  nil,
			}
			handler := NewHandlers(services)

			// Test server
			engine := gin.New()
			engine.POST("/protected", handler.userIdentity, func(ctx *gin.Context) {
				id, _ := ctx.Get(userCtx)
				ctx.String(http.StatusOK, strconv.Itoa(id.(int)))
			})

			// Test request
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/protected", nil)
			request.Header.Set(testCase.headerName, testCase.headerValue)

			// Make request
			engine.ServeHTTP(recorder, request)

			// Assert
			assert.Equal(t, recorder.Code, testCase.expectedStatusCode)
			assert.Equal(t, recorder.Body.String(), testCase.expectedResponseBody)
		})
	}
}
