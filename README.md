![GitHub release](https://img.shields.io/github/v/release/Tagliapietra96/logger)
[![Go Reference](https://pkg.go.dev/badge/Tagliapietra96/logger/path.svg)](https://pkg.go.dev/github.com/Tagliapietra96/logger)
[![Go Report Card](https://goreportcard.com/badge/github.com/Tagliapietra96/logger)](https://goreportcard.com/report/github.com/Tagliapietra96/logger)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

# Table of Contents
1. [Logger - A Lightweight Logging System with Advanced Features](#logger---a-lightweight-logging-system-with-advanced-features)
2. [Features](#features)
   - [SQLite-Based Logging](#sqlite-based-logging)
   - [Flexible Filtering and Configurable Terminal Output](#flexible-filtering-and-configurable-terminal-output)
   - [Export Logs to Multiple Formats](#export-logs-to-multiple-formats)
   - [Alert Log System](#alert-log-system)
3. [Why Choose Logger?](#why-choose-logger)
4. [Usage](#usage)
   - [Install the Package](#install-the-package)
   - [Basic Usage](#basic-usage)
   - [Advanced Configuration](#advanced-configuration)
     - [Setting the Log Storage Folder](#setting-the-log-storage-folder)
     - [Configuring Log Output Format (Inline vs Block)](#configuring-log-output-format-inline-vs-block)
     - [Customizing Caller Information Display](#customizing-caller-information-display)
     - [Configuring Timestamp Display](#configuring-timestamp-display)
     - [Managing Tags for Logs](#managing-tags-for-logs)
     - [Configuring Fatal Notifications](#configuring-fatal-notifications)
     - [Creating a Copy of the Logger Configuration](#creating-a-copy-of-the-logger-configuration)
5. [Log Management Functionality](#log-management-functionality)
   - [Saving Logs to the Database](#saving-logs-to-the-database)
   - [Printing Logs Directly to the Console (Without Persistence)](#printing-logs-directly-to-the-console-without-persistence)
   - [Printing Logs from the Database](#printing-logs-from-the-database)
6. [Export Functionality](#export-functionality)
   - [Key Features](#key-features)
   - [Example Usage](#example-usage)
   - [Use Cases](#use-cases)
7. [Conclusion](#conclusion)
8. [Important Note ⚠️](#important-note-️)
9. [Found a Bug? Have a Suggestion? 💡](#found-a-bug-have-a-suggestion-)
10. [Acknowledgements 🙌](#acknowledgements-)
11. [Support Logger ❤️](#support-logger-️)
12. [License](#license)

---

# Logger - A Lightweight Logging System with Advanced Features
Logger is a Go package that provides a versatile and efficient logging system designed for applications that require structured log management. With its built-in SQLite integration, filtering capabilities, customizable terminal output, and alert notifications, Logger offers a comprehensive solution for tracking and managing application events.
## Features
### SQLite-Based Logging
Logs are stored in a lightweight SQLite database, ensuring persistence and easy retrieval without relying on external systems.
### Flexible Filtering and Configurable Terminal Output
Quickly filter logs by various criteria (e.g., log levels, timestamps, tags) and print them directly to the terminal for debugging or analysis. Using the powerful [`charmbracelet/lipgloss`](https://github.com/charmbracelet/lipgloss) package, Logger provides styled terminal output with predefined, visually appealing formats. You can customize the visibility and detail level of key log components, such as:
- **Caller Information:** Choose whether to display or hide details like file names, line numbers, and function names.
- **Timestamps:** Include or exclude timestamps or adjust their level of precision.
- **Tags:** Display all tags or omit them entirely for cleaner output.

Logs can be printed in two distinct formats:
- **Inline Mode:** Logs are printed as compact, single-line entries for efficient viewing of multiple logs at once.
- **Block Mode:** Each log is displayed as a "card," with detailed, structured output that visually separates it from other logs for enhanced readability.

This flexible configuration ensures that the log output is tailored to the specific needs of the user, whether they require concise summaries or detailed insights.
### Export Logs to Multiple Formats
Easily export log data from the SQLite database in various formats:

- **JSON:** Structured and machine-readable for integration with other tools.
- **CSV:** Convenient for data analysis and spreadsheet manipulation.
- **LOG:** Traditional plain text format for compatibility with other log systems.
### Alert Log System
Create "alert logs" that trigger real-time notifications for critical events, allowing immediate response to important issues. Alerts are powered by the [`gen2brain/beeep`](https://github.com/gen2brain/beeep) package, which provides cross-platform notifications via native system alerts, ensuring you never miss a critical event.


## Why Choose Logger?
Logger combines simplicity with powerful features, making it ideal for developers who want a self-contained logging solution that goes beyond traditional file-based logging. Whether you need detailed logs for debugging, instant alerts for critical events, or stylish terminal output, Logger has you covered.
With Logger, you get a complete, lightweight logging solution that is visually engaging and ready to meet your application's logging needs.

---


## Usage
### Install the Package
Get the `logger` package via `go get`:
```bash
go get github.com/Tagliapietra96/logger
```

### Basic Usage
Create and configure a basic logger:
```go
package main

import "github.com/Tagliapietra96/logger"

func main() {
    // Initialize a logger with the default configuration
    log := logger.New() 

    // Initialize a logger with default tags 'tag1' and 'tag2'
    // All logs created with this logger will include these tags
    logWithTags := logger.New("tag1", "tag2") 
}
```

### Advanced Configuration
Logger provides several options to customize how logs are stored, formatted, and displayed. Below are detailed configurations to tailor the logger to your needs.


#### Setting the Log Storage Folder
By default, the logs are stored in an SQLite database in the binary directory (if the apllication isn't built the default is the current working directory). You can change the storage location with the `Folder` method:
```go
log := logger.New()

// Set the folder where the SQLite database will be stored
log.Folder("~/projects/my-logs/")
```
> **Note:** Ensure the specified folder exists and has the necessary write permissions.


#### Configuring Log Output Format (Inline vs Block)
You can control how logs are printed to the terminal. Logs can be displayed in a compact, single-line format (`inline`), or in a more detailed, block format, where each log is presented as a card-like entry (`block`).

```go
// Print logs in a single-line format
log.Inline(true)  

// Print logs in a block format (default)
log.Inline(false)  
```
> - **Inline Mode:** Suitable for quick, concise debugging.
> - **Block Mode:** Ideal for comprehensive, formatted log displays with better readability.


#### Customizing Caller Information Display
Control how much information about the function calling the logger is shown. You can hide it completely, or display varying levels of detail:

```go
// Hide all caller information
log.Caller(logger.HideCaller)

// Show only the caller's file name (default behavior)
log.Caller(logger.ShowCallerFile)

// Show the file name and the line number where the log was called
log.Caller(logger.ShowCallerLine)

// Show file name, line number, and function name for precise tracking
log.Caller(logger.ShowCallerFunction)
```
> **Use Case:** Showing the caller details is useful for debugging complex applications where knowing the exact source of a log is critical.


#### Configuring Timestamp Display
Decide how much timestamp information you want in your logs. You can hide it entirely or choose from different levels of detail:

```go
// Hide timestamp information entirely
log.Timestamp(logger.HideTimestamp)

// Display only the date in the format "YYYY-MM-DD"
log.Timestamp(logger.ShowDate)  

// Display date and time in the format "YYYY-MM-DD HH:MM:SS" (default)
log.Timestamp(logger.ShowDateTime)

// Display the full timestamp with the day of the week included
log.Timestamp(logger.ShowFullTimestamp)
```
> - **Default Format:** `2006-01-02 15:04:05`
> - **Full Timestamp Example:** `Monday 2006-01-02 15:04:05`


#### Managing Tags for Logs
Tags help categorize logs, making it easier to filter and search. You can add or remove tags dynamically.

```go
// Add tags to the logger instance
log.Tags("tag1", "tag2")

// Add another tag to the existing list
log.Tags("tag3")  // Now the tags are: []string{"tag1", "tag2", "tag3"}

// Override existing tags with a new set
log.SetTags("new-tag1", "new-tag2")  

// Clear all tags
log.SetTags()  // Now the logger has no tags
```


#### Configuring Fatal Notifications
Customize the message and title for critical errors using the `SetFatal` method. This is particularly useful for displaying user-friendly or context-specific messages.

```go
// Set a custom title and message for fatal error notifications
log.SetFatal("MyApp - CRITICAL ERROR", "Oops! Something went wrong. Check the logs.")
```
> - **Default Title:** `"Fatal"`
> - **Default Message:** `"An error occurred, please check the logs for more information"`


These configurations allow you to fully customize your logging setup, ensuring that the logs are both informative and easy to manage, while maintaining flexibility in how they are displayed and stored.

#### Creating a Copy of the Logger Configuration
Logger allows you to create an exact copy of the current logger instance using the `Copy` method. This is useful when you want to reuse the same logging configuration but apply it with different tags or additional settings without modifying the original logger.

```go
log := logger.New()

// Customize the original logger configuration
log.Folder("~/projects/my-logs/")
log.Tags("initial", "setup")
log.Caller(logger.ShowCallerFile)

// Create a copy of the existing logger
log2 := log.Copy()

// Modify the copied logger independently
log2.Tags("new-tag") // Now log2 has ["initial", "setup", "new-tag"]
log2.Caller(logger.HideCaller) // Different caller visibility from original
```
> #### Use Cases:
> - **Independent Logging for Different Contexts:** Use the same core configuration but apply different tags for logging in different parts of your application (e.g., "auth", "database").
> - **Debugging Different Modules Separately:** Keep the original logger for general application logs and use the copied instance to focus on specific modules without altering the primary configuration.
> - **Isolated Fatal Handling:** Customize fatal notifications independently for various modules while retaining a shared base configuration.
> 
> The `Copy` feature enhances flexibility by enabling modular and context-aware logging configurations while maintaining a consistent base setup across different components of your application.


## Log Management Functionality
Logger provides three primary ways to manage logs: saving them to the SQLite database, printing them directly to the console without persistence, and retrieving and printing existing logs from the database. This section details these functionalities, offering examples and explanations for each.

### Saving Logs to the Database
Logs can be saved to the database using various log levels such as `Debug`, `Info`, `Warn`, `Error`, and `Fatal`. Each method formats the provided message and stores it in the SQLite database with optional tags and metadata. The `Fatal` log type also triggers an alert and exits the program.

#### Example Usage:

```go
package main
import (
    "github.com/Tagliapietra96/logger"
    
    "time"
)
func main() {
    log := logger.New()
    msg := "initializing components"
    // Create different log types and save them to the database
    log.Debug("Debug message: %s", msg)
    log.Info("App started at: %s", time.Now().Format("2006-01-02 15:04:05"))
    log.Error("oh no! an error")
    err := myFunc()
    if err != nil {
        // Fatal log - triggers alert and exits
        log.Fatal(err)
    }
    err := log.Warn("Potential issue detected: low disk space")
    if err != nil {
        panic(err)
    }
}
```

#### Key Details:
- **Formatting:** Uses `fmt.Sprintf` for message formatting.
- **Persistence:** Logs are stored in the SQLite database.
- **Error Handling:** Each method returns an error if log creation fails.
- **Alerts:** `Fatal` logs trigger alerts using the `beeep` package and terminate the application.

#### Use Cases:
- **Tracking critical system events** with persistent logs.
- **Debugging issues** by storing logs for later retrieval and analysis.
- **Immediate notification** of unrecoverable errors.


### Printing Logs Directly to the Console (Without Persistence)

For real-time feedback, logs can be printed directly to the terminal using `PrintDebug`, `PrintInfo`, `PrintWarn`, `PrintError`, and `PrintFatal`. These logs are not saved in the database.

#### Example Usage:

```go
package main

import (
    "github.com/Tagliapietra96/logger"
    
    "fmt"
)

func main() {
    log := logger.New()

    // Print logs directly without saving them to the database
    log.PrintDebug("Debugging: %s", "initializing cache")
    log.PrintInfo("Starting process: %s", "data sync")
    log.PrintWarn("Warning: %s", "deprecated API usage")
    
    err := fmt.Errorf("network timeout")
    log.PrintError("Error: %v", err)
    
    // Print fatal log and exit
    log.PrintFatal(fmt.Errorf("Fatal error: %s", "database connection lost"))
}
```

#### Key Details:
- **Real-Time Output:** Logs are printed directly to the terminal.
- **No Persistence:** These logs are not stored in the database.
- **Exit on Fatal:** PrintFatal logs terminate the program after printing.

#### Use Cases:
- **Debugging in development** environments where persistence is unnecessary.
- **Immediate visibility** for application events during runtime.
- **Custom error reporting** without cluttering the database.


### Printing Logs from the Database
Logs stored in the database can be queried and printed using `PrintLogs`. This method supports query options to filter logs based on criteria like level, tags, or date range.

#### Example Usage:

```go
package main

import (
	"github.com/Tagliapietra96/logger"
	"github.com/Tagliapietra96/logger/queries"
)

func main() {
    log := logger.New()

    // Print all logs stored in the database
    log.PrintLogs()

    // Print only "Error" level logs
    log.PrintLogs(queries.LevelEqual(logger.Error))

    // Print only 10 logs
    log.PrintLogs(queries.AddLimit(10))

    // Print the last 10 logs with "Error" level
    log.PrintLogs(
        queries.LevelEqual(logger.Error),
        queries.AddLimit(10),
        queries.SortTimestamp("desc")
    )
}
```

#### Key Details
- **Flexible Querying:** Use `QueryOption` to filter logs by level, tags, or date range. The package also includes the sub-package `github.com/Tagliapietra96/logger/queries`, which provides a comprehensive list of ready-to-use `QueryOption` instances that cover most common use cases, simplifying complex query creation.
- **Database Retrieval:** Retrieves logs from SQLite and prints them in the configured format, ensuring consistency between stored and displayed data.
- **Inline or Block Output:** Logs can be printed inline (single-line log entries) or in block format (each log displayed as a separate card-like structure), allowing for flexible presentation.

#### Use Cases
- **Postmortem Analysis:** Retrieve logs after a critical failure to perform in-depth analysis and identify root causes. The `queries` sub-package can streamline complex queries for these scenarios.
- **Selective Viewing:** View logs filtered by level, tags, or other criteria for efficient debugging and issue tracking. Use predefined `QueryOption` from the `queries` sub-package to expedite setup.
- **Audit Trail Creation:** Query logs from specific date ranges to create detailed audit trails. The `queries` package provides options to filter by time windows or log levels for compliance and reporting.


## Export Functionality
The `Export` method in the `logger` package provides a powerful way to export logs from the SQLite database to different file formats. This feature supports exporting logs based on specified query options, offering flexibility for different use cases such as data analysis, archival, or reporting.

### Key Features

#### Export Formats:
The `ExportType` parameter allows exporting logs in one of the following formats:
- **LOG:** Exports logs in a plain `.log` text file.
- **JSON:** Exports logs in a structured `.json` file, ideal for further data processing or integration with other systems.
- **CSV:** Exports logs in a `.csv` file, suitable for importing into spreadsheet software or databases for analysis.

#### Customizable Queries:
The method accepts `QueryOption` parameters, enabling fine-grained control over which logs to export. You can filter by log level, tags, date ranges, or other criteria. Leverage the `github.com/Tagliapietra96/logger/queries` sub-package for ready-to-use query options.

#### Return Values:
- **File Path:** The method returns the full path to the exported file.
- **Error Handling:** If the export fails, it returns an error describing the issue.

### Example Usage
```go
package main

import (
	"fmt"
	"github.com/Tagliapietra96/logger"
	"github.com/Tagliapietra96/logger/queries"
)

func main() {
	log := logger.New()

	// Export logs as a JSON file filtered by level and tags
	filePath, err := log.Export(logger.JSON, queries.LevelEqual(logger.Info), queries.HasTags("auth"))
	if err != nil {
		fmt.Println("Error exporting logs:", err)
		return
	}
	fmt.Println("Logs exported to:", filePath)

	// Export logs within a specific date range in CSV format
	filePath, err = log.Export(logger.CSV, queries.DateBetween(time.Now().Add(-24*time.Hour), time.Now()))
	if err != nil {
		fmt.Println("Error exporting logs:", err)
		return
	}
	fmt.Println("Logs exported to:", filePath)
}
```

#### Use Cases

- **Backup and Archival:** Periodically export logs for long-term storage in .log or .csv format.
- **Data Integration:** Export logs in .json format for integration with external tools such as ELK Stack or custom data pipelines.
- **Auditing and Reporting:** Generate .csv exports filtered by date range or tags to create detailed audit trails or compliance reports.

## Conclusion
Thank you for exploring **Logger**, a lightweight yet powerful logging system designed to simplify log management for CLI applications. With its user-friendly API, flexible configuration options, and seamless SQLite integration, Logger helps keep your logs organized and accessible. Whether you're building a small utility or a robust command-line tool, Logger offers the essential features to track and analyze application events effectively.

### Important Note ⚠️
Please note that **Logger** has been thoroughly tested only in **CLI applications** running on **Unix-based systems** (Linux, macOS, etc.). While it may work in other environments—such as GUI applications or non-Unix platforms—there’s no guarantee of full compatibility or consistent behavior in these contexts. If you use Logger outside its tested environments, please proceed with caution, and feel free to share your experience.

### Found a Bug? Have a Suggestion? 💡
I am committed to making Logger as reliable and useful as possible, and your feedback is invaluable! If you encounter any bugs, unexpected behavior, or have ideas for new features, please open an **issue** on the [GitHub repository](https://github.com/Tagliapietra96/logger/issues). I’ll do my best to address your concerns and improve the library based on your input.

### Acknowledgements 🙌

Logger wouldn’t be what it is without the amazing contributions of these projects:

- [Charmbracelet Lipgloss](https://github.com/charmbracelet/lipgloss): Provides elegant terminal styling, making the log output visually appealing. (Licensed under the MIT License. See [`THIRD_PARTY_LICENSES.md`](https://github.com/Tagliapietra96/logger/blob/main/THIRD_PARTY_LICENSES.md) for more details.)
- [Gen2brain Beeep](https://github.com/gen2brain/beeep): Enables desktop notifications to ensure you never miss critical log events. (Licensed under the BSD 2-Clause License. See [`THIRD_PARTY_LICENSES.md`](https://github.com/Tagliapietra96/logger/blob/main/THIRD_PARTY_LICENSES.md) for more details.)

I encourage you to visit their repositories, explore their features, and, if you find them useful, consider giving them a ⭐️ to show your support.

### Support Logger ❤️
If Logger has been helpful in your project, I’d be incredibly grateful if you could leave a **star** ⭐️ on the [GitHub repository](https://github.com/Tagliapietra96/logger). Your support means a lot and helps other developers discover and benefit from Logger.

Thank you for choosing Logger, and happy logging! 🚀

**License**
[MIT](https://github.com/Tagliapietra96/logger/blob/main/LICENSE)
