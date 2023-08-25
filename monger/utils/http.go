package utils

var SensitiveHeaders = map[string]bool{
	"X-Numexa-Api-Key": true,
	"Authorization":    true,
	"Organization":     true,
}
