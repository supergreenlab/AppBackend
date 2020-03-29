package feeds

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/minio/minio-go/v6"
	"github.com/spf13/viper"
)

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

func uploadURLHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO check if media exists
	minioClient := createMinioClient()
	expiry := time.Second * 60 * 60
	presignedURL, err := minioClient.PresignedPutObject("feedmedias", fmt.Sprintf("%s-%s", p.ByName("id"), uuid.New().String()), expiry)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully generated presigned URL", presignedURL)
}
