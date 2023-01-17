package builder

const (
	//https://www.mongodb.com/docs/manual/reference/operator/query-comparison/

	//Matches values that are equal to a specified value.
	op_comparison_eq string = "$eq"
	//Matches values that are greater than a specified value.
	op_comparison_gt string = "$gt"
	//Matches values that are greater than or equal to a specified value.
	op_comparison_gte string = "$gte"
	//Matches any of the values specified in an array.
	op_comparison_in string = "$in"
	//Matches values that are less than a specified value.
	op_comparison_lt string = "$lt"
	//Matches values that are less than or equal to a specified value.
	op_comparison_lte string = "$lte"
	//Matches all values that are not equal to a specified value.
	op_comparison_ne string = "$ne"
	//Matches none of the values specified in an array.
	op_comparison_nin string = "$nin"
)

type OpList struct {
}

func (l *OpList) Eq() string {
	return op_comparison_eq
}

func (l *OpList) Gt() string {
	return op_comparison_gt
}

func (l *OpList) Gte() string {
	return op_comparison_gte
}

func (l *OpList) In() string {
	return op_comparison_in
}

func (l *OpList) Lt() string {
	return op_comparison_lt
}

func (l *OpList) Lte() string {
	return op_comparison_lte
}

func (l *OpList) Ne() string {
	return op_comparison_ne
}

func (l *OpList) Nin() string {
	return op_comparison_nin
}
