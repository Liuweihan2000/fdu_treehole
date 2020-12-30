package utils

import "fmt"

func Debug(obj interface{}) {
	fmt.Printf("\n\n%v\n\n", obj)
}
