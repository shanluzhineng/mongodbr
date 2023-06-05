package builder

const (
	//https://www.mongodb.com/docs/manual/reference/operator/query-logical/

	//Joins query clauses with a logical AND returns all documents that match the conditions of both clauses.
	op_comparison_and string = "$and"
	//Inverts the effect of a query expression and returns documents that do not match the query expression.
	op_comparison_not string = "$not"
	//Joins query clauses with a logical NOR returns all documents that fail to match both clauses.
	op_comparison_nor string = "$nor"
	//Joins query clauses with a logical OR returns all documents that match the conditions of either clause.
	op_comparison_or string = "$or"
)

func init() {
	_opList[op_comparison_and] = &Op{name: op_comparison_and}
	_opList[op_comparison_not] = &Op{name: op_comparison_not}
	_opList[op_comparison_nor] = &Op{name: op_comparison_nor}
	_opList[op_comparison_or] = &Op{name: op_comparison_or}
}

func Op_And() *Op {
	return _opList[op_comparison_and]
}

func Op_Not() *Op {
	return _opList[op_comparison_not]
}

func Op_Nor() *Op {
	return _opList[op_comparison_nor]
}

func Op_Or() *Op {
	return _opList[op_comparison_or]
}
