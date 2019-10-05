package vgo

import (
	"fmt"
	"time"
)

var attributes = map[string]string{
	"name":   "نام",
	"age":    "سن",
	"family": "خانواده",
}

var translations = map[string]string{
	"present":            "فیلد %s باید در پارامترهای ارسالی وجود داشته باشد.",
	"required":           "فیلد %s الزامی است.",
	"requiredWith":       "در صورت وجود فیلد %v، فیلد %s نیز الزامی است.",
	"requiredWithAll":    "در صورت وجود فیلدهای %v، فیلد %s نیز الزامی است.",
	"requiredWithout":    "در صورت عدم وجود فیلد %v، فیلد %s الزامی است.",
	"requiredWithoutAll": "در صورت عدم وجود فیلدهای %v، فیلد %s الزامی است.",
	"confirmed":          "%s با فیلد تکرار مطابقت ندارد.",
	"none":               "فیلد %s اشتباه است.",
	"same":               "%s و %s باید همانند هم باشند.",
	"different":          "%s و %s باید از یکدیگر متفاوت باشند.",


	"string.national":     "فیلد %s باید یک کد ملی معتبر باشد.",
	"string.filled":     "فیلد %s باید مقدار داشته باشد.",
	"string.in":         "%s انتخاب شده، معتبر نیست.",
	"string.inArray":    "فیلد %s در لیست %s وجود ندارد.",
	"string.notIn":      "%s انتخاب شده، معتبر نیست.",
	"string.url":        "%s معتبر نمی‌باشد.",
	"string.uuid":       "%s باید یک UUID معتبر باشد.",
	"string.email":      "%s باید یک ایمیل معتبر باشد.",
	"string.mobile":      "%s باید یک شماره موبایل معتبر باشد.",
	"string.phone":      "%s باید یک شماره تلفن معتبر باشد.",
	"string.ip":         "%s باید آدرس IP معتبر باشد.",
	"string.ipv4":       "%s باید یک آدرس معتبر از نوع IPv4 باشد.",
	"string.ipv6":       "%s باید یک آدرس معتبر از نوع IPv6 باشد.",
	"string.json":       "فیلد %s باید یک رشته از نوع JSON باشد.",
	"string.size":       "%s باید برابر با %v کاراکتر باشد.",
	"string.min":        "%s نباید کمتر از %v کاراکتر داشته باشد.",
	"string.max":        "%s نباید بیشتر از %v کاراکتر داشته باشد.",
	"string.between":    "%s باید بین %v و %v کاراکتر باشد.",
	"string.regex":      "فرمت %s معتبر نیست.",
	"string.username":   "%s باید فقط حروف الفبا، اعداد، خط تیره و زیرخط باشد.",
	"string.alphaNum":   "%s باید فقط حروف الفبا و اعداد باشد.",
	"string.persian":    "%s باید فقط حروف الفبای فارسی باشد.",
	"string.alpha":      "%s باید فقط حروف الفبا باشد.",
	"string.startsWith": "%s باید با یکی از این ها شروع شود: %s",
	"string.endsWith":   "فیلد %s باید با یکی از مقادیر زیر خاتمه یابد: %s",
	"string.contains":   "فیلد %s باید شامل یکی از مقادیر زیر باشد: %s",

	"number.digits": "%s باید %v رقم باشد.",
	"number.digitsBetween": "%s باید بین %v و %v رقم باشد.",

	"number.greaterThan": "%s باید بزرگتر از %v باشد.",
	"number.greaterThanOrEqual": "%s باید بزرگتر یا مساوی %v باشد.",
	"number.lessThan": "%s باید کوچکتر از %v باشد.",
	"number.lessThanOrEqual": "%s باید کوچکتر یا مساوی %v باشد.",
	"number.between": "%s باید بین %v و %v باشد.",
	"number.in":         "%s انتخاب شده، معتبر نیست.",

	"date.after": "%s(%v) باید تاریخی بعد از %v باشد.",
	"date.before": "%s(%v) باید تاریخی قبل از %v باشد.",
	"date.between": "%s(%v) باید تاریخی بین %v و %v باشد.",

	"type.string":       "فیلد %s باید رشته باشد.",
	"type.array":        "%s باید آرایه باشد.",
	"type.object":       "%s باید آبجکت باشد.",
	"type.number":       "%s باید عدد یا رشته‌ای از اعداد باشد.",
	"type.bool":         "فیلد %s فقط می‌تواند true و یا false باشد.",
	"type.file":         "%s باید یک فایل معتبر باشد.",
	"type.image":        "%s باید یک تصویر معتبر باشد.",
	"type.date":         "%s یک تاریخ معتبر نیست.",
	"type.none":         "فیلد %s در سرور اشتباه تعریف شده است.",
}

func translate(typ string, args ...interface{}) string {
	trs, ok := translations[typ]
	if !ok {
		trs = translations["none"]
	}
	return fmt.Sprintf(trs, args...)
}

func translateAttribute(name string) string {
	val, ok := attributes[name]
	if ok {
		return val
	}
	return name
}

func formatDate(t time.Time) string {
	return t.UTC().Format("2006-01-02/15:04")
}
func translateAttributes(name ...string) []string {
	for i, item := range name {
		val, ok := attributes[item]
		if ok {
			name[i] = val
		}
	}
	return name
}

var faToEn = []rune{
	'0',
	'1',
	'2',
	'3',
	'4',
	'5',
	'6',
	'7',
	'8',
	'9',
	'۰',
	'۱',
	'۲',
	'۳',
	'۴',
	'۵',
	'۶',
	'۷',
	'۸',
	'۹',
	'٠',
	'١',
	'٢',
	'٣',
	'٤',
	'٥',
	'٦',
	'٧',
	'٨',
	'٩',
}

func convertToNumber(str string) (float64, bool) {
	if str == "" {
		return 0, false
	}
	strict := true
	var out float64 = 0
	dot := false
	var dah float64 = 10
	var neg = false
	for _, char := range str {
		if char == '.' {
			dot = true
		} else if char == '-' && !neg {
			neg = true
		} else {
			fnd := false
			for j, num := range faToEn {
				if num == char {
					var digit = float64(j % 10)
					if !dot {
						out = out*10 + digit
					} else {
						out = out + digit/dah
						dah *= 10
					}
					fnd = true
					break
				}
			}
			if !fnd {
				strict = false
			}
		}
	}
	if neg {
		return -out, strict
	}
	return out, strict
}
