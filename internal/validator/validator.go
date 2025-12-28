package validator

import (
	"regexp"
	"slices"
)

type Validator struct {
	FieldErorrs map[string]string
}

func New() *Validator {
	return &Validator{
		FieldErorrs: make(map[string]string),
	}
}

var EmailRX = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

func (v *Validator) Valid() bool {
	return len(v.FieldErorrs) == 0
}

func (v *Validator) AddError(key, value string) {
	if _, exists := v.FieldErorrs[key]; !exists {
		v.FieldErorrs[key] = value
	}
}

func (v *Validator) Check(ok bool, key, value string) {
	if ok {
		v.AddError(key, value)
	}
}

func CheckEmail(email string, rx *regexp.Regexp) bool {
	return rx.MatchString(email)
}

func (v *Validator) Empty(value string) bool {
	return value == ""
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
