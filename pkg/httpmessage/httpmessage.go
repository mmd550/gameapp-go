package httpmessage

import (
	"fmt"
	"gameapp/pkg/errormessage"
	"gameapp/pkg/richerror"
	"net/http"
)

func Error(err error) (message string, code int) {

	switch err.(type) {
	case richerror.RichError:
		re := err.(richerror.RichError)
		msg := re.Message()

		code := mapKindToHTTPStatusCode(re.Kind())

		if code >= 500 {
			msg = errormessage.ErrorMessageSomethingWentWrong
		}

		return msg, code		
	default:
		if err == nil {
			fmt.Print("encountered unhandled error")
			return errormessage.ErrorMessageSomethingWentWrong, http.StatusInternalServerError
		}
		fmt.Printf("encountered unhandled error: %s", err.Error())
		return errormessage.ErrorMessageSomethingWentWrong, http.StatusInternalServerError
	}

}

func mapKindToHTTPStatusCode(kind richerror.Kind) int {
	switch kind {
	case richerror.KindForbidden:
		return http.StatusUnauthorized
	case richerror.KindInvalid:
		return http.StatusBadRequest
	case richerror.KindNotFound:
		return http.StatusNotFound
	case richerror.KindUnexpected:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
