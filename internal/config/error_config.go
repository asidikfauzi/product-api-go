package config

import "os"

func HandleError(err error) interface{} {
	if os.Getenv("APP_DEBUG") == "true" {
		return err.Error()
	}
	return nil
}
