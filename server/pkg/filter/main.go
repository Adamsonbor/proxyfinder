package filter

import (
	"errors"
)

const (
	FilterCtxKey = "filterOptions"

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
	page    int
	perPage int
	limit   int
	offset  int
	is      bool

	sort   []Sort
	fields []Field
}

type Field struct {
	Name string
	Op   string
	Val  interface{}
	Type string
}

type Sort struct {
	By    string
	Order string
}

func New() Options {
	return &options{
		page:    1,
		perPage: 10,
		limit:   10,
		offset:  0,
		is:      false,
		fields:  []Field{},
		sort:    []Sort{},
	}
}

type Options interface {
	Page() int
	SetPage(int)
	PerPage() int
	SetPerPage(int)
	NextPage()

	Limit() int
	SetLimit(limit int)
	Offset() int
	SetOffset(offset int)
	UpdateLimitAndOffset()

	Sort() []Sort
	AddSort(Sort)
	MapSort(func(*Sort) error) error

	Fields() []Field
	AddField(name string, op string, val interface{}, dtype string) error
	SetField(id int, name string, op string, val interface{}, dtype string)
	MapField(func(*Field) error) error

	Values() []interface{}

	Is() bool
}

func (self *options) Sort() []Sort {
	return self.sort
}

func (self *options) AddSort(sort Sort) {
	if sort.Order != "asc" && sort.Order != "desc" {
		sort.Order = "asc"
	}
	self.sort = append(self.sort, sort)
}

func (self *options) MapSort(fn func(*Sort) error) error {
	for i := range self.sort {
		err := fn(&self.sort[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *options) Page() int {
	return self.page
}

// SetPage sets the current page
// If page is less than 1, it will be set to 1
func (self *options) SetPage(page int) {
	if page < 1 {
		page = 1
	}
	self.page = page
}

func (self *options) PerPage() int {
	return self.perPage
}

// SetPerPage sets the number of items per page
// If perPage is less than 1, it will be set to 1
func (self *options) SetPerPage(perPage int) {
	if perPage < 1 {
		perPage = 10
	}
	self.perPage = perPage
}

// NextPage increments the current page
// and updates the limit and offset
func (self *options) NextPage() {
	self.page++
	self.UpdateLimitAndOffset()
}

func (self *options) UpdateLimitAndOffset() {
	self.limit = self.perPage
	self.offset = (self.page - 1) * self.perPage
}

// AddField adds a field to the options and validates the operator
func (self *options) AddField(name string, op string, val interface{}, dtype string) error {
	err := ValidateOperator(op)
	if err != nil {
		return err
	}

	if !self.is {
		self.is = true
	}

	self.fields = append(self.fields, Field{
		Name: name,
		Op:   op,
		Val:  val,
		Type: dtype,
	})

	return nil
}

func (self *options) SetField(id int, name string, op string, val interface{}, dtype string) {
	self.fields[id] = Field{
		Name: name,
		Op:   op,
		Val:  val,
		Type: dtype,
	}
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

func (self *options) SetLimit(limit int) {
	self.limit = limit
}

func (self *options) SetOffset(offset int) {
	self.offset = offset
}

func (self *options) Limit() int {
	return self.limit
}

func (self *options) Offset() int {
	return self.offset
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
