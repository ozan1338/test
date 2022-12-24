package userhandler

import (
	"net/http"
	"test/dto"
	"test/helpers"
	"test/service"
	"time"

	"test/pkg/jwt"
)

var accesTokenDuration time.Duration = 15 * time.Minute

type userHandler struct {
	userService service.UserServiceInterface
	helpers helpers.HelpersInterface
	JWT jwt.Maker
}

func NewUserHandler(userService service.UserServiceInterface, helpers helpers.HelpersInterface,JWT jwt.Maker) *userHandler {
	return &userHandler{
		userService: userService,
		helpers: helpers,
		JWT: JWT,
	}
}

func (h userHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var user dto.UsersRequest
	if err := h.helpers.ReadJSON(w,r,&user); err != nil {
		h.helpers.WriteResponse(w,err.GetStatus(),err)
		return
	}

	result, err := h.userService.LoginUser(user)
	if err != nil {
		h.helpers.WriteResponse(w,err.GetStatus(),err)
		return
	}

	jwtToken, _, err := h.JWT.CreateToken(result.ID, accesTokenDuration)
	if err != nil {
		h.helpers.WriteResponse(w,err.GetStatus(),err)
		return
	}

	result.JWT = jwtToken

	h.helpers.WriteResponse(w,http.StatusOK, result)
}

func (h userHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user dto.UsersRequest
	if err := h.helpers.ReadJSON(w,r,&user); err != nil {
		h.helpers.WriteResponse(w,err.GetStatus(),err)
		return
	}
	
	result, err := h.userService.CreateUser(user)
	if err != nil {
		h.helpers.WriteResponse(w,err.GetStatus(),err)
		return
	}

	

	jwtToken, _, err :=h.JWT.CreateToken(result.ID,accesTokenDuration)

	if err != nil {
		h.helpers.WriteResponse(w,err.GetStatus(),err)
		return
	}

	result.JWT = jwtToken

	h.helpers.WriteResponse(w,http.StatusOK, result)

}