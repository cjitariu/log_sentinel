# Log Sentinel

Log Sentinel is a Go application designed to process and analyze job execution logs stored in CSV format. It parses log entries, tracks job durations, and provides warnings or errors if job durations exceed specified thresholds.

## Features

- Loads and parses CSV log files.
- Tracks job start and end times, calculating durations.
- Logs warnings for jobs running longer than 5 minutes and errors for jobs running longer than 10 minutes.
- JSON logging using Go's `slog` package.
- Modular code structure for easy extension.
- Unit tests for core logic.

## Project Structure

```
.
├── main.go                # Application entry point
├── main_test.go           # Unit tests for main logic
├── output/
│   └── logs.json          # Processed log output
├── pkg/
│   ├── csvloader/
│   │   └── csvloader.go   # CSV loading utility
│   └── models/
│       └── models.go      # Data model definitions
├── resources/
│   └── logs.log           # Example/input log file
├── go.mod                 # Go module file
└── README.md              # Project documentation
```

## Log Format

Each line in the CSV log file should have the following columns:

```
Timestamp,Job Description,Event Type,PID
```

- **Timestamp**: Time of the event (e.g., `12:00:00`)
- **Job Description**: Name or description of the job (e.g., `JobA`)
- **Event Type**: `START` or `END`
- **PID**: Process ID (e.g., `1234`)

**Example:**
```
12:00:00,JobA,START,1234
12:07:00,JobA,END,1234
```

## Usage

1. Place your CSV log file at `resources/logs.log` using the format above.
2. Run the application:

   ```sh
   go mod tidy
   go run main.go
   ```
3. Build and run the application:

   ```sh
   go build -o log-sentinel
   ./log-sentinel
   ```
4. The application will output structured logs to stdout, only warnings and errors for long-running jobs.

## Configuration

- **Warning threshold:** 5 minutes (jobs running longer will trigger a warning)
- **Error threshold:** 10 minutes (jobs running longer will trigger an error)
- You can adjust these thresholds in `main.go` by changing the `warningThreshold` and `errorThreshold` constants.

## Testing

Run unit tests with:

```sh
go test
```


