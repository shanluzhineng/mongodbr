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

func (l *OpList) And() string {
	return op_comparison_and
}

func (l *OpList) Not() string {
	return op_comparison_not
}

func (l *OpList) Nor() string {
	return op_comparison_nor
}

func (l *OpList) Or() string {
	return op_comparison_or
}
