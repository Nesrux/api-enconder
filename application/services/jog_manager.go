package services

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/Nesrux/api-enconder/application/repositories"
	"github.com/Nesrux/api-enconder/domain"
	"github.com/Nesrux/api-enconder/framework/queue"
	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
)

type JobManager struct {
	Db               *gorm.DB
	Domain           domain.Job
	MessageChannel   chan amqp.Delivery
	JobReturnChannel chan JobWorkerResult
	RabbitMQ         *queue.RabbitMQ
}

type JobNotification struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewJobManager(db *gorm.DB, rabbitMq *queue.RabbitMQ, jobReturnChan chan JobWorkerResult,
	messageChannel chan amqp.Delivery) *JobManager {
	return &JobManager{
		Db:               db,
		Domain:           domain.Job{},
		MessageChannel:   messageChannel,
		JobReturnChannel: jobReturnChan,
		RabbitMQ:         rabbitMq,
	}
}

func (j *JobManager) Start(ch *amqp.Channel) {
	VideoService := NewVideoService()
	VideoService.VideoRepository = repositories.VideoRepositoryDb{Db: j.Db}

	jobService := JobService{
		JobRepository: repositories.JobRepositoryDb{Db: j.Db},
		VideoService:  VideoService,
	}
	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY_WORKERS"))
	if err != nil {
		log.Fatalf("error loading var: Concurrency_Workers")
	}

	for qtdProcess := 0; qtdProcess < concurrency; qtdProcess++ {
		go JobWorker(j.MessageChannel, j.JobReturnChannel, jobService, j.Domain, qtdProcess)
	}

	for jobResult := range j.JobReturnChannel {
		if jobResult.Error != nil {
			err = j.chackParseErrors(jobResult)
		} else {
			err = j.notifySuccess(jobResult, ch)
		}
		if err != nil {
			 jobResult.Message.Reject(false)
		}
	}

}

func (j *JobManager) chackParseErrors(jobResult JobWorkerResult) error {
	if jobResult.Job.ID != "" {
		log.Printf("MessageId #{jobResult.Message.DeliveryTag}. Error parsing job #{jobResult.Job.ID}")
	} else {
		log.Printf("MessageId #{jobResult.Message.DeliveryTag}. Error parsing message #{jobResult.Job.ID}")
	}
	errorMessage := JobNotification{
		Message: string(jobResult.Message.Body),
		Error:   jobResult.Error.Error(),
	}
	jobJson, err := json.Marshal(errorMessage)

	return nil
}

func (j *JobManager) notifySuccess(){
	
}
