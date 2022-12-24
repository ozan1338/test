package dto

import "test/domain"

type JobListResponse struct {
	Job []domain.Job
	Total_Page  int    `json:"total_page"`
}

type JobDetailResponse struct {
	ID          int    `json:"id"`
	City        string `json:"city"`
	Full_Time   bool   `json:"full_time"`
	Description string `json:"description"`
}

func ToJobDto(j domain.Job) JobDetailResponse {
	return JobDetailResponse{
		ID: j.ID,
		City: j.City,
		Full_Time: j.Full_Time,
		Description: j.Description,
	}
}

type JobRequest struct {
	ID          int    `json:"id"`
	Id_User     int    `json:"id_user"`
	City        string `json:"city"`
	Full_Time   bool   `json:"full_time"`
	Description string `json:"description"`
}