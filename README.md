# Logger - A Lightweight Logging System with Advanced Features
Logger is a Go package that provides a versatile and efficient logging system designed for applications that require structured log management. With its built-in SQLite integration, filtering capabilities, customizable terminal output, and alert notifications, Logger offers a comprehensive solution for tracking and managing application events.
## Features
### SQLite-Based Logging
Logs are stored in a lightweight SQLite database, ensuring persistence and easy retrieval without relying on external systems.
### Flexible Filtering and Configurable Terminal Output
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

---

1. **Setting the Log Storage Folder**
By default, the logs are stored in an SQLite database in the current working directory. You can change the storage location with the `Folder` method:
```go
log := logger.New()

// Set the folder where the SQLite database will be stored
log.Folder("~/projects/my-logs/")
```
> **Note:** Ensure the specified folder exists and has the necessary write permissions.

---

2. **Configuring Log Output Format (Inline vs Block)**
You can control how logs are printed to the terminal. Logs can be displayed in a compact, single-line format (`inline`), or in a more detailed, block format, where each log is presented as a card-like entry (`block`).

```go
// Print logs in a single-line format
log.Inline(true)  

// Print logs in a block format (default)
log.Inline(false)  
```
> **Inline Mode:** Suitable for quick, concise debugging.
> **Block Mode:** Ideal for comprehensive, formatted log displays with better readability.

---

3. **Customizing Caller Information Display**
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

---

4. **Configuring Timestamp Display**
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
> **Default Format:** `2006-01-02 15:04:05`
> **Full Timestamp Example:** `Monday 2006-01-02 15:04:05`

---

5. **Managing Tags for Logs**
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

---

6. **Configuring Fatal Notifications**
Customize the message and title for critical errors using the `SetFatal` method. This is particularly useful for displaying user-friendly or context-specific messages.

```go
// Set a custom title and message for fatal error notifications
log.SetFatal("MyApp - CRITICAL ERROR", "Oops! Something went wrong. Check the logs.")
```
> **Default Title:** `"Fatal"`
> **Default Message:** `"An error occurred, please check the logs for more information"`

---

These configurations allow you to fully customize your logging setup, ensuring that the logs are both informative and easy to manage, while maintaining flexibility in how they are displayed and stored.