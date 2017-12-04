/*
 * Copyright (c) 2017. Braedon Wooding
 * Created under LICENSE, see the file LICENSE for information
 */

package core

import (
	"reflect"
)

type RuntimeCore struct {
	stack []interface{}
	// Should always begin at -1 to indicate that the array is empty
	stackPtr  int
	registers []interface{}
}

func (runtime *RuntimeCore) CurrentStackSize() int {
	return runtime.stackPtr + 1
}

func (runtime *RuntimeCore) MaxStackSize() int {
	return len(runtime.stack)
}

func (runtime *RuntimeCore) RegisterSize() int {
	return len(runtime.registers)
}

func (runtime *RuntimeCore) ResizeStack(size int, carryAcross bool) {
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

func (runtime *RuntimeCore) ResizeRegisters(size int, carryAcross bool) {
	if carryAcross {
		temp := runtime.registers
		runtime.registers = make([]interface{}, size)
		copy(temp, runtime.registers)
	} else {
		runtime.registers = make([]interface{}, size)
	}
}

func (runtime *RuntimeCore) SetStackPtr(value int) {
	runtime.stackPtr = value
}

func (runtime *RuntimeCore) IncrementStackPtr(value int) {
	runtime.stackPtr += value
}

func (runtime *RuntimeCore) DecrementStackPtr(value int) {
	runtime.stackPtr -= value
}

func (runtime *RuntimeCore) GetRegister(index int) interface{} {
	return runtime.registers[index]
}

func (runtime *RuntimeCore) SetRegister(index int, obj interface{}) {
	runtime.registers[index] = obj
}

func (runtime *RuntimeCore) PushOntoStack(obj interface{}) {
	runtime.stackPtr += 1
	runtime.stack[runtime.stackPtr] = obj
}

func (runtime *RuntimeCore) PeekFromStack() interface{} {
	return runtime.stack[runtime.stackPtr]
}

func (runtime *RuntimeCore) PopFromStack() interface{} {
	obj := runtime.PeekFromStack()
	runtime.stackPtr -= 1
	return obj
}

func (runtime *RuntimeCore) TopOfStackIsOfType(typeToCheck reflect.Type) bool {
	return reflect.TypeOf(runtime.PeekFromStack()).ConvertibleTo(typeToCheck)
}
