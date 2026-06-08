package uservalidation

import (
	"errors"
	"gameapp/dto"
	"gameapp/pkg/errormessage"
	"gameapp/pkg/richerror"
	"gameapp/validation"

	"github.com/go-playground/validator/v10"
)

type UserRepository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
}

type UserValidator struct {
	validate       *validator.Validate
	userRepository UserRepository
}

func New(validate *validator.Validate, userRepository UserRepository) UserValidator {
	return UserValidator{
		validate:       validate,
		userRepository: userRepository,
	}
}

func (v UserValidator) ValidateRegisterRequest(registerRequest dto.RegisterRequest) (error, map[string]string) {

	err := v.validate.Struct(registerRequest)

	op := "UserValidator.ValidateRegisterRequest"

	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return richerror.New(op).WithErr(err).WithMessage(errormessage.SomethingWentWrong).WithKind(richerror.KindUnexpected), nil
		}
		var validateErrs validator.ValidationErrors

		fieldErrors := make(map[string]string)

		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				fieldErrors[e.Field()] = validation.FieldErrorMessage(e)
			}
		}

		return richerror.New(op).WithMessage(errormessage.BadRequest).WithKind(richerror.KindInvalid).WithErr(err), fieldErrors
	}

	// check uniqueness of phone number
	if isUnique, err := v.userRepository.IsPhoneNumberUnique(registerRequest.PhoneNumber); err != nil || !isUnique {
		if err != nil {
			return richerror.New(op).WithMessage(errormessage.SomethingWentWrong).WithErr(err).WithKind(richerror.KindUnexpected), nil
		}

		return richerror.New(op).WithMessage(errormessage.BadRequest).WithKind(richerror.KindInvalid), map[string]string{
			"phone_number": errormessage.PhoneNumberIsNotUique,
		}
	}

	return nil, nil
}
