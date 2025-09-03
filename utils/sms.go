package utils

import "fmt"

func SendSMS(phone, message string) error {
	fmt.Printf("Sending SMS to %s: %s\n", phone, message)
	return nil
}
