package builder

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func parseDate(s string) time.Time {
	re := regexp.MustCompile(`(\d{4})\-(\d{1,2})-(\d{1,2})`)
	matches := re.FindStringSubmatch(s)
	if len(matches) == 0 {
		return time.Unix(0, 0)
	}
	year, e1 := strconv.Atoi(matches[1])
	month, e2 := strconv.Atoi(matches[2])
	day, e3 := strconv.Atoi(matches[3])
	if e1 != nil || e2 != nil || e3 != nil {
		return time.Unix(0, 0)
	}
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}

func baseName(filename string) string {
	_, file := filepath.Split(filename)
	return file
}

func htmlName(name string) string {
	comma := strings.LastIndex(name, ",")
	if comma != -1 {
		name = strings.TrimSpace(name[comma+1:])
	}
	ext := filepath.Ext(name)
	return strings.TrimSuffix(name, ext) + ".html"
}
