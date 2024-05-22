package services

import (
	"errors"
	"os"
	"strconv"

	"github.com/Nesrux/api-enconder/application/repositories"
	"github.com/Nesrux/api-enconder/domain"
)

const (
	localStoragePath  = "localStoragePath"
	inputBucketName   = "inputBucketName"
	concurrency_upload= "concurrency_upload"
)

type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobRepository
	VideoService  VideoService
}

func (j *JobService) Start() error {
	err := j.changeJobStatus("DOWNLOADING")

	if err != nil {
		return j.failJob(err)
	}
	err = j.VideoService.Downlaod(os.Getenv(inputBucketName))

	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus("FRAGMENTING")

	if err != nil {
		return j.failJob(err)
	}
	err = j.VideoService.Fragment()

	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus("ENCONDING")
	if err != nil {
		return j.failJob(err)
	}
	err = j.VideoService.Enconde()

	if err != nil {
		return j.failJob(err)
	}
	err = j.performUpload()
	if err != nil {
		return j.failJob(err)
	}
	err = j.changeJobStatus("FiNISHING")
	if err != nil {
		return j.failJob(err)
	}
	err = j.VideoService.Finish()
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus("COMPLETED")
	if err != nil {
		return j.failJob(err)
	}

	return nil
}

func (j *JobService) performUpload() error {
	err := j.changeJobStatus("UPLOADING")

	if err != nil {
		return j.failJob(err)
	}

	videoUpload := NewVideoUpload()
	videoUpload.OutputBucket = os.Getenv(inputBucketName)
	videoUpload.VideoPath = os.Getenv(localStoragePath) + "/" + j.VideoService.Video.ID
	concurrency, _ := strconv.Atoi(os.Getenv(concurrency_upload))
	doneUpload := make(chan string)

	go videoUpload.ProcessUpload(concurrency, doneUpload)

	var uploadResult string
	uploadResult = <-doneUpload

	if uploadResult != "upload completed" {
		return j.failJob(errors.New(uploadResult))
	}

	return err
}

func (j *JobService) changeJobStatus(status string) error {
	var err error

	j.Job.Status = status
	j.Job, err = j.JobRepository.Update(j.Job)

	if err != nil {
		return j.failJob(err)
	}

	return nil

}
func (j *JobService) failJob(error error) error {
	j.Job.Status = "FAILED"
	j.Job.Error = error.Error()

	_, err := j.JobRepository.Update(j.Job)
	if err != nil {
		return err
	}
	return error
}
