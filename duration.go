package iso8601

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var (
	// ErrBadFormat is returned when parsing fails
	ErrBadFormat = errors.New("bad format string")

	// ErrNoMonth is raised when a month is in the format string
	ErrNoMonth = errors.New("no months allowed")

	full = regexp.MustCompile(`P((?P<year>\d+)Y)?((?P<month>\d+)M)?((?P<day>\d+)D)?(T((?P<hour>\d+)H)?((?P<minute>\d+)M)?((?P<second>\d+)S)?)?`)
	week = regexp.MustCompile(`P((?P<week>\d+)W)`)
)

type Duration struct {
	Years   int
	Months  int
	Weeks   int
	Days    int
	Hours   int
	Minutes int
	Seconds int
}

// adapted from https://github.com/BrianHicks/finch/duration
func ParseDuration(value string) (time.Duration, *Duration, error) {
	var match []string
	var regex *regexp.Regexp
	result := Duration{
		Years:   0,
		Months:  0,
		Weeks:   0,
		Days:    0,
		Hours:   0,
		Minutes: 0,
		Seconds: 0,
	}

	if week.MatchString(value) {
		match = week.FindStringSubmatch(value)
		regex = week
	} else if full.MatchString(value) {
		match = full.FindStringSubmatch(value)
		regex = full
	} else {
		return time.Duration(0), nil, ErrBadFormat
	}

	d := time.Duration(0)
	day := time.Hour * 24
	week := day * 7
	year := day * 365

	for i, name := range regex.SubexpNames() {
		part := match[i]
		if i == 0 || name == "" || part == "" {
			continue
		}

		value, err := strconv.Atoi(part)
		if err != nil {
			return time.Duration(0), nil, err
		}
		switch name {
		case "year":
			result.Years = value
			d += year * time.Duration(value)
		case "month":
			result.Months = value
			if value != 0 {
				return time.Duration(0), nil, ErrNoMonth
			}
		case "week":
			result.Weeks = value
			d += week * time.Duration(value)
		case "day":
			result.Days = value
			d += day * time.Duration(value)
		case "hour":
			result.Hours = value
			d += time.Hour * time.Duration(value)
		case "minute":
			result.Minutes = value
			d += time.Minute * time.Duration(value)
		case "second":
			result.Seconds = value
			d += time.Second * time.Duration(value)
		}
	}

	return d, &result, nil
}

func FormatDuration(duration time.Duration) string {
	// we're not doing negative durations
	if duration.Seconds() <= 0 {
		return "PT0S"
	}

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) - (hours * 60)
	seconds := int(duration.Seconds()) - (hours*3600 + minutes*60)

	// we're not doing Y,M,W
	s := "PT"
	if hours > 0 {
		s = fmt.Sprintf("%s%dH", s, hours)
	}
	if minutes > 0 {
		s = fmt.Sprintf("%s%dM", s, minutes)
	}
	if seconds > 0 {
		s = fmt.Sprintf("%s%dS", s, seconds)
	}

	return s
}
