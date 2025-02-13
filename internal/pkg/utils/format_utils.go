package utils

import (
	"fmt"
	"product-api-go/internal/config"
	"product-api-go/internal/pkg/constant"
	"strings"
	"time"
)

func FormatFieldName(fieldName string) string {
	var formatted strings.Builder
	runes := []rune(fieldName)

	for i, r := range runes {
		if i > 0 && r >= 'A' && r <= 'Z' {
			formatted.WriteRune(' ')
		}
		formatted.WriteRune(r)
	}

	return strings.ToLower(formatted.String())
}

func FormatTimeWithTimezone(utcTime time.Time) (string, error) {
	timezone := config.Env("APP_TIMEZONE")

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return "", fmt.Errorf(constant.FailedToLoadTimeZone, err)
	}

	localTime := utcTime.In(location)
	return localTime.Format("02-01-2006 15:04:05"), nil
}
