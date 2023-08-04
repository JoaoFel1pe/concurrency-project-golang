# XKCD Comic Scraper

This is a Go program that scrapes XKCD comics and saves the data in JSON format. It fetches comic information from the XKCD website and stores it in a JSON file for further analysis or use.

## How it works

1. The program starts by allocating a set of jobs to be done. It sends job numbers to the `jobs` channel.

2. It then creates a worker pool to concurrently process the jobs. The number of workers can be adjusted by changing the `noOfWorkers` constant.

3. Each worker fetches a comic using the provided job number. The comics are fetched from XKCD's website in JSON format.

4. The fetched comic data is stored in the `results` channel for further processing.

5. The results are collected and stored in a slice called `resultCollection`.

6. Once all the jobs are completed, the program encodes the `resultCollection` slice to JSON format and writes it to a file named `xkcd.json`.

## Dependencies

This program uses the following dependencies:

- `net/http` for making HTTP requests
- `encoding/json` for encoding and decoding JSON data
- `sync` for handling synchronization between goroutines
- `log` for logging errors
- `os` for file operations

## How to use

1. Make sure you have Go installed on your machine.

2. Clone the repository or copy the code into a file named `main.go`.

3. Run the program using the following command:

   ```bash
   go run main.go
   ```

4. The program will start fetching XKCD comics concurrently and save the data to `xkcd.json` in the same directory.

5. The program will output the status of each comic being fetched, including the comic number.

## Notes

- The program is designed to fetch 3000 comics concurrently. You can adjust the `noOfJobs` and `noOfWorkers` constants based on your requirements.

- Be mindful of the XKCD website's terms of use when scraping data. Frequent and aggressive scraping might violate their policies.

- Remember to handle errors properly and implement error checking as needed, especially if you plan to use this code in a production environment.