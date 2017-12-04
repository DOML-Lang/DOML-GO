/*
 * Copyright (c) 2017. Braedon Wooding
 * Created under LICENSE, see the file LICENSE for information
 */
package tests

import (
	"doml/runtime"
	"testing"
)

func TestPush(t *testing.T) {
	t.Log("Testing Push Functionality")
	t.Log("No need to resize and allowing resize")
	runtimeInstance := runtime.NewRuntime(1, 0, nil)
	err := runtimeInstance.Push(2, true)
	if err != nil {
		t.Error(err.Error())
	}

	t.Log("Resetting")
	ResetStack(&runtimeInstance, 1)

	t.Log("No need to resize and not allowing resize")
	err = runtimeInstance.Push(2, false)
	if err != nil {
		t.Error(err.Error())
	}

	t.Log("Resetting")
	ResetStack(&runtimeInstance, 0)

	t.Log("With need to resize and allowing resize")
	err = runtimeInstance.Push(2, true)
	if err != nil {
		t.Error(err.Error())
	}

	t.Log("Resetting")
	ResetStack(&runtimeInstance, 0)

	t.Log("With need to resize and not allowing resize")
	err = runtimeInstance.Push(2, false)
	if err == nil {
		t.Error("Was expecting an error")
	}
}

func TestPop(t *testing.T) {
	t.Log("Testing Pop Functionality")
	t.Log("Testing pop with object")
	runtimeInstance := runtime.NewRuntime(1, 0, nil)
	err := runtimeInstance.Push(2, false)
	if err != nil {
		t.Error(err.Error())
	}

	obj, err := runtimeInstance.Pop()
	if err != nil {
		t.Error(err.Error())
	}

	if obj != 2 {
		t.Error("popped object doesn't equal pushed object")
	}

	ResetStack(&runtimeInstance, 0)
	t.Log("Testing pop without object")
	obj, err = runtimeInstance.Pop()
	if err == nil || obj != nil {
		t.Error("Was expecting an error and a nil object")
	}
}

func TestPopNoReturn(t *testing.T) {
	t.Log("Testing Pop With No Return Functionality")
	t.Log("Testing pop with object")
	runtimeInstance := runtime.NewRuntime(1, 0, nil)
	err := runtimeInstance.Push(2, false)
	if err != nil {
		t.Error(err.Error())
	}

	err = runtimeInstance.PopNoReturn()
	if err != nil {
		t.Error(err.Error())
	}

	ResetStack(&runtimeInstance, 0)
	t.Log("Testing pop without object")
	err = runtimeInstance.PopNoReturn()
	if err == nil {
		t.Error("Was expecting an error")
	}
}
