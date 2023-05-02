package mongodbr

type IValidation interface {
	Validate() error
}

// validate object if object implement IValidation interface
func Validate(v interface{}) error {
	validation := v.(IValidation)
	if validation == nil {
		return nil
	}
	return validation.Validate()
}
