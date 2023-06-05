package builder

const (
	//https://www.mongodb.com/docs/manual/reference/operator/query-element/

	//Matches documents that have the specified field.
	op_comparison_exists string = "$exists"
	//Selects documents if a field is of the specified type.
	op_comparison_type string = "$type"
)

func init() {
	_opList[op_comparison_exists] = &Op{name: op_comparison_exists}
	_opList[op_comparison_type] = &Op{name: op_comparison_type}
}

func Op_Exists() *Op {
	return _opList[op_comparison_exists]
}

func Op_Type() *Op {
	return _opList[op_comparison_type]
}
