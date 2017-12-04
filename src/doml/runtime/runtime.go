package runtime

// This file exists as the runtime for DOML
// Allows it to effectively run IR
import (
	"errors"
	"reflect"
)

type runtime struct {
	stack []interface{}
	// Should always begin at -1 to indicate that the array is empty
	stackPtr  int
	registers []interface{}
}

func NewRuntime(stackSize, registerSize int) *runtime {
	return &runtime{make([]interface{}, stackSize), -1, make([]interface{}, registerSize)}
}

func NewEmptyRuntime() *runtime {
	return &runtime{stackPtr: -1}
}

func (runtime *runtime) CurrentStackSize() int {
	return runtime.stackPtr + 1
}

func (runtime *runtime) MaxStackSize() int {
	return len(runtime.stack)
}

func (runtime *runtime) RegisterSize() int {
	return len(runtime.registers)
}

// ClearSpace() -> Clears the stack
// Returns true if stack was cleared, if stack already empty returns false
func (runtime *runtime) ClearSpace() bool {
	if runtime.CurrentStackSize() > 0 {
		runtime.stack = make([]interface{}, runtime.MaxStackSize())
		return true
	}

	runtime.stackPtr = 0
	return false
}

func (runtime *runtime) ClearRegisters() bool {
	if runtime.RegisterSize() > 0 {
		runtime.registers = make([]interface{}, runtime.RegisterSize())
		return true
	}

	return false
}

func (runtime *runtime) ReserveRegister(space int) bool {
	if space > runtime.RegisterSize() {
		runtime.registers = make([]interface{}, space)
		return true
	}

	return false
}

func (runtime *runtime) ReserveSpace(space int) bool {
	if space > runtime.MaxStackSize() {
		runtime.stack = make([]interface{}, space)
		return true
	}

	return false
}

func (runtime *runtime) GetObject(index int) (obj interface{}, err error) {
	if index >= 0 && index <= runtime.RegisterSize() {
		obj = runtime.registers[index]
		if obj == nil {
			err = errors.New("object is null")
		}
	} else {
		obj = nil
		err = errors.New("index out of range")
	}

	return
}

func (runtime *runtime) SetObject(index int, obj interface{}) error {
	if index >= 0 && index <= runtime.RegisterSize() && obj != nil {
		runtime.registers[index] = obj
		return nil
	} else {
		return errors.New("object is null or index is out of range")
	}
}

func (runtime *runtime) RemoveObject(index int) error {
	if index >= 0 && index <= runtime.RegisterSize() {
		runtime.registers[index] = nil
		return nil
	} else {
		return errors.New("index out of range")
	}
}

func (runtime *runtime) Push(obj interface{}, resizeIfNoSpace bool) error {
	if runtime.CurrentStackSize() == runtime.MaxStackSize() {
		if resizeIfNoSpace {
			temp := runtime.stack
			runtime.stack = make([]interface{}, runtime.MaxStackSize()+1)
			copy(temp, runtime.stack)
		} else {
			return errors.New("no space on stack for object")
		}
	}

	runtime.stackPtr += 1
	runtime.stack[runtime.stackPtr] = obj
	return nil
}

func (runtime *runtime) PopNoReturn() error {
	if runtime.CurrentStackSize() > 0 {
		if runtime.stack[runtime.stackPtr] == nil {
			return errors.New("object is null")
		}
	} else {
		return errors.New("nothing to pop off stack")
	}

	runtime.stackPtr -= 1
	return nil
}

func (runtime *runtime) Pop() (obj interface{}, err error) {
	obj, err = runtime.Peek()
	if err == nil {
		runtime.stackPtr -= 1
	}
	return
}

func (runtime *runtime) Peek() (obj interface{}, err error) {
	err = nil
	obj = nil
	if runtime.CurrentStackSize() > 0 {
		obj = runtime.stack[runtime.stackPtr]
		if obj == nil {
			err = errors.New("object is null")
		}
	} else {
		err = errors.New("nothing to pop off stack")
	}
	return
}

func (runtime *runtime) TopIsOfType(typeToCheck reflect.Type) (isOfType bool, err error) {
	obj, err := runtime.Peek()
	if err != nil {
		isOfType = false
		return
	}

	return reflect.TypeOf(obj) == typeToCheck, nil
}
