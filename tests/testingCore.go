/*
 * Copyright (c) 2017. Braedon Wooding
 * Created under LICENSE, see the file LICENSE for information
 */

package tests

import (
	"doml/runtime"
)

// ResetStack just properly resets stack
func ResetStack(runtime *runtime.Runtime, newSize int) {
	runtime.Unsafe.ResizeStack(newSize, false)
	runtime.Unsafe.SetStackPtr(-1)
}

// ResetRegisters just properly resets registers
func ResetRegisters(runtime *runtime.Runtime, newSize int) {
	runtime.Unsafe.ResizeRegisters(newSize, false)
}
