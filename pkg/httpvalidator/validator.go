package httpvalidator

import "fmt"

const (
	requiredMessage = "field is required"
)

type Rule struct {
	Description string
	Validate    func(string) bool
}

type BodyField struct {
	Required bool
	Rules    []Rule
}

type Fields map[string]BodyField

type RequestBody struct {
	Fields Fields
}

type Bodies map[string]RequestBody

type PathValues map[string][]Rule

type Validator struct {
	BodyTemplates      Bodies
	PathValueTemplates PathValues
}

func NewValidator() Validator {
	return Validator{
		BodyTemplates:      make(map[string]RequestBody),
		PathValueTemplates: make(map[string][]Rule),
	}
}

func (v *Validator) AddBodyTemplate(templateName string, template RequestBody) {
	v.BodyTemplates[templateName] = template
}

func (v *Validator) AddPathValueTemplate(templateName string, rules []Rule) {
	v.PathValueTemplates[templateName] = rules
}

type ValidationError struct {
	Location string
	Param    string
	Value    string
	Message  string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error: location: %s, param: %s, value: %s, message: %s", e.Location, e.Param, e.Value, e.Message)
}

func (v *Validator) ValidateBody(bodyName string, body map[string]string) []ValidationError {
	response := make([]ValidationError, 0)

	for param, field := range v.BodyTemplates[bodyName].Fields {
		value, ok := body[param]
		if !ok && field.Required {
			response = append(response, ValidationError{Location: "body", Param: param, Value: value, Message: requiredMessage})
			continue
		}
		for _, rule := range field.Rules {
			if !rule.Validate(value) {
				response = append(response, ValidationError{Location: "body", Param: param, Value: value, Message: rule.Description})
			}
		}
	}

	return response
}

func (v *Validator) ValidatePathValue(param string, value string) []ValidationError {
	response := make([]ValidationError, 0)

	for _, rule := range v.PathValueTemplates[param] {
		if !rule.Validate(value) {
			response = append(response, ValidationError{Location: "path", Param: param, Value: value, Message: rule.Description})
		}
	}

	return response
}
