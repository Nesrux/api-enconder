package services_test

import (
	"log"
	"testing"

	"github.com/Nesrux/api-enconder/application/repositories"
	"github.com/Nesrux/api-enconder/application/services"
	"github.com/Nesrux/api-enconder/domain"
	"github.com/Nesrux/api-enconder/framework/database"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func Test_VideoService_download(t *testing.T) {
	video, repo := prepare()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Downlaod("bucket-encoder")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

}

func prepare() (*domain.Video, repositories.VideoRepositoryDb) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.FilePath = "videotest.mp4"

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	return video, repo
}
