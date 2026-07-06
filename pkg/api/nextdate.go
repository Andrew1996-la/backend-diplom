package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat is empty")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid date: %w", err)
	}

	switch {
	case repeat == "y":
		for !date.After(now) {
			date = date.AddDate(1, 0, 0)
		}

	case strings.HasPrefix(repeat, "d "):
		parts := strings.Split(repeat, " ")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid repeat format")
		}

		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("invalid interval")
		}

		if interval < 1 || interval > 400 {
			return "", fmt.Errorf("interval must be between 1 and 400")
		}

		for !date.After(now) {
			date = date.AddDate(0, 0, interval)
		}

	default:
		return "", fmt.Errorf("unsupported repeat format")
	}

	return date.Format(dateFormat), nil
}
