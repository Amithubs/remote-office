package utils

import (
	"cloud.google.com/go/storage"
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

var generator *shortid.Shortid

type clientError struct {
	ID            string `json:"id"`
	MessageToUser string `json:"messageToUser"`
	DeveloperInfo string `json:"developerInfo"`
	Err           string `json:"error"`
	StatusCode    int    `json:"statusCode"`
	IsClientError bool   `json:"isClientError"`
}

func init() {
	g, err := shortid.New(1, shortid.DefaultABC, rand.Uint64())
	if err != nil {
		logrus.Panicf("Failed to initialize utils package with error: %+v", err)
	}
	generator = g
}

// ParseBody parses the values from io reader to a given interface
func ParseBody(body io.Reader, out interface{}) error {
	err := json.NewDecoder(body).Decode(out)
	if err != nil {
		return err
	}
	return nil
}

// EncodeJSONBody writes the JSON body to response writer
func EncodeJSONBody(resp http.ResponseWriter, data interface{}) error {
	return json.NewEncoder(resp).Encode(data)
}

// RespondJSON sends the interface as a JSON
func RespondJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.WriteHeader(statusCode)
	if body != nil {
		if err := EncodeJSONBody(w, body); err != nil {
			logrus.Errorf("Failed to respond JSON with error: %+v", err)
		}
	}
}

// newClientError creates structured client error response message
func newClientError(err error, statusCode int, messageToUser string, additionalInfoForDevs ...string) *clientError {
	additionalInfoJoined := strings.Join(additionalInfoForDevs, "\n")
	if len(additionalInfoJoined) == 0 {
		additionalInfoJoined = messageToUser
	}

	errorID, _ := generator.Generate()
	var errString string
	if err != nil {
		errString = err.Error()
	}
	return &clientError{
		ID:            errorID,
		MessageToUser: messageToUser,
		DeveloperInfo: additionalInfoJoined,
		Err:           errString,
		StatusCode:    statusCode,
		IsClientError: true,
	}
}

// RespondError sends an error message to the API caller and logs the error
func RespondError(w http.ResponseWriter, statusCode int, err error, messageToUser string, additionalInfoForDevs ...string) {
	logrus.Errorf("status: %d, message: %s", statusCode, messageToUser)
	clientError := newClientError(err, statusCode, messageToUser, additionalInfoForDevs...)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(clientError); err != nil {
		logrus.Errorf("Failed to send error to caller with error: %+v", err)
	}
}

// HashAndSaltPassword create a hash for given password
func HashAndSaltPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePasswords compares a hashed password with plain text password
func ComparePasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPwd))
	if err != nil {
		logrus.Errorf("Failed to compare password with error: %+v", err)
		return false
	}

	return true
}

// HashString generates SHA256 for a given string
func HashString(toHash string) string {
	sha := sha512.New()
	sha.Write([]byte(toHash))
	return hex.EncodeToString(sha.Sum(nil))
}
func FileFromRequest(r *http.Request, key string) (multipart.File, *multipart.FileHeader, error) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return nil, nil, err
	}
	file, handler, err := r.FormFile(key)
	if err != nil {
		return nil, nil, err
	}
	return file, handler, err
}
// UploadFile upload the given multipart file to firebase storage with the given name
func UploadFile(file multipart.File, filename, bucket string) (string, error) {

	u := uuid.New()
	uniqueFilename := fmt.Sprintf("%s-%s", u, filename)

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(os.Getenv("FIREBASE_CONFIG_FILE_PATH")))
	if err != nil {
		return "", err
	}
	defer func() {
		if err = client.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	wc := client.Bucket(bucket).Object(uniqueFilename).NewWriter(ctx)
	if _, err = io.Copy(wc, file); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}
	return uniqueFilename, nil
}
// GenerateURL return signed url for the given file name
func GenerateURL(filename, bucket, method string, expires time.Time) (string, error) {

	configFile, err := ioutil.ReadFile(os.Getenv("FIREBASE_CONFIG_FILE_PATH"))
	if err != nil {
		return "", err
	}

	cfg, err := google.JWTConfigFromJSON(configFile)
	if err != nil {
		return "", err
	}

	url, err := storage.SignedURL(bucket, filename, &storage.SignedURLOptions{
		GoogleAccessID: cfg.Email,
		PrivateKey:     cfg.PrivateKey,
		Method:         method,
		Expires:        expires,
	})
	if err != nil {
		return "", err
	}
	return url, nil
}