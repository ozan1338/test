package router

import (
	"net/http"
	"test/handlers/pinghandler"
	"test/handlers/userhandler"
	"test/helpers"
	"test/pkg/jwt"
	"test/pkg/postgresql"
	"test/repo"
	"test/service"

	"github.com/gorilla/mux"
)

func userRouter(r *mux.Router) {
	//TODO
	//1.add repository db
	db := postgresql.Conn
	repo := repo.NewRepoUser(db)

	//2. add helper
	helper := helpers.NewHelper()

	//4. add service
	userService := service.NewUserService(repo)

	//5. add jwt
	jwtMaker := jwt.NewJWTMaker("secret")

	//6.add handlers for user
	h := pinghandler.PingHandler(helper)
	u := userhandler.NewUserHandler(userService, helper, jwtMaker)

	//serve the route
	r.HandleFunc("/", h.Ping).Methods(http.MethodGet)
	r.HandleFunc("/login", u.LoginUser).Methods(http.MethodPost)
	r.HandleFunc("/register", u.Register).Methods(http.MethodPost)
}