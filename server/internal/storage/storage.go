package storage

import (
	"bytes"
	"fmt"
	"proxyfinder/pkg/options"
)

var (
	ErrRecordNotFound = fmt.Errorf("Recod not found")
	ErrEmptyOptions   = fmt.Errorf("Empty options")
	ErrInvalidId      = fmt.Errorf("Invalid id")
	ErrInvalidPage    = fmt.Errorf("Invalid page")
	ErrInvalidPerPage = fmt.Errorf("Invalid perPage")

	DefaultLimit  = 40
	DefaultOffset = 0
)

type QueryBuilder struct {
	query   bytes.Buffer
	page    int
	perPage int
	limit   int
	offset  int
	isWhere bool
	isSort  bool
	where   options.Options
	sort    options.Options
	values  []interface{}
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		limit:  DefaultLimit,
		offset: DefaultOffset,
	}
}

func (self *QueryBuilder) Filter(filter options.Options) error {
	var (
		ok bool
	)
	if filter == nil {
		return ErrEmptyOptions
	}
	for _, field := range filter.Fields() {
		if field.Name == "page" {
			self.page, ok = field.Val.(int)
			if !ok {
				return ErrInvalidPage
			}
			continue
		}
		if field.Name == "perPage" {
			self.perPage, ok = field.Val.(int)
			if !ok {
				return ErrInvalidPerPage
			}
			continue
		}
		self.where = filter
	}

	return nil
}

func (self *QueryBuilder) Sort(sort options.Options) {
	if sort == nil {
		return
	}
	self.sort = sort
}

func (self *QueryBuilder) AddLimit(limit int) {
	self.limit = limit
}

func (self *QueryBuilder) AddOffset(offset int) {
	self.offset = offset
}

func (self *QueryBuilder) SetFilter(filter options.Options) {
	self.where = filter
}

func (self *QueryBuilder) SetSort(sort options.Options) {
	self.sort = sort
}

func (self *QueryBuilder) buildWhere() {
	if self.where == nil {
		return
	}
	for _, v := range self.where.Fields() {
		if v.Name == "page" || v.Name == "perPage" {
			continue
		}
		switch v.Val.(type) {
		case []string:
			self.addFilteSlice(v)
			break
		default:
			self.addFilterString(v)
		}
	}
}

func (self *QueryBuilder) addFilterString(field options.Field) {
	if self.isWhere {
		self.query.WriteString(fmt.Sprintf("AND %s %s ?", field.Name, field.Op))
	} else {
		self.query.WriteString(fmt.Sprintf("WHERE %s %s ?", field.Name, field.Op))
		self.isWhere = true
	}
	self.values = append(self.values, field.Val)
}

func (self *QueryBuilder) addFilteSlice(field options.Field) {
	if self.isWhere {
		self.query.WriteString(fmt.Sprintf("AND %s %s (", field.Name, field.Op))
	} else {
		self.query.WriteString(fmt.Sprintf("WHERE %s %s (", field.Name, field.Op))
		self.isWhere = true
	}
	for i, str := range field.Val.([]string) {
		self.values = append(self.values, str)
		if i > 0 {
			self.query.WriteString(", ")
		}
		self.query.WriteString(`?`)
	}
	self.query.WriteString(")")
}

func (self *QueryBuilder) buildSort() {
	for _, v := range self.sort.Fields() {
		if self.isSort {
			self.query.WriteString(fmt.Sprintf(", %s %s", v.Name, v.Val.(string)))
		} else {
			self.query.WriteString(fmt.Sprintf(" ORDER BY %s %s", v.Name, v.Val.(string)))
			self.isSort = true
		}
	}
}

func (self *QueryBuilder) buildLimit() {
	if self.limit == 0 {
		if self.perPage == 0 {
			return
		}
		self.limit = self.perPage
	}

	self.query.WriteString(fmt.Sprintf(" LIMIT %d", self.limit))
}

func (self *QueryBuilder) buildOffset() {
	if self.limit == 0 {
		return
	}
	if self.offset == 0 {
		if self.page == 0 {
			return
		}
		self.offset = (self.page - 1) * self.perPage
	}

	self.query.WriteString(fmt.Sprintf(" OFFSET %d", self.offset))
}

func (self *QueryBuilder) BuildQuery(query string) string {
	self.buildWhere()
	self.buildSort()
	self.buildLimit()
	self.buildOffset()
	// fmt.Println(fmt.Sprintf("%s %s", query, self.query.String()))
	return fmt.Sprintf("%s %s", query, self.query.String())
}

func (self *QueryBuilder) Values() []interface{} {
	return self.values
}
