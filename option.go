package pdf

import (
	"fmt"
	"strconv"
)

type optType int

const (
	optTypeString optType = iota + 1
	optTypeBool
	optTypeInt
	optTypeUint
	optTypeFloat
)

type setterFunc func(name, value string) error

type setOp struct {
	name     string
	value    interface{}
	typ      optType
	setter   setterFunc
	setEmpty bool
}

func newSetOp(name string, value interface{}, typ optType, setter setterFunc, setEmpty bool) *setOp {
	return &setOp{
		name:     name,
		value:    value,
		typ:      typ,
		setter:   setter,
		setEmpty: setEmpty,
	}
}

func (op *setOp) execute() error {
	switch op.typ {
	case optTypeString:
		if val := op.value.(string); op.setEmpty || val != "" {
			return op.setter(op.name, val)
		}
	case optTypeBool:
		if val := op.value.(bool); op.setEmpty || val {
			return op.setter(op.name, strconv.FormatBool(val))
		}
	case optTypeInt:
		if val := op.value.(int64); op.setEmpty || val > 0 {
			return op.setter(op.name, strconv.FormatInt(val, 10))
		}
	case optTypeUint:
		if val := op.value.(uint64); op.setEmpty || val > 0 {
			return op.setter(op.name, strconv.FormatUint(val, 10))
		}
	case optTypeFloat:
		if val := op.value.(float64); op.setEmpty || uint64(val) > 0 {
			return op.setter(op.name, strconv.FormatFloat(val, 'E', -1, 64))
		}
	default:
		return fmt.Errorf("invalid option type: %d", op.typ)
	}

	return nil
}
