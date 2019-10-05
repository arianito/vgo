package vgo

import (
	"strings"
)

type File struct {
	MimeType string
	Buffer []byte
}

func checkEmptiness(obj interface{}, nullable bool) bool {
	switch obj.(type) {
	case nil:
		if nullable {
			return false
		}
		return true
	case string:
		return len(obj.(string)) < 1
	case []interface{}:
		return len(obj.([]interface{})) < 1
	case map[string]interface{}:
		return len(obj.(map[string]interface{})) < 1
	}
	return false
}

var sharedOperators = map[string]validatorFunc{
	"nullable": func(context *phaseContext, obj subjectObj) error {
		context.nullable = true
		return nil
	},
	"present": func(context *phaseContext, obj subjectObj) error {
		if _, ok := obj[context.name]; !ok {
			context.hasError = true
			context.err =  translate("present", translateAttribute(context.name))
		}
		return nil
	},
	"required": func(context *phaseContext, obj subjectObj) error {
		if val, ok := obj[context.name]; !ok || checkEmptiness(val, context.nullable) {
			context.hasError = true
			context.err = translate("required", translateAttribute(context.name))
		}
		return nil
	},
	"requiredWith": func(context *phaseContext, obj subjectObj) error {
		allFieldsExists := true
		for _, item := range context.args {
			if val, ok := obj[item]; !ok || checkEmptiness(val, false) {
				allFieldsExists = false
				break
			}
		}
		if allFieldsExists {
			if val, ok := obj[context.name]; !ok || checkEmptiness(val, context.nullable) {
				context.hasError = true
				if len(context.args) > 1 {
					context.err =  translate("requiredWithAll", strings.Join(translateAttributes(context.args...), "|"), translateAttribute(context.name))
				}else {
					context.err =  translate("requiredWith", translateAttribute(context.args[0]), translateAttribute(context.name))
				}
			}
		}
		return nil
	},
	"requiredWithout": func(context *phaseContext, obj subjectObj) error {
		allFieldsExists := true
		for _, item := range context.args {
			if val, ok := obj[item]; !ok || checkEmptiness(val, false) {
				allFieldsExists = false
				break
			}
		}
		if !allFieldsExists {
			if val, ok := obj[context.name]; !ok || checkEmptiness(val, context.nullable) {
				context.hasError = true
				if len(context.args) > 1 {
					context.err =  translate("requiredWithoutAll", strings.Join(translateAttributes(context.args...), "|"), translateAttribute(context.name))
				}else {
					context.err =  translate("requiredWithout", translateAttribute(context.args[0]), translateAttribute(context.name))
				}
			}
		}
		return nil
	},
	"confirmed": func(context *phaseContext, obj subjectObj) error {
		arg := context.name + "Confirmation"
		if len(context.args) > 0 {
			arg = context.args[0]
		}
		a, aOk := obj[context.name]
		b, bOk := obj[arg]
		if !aOk || !bOk || a != b {
			context.hasError = true
			context.err =  translate("confirmed", translateAttribute(context.name))
		}
		return nil
	},
}
