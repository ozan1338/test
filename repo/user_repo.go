package repo

import (
	"database/sql"
	"strings"
	"test/domain"
	"test/log"
	resError "test/util/errors_response"
)

//go:generate mockgen -destination=../mocks/repo/mockUserRepo.go -package=repo test/repo UserRepo
type UserRepo interface {
	GetUserByEmail(*domain.Users) (*domain.Users,resError.RespError)
	CreateUser(u domain.Users) (int, resError.RespError)
}

const (
	queryGetUserByEmail = `select id, password, name from public.Users where email = $1;`
	queryCreateUser = `insert into public.Users(name,email,password) values($1,$2,$3) returning id;`
)

func NewRepoUser(db *sql.DB) UserRepo {
	return &repo{db: db}
}

func (r *repo) GetUserByEmail(u *domain.Users) (*domain.Users,resError.RespError) {
	stmt, err := r.db.Prepare(queryGetUserByEmail)
	if err != nil {
		log.Error("error when trying to prepare get user by email", err)
		return nil,resError.NewBadRequestError("database error")
	}
	defer stmt.Close()

	if err := stmt.QueryRow(u.Email).Scan(&u.ID,&u.Password,&u.Name); err != nil {
		if strings.Contains(err.Error(),sql.ErrNoRows.Error()) {
			return nil,resError.NewBadRequestError("user not found")
		}
		log.Error("error when trying scan get user by email", err)
		return nil,resError.NewBadRequestError("database error")
	}

	return u,nil
}

func (r *repo) CreateUser(u domain.Users) (int, resError.RespError) {
	stmt, err := r.db.Prepare(queryCreateUser)
	if err != nil {
		log.Error("error when trying to prepare create user statement", err)
		return 0, resError.NewBadRequestError("database error")
	}
	defer stmt.Close()

	if err := stmt.QueryRow(u.Name,u.Email,u.Password).Scan(&u.ID); err != nil {
		log.Error("error when trying to scan create user", err)
		return 0, resError.NewBadRequestError("database error")
	}
	
	return u.ID, nil
}