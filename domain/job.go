package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Job struct {
	ID               string    `valid:"uuid"`
	OutputBucketPath string    `valid:"notnull"`
	Status           string    `valid:"notnull"`
	Video            *Video    `valid:"-"`
	VideoId          string    `valid:"-"`
	Error            string    `valid:"-"`
	CreatedAt        time.Time `valid:"notnull"`
	UpdatedAt        time.Time `valid:"notnull"`
}

func NewJob(output, status string, video *Video) (*Job, error) {
	job := Job{
		OutputBucketPath: output,
		Status:           status,
		Video:            video,
	}
	job.prepare()

	err := job.Validate()

	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (j *Job) prepare() {
	j.ID = uuid.NewString()
	j.CreatedAt = time.Now()
	j.UpdatedAt = time.Now()
}

func (j *Job) Validate() error {
	_, err := govalidator.ValidateStruct(j)
	if err != nil {
		return err
	}
	return nil
}
