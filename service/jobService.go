package service

import (
	"test/domain"
	"test/dto"
	"test/repo"
	resError "test/util/errors_response"
)

type JobServiceInterface interface {
	GetListJobs(string,  string,  string,  int) (*dto.JobListResponse,resError.RespError)
	GetDetailJob(int) (*dto.JobDetailResponse,resError.RespError)
	InsertJob(dto.JobRequest) (*dto.JobDetailResponse,resError.RespError)
}

type JobService struct {
	repo repo.JobRepo
}

func NewJobService(jobRepo repo.JobRepo) JobServiceInterface {
	return &JobService{repo: jobRepo}
}

func (s JobService) InsertJob(j dto.JobRequest) (*dto.JobDetailResponse,resError.RespError) {
	var initialJob  domain.Job

	initialJob.City = j.City
	initialJob.Description = j.Description
	initialJob.Full_Time = j.Full_Time
	initialJob.Id_User = j.Id_User

	job, err := s.repo.InsertJob(initialJob)
	if err != nil {
		return nil,err
	}

	result := dto.ToJobDto(*job)

	return &result,nil
}

func (s JobService) GetDetailJob(id int) (*dto.JobDetailResponse,resError.RespError) {
	job, err := s.repo.GetJob(id)
	if err != nil {
		return nil,err
	}

	result := dto.ToJobDto(*job)

	return &result,nil
}

func (s JobService) GetListJobs(description,city string, full_time string, page int) (*dto.JobListResponse,resError.RespError) {
	var query domain.Job = domain.Job{
		Description: description,
		City: city,
	}

	resultQuery, err := s.repo.GetListJob(query,page,full_time)
	if err != nil {
		return nil,err
	}

	// var job dto.JobResponse
	var response dto.JobListResponse

	response.Job = append(response.Job, resultQuery...)

	totalPage, err := s.repo.CountListUser(query,full_time)
	if err != nil {
		return nil,err
	}
	response.Total_Page = totalPage
	
	return &response,nil
}