package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	return date.Format(dateFormat) > now.Format(dateFormat)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if len(repeat) == 0 {
		return "", fmt.Errorf("empty repeat rule")
	}

	start, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Split(repeat, " ")
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid d format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("invalid days count")
		}
		date := start
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(dateFormat), nil

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("invalid y format")
		}
		date := start
		for {
			year := date.Year() + 1
			month := start.Month()
			day := start.Day()

			date = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

			if date.Month() != month {
				date = time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
			}
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(dateFormat), nil

	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid w format")
		}
		days := strings.Split(parts[1], ",")
		var validDays [8]bool
		for _, d := range days {
			n, err := strconv.Atoi(strings.TrimSpace(d))
			if err != nil || n < 1 || n > 7 {
				return "", fmt.Errorf("invalid weekday")
			}
			validDays[n] = true
		}
		date := start
		for {
			wd := int(date.Weekday())
			if wd == 0 {
				wd = 7
			}
			if validDays[wd] && afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
			date = date.AddDate(0, 0, 1)
			if date.Year() > now.Year()+100 {
				return "", fmt.Errorf("date not found")
			}
		}

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf("invalid m format")
		}
		dayParts := strings.Split(parts[1], ",")
		var validDays []int
		for _, d := range dayParts {
			n, err := strconv.Atoi(strings.TrimSpace(d))
			if err != nil || n < -2 || n == 0 || n > 31 {
				return "", fmt.Errorf("invalid month day")
			}
			validDays = append(validDays, n)
		}

		var validMonths [13]bool
		if len(parts) == 3 {
			monthParts := strings.Split(parts[2], ",")
			for _, m := range monthParts {
				n, err := strconv.Atoi(strings.TrimSpace(m))
				if err != nil || n < 1 || n > 12 {
					return "", fmt.Errorf("invalid month")
				}
				validMonths[n] = true
			}
		} else {
			for i := 1; i <= 12; i++ {
				validMonths[i] = true
			}
		}

		date := start
		for {
			month := int(date.Month())
			if !validMonths[month] {
				date = date.AddDate(0, 0, 1)
				continue
			}

			day := date.Day()
			lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

			valid := false
			for _, d := range validDays {
				if d > 0 && day == d {
					valid = true
					break
				}
				if d == -1 && day == lastDay {
					valid = true
					break
				}
				if d == -2 && day == lastDay-1 {
					valid = true
					break
				}
			}

			if valid && afterNow(date, now) {
				return date.Format(dateFormat), nil
			}

			date = date.AddDate(0, 0, 1)
			if date.Year() > now.Year()+1000 {
				return "", fmt.Errorf("date not found")
			}
		}

	default:
		return "", fmt.Errorf("unsupported repeat format")
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}
