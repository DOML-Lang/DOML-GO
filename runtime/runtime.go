/*
 * Copyright (c) 2017. Braedon Wooding
 * Created under LICENSE, see the file LICENSE for information
 */

package runtime

// This file exists as the runtime for DOML
// Allows it to effectively run IR
import (
	"doml/internal/core"
	"errors"
	"fmt"
	"reflect"
)

type Runtime struct {
	Unsafe       core.RuntimeCore
	Instructions *[]Instruction
}

func PrintInstruction(instruction Instruction) {
	println("Opcode: ", fmt.Sprintf("%02d", instruction.Opcode), "  Parameter: ", fmt.Sprint(instruction.Parameter))
}

func NewRuntime(stackSize, registerSize int, instructions *[]Instruction) Runtime {
	runtime := Runtime{}
	runtime.Unsafe.ResizeStack(stackSize, false)
	runtime.Unsafe.ResizeRegisters(registerSize, false)
	runtime.Unsafe.SetStackPtr(-1)
	runtime.Instructions = instructions
	return runtime
}

func NewEmptyRuntime(instructions *[]Instruction) Runtime {
	runtime := Runtime{}
	runtime.Unsafe.SetStackPtr(-1)
	runtime.Instructions = instructions
	return runtime
}

// ClearSpace() -> Clears the stack
// Returns true if stack was cleared, if stack already empty returns false
func (runtime *Runtime) ClearSpace() bool {
	if runtime.Unsafe.CurrentStackSize() > 0 {
		runtime.Unsafe.ResizeStack(runtime.Unsafe.MaxStackSize(), false)
		return true
	}

	runtime.Unsafe.SetStackPtr(-1)
	return false
}

func (runtime *Runtime) ClearRegisters() bool {
	if runtime.Unsafe.RegisterSize() > 0 {
		runtime.Unsafe.ResizeRegisters(runtime.Unsafe.MaxStackSize(), false)
		return true
	}

	return false
}

func (runtime *Runtime) ReserveRegister(space int) bool {
	if space > runtime.Unsafe.RegisterSize() {
		runtime.Unsafe.ResizeRegisters(space, false)
		return true
	}

	return false
}

func (runtime *Runtime) ReserveSpace(space int) bool {
	if space > runtime.Unsafe.MaxStackSize() {
		runtime.Unsafe.ResizeStack(space, false)
		return true
	}

	return false
}

func (runtime *Runtime) GetObject(index int) (obj interface{}, err error) {
	if index >= 0 && index <= runtime.Unsafe.RegisterSize() {
		obj = runtime.Unsafe.GetRegister(index)
		if obj == nil {
			err = errors.New("object is null")
		}
	} else {
		obj = nil
		err = errors.New("index out of range")
	}

	return
}

func (runtime *Runtime) SetObject(index int, obj interface{}) error {
	if index >= 0 && index <= runtime.Unsafe.RegisterSize() && obj != nil {
		runtime.Unsafe.SetRegister(index, obj)
		return nil
	} else {
		return errors.New("object is null or index is out of range")
	}
}

func (runtime *Runtime) RemoveObject(index int) error {
	if index >= 0 && index <= runtime.Unsafe.RegisterSize() {
		runtime.Unsafe.SetRegister(index, nil)
		return nil
	} else {
		return errors.New("index out of range")
	}
}

func (runtime *Runtime) Push(obj interface{}, resizeIfNoSpace bool) error {
	if runtime.Unsafe.CurrentStackSize() == runtime.Unsafe.MaxStackSize() {
		if resizeIfNoSpace {
			runtime.Unsafe.ResizeStack(runtime.Unsafe.MaxStackSize()+1, true)
		} else {
			return errors.New("no space on stack for object")
		}
	}

	runtime.Unsafe.PushOntoStack(obj)
	return nil
}

func (runtime *Runtime) PopNoReturn() error {
	if runtime.Unsafe.CurrentStackSize() > 0 {
		if runtime.Unsafe.PeekFromStack() == nil {
			return errors.New("object is null")
		}
	} else {
		return errors.New("nothing to pop off stack")
	}

	runtime.Unsafe.DecrementStackPtr(1)
	return nil
}

func (runtime *Runtime) Pop() (obj interface{}, err error) {
	obj, err = runtime.Peek()
	if err == nil {
		runtime.Unsafe.DecrementStackPtr(1)
	}
	return
}

func (runtime *Runtime) Peek() (obj interface{}, err error) {
	err = nil
	obj = nil
	if runtime.Unsafe.CurrentStackSize() > 0 {
		obj = runtime.Unsafe.PeekFromStack()
		if obj == nil {
			err = errors.New("object is null")
		}
	} else {
		err = errors.New("nothing to pop off stack")
	}
	return
}

func (runtime *Runtime) TopIsOfType(typeToCheck reflect.Type) (isOfType bool, err error) {
	if runtime.Unsafe.CurrentStackSize() > 0 {
		if runtime.Unsafe.PeekFromStack() != nil {
			isOfType = runtime.Unsafe.TopOfStackIsOfType(typeToCheck)
		} else {
			err = errors.New("object is null")
		}
	} else {
		err = errors.New("object is null")
	}
	return
}
