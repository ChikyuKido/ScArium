package jobs

import (
	"ScArium/common/config"
	"ScArium/common/log"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type JobStatus string

const (
	Pending JobStatus = "pending"
	Running JobStatus = "running"
	Success JobStatus = "success"
	Failed  JobStatus = "failed"
)

type Job struct {
	ID             string
	Name           string
	AdditionalName string
	Params         map[string]string
	Status         JobStatus
	StartedAt      time.Time
	EndedAt        time.Time
	JobFunc        func(job Job) error
	JobFile        *os.File
}

var JobQueue = make(chan Job, 100)

func CreateJob(name string, additionalName string, params map[string]string, jobFunc func(job Job) error) error {
	job := Job{
		ID:             uuid.New().String(),
		Name:           name,
		AdditionalName: additionalName,
		Params:         params,
		Status:         Pending,
		StartedAt:      time.Time{},
		EndedAt:        time.Time{},
		JobFunc:        jobFunc,
	}

	err := job.createJobFile()
	if err != nil {
		return fmt.Errorf("failed to create job: %v", err)
	}
	JobQueue <- job

	return nil
}

func Worker() {
	for job := range JobQueue {
		job.Status = Running
		job.StartedAt = time.Now()
		job.AppendLog(logrus.InfoLevel, "Job started")

		err := job.JobFunc(job)
		if err != nil {
			job.Status = Failed
			job.AppendLog(logrus.ErrorLevel, "Job failed: "+err.Error())
		} else {
			job.Status = Success
			job.AppendLog(logrus.InfoLevel, "Job completed successfully")
		}
		job.EndedAt = time.Now()
		err = job.finalizeJob()
		if err != nil {
			log.E.Errorf("Failed to finalize job: %v", err)
		}
	}
}

func (job *Job) createJobFile() error {
	file, err := os.OpenFile(fmt.Sprintf("%s/%s.job", config.JOBS_PATH, job.ID), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	metadata := fmt.Sprintf("%s,%s,%s,%s,%d\n", job.ID, job.Name, job.AdditionalName, job.Params, time.Now().Unix())
	if _, err := file.WriteString(metadata); err != nil {
		return err
	}
	job.JobFile = file
	return nil
}

func (job *Job) AppendLog(level logrus.Level, logMsg string) {
	if _, err := job.JobFile.WriteString(fmt.Sprintf("[%s] %s", level.String(), logMsg) + "\n"); err != nil {
		log.I.Errorf("Could not write job log file: %v", err)
	}
}

func (job *Job) finalizeJob() error {
	metadata := fmt.Sprintf("%d,%d,%s\n", job.StartedAt.Unix(), job.EndedAt.Unix(), job.Status)
	if _, err := job.JobFile.WriteString(metadata + "\n"); err != nil {
		return err
	}
	return job.JobFile.Close()
}
