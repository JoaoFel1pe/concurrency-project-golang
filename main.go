package main

import (
	"concurrency/workerpool"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	// allocate jobs
	noOfJobs := 3000
	go workerpool.AllocateJobs(noOfJobs) // create a goroutine to execute the allocateJobs function that sends jobs to the jobs communication channel

	// get results
	done := make(chan bool)        // create a communication channel to receive a boolean value (done) which will be used to finalize the execution of the goroutine that will be created later to receive the job results
	go workerpool.GetResults(done) // create a goroutine to execute the getResults function that receives the job results

	// create worker pool
	noOfWorkers := 1000
	workerpool.CreateWorkerPool(noOfWorkers) // create a goroutine to execute the createWorkerPool function that creates goroutines to execute the jobs received from the jobs communication channel

	<-done // wait for the done communication channel to receive true

	data, err := json.MarshalIndent(workerpool.ResultCollection, "", "  ") // encode the resultCollection slice to JSON format

	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}

	err = writeToFile(data)
	if err != nil {
		log.Fatalf("Error writing to file: %s", err)
	}
}

func writeToFile(data []byte) error {
	f, err := os.Create("xkcd.json") // create a file named xkcd.json
	if err != nil {
		return fmt.Errorf("create file: %v", err)
	}

	defer f.Close() // close the created file

	_, err = f.Write(data) // write the encoded result slice to the created file
	if err != nil {
		return fmt.Errorf("write file: %v", err)
	}

	return nil
}
