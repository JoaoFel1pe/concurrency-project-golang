package workerpool

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var jobs = make(chan Job, 100)        // create a communication channel for sending jobs
var Results = make(chan Result, 100) // create a communication channel for receiving job Results
var ResultCollection []Result        // create a slice to store the job Results

const Url = "https://xkcd.com"

type Result struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	Safe_title string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
}

type Job struct {
	number int
}

func Fetch(n int) (*Result, error) {

	client := &http.Client{
		Timeout: 5 * time.Minute, // create an HTTP client with a timeout of 5 minutes for the HTTP request to be made
	}

	url := strings.Join([]string{Url, fmt.Sprintf("%d", n), "info.0.json"}, "/") // construct the URL for the HTTP request with the comic number and the response format (json) at the end of the URL (e.g., https://xkcd.com/614/info.0.json)

	req, err := http.NewRequest("GET", url, nil) // create an HTTP GET request for the constructed URL with no request body (nil)

	// check for any errors during the creation of the HTTP request
	if err != nil {
		return nil, fmt.Errorf("http request: %v", err)
	}

	resp, err := client.Do(req) // execute the created HTTP request

	if err != nil {
		return nil, fmt.Errorf("http request: %v", err)
	}

	var data Result

	if resp.StatusCode != http.StatusOK { // check if the HTTP response status code is different from 200

		data = Result{} // if it's different from 200, initialize the previously created Result struct with empty values

	} else {
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil { // decode the body of the HTTP response into the previously created Result struct
			return nil, fmt.Errorf("decode: %v", err)
		}
	}

	resp.Body.Close() // close the HTTP response body

	return &data, nil // return the previously created Result struct
}
func CreateWorkerPool(noOfWorkers int) {
	var wg sync.WaitGroup // create a WaitGroup to control the execution of goroutines

	for i := 0; i <= noOfWorkers; i++ {
		wg.Add(1)      // increase the WaitGroup counter by 1
		go Worker(&wg) // create a goroutine to execute the worker function
	}

	wg.Wait()      // wait for all goroutines to finish
	close(Results) // close the Results communication channel
}

func Worker(wg *sync.WaitGroup) { // function to be executed by the goroutines that will be created later to execute the jobs received from the jobs communication channel

	for job := range jobs { // receive jobs from the jobs communication channel

		Result, err := Fetch(job.number) // fetch the comic with the received job number
		if err != nil {
			log.Printf("error in fetching: %v\n", err)
		}

		Results <- *Result // send the fetched Result to the Results communication channel
	}

	wg.Done() // decrease the WaitGroup counter by 1
}

func AllocateJobs(noOfJobs int) {

	for i := 1; i <= noOfJobs; i++ { // send jobs to the jobs communication channel
		jobs <- Job{i}
	}

	close(jobs) // close the jobs communication channel
}

func GetResults(done chan bool) { // function to be executed by the goroutine that will be created later to receive the job Results

	for Result := range Results {
		if Result.Num != 0 { // check if the comic number is different from 0
			fmt.Printf("Retrieving issue #%d\n", Result.Num)    // print the comic number being fetched
			ResultCollection = append(ResultCollection, Result) // add the fetched Result to the ResultCollection slice
		}
	}

	done <- true // send true to the done communication channel
}
