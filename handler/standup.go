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

func CreateStandUp(w http.ResponseWriter,r *http.Request){
	userCT:=middleware.UserContext(r)
	userID:= userCT.ID
	reqBody:= struct {
		StandUp string `json:"data"`
		Date time.Time  `json:"date"`
	}{}
	if err := utils.ParseBody(r.Body, &reqBody); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Failed to decode request body")
		return
	}

	if err:= dbHelper.InsertStandUp(userID,reqBody.StandUp,reqBody.Date); err!=nil{
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to store standup entry")
	}
	w.WriteHeader(http.StatusCreated)
	}

func ModifyStandUp(w http.ResponseWriter, r *http.Request){
	userCT:=middleware.UserContext(r)
	userID:=userCT.ID
	standUpID,err:=strconv.Atoi(chi.URLParam(r,"id"))
	if err!=nil{
		utils.RespondError(w ,http.StatusBadRequest,err,"Failed to convert standUp ID to int")
		return
	}
	reqBody:= struct {
		StandUp string `json:"data"`
		Date time.Time  `json:"date"`
	}{}
	if err := utils.ParseBody(r.Body, &reqBody); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Failed to decode request body")
		return
	}
	if err:= dbHelper.ModifyStandUp(userID,reqBody.StandUp,reqBody.Date,standUpID); err!=nil{
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to store standUp entry")
	}
	w.WriteHeader(http.StatusCreated)
}

func GetStandUp(w http.ResponseWriter, r *http.Request){
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

	leaves, err:= dbHelper.GetStandUp(userID,offSet,limit)
	if  err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to get leave entry")
		return
	}
	utils.RespondJSON(w,http.StatusOK,leaves)
}