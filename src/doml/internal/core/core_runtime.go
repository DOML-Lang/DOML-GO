package core

import (
	"reflect"
)

type Runtime struct {
	stack []interface{}
	// Should always begin at -1 to indicate that the array is empty
	stackPtr  int
	registers []interface{}
}

func (runtime *Runtime) CurrentStackSize() int {
	return runtime.stackPtr + 1
}

func (runtime *Runtime) MaxStackSize() int {
	return len(runtime.stack)
}

func (runtime *Runtime) RegisterSize() int {
	return len(runtime.registers)
}

func (runtime *Runtime) ResizeStack(size int, carryAcross bool) {
	if carryAcross {
		temp := runtime.stack
		runtime.stack = make([]interface{}, size)
		copy(temp, runtime.stack)
		print(len(runtime.stack))
		print(size)
	} else {
		runtime.stack = make([]interface{}, size)
	}
}

func (runtime *Runtime) ResizeRegisters(size int, carryAcross bool) {
	if carryAcross {
		temp := runtime.registers
		runtime.registers = make([]interface{}, size)
		copy(temp, runtime.registers)
	} else {
		runtime.registers = make([]interface{}, size)
	}
}

func (runtime *Runtime) SetStackPtr(value int) {
	runtime.stackPtr = value
}

func (runtime *Runtime) IncrementStackPtr(value int) {
	runtime.stackPtr += value
}

func (runtime *Runtime) DecrementStackPtr(value int) {
	runtime.stackPtr -= value
}

func (runtime *Runtime) GetRegister(index int) interface{} {
	return runtime.registers[index]
}

func (runtime *Runtime) SetRegister(index int, obj interface{}) {
	runtime.registers[index] = obj
}

func (runtime *Runtime) PushOntoStack(obj interface{}) {
	runtime.stackPtr += 1
	runtime.stack[runtime.stackPtr] = obj
}

func (runtime *Runtime) PeekFromStack() interface{} {
	return runtime.stack[runtime.stackPtr]
}

func (runtime *Runtime) PopFromStack() interface{} {
	obj := runtime.PeekFromStack()
	runtime.stackPtr -= 1
	return obj
}

func (runtime *Runtime) TopOfStackIsOfType(typeToCheck reflect.Type) bool {
	return reflect.TypeOf(runtime.PeekFromStack()).ConvertibleTo(typeToCheck)
}
