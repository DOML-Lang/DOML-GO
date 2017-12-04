/*
 * Copyright (c) 2017. Braedon Wooding
 * Created under LICENSE, see the file LICENSE for information
 */

package core

type InstructionCore struct {
	Opcode    byte
	Parameter interface{}
}

//noinspection GoUnusedConst
const (
	/* ------ System Instructions ------ */

	Nop byte = iota
	Comment
	MakeSpace
	MakeReg

	/* ------ Call Instructions ------ */

	Set
	Call
	New

	/* ------ Push Instructions ------ */

	RegObj
	UnRegObj
	Copy
	Pop
	PushObj
	PushInt
	PushNum
	PushDec
	PushStr
	PushBool
	Push
	CountOfInstructions
)