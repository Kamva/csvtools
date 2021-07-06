package main

import (
	"fmt"

	"github.com/jszwec/csvutil"
	"github.com/kamva/csvtools"
)

//--------------------------------
// Define groups
//--------------------------------

var BookAuthorGroup = csvtools.NewGroup(csvtools.GroupOpts{
	SetVal:     func(r csvtools.Record, val interface{}) { r.(*CSVBook).Author = val.(*CSVAuthor) },
	Val:        func(r csvtools.Record) interface{} { return r.(*CSVBook).Author },
	ValIsEmpty: func(r csvtools.Record) bool { return r.(*CSVBook).Author.Name == "" },
	List:       func(r csvtools.Record) []interface{} { return csvtools.InterfaceToSlice(r.(*CSVBook).Authors) },
	SetList: func(r csvtools.Record, values []interface{}) {
		l := make([]*CSVAuthor, len(values))
		for i, v := range values {
			l[i] = v.(*CSVAuthor)
		}
		r.(*CSVBook).Authors = l
	},
})

var bookAgeGroup = csvtools.NewGroup(csvtools.GroupOpts{
	SetVal:     func(r csvtools.Record, val interface{}) { r.(*CSVBook).Age = val.(int) },
	Val:        func(r csvtools.Record) interface{} { return r.(*CSVBook).Age },
	ValIsEmpty: func(r csvtools.Record) bool { return r.(*CSVBook).Age == 0 },
	List:       func(r csvtools.Record) []interface{} { return csvtools.InterfaceToSlice(r.(*CSVBook).Ages) },
	SetList: func(r csvtools.Record, values []interface{}) {
		l := make([]int, len(values))
		for i, v := range values {
			l[i] = v.(int)
		}
		r.(*CSVBook).Ages = l
	},
})

type CSVBook struct {
	Name  string `json:"name" csv:"name"` // Key must be pointer.
	Color string `json:"color" csv:"color,omitempty"`
	Ages  []int  `json:"ages" csv:"-"`
	Age   int    `json:"age" csv:"age"`

	// Per each Array we must have two fields: array must be ignored.
	Authors []*CSVAuthor `json:"authors" csv:"-"`             // convert array to interface.
	Author  *CSVAuthor   `json:"author" csv:"author.,inline"` // convert array to interface.
}

type CSVAuthor struct {
	Name  string `json:"name" csv:"name"`   // key must be a pointer.
	Score *int   `json:"score" csv:"score"` // change to pointer to prevent form printing zero av value.
}

func (c *CSVBook) CloneEmptyRecord() csvtools.Record {
	return &CSVBook{Name: c.Name}
}

func (c *CSVBook) CSVKey() interface{} {
	return c.Name
}

func (c *CSVBook) CSVGroups() []csvtools.Group {
	return []csvtools.Group{BookAuthorGroup, bookAgeGroup} // just child groups.
}

func bookRecordsToStruct(records []csvtools.Record) []*CSVBook {
	l := make([]*CSVBook, len(records))
	for i, v := range records {
		l[i] = v.(*CSVBook)
	}
	return l
}

func main() {
	books := []*CSVBook{
		{
			Name:  "book a",
			Color: "white",
			Ages:  []int{1, 2, 10},
			//Age:  []int{1, 2, 3},
			Authors: []*CSVAuthor{
				{
					Name:  "ali",
					Score: newInt(3),
				},
				{
					Name:  "reza",
					Score: newInt(4),
				},
			},
		},

		{
			Name:  "book b",
			Color: "red",
			Ages:  []int{2, 4, 6},
			//Age:  []int{1, 2, 3},
			Authors: []*CSVAuthor{
				{
					Name:  "John",
					Score: newInt(5),
				},
				{
					Name:  "Jessy",
					Score: newInt(6),
				},
			},
		},
	}

	//--------------------------------
	// Marshal
	//--------------------------------

	records, err := csvtools.DecomposeRecords(csvtools.ToRecordsSlice(books))
	if err != nil {
		panic(err)
	}
	bytes, err := csvutil.Marshal(bookRecordsToStruct(records))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))

	//--------------------------------
	// Unmarshal
	//--------------------------------

	books = make([]*CSVBook, 0) // to keep types, unmarshal to its original type first.
	if err := csvutil.Unmarshal(bytes, &books); err != nil {
		panic(err)
	}
	bookRecords, err := csvtools.ComposeRecords(csvtools.ToRecordsSlice(books))
	if err != nil {
		panic(err)
	}
	fmt.Println(bookRecordsToStruct(bookRecords))
}

var _ csvtools.Record = &CSVBook{}

// util
func newInt(v int) *int {
	return &v
}
