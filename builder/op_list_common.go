package builder

func init() {
	_opList[setKey] = &Op{name: setKey}
}

func Op_Set() *Op {
	return _opList[setKey]
}
