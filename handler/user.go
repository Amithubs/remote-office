package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/tejashwikalptaru/remote-office/dbHelper"
	"github.com/tejashwikalptaru/remote-office/middleware"
	"github.com/tejashwikalptaru/remote-office/models"
	"github.com/tejashwikalptaru/remote-office/utils"
	"net/http"
	"os"
	"time"
)


func SignUp(w http.ResponseWriter, r *http.Request) {
	reqBody := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := utils.ParseBody(r.Body, &reqBody); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Failed to decode request body")
		return
	}
	exist, err := dbHelper.IsUserExist(reqBody.Email)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to find user")
		return
	}
	if exist > 0 {
		utils.RespondError(w, http.StatusBadRequest, nil, "User already exist")
		return
	}
	hash, err := utils.HashAndSaltPassword(reqBody.Password)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to secure password")
		return
	}
	if err := dbHelper.InsertUser(reqBody.Name, reqBody.Email, hash,models.Employee); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to create user")
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	reqBody := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := utils.ParseBody(r.Body, &reqBody); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Failed to decode request body")
		return
	}

	password, err := dbHelper.GetUserPasswordByEmail(reqBody.Email)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Failed to find user")
		return
	}
	if !utils.ComparePasswords(password, reqBody.Password) {
		utils.RespondError(w, http.StatusUnauthorized, nil, "Incorrect password")
		return
	}

	userID, err := dbHelper.IsUserExist(reqBody.Email)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to find user")
		return
	}
	token := utils.HashString(fmt.Sprintf("%s-%s", reqBody.Email, time.Now().String()))
	if err := dbHelper.StoreUserSession(userID, token); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to store token")
		return
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	utils.RespondJSON(w, http.StatusOK, response)
}



func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.UserContext(r)
	utils.RespondJSON(w, http.StatusOK, userCtx)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	userCT := middleware.UserContext(r)
	userID := userCT.ID
	token := userCT.Token

	allSessions := chi.URLParam(r, "allSessions")

	if allSessions == "true" {
		if err := dbHelper.DeleteUserSessionAll(userID); err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err, "Failed to logout all sessions")
			return
		}
	} else {
		if err := dbHelper.DeleteUserSession(token, userID); err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err, "Failed to logout all sessions")
			return
		}

	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	//todo update user get ID from user context
	userCTX := middleware.UserContext(r)
	reqBody := struct {
		Phone    string `json:"Phone" db:"Phone"`
		Position string `json:"Position" db:"Position"`
	}{}
	if err := utils.ParseBody(r.Body, &reqBody); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to decode request body")
		return
	}
	if reqBody.Phone == "" {
		utils.RespondError(w, http.StatusInternalServerError, fmt.Errorf("invalid value of field Phone"), "Phone is empty")
		return
	}
	if reqBody.Position == "" {
		utils.RespondError(w, http.StatusInternalServerError, fmt.Errorf("invalid value of field Position"), "Position is empty")
		return
	}
	if err := dbHelper.UpdateUserInfo(userCTX.ID, reqBody.Phone, reqBody.Position); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "something went wrong")
		return
	}

	user, err := dbHelper.GetUserByToken(userCTX.Token)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "unable to fetch updated record")
		return
	}
	utils.RespondJSON(w, http.StatusOK, user)
}
// UploadProfileImage uploads profile image of the user to firebase and return image link
func UploadProfileImage(w http.ResponseWriter, r *http.Request) {
	userCTX := middleware.UserContext(r)
	file, handler, err := utils.FileFromRequest(r, "image")
	defer func() {
		if err=file.Close(); err != nil {
			logrus.Error(err)
		}
	}()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, err.Error())
		return
	}
	fileName, err := utils.UploadFile(file, handler.Filename, os.Getenv("FIREBASE_BUCKET_NAME"))
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, err.Error())
		return
	}
	if err := dbHelper.InsertImageAndUpdateUsers("profile", os.Getenv("FIREBASE_BUCKET_NAME"), fileName, userCTX.ID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, err.Error())
		return
	}

	url, err := utils.GenerateURL(fileName, os.Getenv("FIREBASE_BUCKET_NAME"), "GET", time.Now().Add(time.Minute*60))
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		URL string json:"url"
	}{
		url,
	})
}