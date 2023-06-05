package builder

const (
	//https://www.mongodb.com/docs/manual/reference/operator/update-array/

	// Adds elements to an array only if they do not already exist in the set.
	op_array_addToSet string = "$addToSet"
	// Removes the first or last item of an array.
	op_array_pop string = "$pop"
	// Removes all array elements that match a specified query.
	op_array_pull string = "$pull"
	// Adds an item to an array.
	op_array_push string = "$push"
	// Removes all matching values from an array.
	op_array_pullAll string = "$pullAll"

	//The $elemMatch operator matches documents that contain an array field with at least one element
	// that matches all the specified query criteria.
	op_array_elemMatch string = "$elemMatch"
)

func init() {
	_opList[op_array_addToSet] = &Op{name: op_array_addToSet}
	_opList[op_array_pop] = &Op{name: op_array_pop}
	_opList[op_array_pull] = &Op{name: op_array_pull}
	_opList[op_array_push] = &Op{name: op_array_push}
	_opList[op_array_pullAll] = &Op{name: op_array_pullAll}
	_opList[op_array_elemMatch] = &Op{name: op_array_elemMatch}
}

func Op_AddToSet() *Op {
	return _opList[op_array_addToSet]
}

func Op_Pop() *Op {
	return _opList[op_array_pop]
}

func Op_Pull() *Op {
	return _opList[op_array_pull]
}

func Op_Push() *Op {
	return _opList[op_array_push]
}

func Op_PullAll() *Op {
	return _opList[op_array_pullAll]
}

func Op_ElemMatch() *Op {
	return _opList[op_array_elemMatch]
}
