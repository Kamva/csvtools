package csvtools

import "reflect"

type Record interface {
	CSVKey() interface{}
	CSVGroups() []Group
	CloneEmptyRecord() Record
}

type Group interface {
	SetVal(r Record, val interface{})
	Val(r Record) interface{}
	ValIsEmpty(r Record) bool
	List(r Record) []interface{}
	SetList(r Record, values []interface{})
}

type GroupOpts struct {
	SetVal     func(r Record, val interface{})
	Val        func(r Record) interface{}
	ValIsEmpty func(r Record) bool
	List       func(r Record) []interface{}
	SetList    func(r Record, values []interface{})
}

type fnGroup struct {
	setVal     func(r Record, val interface{})
	val        func(r Record) interface{}
	valIsEmpty func(r Record) bool
	list       func(r Record) []interface{}
	setList    func(r Record, values []interface{})
}

func (f *fnGroup) SetVal(r Record, val interface{}) {
	f.setVal(r, val)
}

func (f *fnGroup) Val(r Record) interface{} {
	return f.val(r)
}

func (f *fnGroup) ValIsEmpty(r Record) bool {
	return f.valIsEmpty(r)
}

func (f *fnGroup) List(r Record) []interface{} {
	return f.list(r)
}

func (f *fnGroup) SetList(r Record, l []interface{}) {
	f.setList(r, l)
}

func NewGroup(o GroupOpts) Group {
	return &fnGroup{
		setVal:     o.SetVal,
		val:        o.Val,
		valIsEmpty: o.ValIsEmpty,
		list:       o.List,
		setList:    o.SetList,
	}
}

type decomposedItems struct {
	group Group
	list  []interface{}
}

func DecomposeRecords(l []Record) ([]Record, error) {
	res := make([]Record, 0)
	for _, r := range l {
		records, err := DecomposeRecord(r)
		if err != nil {
			return nil, err
		}

		res = append(res, records...)
	}
	return res, nil
}

func DecomposeRecord(parent Record) ([]Record, error) { // parent should be an interface.
	var records []Record // []CSVBook // Record
	records = append(records, parent)

	var decomposed []*decomposedItems
	for _, g := range parent.CSVGroups() { // replace with Record.CSVGroups()
		_, childIsRecord := g.Val(parent).(Record)
		if !childIsRecord {
			decomposed = append(decomposed, &decomposedItems{group: g, list: g.List(parent)})
			continue
		}

		for _, v := range g.List(parent) { // loop over structs
			l, err := DecomposeRecord(v.(Record)) // check if it has been implemented Record interface
			if err != nil {
				return nil, err
			}
			decomposed = append(decomposed, &decomposedItems{group: g, list: InterfaceToSlice(l)})
		}
	}

	return cloneDecomposedRecords(parent, decomposed), nil
}

func ComposeRecords(records []Record) ([]Record, error) {
	if len(records) == 0 {
		return nil, nil
	}
	g, err := groupRecordsByKey(records)
	if err != nil {
		return nil, err
	}
	res := make([]Record, len(g))
	for i, recordsGroup := range g {
		var err error
		if res[i], err = ComposeGroupedRecords(recordsGroup); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func ComposeGroupedRecords(records []Record) (Record, error) {
	if len(records) == 0 {
		return nil, nil
	}

	first := records[0]
	for _, g := range first.CSVGroups() {
		var values []interface{}
		children := extractChildren(records, g)
		_, childrenAreRecord := g.Val(records[0]).(Record)

		// if children are not group records, so every record is simply a value of list:
		if !childrenAreRecord {
			g.SetList(first, children)
			continue
		}

		recordGroups, err := groupRecordsByKey(ToRecordsSlice(children))
		if err != nil {
			return nil, err
		}
		for _, rg := range recordGroups {
			m, err := ComposeGroupedRecords(rg)
			if err != nil {
				return nil, err
			}
			if m != nil {
				values = append(values, m)
			}
		}
		g.SetList(first, values)

	}
	return first, nil
}

func cloneDecomposedRecords(parent Record, decomposed []*decomposedItems) []Record {
	var records []Record

	reachedEnd := false
	for i := 0; !reachedEnd; i++ {
		reachedEnd = true
		clone := parent.CloneEmptyRecord()
		if i == 0 {
			clone = parent
		}

		for _, d := range decomposed {
			if i < len(d.list) {
				d.group.SetVal(clone, d.list[i])
			}

			if i+1 < len(d.list) { // If this is not last index for this decomposed list, continue
				reachedEnd = false
			}
		}
		records = append(records, clone)
	}

	return records
}

func extractChildren(records []Record, g Group) []interface{} {
	var children []interface{}
	for _, r := range records {
		if !g.ValIsEmpty(r) {
			children = append(children, g.Val(r))
		}
	}
	return children
}

// groupRecordsByKey groups by key. if key is empty, so every record is a record group.
func groupRecordsByKey(records []Record) ([][]Record, error) {
	if len(records) == 0 {
		return nil, nil
	}

	recordGroups := make([][]Record, 0)
	var rg []Record
	var prev = records[0].CSVKey()
	for _, r := range records {
		// if current key is different than previous, we have a new record group.
		if !reflect.DeepEqual(prev, r.CSVKey()) {
			recordGroups = append(recordGroups, rg)
			rg = make([]Record, 0)
		}

		rg = append(rg, r)
		prev = r.CSVKey()
	}
	recordGroups = append(recordGroups, rg) // Add last record group

	return recordGroups, nil
}
