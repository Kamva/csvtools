`csvtools` Provide tools to work with csv files.

#### Available tools:
- __Flattener__: Flat your structs to csv files and merge them after unmarshal.




## Flattener:

#### How to use?

__Marshal__
- Each struct that contains arrays must implement the `Record` interface.
- For each array field in your struct add another field with type of your array field's type and also set tag of that array field to dash(`-`) to ignore for csv. e.g.,
```go
struct Book{
	Authors `csv:"-"` // please note csv tag value is dash.
	Author `csv:"author_,inline"`
}
```
  
- For each Array field in your struct create a new instance of `Group` and return created groups in the struct as return param of `CSVGroups()` method.
  
- call to `DecomposeRecords` do decompose your records.
- Simply call to the `csvutil.Marshal` method to marshal your list.

__Unmarshal__
- Simply use `csvutil.Unmarshal` to unmarshal your csv records to your struct list.

- Call to the `ComposeRecords` to group records and then merge them. it will return array of merged records.

