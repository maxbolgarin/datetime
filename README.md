# datetime

[![Go Version][version-img]][doc] [![GoDoc][doc-img]][doc] [![Build][ci-img]][ci] [![GoReport][report-img]][report]

`datetime` — package for handling datetime in UX friendly manner. Date is stored without time in `yyyy-mm-dd` format. Time is stored in `HH:MM` format without seconds. Timezone is stored in `UTC(+|-)HH:MM` format.

```
go get -u github.com/maxbolgarin/datetime
```

## Overview

The `datetime` package offers the following main functionalities:

- **Time**: Represents a specific time of day (hours and minutes) without a date and timezone context.
  - Functions to create new time instances from integers or strings.
  - Methods to add/subtract time, compare, and convert to/from strings/JSON.

- **Date**: Represents a specific calendar date without attached time.
  - Functions to create dates from integers or strings.
  - Methods to navigate days, compare, and convert to/from strings/JSON.

- **Timezone**: Represents time zones as UTC offsets.
  - Functions for creating time zones from predefined locations or custom strings.
  - Supports JSON serialization.


### Advantages

- **Lightweight**: Minimal overhead for handling common date/time operations.
- **Flexible Parsing**: Ability to parse various common string formats for both dates and times.
- **JSON Support**: Provides seamless JSON marshaling and unmarshaling.
- **Discrete Time and Date Handling**: Time and Date are considered separate entities, which simplifies some specific usage contexts.

### Disadvantages

- **Limited Scope**: Primarily intended for handling basic operations and may not suit complex calendrical computations.
- **UTC Focused Time Representations**: Time representations are tied to UTC, which can require additional handling for local or non-standard time zones.


## Usage Examples

### Time Management

```go
package main

import (
    "fmt"
    "yourmodule/datetime"
)

func main() {
    // Create a new Time instance at 10:30
    timeInstance := datetime.NewTime(10, 30)
    
    // Print the Time string representation
    fmt.Println("Time:", timeInstance.String())

    // Parse a time string
    parsedTime, err := datetime.NewTimeFromString("14:45")
    if err != nil {
        fmt.Println("Error parsing time:", err)
    }
    fmt.Println("Parsed Time:", parsedTime.String())
    
    // Calculate time difference
    duration := timeInstance.Range(parsedTime)
    fmt.Println("Time Difference:", duration)
}
```

### Date Management

```go
package main

import (
    "fmt"
    "yourmodule/datetime"
)

func main() {
    // Create a new Date instance for January 1, 2022
    dateInstance := datetime.NewDate(2022, 1, 1)
    
    // Print the Date string representation
    fmt.Println("Date:", dateInstance.String())

    // Parse a date string
    parsedDate, err := datetime.NewDateFromString("2023-08-15")
    if err != nil {
        fmt.Println("Error parsing date:", err)
    }
    
    fmt.Println("Parsed Date:", parsedDate.String())
    
    // Determine if the given date is today
    timezone := datetime.NewTimezoneFromTime(time.Now())
    fmt.Println("Is Today:", parsedDate.IsToday(datetime.EmptyTime, timezone.Loc()))
}
```

### Timezone Handling

```go
package main

import (
    "fmt"
    "yourmodule/datetime"
    "time"
)

func main() {
    // Create Timezone instance for UTC+5:00
    tz, err := datetime.ParseTimezone("UTC+5:00")
    if err != nil {
        fmt.Println("Error parsing timezone:", err)
    }
    
    fmt.Println("Timezone Location:", tz.Loc().String())
    fmt.Println("Timezone Offset (hours):", tz.OffsetHours())
}
```

## Contributions and Issues

Feel free to contribute to this package or report issues you encounter during usage. Collaborative improvements are welcome to refine and extend the usability of the `datetime` package.

## License

This project is licensed under the terms of the [MIT License](LICENSE).

[MIT License]: LICENSE.txt
[version-img]: https://img.shields.io/badge/Go-%3E%3D%201.13-%23007d9c
[doc-img]: https://pkg.go.dev/badge/github.com/maxbolgarin/datetime
[doc]: https://pkg.go.dev/github.com/maxbolgarin/datetime
[ci-img]: https://github.com/maxbolgarin/datetime/actions/workflows/go.yml/badge.svg
[ci]: https://github.com/maxbolgarin/datetime/actions
[report-img]: https://goreportcard.com/badge/github.com/maxbolgarin/datetime
[report]: https://goreportcard.com/report/github.com/maxbolgarin/datetime
