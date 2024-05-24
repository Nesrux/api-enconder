package services_test

import (
	"log"
	"os"
	"testing"

	"github.com/Nesrux/api-enconder/application/services"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func Test_videoservice_upload(t *testing.T) {
	video, repo := prepare()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("bucket-encoder")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	videoUpload := services.NewVideoUpload()
	videoUpload.OutputBucket = "bucket-encoder"
	videoUpload.VideoPath = os.Getenv(services.LOCAL_STORAGE_PATH) + "/" + video.ID

	doneUpload := make(chan string)
	go videoUpload.ProcessUpload(30, doneUpload)

	result := <-doneUpload
	require.Equal(t, result, "upload completed")
}
