package jobhandler

import (
	"net/http"
	"strconv"
	"test/dto"
	"test/helpers"
	"test/pkg/jwt"
	"test/service"
	resError "test/util/errors_response"

	"github.com/gorilla/mux"
)

type jobHandler struct {
	jobService service.JobServiceInterface
	helpers helpers.HelpersInterface
	JWT     jwt.Maker
}

func NewJobHandler(jobService service.JobServiceInterface,helpers helpers.HelpersInterface, JWT jwt.Maker) *jobHandler {
	return &jobHandler{
		jobService: jobService,
		helpers:     helpers,
		JWT:         JWT,
	}
}

func (h jobHandler) GetJobList(w http.ResponseWriter, r *http.Request) {
	// create variable to hold query http
	// fmt.Println("kepanggil")
	description := r.URL.Query().Get("description")
	city := r.URL.Query().Get("city")
	page := r.URL.Query().Get("page")
	full_time := r.URL.Query().Get("full_time")

	var offset int
	var err error

	if page != "" {
		offset, err = strconv.Atoi(page)
		if err != nil {
			errConvert := resError.NewBadRequestError("can't convert")
			h.helpers.WriteResponse(w,errConvert.GetStatus(),errConvert)
			return
		}

		if offset > 1 {
			offset = (offset - 1) * 5;
		}
		
	} else {
		offset = 0
	}

	// call service to get list
	result, getErr := h.jobService.GetListJobs(description,city,full_time,offset)
	if getErr != nil {
		h.helpers.WriteResponse(w,getErr.GetStatus(),getErr)
		return
	}


	//send the response
	h.helpers.WriteResponse(w,http.StatusOK,result)
}

func (h jobHandler) GetDetailJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["job_id"]

	id_job,_ := strconv.Atoi(id)

	result, err := h.jobService.GetDetailJob(id_job)
	if err != nil {
		h.helpers.WriteResponse(w,err.GetStatus(),err)
		return
	}


	h.helpers.WriteResponse(w,http.StatusOK,result)
}

func (h jobHandler) InsertJob(w http.ResponseWriter, r *http.Request) {
	var job dto.JobRequest

	if err := h.helpers.ReadJSON(w,r,&job); err != nil {
		h.helpers.WriteResponse(w,err.GetStatus(),err)
		return
	}

	result, err := h.jobService.InsertJob(job)
	if err != nil {
		h.helpers.WriteResponse(w,err.GetStatus(),err)
		return
	}

	h.helpers.WriteResponse(w,http.StatusOK,result)
}