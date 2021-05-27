package handler

import (
	"fmt"
	"github.com/tejashwikalptaru/remote-office/middleware"
	"net/http"
)

func AdminRights(w http.ResponseWriter, r *http.Request){
	userCT:=middleware.UserContext(r)
	userName:=userCT.Name
	fmt.Sprintln(userName)
	w.WriteHeader(http.StatusOK)
}
