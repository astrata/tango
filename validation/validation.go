package validation

import (
	"fmt"
	"github.com/astrata/tango"
	"regexp"
	"strconv"
)

// Validation function.
type Rule func(string) error

// A set of rules to be applied on a single variable.
type Constraint struct {
	All     []Rule
	Message string
}

// A set of validation rules.
type Rules struct {
	Map map[string]Constraint
}

// Returns a new set of validation rules.
func New() *Rules {
	self := &Rules{}
	self.Map = make(map[string]Constraint)
	return self
}

// Adds a rule to a set of constraints.
func (self *Constraint) Add(rule Rule) {
	self.All = append(self.All, rule)
}

// Validates input data.
func (self *Rules) Validate(params tango.Value) (bool, map[string][]string) {
	valid := true
	messages := map[string][]string{}

	for key, _ := range params {
		value := params.GetString(key)
		if constraint, ok := self.Map[key]; ok == true {
			errors := []string{}
			passed := true
			for _, rule := range constraint.All {
				test := rule(value)
				if test != nil {
					passed = false
					errors = append(errors, test.Error())
				}
			}
			if passed == false {
				messages[key] = errors
				valid = false
			}
		}
	}

	return valid, messages
}

func (self *Rules) Add(name string, rule Rule, message string) {
	var constraint Constraint
	var ok bool

	if constraint, ok = self.Map[name]; ok == false {
		constraint = Constraint{}
		constraint.All = []Rule{}
		constraint.Message = "An error ocurred"
	}

	constraint.Message = message
	constraint.Add(rule)

	self.Map[name] = constraint
}

// A rule that returns error if the value is empty.
func NotEmpty(value string) error {
	if value == "" {
		return fmt.Errorf("This value is required")
	}
	return nil
}

// A rule that returns error if the value is not an URL.
func Url(value string) error {
	match := MatchExpr(value, `(?i)^[a-z]+:\/\/[a-z0-9][a-z0-9\-\.]*`)
	if match == nil {
		return nil
	}
	return fmt.Errorf("Value must be an URL.")
}

// A rule that returns error if the value is not a BSON ObjectId.
func ObjectId(value string) error {
	match := MatchExpr(value, `^[a-f0-9]{24}$`)
	if match == nil {
		return nil
	}
	return fmt.Errorf("Expecting an ObjectId.")
}

// A rule that returns error if the value is not a-zA-Z0-9.
func Alpha(value string) error {
	match := MatchExpr(value, `(?i)^[a-z0-9]+$`)
	if match == nil {
		return nil
	}
	return fmt.Errorf("Value must be a number or a letter from A to Z (case does not matter).")
}

// A rule that returns error if the value is not an e-mail.
func Email(value string) error {
	passed := MatchExpr(value, `(?i)^[a-z0-9][a-z0-9\.\-+_]*@[a-z0-9\-\.]+.[a-z]+$`)
	if passed == nil {
		return nil
	}
	return fmt.Errorf("Value must be an e-mail address.")
}

// A rule that returns error if the value is not numeric.
func Numeric(value string) error {
	passed := MatchExpr(value, `(?i)^[0-9]+$`)
	if passed == nil {
		return nil
	}
	return fmt.Errorf("Value must be a number.")
}

// A rule that returns error if the value does not match a pattern.
func MatchExpr(value string, expr string) error {
	match, _ := regexp.MatchString(expr, value)
	if match == true {
		return nil
	}
	return fmt.Errorf("Value does not match pattern %s.", expr)
}
