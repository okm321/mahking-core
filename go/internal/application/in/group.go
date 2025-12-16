package in

import "github.com/go-playground/validator/v10"

var validate = validator.New()

// CreateGroupCommand represents input for creating a group.
// Required / length checks are handled here (application-layer input).
type CreateGroupCommand struct {
	Name string `json:"name" validate:"required,max=100"`
}

// Validate applies input-level validation.
func (c CreateGroupCommand) Validate() error {
	return validate.Struct(c)
}
