package options

import (
	"errors"
)

const (
	OptionsCtxKey = "optionsCtxKey"

	OpEq      = "="
	OpLt      = "<"
	OpLte     = "<="
	OpGt      = ">"
	OpGte     = ">="
	OpNe      = "!="
	OpLike    = "like"
	OpIn      = "in"
	OpNotIn   = "not in"
	OpBetween = "between"

	ErrInvalidOp = "Invalid operator"
)

type options struct {
	is        bool
	fields    []Field
}

type Field struct {
	Name string
	Op   string
	Val  interface{}
}

func New() Options {
	return &options{
		is:        false,
		fields:    []Field{},
	}
}

type Options interface {
	Fields() []Field
	AddField(name string, op string, val interface{}) error
	SetField(name string, op string, val interface{}) error
	GetField(name string) (Field, error)
	MapField(func(*Field) error) error
	Values() []interface{}

	Is() bool
}

// AddField adds a field to the options and validates the operator
func (self *options) AddField(name string, op string, val interface{}) error {
	if name == "" {
		return errors.New("field name cannot be empty")
	}

	err := ValidateOperator(op)
	if err != nil {
		return err
	}

	if !self.is {
		self.is = true
	}

	newField := Field{
		Name: name,
		Op:   op,
		Val:  val,
	}

	self.fields = append(self.fields, newField)

	return nil
}

func (self *options) SetField(name string, op string, val interface{}) error {
	err := ValidateOperator(op)
	if err != nil {
		return err
	}

	for i := range self.fields {
		if self.fields[i].Name == name {
			self.fields[i].Op = op
			self.fields[i].Val = val
			return nil
		}
	}

	return nil
}

func (self *options) GetField(name string) (Field, error) {
	for i := range self.fields {
		if self.fields[i].Name == name {
			return self.fields[i], nil
		}
	}

	return Field{}, errors.New("field not found")
}

func (self *options) MapField(fn func(*Field) error) error {
	for i := range self.fields {
		err := fn(&self.fields[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *options) Values() []interface{} {
	var values []interface{}
	for _, field := range self.fields {
		values = append(values, field.Val)
	}
	return values
}

func (self *options) Is() bool {
	return self.is
}

func (self *options) Fields() []Field {
	return self.fields
}

func ValidateOperator(op string) error {
	switch op {
	case OpEq, OpLt, OpLte, OpGt, OpGte, OpNe, OpLike, OpIn, OpNotIn:
		return nil
	default:
		return errors.New(ErrInvalidOp)
	}
}
