package builder

const (
	//https://www.mongodb.com/docs/manual/reference/operator/query-element/

	//Matches documents that have the specified field.
	op_comparison_exists string = "$exists"
	//Selects documents if a field is of the specified type.
	op_comparison_type string = "$type"
)

func (l *OpList) Exists() string {
	return op_comparison_exists
}

func (l *OpList) Type() string {
	return op_comparison_type
}
