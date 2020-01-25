package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	units = []string{"years", "months", "weeks", "days", "hours", "minutes", "seconds"}
)

// FormatByteSize formats byte size into a human readable format.
func FormatByteSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

// FormatDuration parses duration into a human readable format.
func FormatDuration(duration time.Duration, limit int) string {
	log.Tracef("Formatting duration %d", duration)

	result := ""

	if duration < 0 {
		result += "-"
		duration = -duration
	}

	durationMap := getDurationMap(duration)
	log.Tracef("Formatting with duration map %v", durationMap)

	for i := range units {
		u := units[i]
		v := durationMap[u]
		strval := strconv.FormatInt(v, 10)
		switch {

		case v > 1:
			result += strval + " " + u + " "
		case v == 1:
			result += strval + " " + strings.TrimRight(u, "s") + " "
		case v == 0:
			continue
		}
	}

	result = strings.TrimSpace(result)
	if limit > 0 {
		parts := strings.Split(result, " ")
		if len(parts) > limit*2 {
			result = strings.Join(parts[:limit*2], " ")
		}
	}

	return result
}

func getDurationMap(duration time.Duration) map[string]int64 {

	seconds := int64(duration.Seconds()) % 60
	minutes := int64(duration.Minutes()) % 60
	hours := int64(duration.Hours()) % 24

	totalDays := int64(duration / (24 * time.Hour))
	years := totalDays / 365
	months := totalDays % 365 / 30
	weeks := totalDays % 365 % 30 / 7
	days := totalDays % 365 % 30 % 7

	return map[string]int64{
		"seconds": seconds,
		"minutes": minutes,
		"hours":   hours,
		"days":    days,
		"weeks":   weeks,
		"months":  months,
		"years":   years,
	}
}
