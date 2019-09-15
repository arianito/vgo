package vgo

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"reflect"
)

type subjectObj = map[string]interface{}

type validatorFunc = func(context *phaseContext, obj subjectObj) error

type phaseContext struct {
	hasType  bool
	name     string
	typ      string
	rule     string
	value    interface{}
	err      string
	args     []string
	hasError bool
	nullable bool
	mime     string
	required bool
}

func checkInternalTypes(context *phaseContext) bool {
	if context.value == nil {
		return false
	}
	switch context.typ {
	case "string":
		if reflect.TypeOf(context.value).Kind() != reflect.String {
			context.hasError = true
			context.err = translate("type.string", translateAttribute(context.name))
			return false
		}
		break
	case "array":
		if reflect.TypeOf(context.value).Kind() != reflect.Slice {
			context.hasError = true
			context.err = translate("type.array", translateAttribute(context.name))
			return false
		}
		break
	case "number":
		field := reflect.TypeOf(context.value)
		kind := field.Kind()
		if !(kind == reflect.String || kind == reflect.Int8 || kind == reflect.Int || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64 ||
			kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64 ||
			kind == reflect.Float32 || kind == reflect.Float64) {
			context.hasError = true
			context.err = translate("type.number", translateAttribute(context.name))
			return false
		}
		break
	case "object":
		if reflect.TypeOf(context.value).Kind() != reflect.Interface {
			context.hasError = true
			context.err = translate("type.object", translateAttribute(context.name))
			return false
		}
		break
	case "date":
		if reflect.TypeOf(context.value).Kind() != reflect.String {
			context.hasError = true
			context.err = translate("type.date", translateAttribute(context.name))
			return false
		}
		break
	case "image":
		if reflect.TypeOf(context.value).Kind() != reflect.String {
			context.hasError = true
			context.err = translate("type.image", translateAttribute(context.name))
			return false
		}
		break
	case "file":
		if reflect.TypeOf(context.value).Kind() != reflect.String {
			context.hasError = true
			context.err = translate("type.file", translateAttribute(context.name))
			return false
		}
		break
	case "bool":
		if reflect.TypeOf(context.value).Kind() != reflect.Bool {
			context.hasError = true
			context.err = translate("type.bool", translateAttribute(context.name))
			return false
		}
		break
	default:
		context.hasError = true
		context.err = translate("type.none", translateAttribute(context.name))
		return false
	}
	return true
}

func convertInternalTypes(context *phaseContext) bool {
	if context.value == nil {
		return false
	}
	switch context.typ {
	case "number":
		field := reflect.TypeOf(context.value)
		kind := field.Kind()
		var strict = true
		if kind == reflect.Int8 || kind == reflect.Int || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64 ||
			kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64 ||
			kind == reflect.Float32 || kind == reflect.Float64 {
			i := reflect.Indirect(reflect.ValueOf(context.value))
			context.value = i.Convert(reflect.TypeOf(float64(0))).Float()
			return true
		} else if kind == reflect.String {
			context.value, strict = convertToNumber(context.value.(string))
		} else {
			strict = false
		}
		if !strict {
			context.hasError = true
			context.err = translate("type.number", translateAttribute(context.name))
			return false
		}
		break
	case "object":
		if reflect.TypeOf(context.value).Kind() != reflect.Interface {
			context.hasError = true
			context.err = translate("type.object", translateAttribute(context.name))
			return false
		}
		break
	case "date":
		if reflect.TypeOf(context.value).Kind() != reflect.String {
			context.hasError = true
			context.err = translate("type.date", translateAttribute(context.name))
			return false
		}
		tm, err := parseDate(context.value.(string))
		if err != nil {
			context.hasError = true
			context.err = translate("type.date", translateAttribute(context.name))
			return false
		}
		context.value = tm
		break
	case "image":
		if reflect.TypeOf(context.value).Kind() != reflect.String {
			context.hasError = true
			context.err = translate("type.image", translateAttribute(context.name))
			return false
		}
		val := context.value.(string)
		start := 0
		if val[0:5] == "data:" {
			start = 5
			for i, char := range val {
				if char == ';' {
					context.mime = val[start:i]
					start = i + 8
				}
			}
		}
		data, err := base64.StdEncoding.DecodeString(val[start:])
		if err != nil {
			context.hasError = true
			context.err = translate("type.image", translateAttribute(context.name))
			return false
		}
		context.value = data
		break
	case "file":
		if reflect.TypeOf(context.value).Kind() != reflect.String {
			context.hasError = true
			context.err = translate("type.file", translateAttribute(context.name))
			return false
		}
		data, err := base64.StdEncoding.DecodeString(context.value.(string))
		if err != nil {
			context.hasError = true
			context.err = translate("type.file", translateAttribute(context.name))
			return false
		}
		context.value = data
		break
	case "bool":
		if reflect.TypeOf(context.value).Kind() != reflect.Bool {
			context.hasError = true
			context.err = translate("type.bool", translateAttribute(context.name))
			return false
		}
		break
	}
	return true
}

func Validate(body string, rules []string) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		return nil, errors.New("malformed request")
	}
	value, pass := validate(data, rules)
	if pass {
		return value, nil
	}else {
		return value, errors.New("validation failed")
	}

}
func validate(obj subjectObj, rules []string) (map[string]interface{}, bool) {
	var values = make(map[string]interface{})
	var errors = make(map[string]interface{})
	err := false
	for _, rule := range rules {
		context := &phaseContext{
			typ: "any",
		}
		evalRuleChain(rule, func(name string, args ...string) bool {
			context.args = args
			if !context.hasType {
				context.name = name
				if len(args) > 0 {
					context.typ = args[0]
				}
				context.hasType = true
				context.value = obj[name]
				checkInternalTypes(context)
				convertInternalTypes(context)
				return true
			}
			context.rule = name
			fn, ok := SharedOperators[context.rule]
			if ok {
				_ = fn(context, obj)
				if context.hasError {
					return false
				}
			}
			//if !checkInternalTypes(context) {
			//	return false
			//}
			vld, ok := validators[context.typ]
			if ok {
				fn, ok = vld.(map[string]validatorFunc)[context.rule]
				if ok {
					_ = fn(context, obj)
				}
			}
			if context.hasError {
				if context.err == "" {
					context.err = translate("none", translateAttribute(context.name))
				}
				return false
			}
			return true
		})
		if context.hasError {
			errors[context.name] = context.err
			err = true
		} else {
			values[context.name] = context.value
		}
	}
	if err {
		return errors, false
	}
	return values, true
}