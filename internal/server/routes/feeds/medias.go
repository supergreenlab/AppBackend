package feeds

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/minio/minio-go/v6"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func setupBucket(name string) {
	minioClient := createMinioClient()
	exists, err := minioClient.BucketExists(name)
	if err != nil {
		log.Fatalln(err)
	}
	if exists {
		log.Printf("Already created bucket: %s\n", name)
		return
	}
	err = minioClient.MakeBucket(name, "")
	if err != nil {
		log.Fatalln(err)
	}
}

func initStorage() {
	setupBucket("feedmedias")
}

func createMinioClient() *minio.Client {
	accessKey := viper.GetString("S3AccessKey")
	secretKey := viper.GetString("S3SecretKey")
	minioClient, err := minio.New("minio:9000", accessKey, secretKey, false)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return minioClient
}

type feedMediaUploadURLParams struct {
	FileName string `json:"fileName"`
}

type feedMediaUploadURLResult struct {
	FilePath      string `json:"filePath"`
	ThumbnailPath string `json:"thumbnailPath"`
}

func feedMediaUploadURLHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmup := feedMediaUploadURLParams{}
	if err := decodeJSONBody(w, r, &fmup); err != nil {
		logrus.Errorln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path := ""
	if strings.HasSuffix(fmup.FileName, ".mp4") {
		path = fmt.Sprintf("videos-%s.mp4", uuid.New().String())
	} else if strings.HasSuffix(fmup.FileName, ".jpg") {
		path = fmt.Sprintf("pictures-%s.jpg", uuid.New().String())
	} else {
		http.Error(w, "Unknown file type", http.StatusBadRequest)
		return
	}

	res := feedMediaUploadURLResult{}
	minioClient := createMinioClient()
	expiry := time.Second * 60 * 60

	url1, err := minioClient.PresignedPutObject("feedmedias", path, expiry)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res.FilePath = url1.RequestURI()

	path = fmt.Sprintf("thumbnail-%s.jpg", uuid.New().String())
	url2, err := minioClient.PresignedPutObject("feedmedias", path, expiry)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res.ThumbnailPath = url2.RequestURI()

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
