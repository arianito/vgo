package vgo

import (
	"strings"
	"time"
)

func removeSlashes(val string) string {
	val = strings.ReplaceAll(val, "\\)", ")")
	val = strings.ReplaceAll(val, "\\(", "(")
	val = strings.ReplaceAll(val, "\\,", ",")
	return strings.Trim(val, " ")
}
func parseDate(val string) (time.Time, error)  {
	if val == "now" {
		return time.Now().UTC(), nil
	}
	if val == "today" {
		t := time.Now().UTC()
		year, month, day := t.Date()
		return time.Date(year, month, day, 0, 0, 0, 0, time.UTC), nil
	}
	if val == "yesterday" {
		t := time.Now().UTC().Add(time.Hour*-24)
		year, month, day := t.Date()
		return time.Date(year, month, day, 0, 0, 0, 0, time.UTC), nil
	}
	if val == "tomorrow" {
		t := time.Now().UTC().Add(time.Hour*24)
		year, month, day := t.Date()
		return time.Date(year, month, day, 0, 0, 0, 0, time.UTC), nil
	}
	tme, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return tme, err
	}
	return tme.UTC(), nil
}

func evalRuleChain(rule string, call func(name string, args ...string) bool) {
	scopeBlock := 0
	cut := 0
	inBlock := false
	before := ' '
	ln := len(rule)
	paramStart := 0
	argStart := 0
	haveArguments := false
	var args []string
	for i, char := range rule {
		if before != '\\' {
			if char == '(' {
				if scopeBlock == 0 {
					haveArguments = true
					paramStart = i
					argStart = i + 1
					inBlock = true
				}
				scopeBlock++
			} else if char == ')' {
				scopeBlock--
				if scopeBlock == 0 {
					args = append(args, removeSlashes(rule[argStart:i]))
					argStart = i
					inBlock = false
					name := rule[cut-1 : paramStart]
					if !call(name, args...) {
						return
					}
					args = []string{}
				}
			} else if char == ',' && inBlock {
				args = append(args, removeSlashes(rule[argStart:i]))
				argStart = i + 1
			}
		}
		if !inBlock {
			if before != ' ' {
				if char == ' ' || i == ln-1 {
					begin := cut - 1
					if i == ln-1 {
						cut = ln + 1
					} else {
						cut = i + 1
					}
					if !haveArguments {
						name := rule[begin : cut-1]
						if !call(name) {
							return
						}
					}
					haveArguments = false
				}
			} else {
				cut++
			}
		}
		before = char
	}
}

