package client

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/fs1g17/MiniQ/internal/store"
)

type JobResponse struct {
	Job store.Job `json:"job"`
}

func PollForJob(serverURL string) {
	client := &http.Client{
		Timeout: 35 * time.Second,
	}
	var success bool = false
	var jobRequest *JobResponse

	defer func() {
		if r := recover(); r != nil {
			if jobRequest != nil {
				//something in the job failed
				NotifyJobEnd(serverURL, jobRequest.Job.ID, false)
			}
		}
	}()

	for {
		log.Printf("making request to %s\n", serverURL+"/pollJob")
		resp, err := client.Get(serverURL + "/pollJob")
		// if some error, retry
		if err != nil {
			log.Printf("Poll error: %v, retrying in 5s", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// if 204, retry
		if resp.StatusCode == http.StatusNoContent {
			log.Println("no content")
			resp.Body.Close()
			continue
		}

		if resp.StatusCode == http.StatusOK {
			//TODO: process job
			if err := json.NewDecoder(resp.Body).Decode(&jobRequest); err != nil {
				log.Printf("Decode error: %v", err)
			} else {
				log.Printf("Received: %+v", jobRequest)
				log.Println("working on job...")
				time.Sleep(5 * time.Second)
				log.Println("finished working on job")
				success = true
			}
			NotifyJobEnd(serverURL, jobRequest.Job.ID, success)
		}

		resp.Body.Close()
	}
}

type NotifyJobEndData struct {
	JobID   int  `json:"jobID"`
	Success bool `json:"success"`
}

func NotifyJobEnd(serverURL string, jobID int, success bool) {
	client := &http.Client{
		Timeout: 35 * time.Second,
	}

	url := serverURL + "/completeJob"

	data := NotifyJobEndData{JobID: jobID, Success: success}
	b, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to serialize data")
	}

	body := bytes.NewBuffer(b)
	client.Post(url, "application/json", body)
}
