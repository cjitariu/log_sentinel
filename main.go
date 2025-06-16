package main

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"log_sentinel/pkg/csvloader"
	"log_sentinel/pkg/models"
)

const (
	warningThreshold = 5 * time.Minute
	errorThreshold   = 10 * time.Minute
	filePath         = "resources/logs.log"
)

func main() {
	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// Load CSV file
	// Ensure the CSV file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logger.Error("CSV file does not exist", "file", filePath)
		return
	}
	// Load the CSV file
	csvLoader := csvloader.NewCSVLoader(filePath)
	records, err := csvLoader.Load()
	if err != nil {
		logger.Error("Error loading CSV", "error", err)
	}

	jobs := make(map[string]models.Job)
	for _, record := range records {
		// Assuming the CSV has the following columns: Timestamp , Job decription, Event type(start, end), Pid
		//Validate the record length
		if len(record) != 4 {
			logger.Error("Skipping record with insufficient fields", slog.String("record", strings.Join(record, ",")))
			continue
		}
		// Check if the record exists in the jobs map, meaning of the the events has already been processed
		if exists, ok := jobs[record[3]+"_"+strings.ReplaceAll(record[1], " ", "_")]; ok {
			if strings.TrimSpace(record[2]) == "END" {
				exists.EndTime, err = time.Parse("15:04:05", record[0])
				if err != nil {
					logger.Error("Error parsing end time:", "error", err.Error(), "record", strings.Join(record, ","))
					continue
				}
				exists.Duration, err = time.ParseDuration(exists.EndTime.Sub(exists.StartTime).String())
				if err != nil {
					logger.Error("Error calculating duration:", "error", err.Error(), "record", strings.Join(record, ","))
					continue
				}
				// Check the job duration
				checkJobDuration(exists, logger)

				// Update the job in the map

				jobs[record[3]+"_"+strings.ReplaceAll(record[1], " ", "_")] = exists
			} else {
				exists.StartTime, err = time.Parse("15:04:05", record[0])
				if err != nil {
					logger.Error("Error parsing start time:", "error", err.Error(), "record", strings.Join(record, ","))
					continue
				}
				// Check job duration
				checkJobDuration(exists, logger)
				// Update the job in the map
				jobs[record[3]+"_"+strings.ReplaceAll(record[1], " ", "_")] = exists
			}
		} else {
			//check the event type
			if strings.TrimSpace(record[2]) == "START" {
				startTime, err := time.Parse("15:04:05", record[0])
				if err != nil {
					logger.Error("Error parsing start time:", "error", err.Error(), "record", strings.Join(record, ","))
					continue
				}
				// Initialize the job with the start time
				jobs[record[3]+"_"+strings.ReplaceAll(record[1], " ", "_")] = models.Job{
					StartTime: startTime,
					Pid:       record[3],
					Name:      record[1],
				}
			} else if strings.TrimSpace(record[2]) == "END" {

				endTime, err := time.Parse("15:04:05", record[0])
				if err != nil {
					logger.Error("Error parsing end time:", "error", err.Error(), "record", strings.Join(record, ","))
					continue
				}
				// Initialize the job with the end time
				jobs[record[3]+"_"+strings.ReplaceAll(record[1], " ", "_")] = models.Job{
					EndTime: endTime,
					Pid:     record[3],
					Name:    record[1],
				}
			} else {
				logger.Error("Unknown event type in record:", "record", strings.Join(record, ","))
				continue

			}

		}
	}

	// Do something with the jobs slice
	//fmt.Println(jobs)
}

func checkJobDuration(job models.Job, logger *slog.Logger) {
	if job.Duration > warningThreshold && job.Duration < errorThreshold {
		// Log a warning if the job duration exceeds warning threshold
		// but is less than error threshold
		logger.Warn("Job duration exceeds warning threshold", "job", job, "minutes", job.Duration.Minutes())
	}
	if job.Duration > errorThreshold {
		// Log an error if the job duration exceeds error threshold
		logger.Error("Job duration exceeds error threshold", "job", job, "minutes", job.Duration.Minutes())
	}
}
