package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/tejashwikalptaru/remote-office/dbHelper"
	"github.com/tejashwikalptaru/remote-office/middleware"
	"github.com/tejashwikalptaru/remote-office/utils"
	"net/http"
	"strconv"
	"time"
)

func CreateLeave(w http.ResponseWriter, r *http.Request){
	userCT:=middleware.UserContext(r)
	userID:=userCT.ID
	reqBody:= struct {
		LeaveFrom  time.Time `json:"leaveFrom"`
		LeaveTo  time.Time `json:"leaveTo"`
		Reason  string `json:"reason"`
		LeaveType string `json:"leaveType"`
	}{}

	if err := utils.ParseBody(r.Body, &reqBody); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Failed to decode request body")
		return
	}
	if err:= dbHelper.InsertLeave(userID, reqBody.LeaveFrom, reqBody.LeaveTo, reqBody.Reason, reqBody.LeaveType); err != nil {

		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to store leave entry")
	}
	w.WriteHeader(http.StatusCreated)
}
func ModifyLeave(w http.ResponseWriter, r *http.Request){
	userCT:=middleware.UserContext(r)
	userID:=userCT.ID
	leaveID,err:=strconv.Atoi(chi.URLParam(r,"id"))
	if err!=nil{
		utils.RespondError(w ,http.StatusBadRequest,err,"Failed to convert leaveid to int")
		return
	}
	reqBody:= struct {
		LeaveFrom  time.Time `json:"leaveFrom"`
		LeaveTo  time.Time `json:"leaveTo"`
		Reason  string `json:"reason"`
		LeaveType string `json:"leaveType"`
	}{}
	if err:=utils.ParseBody(r.Body,&reqBody); err != nil {
		utils.RespondError(w ,http.StatusBadRequest,err,"Failed to decode request Body")
		return
	}
	if err:= dbHelper.ModifyLeave(userID,reqBody.LeaveFrom, reqBody.LeaveTo, reqBody.Reason, reqBody.LeaveType,leaveID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to store leave entry")
	}
	w.WriteHeader(http.StatusCreated)
}

func GetLeave(w http.ResponseWriter, r *http.Request){
	var (
		offSet =0
		limit =20
		err error =nil
	)
	userCT:=middleware.UserContext(r)
	userID:=userCT.ID
	var offsetStr ,limitStr string
	offsetStr=chi.URLParam(r,"offset")
	if offsetStr != "" {
		offSet, err = strconv.Atoi(offsetStr)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, err, "Failed to convert q-offset to int")
			return
		}
	}
	if limitStr=chi.URLParam(r,"limit");limitStr!="" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, err, "Failed to convert limit to int")
			return
		}
	}

	leaves, err:= dbHelper.GetLeave(userID,offSet,limit)
	if  err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to get leave entry")
		return
	}
	utils.RespondJSON(w,http.StatusOK,leaves)
}
func GetLeaveStats(w http.ResponseWriter, r *http.Request){
	userCT:=middleware.UserContext(r)
	userID:=userCT.ID
	year := time.Now().Year()

	leaveStats, err := dbHelper.GetLeaveStats(userID, year)
	if  err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to get leave Stats")
		return
	}
	utils.RespondJSON(w,http.StatusOK,leaveStats)
}