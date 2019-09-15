package vgo

import (
	"encoding/json"
	"github.com/google/uuid"
	"math"
	"net"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var emailValidator = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var regexUsername = regexp.MustCompile("^[a-zA-Z]+[\\-_a-zA-Z0-9]+$")
var alphaNumeric = regexp.MustCompile("^[a-zA-Z0-9\\s.\\-]+$")
var persian = regexp.MustCompile("^[\u0600-\u06FF\\s]+$")
var alphaPersian = regexp.MustCompile("^[a-zA-Z0-9\u0600-\u06FF\\s]+$")
var alpha = regexp.MustCompile("^[a-zA-Z0-9\\s]+$")

func contains(val string, args []string) bool {
	for _, item := range args {
		if val == item {
			return true
		}
	}
	return false
}

var validators = map[string]interface{}{
	"date": map[string]validatorFunc{
		"after": func(context *phaseContext, obj subjectObj) error {
			a, _ := parseDate(context.args[0])
			v := context.value.(time.Time)
			if !v.After(a) {
				context.hasError = true
				context.err = translate("date.after", translateAttribute(context.name), formatDate(v), formatDate(a))
			}
			return nil
		},
		"before": func(context *phaseContext, obj subjectObj) error {
			a, _ := parseDate(context.args[0])
			v := context.value.(time.Time)
			if !v.Before(a) {
				context.hasError = true
				context.err = translate("date.before", translateAttribute(context.name), formatDate(v), formatDate(a))
			}
			return nil
		},
		"between": func(context *phaseContext, obj subjectObj) error {
			a, _ := parseDate(context.args[0])
			b, _ := parseDate(context.args[1])
			v := context.value.(time.Time)
			if v.Before(a) || v.After(b) {
				context.hasError = true
				context.err = translate("date.between", translateAttribute(context.name), formatDate(v), formatDate(a), formatDate(b))
			}
			return nil
		},
	},
	"number": map[string]validatorFunc{
		"digits": func(context *phaseContext, obj subjectObj) error {
			a, _ := strconv.Atoi(context.args[0])
			v := (int64)(math.Floor(context.value.(float64)))
			k:=1
			for i := v; i > 10; i/=10 {
				k++
			}
			if a != k {
				context.hasError = true
				context.err = translate("number.digits", translateAttribute(context.name), a)
			}
			return nil
		},
		"digitsBetween": func(context *phaseContext, obj subjectObj) error {
			a, _ := strconv.Atoi(context.args[0])
			b, _ := strconv.Atoi(context.args[1])
			v := (int64)(math.Floor(context.value.(float64)))
			k:=1
			for i := v; i > 10; i/=10 {
				k++
			}
			if k < a || k > b {
				context.hasError = true
				context.err = translate("number.digitsBetween", translateAttribute(context.name), a, b)
			}
			return nil
		},
		"integer": func(context *phaseContext, obj subjectObj) error {
			i := reflect.Indirect(reflect.ValueOf(context.value))
			context.value = i.Convert(reflect.TypeOf(0)).Int()
			return nil
		},
		"greaterThan": func(context *phaseContext, obj subjectObj) error {
			a, _ := strconv.ParseFloat(context.args[0], 64)
			i := reflect.Indirect(reflect.ValueOf(context.value))
			val := i.Convert(reflect.TypeOf(float64(0))).Float()
			if val <= a {
				context.hasError = true
				context.err = translate("number.greaterThan", translateAttribute(context.name), a)
			}
			return nil
		},
		"greaterThanOrEqual": func(context *phaseContext, obj subjectObj) error {
			a, _ := strconv.ParseFloat(context.args[0], 64)
			i := reflect.Indirect(reflect.ValueOf(context.value))
			val := i.Convert(reflect.TypeOf(float64(0))).Float()
			if val < a {
				context.hasError = true
				context.err = translate("number.greaterThanOrEqual", translateAttribute(context.name), a)
			}
			return nil
		},
		"lessThan": func(context *phaseContext, obj subjectObj) error {
			a, _ := strconv.ParseFloat(context.args[0], 64)
			i := reflect.Indirect(reflect.ValueOf(context.value))
			val := i.Convert(reflect.TypeOf(float64(0))).Float()
			if val >= a {
				context.hasError = true
				context.err = translate("number.lessThan", translateAttribute(context.name), a)
			}
			return nil
		},
		"lessThanOrEqual": func(context *phaseContext, obj subjectObj) error {
			a, _ := strconv.ParseFloat(context.args[0], 64)
			i := reflect.Indirect(reflect.ValueOf(context.value))
			val := i.Convert(reflect.TypeOf(float64(0))).Float()
			if val > a {
				context.hasError = true
				context.err = translate("number.lessThanOrEqual", translateAttribute(context.name), a)
			}
			return nil
		},
		"between": func(context *phaseContext, obj subjectObj) error {
			a, _ := strconv.ParseFloat(context.args[0], 64)
			b, _ := strconv.ParseFloat(context.args[1], 64)
			i := reflect.Indirect(reflect.ValueOf(context.value))
			val := i.Convert(reflect.TypeOf(float64(0))).Float()
			if val < a || val > b {
				context.hasError = true
				context.err = translate("number.between", translateAttribute(context.name), a, b)
			}
			return nil
		},
	},
	"string": map[string]validatorFunc{
		"filled": func(context *phaseContext, obj subjectObj) error {
			str := context.value.(string)
			if len(str) < 1 {
				context.hasError = true
				context.err = translate("string.filled", translateAttribute(context.name))
			}
			return nil
		},
		"json": func(context *phaseContext, obj subjectObj) error {
			str := context.value.(string)
			var data map[string]interface{}
			err := json.Unmarshal([]byte(str), &data)
			if err != nil {
				context.hasError = true
				context.err = translate("string.json", translateAttribute(context.name))
			}
			return nil
		},
		"url": func(context *phaseContext, obj subjectObj) error {
			_, err := url.ParseRequestURI(context.value.(string))
			if err != nil {
				context.hasError = true
				context.err = translate("string.url", translateAttribute(context.name))
			}
			return nil
		},
		"uuid": func(context *phaseContext, obj subjectObj) error {
			_, err := uuid.Parse(context.value.(string))
			if err != nil {
				context.hasError = true
				context.err = translate("string.uuid", translateAttribute(context.name))
			}
			return nil
		},
		"ip": func(context *phaseContext, obj subjectObj) error {
			test := net.ParseIP(context.value.(string))
			if test.To4() == nil || test.To16() == nil {
				context.hasError = true
				context.err = translate("string.ip", translateAttribute(context.name))
			}
			return nil
		},
		"ipv4": func(context *phaseContext, obj subjectObj) error {
			test := net.ParseIP(context.value.(string))
			if test.To4() == nil {
				context.hasError = true
				context.err = translate("string.ipv4", translateAttribute(context.name))
			}
			return nil
		},
		"ipv6": func(context *phaseContext, obj subjectObj) error {
			test := net.ParseIP(context.value.(string))
			if test.To16() == nil {
				context.hasError = true
				context.err = translate("string.ipv6", translateAttribute(context.name))
			}
			return nil
		},
		"email": func(context *phaseContext, obj subjectObj) error {
			if !emailValidator.MatchString(context.value.(string)) {
				context.hasError = true
				context.err = translate("string.email", translateAttribute(context.name))
			}
			return nil
		},
		"mobile": func(context *phaseContext, obj subjectObj) error {
			if !regexp.MustCompile("^[0][9][0-9]{9}$").MatchString(context.value.(string)) {
				context.hasError = true
				context.err = translate("string.mobile", translateAttribute(context.name))
			}
			return nil
		},
		"phone": func(context *phaseContext, obj subjectObj) error {
			if !regexp.MustCompile("^[0][1-8][0-9]{9}$").MatchString(context.value.(string)) {
				context.hasError = true
				context.err = translate("string.phone", translateAttribute(context.name))
			}
			return nil
		},
		"in": func(context *phaseContext, obj subjectObj) error {
			for _, item := range context.args {
				if item == context.value {
					return nil
				}
			}
			context.hasError = true
			context.err = translate("string.in", translateAttribute(context.name))
			return nil
		},
		"inArray": func(context *phaseContext, obj subjectObj) error {
			val, ok := obj[context.args[0]]
			if ok && reflect.TypeOf(val).Kind() == reflect.Slice {
				for _, item := range val.([]interface{}) {
					if item == context.value {
						return nil
					}
				}
				context.hasError = true
				context.err = translate("string.inArray", translateAttribute(context.name), translateAttribute(context.args[0]))
			}
			return nil
		},
		"notIn": func(context *phaseContext, obj subjectObj) error {
			for _, item := range context.args {
				if item == context.value {
					context.hasError = true
					context.err = translate("string.notIn", translateAttribute(context.name))
					return nil
				}
			}
			return nil
		},
		"size": func(context *phaseContext, obj subjectObj) error {
			val, _ := strconv.Atoi(context.args[0])
			str, ok :=context.value.(string)
			if ok && len(str) != val {
				context.hasError = true
				context.err = translate("string.size", translateAttribute(context.name), val)
				return nil
			}
			return nil
		},
		"min": func(context *phaseContext, obj subjectObj) error {
			a, _ := strconv.Atoi(context.args[0])
			c := len(context.value.(string))
			if c < a {
				context.hasError = true
				context.err = translate("string.min", translateAttribute(context.name), a)
				return nil
			}
			return nil
		},
		"max": func(context *phaseContext, obj subjectObj) error {
			b, _ := strconv.Atoi(context.args[0])
			c := len(context.value.(string))
			if c > b {
				context.hasError = true
				context.err = translate("string.max", translateAttribute(context.name), b)
				return nil
			}
			return nil
		},
		"between": func(context *phaseContext, obj subjectObj) error {
			a, _ := strconv.Atoi(context.args[0])
			b, _ := strconv.Atoi(context.args[1])
			c := len(context.value.(string))
			if c < a || c > b {
				context.hasError = true
				context.err = translate("string.between", translateAttribute(context.name), a, b)
				return nil
			}
			return nil
		},
		"username": func(context *phaseContext, obj subjectObj) error {
			if !regexUsername.MatchString(context.value.(string)) {
				context.hasError = true
				context.err = translate("string.username", translateAttribute(context.name))
				return nil
			}
			return nil
		},
		"alphaNum": func(context *phaseContext, obj subjectObj) error {
			if !alphaNumeric.MatchString(context.value.(string)) {
				context.hasError = true
				context.err = translate("string.alphaNum", translateAttribute(context.name))
				return nil
			}
			return nil
		},
		"alpha": func(context *phaseContext, obj subjectObj) error {

			hasFa := contains("fa", context.args)
			hasEn := contains("en", context.args)
			if hasFa && hasEn {
				if !alphaPersian.MatchString(context.value.(string)) {
					context.hasError = true
					context.err = translate("string.alpha", translateAttribute(context.name))
					return nil
				}
			} else if hasFa {
				if !persian.MatchString(context.value.(string)) {
					context.hasError = true
					context.err = translate("string.persian", translateAttribute(context.name))
					return nil
				}
			} else if hasEn {
				if !alpha.MatchString(context.value.(string)) {
					context.hasError = true
					context.err = translate("string.alpha", translateAttribute(context.name))
					return nil
				}
			}
			return nil
		},
		"regex": func(context *phaseContext, obj subjectObj) error {
			re := regexp.MustCompile(context.args[0]).MatchString(context.value.(string))
			if !re {
				context.hasError = true
				context.err = translate("string.regex", translateAttribute(context.name))
				return nil
			}
			return nil
		},
		"notRegex": func(context *phaseContext, obj subjectObj) error {
			re := regexp.MustCompile(context.args[0]).MatchString(context.value.(string))
			if re {
				context.hasError = true
				context.err = translate("string.regex", translateAttribute(context.name))
				return nil
			}
			return nil
		},
		"contains": func(context *phaseContext, obj subjectObj) error {
			val := context.value.(string)
			for _, arg := range context.args {
				if strings.Index(val, arg) > -1 {
					return nil
				}
			}
			context.hasError = true
			context.err = translate("string.contains", translateAttribute(context.name), strings.Join(context.args, ","))
			return nil
		},
		"startsWith": func(context *phaseContext, obj subjectObj) error {
			val := context.value.(string)
			for _, arg := range context.args {
				if strings.Index(val, arg) == 0 {
					return nil
				}
			}
			context.hasError = true
			context.err = translate("string.startsWith", translateAttribute(context.name), strings.Join(context.args, ","))
			return nil
		},
		"endsWith": func(context *phaseContext, obj subjectObj) error {
			val := context.value.(string)
			nv := len(val)
			for _, arg := range context.args {
				if strings.Index(val, arg) == nv-len(arg) {
					return nil
				}
			}
			context.hasError = true
			context.err = translate("string.endsWith", translateAttribute(context.name), strings.Join(context.args, ","))
			return nil
		},
		"same": func(context *phaseContext, obj subjectObj) error {
			arg := context.args[0]
			a, aOk := obj[context.name]
			b, bOk := obj[arg]
			if !aOk || !bOk {
				context.hasError = true
				context.err = translate("none", translateAttribute(context.name))
			}
			if a != b {
				context.hasError = true
				context.err = translate("same", translateAttribute(context.name), translateAttribute(arg))
			}
			return nil
		},
		"different": func(context *phaseContext, obj subjectObj) error {
			arg := context.args[0]
			a, aOk := obj[context.name]
			b, bOk := obj[arg]
			if !aOk || !bOk {
				context.hasError = true
				context.err = translate("none", translateAttribute(context.name))
			}
			if a == b {
				context.hasError = true
				context.err = translate("different", translateAttribute(context.name), translateAttribute(arg))
			}
			return nil
		},
	},
}
