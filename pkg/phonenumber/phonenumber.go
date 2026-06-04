package phonenumber

import "regexp"

var iranianMobileRegex = regexp.MustCompile(`^(\+98|0)?9\d{9}$`)

func IsValid(number string) bool {
	return iranianMobileRegex.MatchString(number)
}
