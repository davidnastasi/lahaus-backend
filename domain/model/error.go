package model

import "fmt"

type DomainError struct {
	Code        int64  `json:"code"`
	Description string `json:"description"`
	Details     string `json:"details"`
}

func (d *DomainError) Error() string {
	return fmt.Sprintf("code: %d, description: %s, details: %s", d.Code, d.Description, d.Details)
}

func NewDomainError(err error) *DomainError {
	return &DomainError{
		Code:        1,
		Description: "Domain error",
		Details:     err.Error(),
	}
}

type EntityNotFoundError struct {
	Code        int64  `json:"code"`
	Description string `json:"description"`
	Details     string `json:"details"`
}

func (d *EntityNotFoundError) Error() string {
	return fmt.Sprintf("code: %d, description: %s, details: %s", d.Code, d.Description, d.Details)
}

func NewEntityNotFoundError(err error) *EntityNotFoundError {
	return &EntityNotFoundError{
		Code:        10,
		Description: "Entity not found",
		Details:     err.Error(),
	}
}

type InternalServerError struct {
	Code        int64  `json:"code"`
	Description string `json:"description"`
	Details     string `json:"details"`
}

func (d *InternalServerError) Error() string {
	return fmt.Sprintf("code: %d, description: %s, details: %s", d.Code, d.Description, d.Details)
}

func NewInternalServerError(err error) *InternalServerError {
	return &InternalServerError{
		Code:        20,
		Description: "Server error",
		Details:     err.Error(),
	}
}

type UnauthorizedError struct {
	Code        int64  `json:"code"`
	Description string `json:"description"`
	Details     string `json:"details"`
}

func (d *UnauthorizedError) Error() string {
	return fmt.Sprintf("code: %d, description: %s, details: %s", d.Code, d.Description, d.Details)
}

func NewUnauthorizedError(err error) *UnauthorizedError {
	return &UnauthorizedError{
		Code:        30,
		Description: "Security issue",
		Details:     err.Error(),
	}
}
