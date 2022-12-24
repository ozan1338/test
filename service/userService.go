package service

import (
	"net/http"
	"test/domain"
	"test/dto"
	"test/repo"
	resError "test/util/errors_response"
)

//go:generate mockgen -destination=../mocks/service/mockUserService.go -package=service test/service UserServiceInterface
type UserServiceInterface interface {
	LoginUser(dto.UsersRequest) (*dto.UsersResponse, resError.RespError)
	CreateUser(dto.UsersRequest) (*dto.UsersResponse, resError.RespError)
}

type UserService struct {
	repo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) UserServiceInterface {
	return &UserService{repo: userRepo}
}

func (s UserService) CreateUser(u dto.UsersRequest) (*dto.UsersResponse, resError.RespError) {

	if err := u.Validate(true); err != nil {
		return nil, err
	}

	var user domain.Users
	user.Email = u.Email
	user.Name = u.Name
	user.Password = u.Password
	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	id, err := s.repo.CreateUser(user);
	if err != nil {
		return nil, err
	}
	
	user.ID = id

	result := dto.ToDto(user)

	return &result, nil
}

func (s UserService) LoginUser(ur dto.UsersRequest) (*dto.UsersResponse, resError.RespError) {
	if err := ur.Validate(false); err != nil {
		return nil, err
	}

	var u domain.Users
	u.Email = ur.Email
	// user.Password = u.Password

	user,getErr := s.repo.GetUserByEmail(&u);
	if getErr != nil {
		return nil, getErr
	}

	if match := user.CheckPassword(ur.Password); !match {
		return nil, resError.NewRespError("password not match", http.StatusUnauthorized, "unauthorized")
	}

	result := dto.ToDto(*user)

	return &result,nil
}