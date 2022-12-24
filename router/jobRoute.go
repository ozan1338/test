package router

import (
	"fmt"
	"net/http"
	"strings"
	"test/handlers/jobhandler"
	"test/helpers"
	"test/pkg/jwt"
	"test/pkg/postgresql"
	"test/repo"
	"test/service"

	resError "test/util/errors_response"

	"github.com/gorilla/mux"
)

var (
	authorization = "Authorization"
	helper = helpers.NewHelper()
	unauthorized = "unauthorized"
	jwtMaker = jwt.NewJWTMaker("secret")
)

func jobRouter(r *mux.Router) {
	//TODO
	//2.add repository db
	db := postgresql.Conn
	repo := repo.NewJobRepo(db)

	//3.add service
	s := service.NewJobService(repo)

	//6.add handlers
	h := jobhandler.NewJobHandler(s,helper,jwtMaker)

	//7.serve
	subRouter := r.PathPrefix("/job").Subrouter()
	subRouter.Use(auth)
	subRouter.HandleFunc("/get-list",h.GetJobList).Methods(http.MethodGet)
	subRouter.HandleFunc("/get-detail/{job_id:[0-9]+}",h.GetDetailJob).Methods(http.MethodGet)
	subRouter.HandleFunc("/create-job",h.InsertJob).Methods(http.MethodPost)
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the authorization header
		authHeader := r.Header.Get(authorization)

		//check
		if err := checkBearerAuth(authHeader,r); err != nil {
			fmt.Println(err)
			helper.WriteResponse(w,err.GetStatus(),err)
			return
		}
		next.ServeHTTP(w,r)
	})
}

func checkBearerAuth(authHeader string, r *http.Request) resError.RespError {
	//sanity check

	if authHeader == "" {
		return resError.NewRespError("no auth header", http.StatusUnauthorized, unauthorized)
	}

	//split header on spaces
	headersSplit := strings.Split(authHeader, " ")
	if len(headersSplit) != 2 {
		return resError.NewRespError("invalid auth header", http.StatusUnauthorized, unauthorized)
	}

	//check to see if we have word "Bearer"
	if headersSplit[0] != "Bearer" {
		return resError.NewRespError("unauthorized: no Bearer", http.StatusUnauthorized, unauthorized)
	}

	token := headersSplit[1]

	_ , err := jwtMaker.VerifyToken(token)
	if err != nil {
		return err
	}

	return nil
}