package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("repeat rule is empty")
	}

	switch parts[0] {
	case "d":
		return nextByDays(now, date, parts)
	case "y":
		return nextByYear(now, date, parts)
	case "w":
		return nextByWeekdays(now, date, parts)
	case "m":
		return nextByMonthDays(now, date, parts)
	default:
		return "", fmt.Errorf("unsupported repeat rule: %s", parts[0])
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}

	now := time.Now()
	if value := r.FormValue("now"); value != "" {
		parsed, err := time.Parse(DateFormat, value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		now = parsed
	}

	next, err := NextDate(now, r.FormValue("date"), r.FormValue("repeat"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(next))
}

func nextByDays(now time.Time, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("invalid day repeat rule")
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil || interval < 1 || interval > 400 {
		return "", errors.New("invalid day interval")
	}

	for {
		date = date.AddDate(0, 0, interval)
		if after(date, now) {
			return date.Format(DateFormat), nil
		}
	}
}

func nextByYear(now time.Time, date time.Time, parts []string) (string, error) {
	if len(parts) != 1 {
		return "", errors.New("invalid year repeat rule")
	}

	for {
		date = date.AddDate(1, 0, 0)
		if after(date, now) {
			return date.Format(DateFormat), nil
		}
	}
}

func nextByWeekdays(now time.Time, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("invalid weekday repeat rule")
	}

	weekdays, err := parseWeekdays(parts[1])
	if err != nil {
		return "", err
	}

	start := laterDate(date, now).AddDate(0, 0, 1)
	for i := 0; i < 7; i++ {
		if weekdays[weekdayNumber(start)] {
			return start.Format(DateFormat), nil
		}
		start = start.AddDate(0, 0, 1)
	}

	return "", errors.New("next weekday was not found")
}

func nextByMonthDays(now time.Time, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 && len(parts) != 3 {
		return "", errors.New("invalid month repeat rule")
	}

	days, err := parseMonthDays(parts[1])
	if err != nil {
		return "", err
	}

	months := [13]bool{}
	if len(parts) == 3 {
		months, err = parseMonths(parts[2])
		if err != nil {
			return "", err
		}
	} else {
		for month := 1; month <= 12; month++ {
			months[month] = true
		}
	}

	current := laterDate(date, now).AddDate(0, 0, 1)
	for i := 0; i < 366*5; i++ {
		month := int(current.Month())
		day := current.Day()
		last := lastDayOfMonth(current)

		if months[month] && (days[day] || days[31] && day == 31 || days[30] && day == 30 || days[-1] && day == last || days[-2] && day == last-1) {
			return current.Format(DateFormat), nil
		}
		current = current.AddDate(0, 0, 1)
	}

	return "", errors.New("next month day was not found")
}

func parseWeekdays(value string) ([8]bool, error) {
	var weekdays [8]bool
	for _, item := range strings.Split(value, ",") {
		weekday, err := strconv.Atoi(item)
		if err != nil || weekday < 1 || weekday > 7 {
			return weekdays, errors.New("invalid weekday")
		}
		weekdays[weekday] = true
	}
	return weekdays, nil
}

func parseMonthDays(value string) (map[int]bool, error) {
	days := make(map[int]bool)
	for _, item := range strings.Split(value, ",") {
		day, err := strconv.Atoi(item)
		if err != nil || day == 0 || day < -2 || day > 31 {
			return nil, errors.New("invalid month day")
		}
		days[day] = true
	}
	return days, nil
}

func parseMonths(value string) ([13]bool, error) {
	var months [13]bool
	for _, item := range strings.Split(value, ",") {
		month, err := strconv.Atoi(item)
		if err != nil || month < 1 || month > 12 {
			return months, errors.New("invalid month")
		}
		months[month] = true
	}
	return months, nil
}

func weekdayNumber(date time.Time) int {
	if date.Weekday() == time.Sunday {
		return 7
	}
	return int(date.Weekday())
}

func laterDate(a time.Time, b time.Time) time.Time {
	if after(a, b) {
		return a
	}
	return b
}

func after(date time.Time, compare time.Time) bool {
	left := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	right := time.Date(compare.Year(), compare.Month(), compare.Day(), 0, 0, 0, 0, time.UTC)
	return left.After(right)
}

func lastDayOfMonth(date time.Time) int {
	return time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
