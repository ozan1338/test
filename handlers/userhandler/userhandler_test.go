package userhandler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"test/dto"
	"testing"
	"time"

	"test/helpers"
	mockJwt "test/mocks/pkg/jwt"
	mockService "test/mocks/service"
	_ "test/pkg/jwt"
	resError "test/util/errors_response"

	"github.com/golang/mock/gomock"
)

var (
	s *mockService.MockUserServiceInterface
	h *userHandler
	jwtMaker *mockJwt.MockMaker
	// router *mux.Router
)

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	s = mockService.NewMockUserServiceInterface(ctrl)
	helper := helpers.NewHelper()
	// jwtMaker := jwt.NewJWTMaker("secret")
	jwtMaker  = mockJwt.NewMockMaker(ctrl)
	h = NewUserHandler(s,helper,jwtMaker)
	// router = mux.NewRouter()
	return func ()  {
		// router = nil
		defer ctrl.Finish()
	}
}

var users = []dto.UsersResponse{
	{ID: 1,Email: "test@mail.com",Name: "ozan", JWT: ""},
}

var usersRequest = dto.UsersRequest{
	Name: "ozan",
	Email: "test@mail.com",
	Password: "123",
}

func TestUserHandlerLoginUser(t *testing.T) {
	//arrange 
	var teardown = setup(t)
	defer teardown()

	test := []struct{
		name string
		json string
		expectedStatus int
		stubService func() *gomock.Call
		stubJwt func() *gomock.Call
	} {
		{
			name:"login user ok",
			json: `{
				"email":"test@mail.com",
				"password":"123"
			}`,
			expectedStatus: http.StatusOK,
			stubService: func() *gomock.Call {
				var testUser = usersRequest
				testUser.Name = ""
				return s.EXPECT().LoginUser(testUser).Return(&users[0],nil)
			},
			stubJwt: func() *gomock.Call {
				return jwtMaker.EXPECT().CreateToken(users[0].ID, 15 * time.Minute)
			},
		},
		{
			name:"login user bad json",
			json: `{
				"email":"test@mail.com",
				"password":123
			}`,
			expectedStatus: http.StatusBadRequest,
			stubService: func() *gomock.Call {
				return nil
			},
			stubJwt: func() *gomock.Call {
				return nil
			},
		},
		{
			name:"login user error",
			json: `{
				"email":"test@mail.com",
				"password":"123"
			}`,
			expectedStatus: http.StatusBadRequest,
			stubService: func() *gomock.Call {
				var testUser = usersRequest
				testUser.Name = ""
				return s.EXPECT().LoginUser(testUser).Return(nil,resError.NewBadRequestError("some error"))
			},
			stubJwt: func() *gomock.Call {
				return nil
			},
		},
		{
			name:"login user jwt error",
			json: `{
				"email":"test@mail.com",
				"password":"123"
			}`,
			expectedStatus: http.StatusUnauthorized,
			stubService: func() *gomock.Call {
				var testUser = usersRequest
				testUser.Name = ""
				return s.EXPECT().LoginUser(testUser).Return(&users[0],nil)
			},
			stubJwt: func() *gomock.Call {
				return jwtMaker.EXPECT().CreateToken(users[0].ID, 15 * time.Minute).Return("",nil,resError.NewRespError("some error", http.StatusUnauthorized,"unauthorized"))
			},
		},
	}

	for _, item := range test {
		item.stubJwt()
		item.stubService()

		var req *http.Request
		req, _ = http.NewRequest(http.MethodPost,"/", strings.NewReader(item.json))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(h.LoginUser)
		handler.ServeHTTP(rr,req)

		if rr.Code != item.expectedStatus {
			t.Errorf("%s : expected %d but got %d",item.name,item.expectedStatus,rr.Code)
		}
	}

}


func TestUserHandlerRegisterUser(t *testing.T) {
	//arrange 
	var teardown = setup(t)
	defer teardown()

	test := []struct{
		name string
		json string
		expectedStatus int
		stubService func() *gomock.Call
		stubJwt func() *gomock.Call
	} {
		{"create user ok",
		`{
			"name":"ozan",
			"email":"test@mail.com",
			"password":"123"
		}`, http.StatusOK, 
		func() *gomock.Call { 
			return s.EXPECT().CreateUser(usersRequest).Return(&users[0],nil) 
		},
		func() *gomock.Call {
			return jwtMaker.EXPECT().CreateToken(users[0].ID,15 * time.Minute)
		}},
		{"create user bad json",
		`{
			"name":"ozan",
			"email":"test@mail.com",
			"password":123
		}`, http.StatusBadRequest, 
		func() *gomock.Call { 
			return nil
		},func() *gomock.Call {
			return nil
		}},
		{"create user error",
		`{
			"name":"ozan",
			"email":"test@mail.com",
			"password":"123"
		}`, http.StatusBadRequest, 
		func() *gomock.Call { 
			return s.EXPECT().CreateUser(usersRequest).Return(nil,resError.NewBadRequestError("some error"))
		},func() *gomock.Call {
			return nil
		}},
		{"create user jwt error",
		`{
			"name":"ozan",
			"email":"test@mail.com",
			"password":"123"
		}`, http.StatusUnauthorized, 
		func() *gomock.Call { 
			return s.EXPECT().CreateUser(usersRequest).Return(&users[0],nil)
		},func() *gomock.Call {
			return jwtMaker.EXPECT().CreateToken(users[0].ID,15 * time.Minute).Return("",nil,resError.NewRespError("some error", http.StatusUnauthorized, "jwt error"))
		}},
	}

	//act
	for _, item := range test{
		item.stubService()
		item.stubJwt()

		var req *http.Request
		req, _ = http.NewRequest(http.MethodPost,"/", strings.NewReader(item.json))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(h.Register)
		handler.ServeHTTP(rr,req)

		if rr.Code != item.expectedStatus {
			t.Errorf("%s: wrong status return; expected %d but got %d", item.name, item.expectedStatus, rr.Code)
		}
	}
}