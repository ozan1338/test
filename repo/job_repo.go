package repo

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"test/domain"
	"test/log"
	resError "test/util/errors_response"
)

type JobRepo interface {
	GetListJob(domain.Job,int,string) ([]domain.Job,resError.RespError)
	CountListUser(domain.Job, string) (int,resError.RespError)
	GetJob(int) (*domain.Job,resError.RespError)
	InsertJob(domain.Job) (*domain.Job,resError.RespError)
}

var (
	queryGetListJobsWithoutFullTime = `select id,city,full_time,description
	from public.job j
	where j.city like $1 and j.description like $2 
	limit 5 offset $3;
	`

	queryGetListJobsWithFullTime = `
	select id,city,full_time,description
	from public.job j
	where j.city like $1 and j.description like $2 and j.full_time = $3
	limit 5 offset $4;
	`

	queryCountListJobsWithFullTime = `
	select count(*)
	from public.job j
	where j.city ilike $1 and j.description ilike $2 and j.full_time = $3
	`

	queryCountListJobsWithoutFullTime = `
	select count(*)
	from public.job j
	where j.city ilike $1 and j.description ilike $2
	`

	queryGetJobById = `
	select id,id_user ,city,full_time,description
	from public.job j 
	where j.id = $1;
	`

	queryInsertJob =`
	insert into public.job(id_user,city,full_time,description) values($1,$2,$3,$4) returning id;
	`
)

func NewJobRepo(db *sql.DB) JobRepo {
	return &repo{db: db}
}



func (r *repo) GetListJob(j domain.Job,page int, hasFullTime string) ([]domain.Job,resError.RespError) {
	var stmt *sql.Stmt
	var err error

	if hasFullTime != "" {
		stmt, err = r.db.Prepare(queryGetListJobsWithFullTime)
	} else {
		stmt, err = r.db.Prepare(queryGetListJobsWithoutFullTime)
	}

	if err != nil {
		log.Error("error when trying to prepare get list job", err)
		return nil,resError.NewBadRequestError("database error")
	}
	defer stmt.Close()

	var city string
	var description string

	if j.City != "" {
		city = j.City
	} else {
		city = fmt.Sprintf("%%%%")
	}

	if j.Description != "" {
		description = j.Description
	} else {
		description = fmt.Sprintf("%%%%")
	}

	var full_time bool

	
	var rows *sql.Rows
	
	if hasFullTime != "" {
		
		full_time, err = strconv.ParseBool(hasFullTime)
		if err != nil {
			return nil,resError.NewBadRequestError("can't convert full time")
		}

		rows, err = stmt.Query(city,description,full_time,page)
	} else {
		rows, err = stmt.Query(city,description,page)

	}

	if err != nil {
		log.Error("error when trying to get list job", err)
		return nil,resError.NewBadRequestError("database error")
	}
	defer rows.Close()

	result := make([]domain.Job, 0)
	for rows.Next() {
		var job domain.Job
		if err := rows.Scan(&job.ID,&job.City,&job.Full_Time,&job.Description); err != nil {
			log.Error("error when trying to scan user", err)
			return nil,resError.NewBadRequestError("database error")
		}
		result = append(result, job)
	}

	if len(result) == 0 {
		return nil,resError.NewBadRequestError("no user found")
	}

	return result,nil
}

func (r *repo) CountListUser(j domain.Job,hasFullTime string) (int,resError.RespError) {
	var stmt *sql.Stmt
	var err error
	
	if hasFullTime != "" {
		stmt, err = r.db.Prepare(queryCountListJobsWithFullTime)
	}else {
		stmt, err = r.db.Prepare(queryCountListJobsWithoutFullTime)
	}

	if err != nil {
		log.Error("error when trying to prepare count list job", err)
		return 0,resError.NewBadRequestError("database error")
	}
	defer stmt.Close()

	var city string
	var description string

	if j.City != "" {
		city = j.City
	} else {
		city = fmt.Sprintf("%%%%")
	}

	if j.Description != "" {
		description = j.Description
	} else {
		description = fmt.Sprintf("%%%%")
	}

	var full_time bool

	var row *sql.Row
	
	if hasFullTime != "" {
		
		full_time, err = strconv.ParseBool(hasFullTime)
		if err != nil {
			return 0,resError.NewBadRequestError("can't convert full time")
		}

		row = stmt.QueryRow(city,description,full_time)
	} else {
		row = stmt.QueryRow(city,description)

	}

	var totalPage int
	
	if err = row.Scan(&totalPage); err != nil {
		log.Error("error when trying scan count",err)
		return 0, resError.NewBadRequestError("database error")
	}

	return totalPage,nil
}


func (r *repo) GetJob(id int) (*domain.Job, resError.RespError) {
	var job domain.Job

	stmt, err := r.db.Prepare(queryGetJobById)
	if err != nil {
		log.Error("error when trying to prepare get job", err)
		return nil,resError.NewBadRequestError("database error")
	}
	defer stmt.Close()

	if err := stmt.QueryRow(id).Scan(&job.ID,&job.Id_User,&job.City,&job.Full_Time,&job.Description); err != nil {
		if strings.Contains(err.Error(),sql.ErrNoRows.Error()) {
			return nil,resError.NewBadRequestError("user not found")
		}
		log.Error("error when trying scan get job", err)
		return nil,resError.NewBadRequestError("database error")
	}

	return &job,nil
}

func (r *repo) InsertJob(j domain.Job) (*domain.Job,resError.RespError) {
	stmt, err := r.db.Prepare(queryInsertJob)
	if err != nil {
		log.Error("error when trybg to prepare insert job", err)
		return nil,resError.NewBadRequestError("database error")
	}

	defer stmt.Close()

	if err := stmt.QueryRow(j.Id_User,j.City,j.Full_Time,j.Description).Scan(&j.ID); err != nil {
		log.Error("error when trying to scan id", err)
		return nil, resError.NewBadRequestError("database error")
	}

	return &j,nil
}