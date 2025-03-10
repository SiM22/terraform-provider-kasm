package validators

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
)

// StringValidator is a custom string validator
type StringValidator struct {
	Desc       string
	ValidateFn func(string) bool
	ErrMessage string
}

func (v StringValidator) Description(ctx context.Context) string {
	return v.Desc
}

func (v StringValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v StringValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if !v.ValidateFn(req.ConfigValue.ValueString()) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid value",
			v.ErrMessage,
		)
	}
}

// StringOneOf returns a validator which ensures that any configured string
// value matches one of the given values exactly.
func StringOneOf(validValues ...string) validator.String {
	return StringValidator{
		Desc: fmt.Sprintf("value must be one of: %v", validValues),
		ValidateFn: func(value string) bool {
			for _, validValue := range validValues {
				if value == validValue {
					return true
				}
			}
			return false
		},
		ErrMessage: fmt.Sprintf("value must be one of: %v", validValues),
	}
}

// Int64Validator is a custom int64 validator
type Int64Validator struct {
	Desc       string
	ValidateFn func(int64) bool
	ErrMessage string
}

func (v Int64Validator) Description(ctx context.Context) string {
	return v.Desc
}

func (v Int64Validator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v Int64Validator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if !v.ValidateFn(req.ConfigValue.ValueInt64()) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid value",
			v.ErrMessage,
		)
	}
}

// Float64Validator is a custom float64 validator
type Float64Validator struct {
	Desc       string
	ValidateFn func(float64) bool
	ErrMessage string
}

func (v Float64Validator) Description(ctx context.Context) string {
	return v.Desc
}

func (v Float64Validator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v Float64Validator) ValidateFloat64(ctx context.Context, req validator.Float64Request, resp *validator.Float64Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if !v.ValidateFn(req.ConfigValue.ValueFloat64()) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid value",
			v.ErrMessage,
		)
	}
}

// Common validator functions
func ValidateURL() validator.String {
	return StringValidator{
		Desc: "must be a valid URL",
		ValidateFn: func(val string) bool {
			matched, _ := regexp.MatchString(`^https?://`, val)
			return matched
		},
		ErrMessage: "URL must start with http:// or https://",
	}
}

func Int64AtLeast(min int64) validator.Int64 {
	return Int64Validator{
		Desc: fmt.Sprintf("must be at least %d", min),
		ValidateFn: func(val int64) bool {
			return val >= min
		},
		ErrMessage: fmt.Sprintf("value must be at least %d", min),
	}
}
