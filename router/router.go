package router

import "github.com/gorilla/mux"

func RouterInit(r *mux.Router) {
	userRouter(r)
	jobRouter(r)
}