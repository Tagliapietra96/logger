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
