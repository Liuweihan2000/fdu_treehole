package utils

import "fmt"

func FatalErrorHandle(err error, msg string) {
	if err != nil {
		panic(fmt.Errorf("Fatal error: %s \n %s \n", msg, err))
	}
}
