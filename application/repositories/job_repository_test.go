package repositories_test

import (
	"testing"

	"github.com/Nesrux/api-enconder/application/repositories"
	"github.com/Nesrux/api-enconder/domain"
	"github.com/Nesrux/api-enconder/framework/database"
	"github.com/stretchr/testify/require"
)

func Test_JobRepository_DbInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.FilePath = "path"

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	job, err := domain.NewJob("output_path", "Pending", video)
	require.Nil(t, err)

	repoJob := repositories.JobRepositoryDb{Db: db}
	repoJob.Insert(job)

	j, err := repoJob.Find(job.ID)
	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, j.ID, job.ID)
	require.Equal(t, j.VideoID, video.ID)

}
