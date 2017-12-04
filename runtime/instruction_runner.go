/*
 * Copyright (c) 2017. Braedon Wooding
 * Created under LICENSE, see the file LICENSE for information
 */

package runtime

import (
	"doml/internal/core"
	"errors"
	"github.com/shopspring/decimal"
)

type Instruction core.InstructionCore

var invalidParameterError = errors.New("invalid parameter")

func(runtime *Runtime) HandleInstructions() error {
	for _, instruction := range *runtime.Instructions {
		switch instruction.Opcode {
		/* ------ System Instructions ------ */
		case core.Nop:
		case core.Comment:
			// Do Explicitly nothing
			return nil
		case core.MakeSpace:
			if space, ok := instruction.Parameter.(int); ok {
				runtime.ReserveSpace(space)
			} else {
				return invalidParameterError
			}
		case core.MakeReg:
			if space, ok := instruction.Parameter.(int); ok {
				runtime.ReserveRegister(space)
			} else {
				return invalidParameterError
			}
			/* ------ Call Instructions ------ */
		case core.Set:
			// Not Done Yet
			return errors.New("not implemented yet")
		case core.Call:
			// Not Done Yet
			return errors.New("not implemented yet")
		case core.New:
			// Not Done Yet
			return errors.New("not implemented yet")

			/* ------ Push Instructions ------ */
		case core.RegObj:
			if index, ok := instruction.Parameter.(int); ok {
				if obj, err := runtime.Pop(); err == nil {
					runtime.SetObject(index, obj)
				} else {
					return err
				}
			} else {
				return invalidParameterError
			}
		case core.UnRegObj:
			if index, ok := instruction.Parameter.(int); ok {
				runtime.RemoveObject(index)
			} else {
				return invalidParameterError
			}
		case core.Copy:
			if index, ok := instruction.Parameter.(int); ok {
				if obj, err := runtime.Peek(); err == nil {
					for ; index > 0; index-- {
						if err = runtime.Push(obj, true); err != nil {
							return err
						}
					}
				} else {
					return err
				}
			} else {
				return invalidParameterError
			}
		case core.Pop:
			if index, ok := instruction.Parameter.(int); ok {
				for ; index > 0; index-- {
					if err := runtime.PopNoReturn(); err != nil {
						return err
					}
				}
			} else {
				return invalidParameterError
			}
		case core.PushObj:
			if index, ok := instruction.Parameter.(int); ok {
				if obj, err := runtime.GetObject(index); err == nil {
					if err = runtime.Push(obj, true); err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				return invalidParameterError
			}
		case core.PushInt:
			if value, ok := instruction.Parameter.(int64); ok {
				if err := runtime.Push(value, true); err != nil {
					return err
				}
			}
		case core.PushNum:
			if value, ok := instruction.Parameter.(float64); ok {
				if err := runtime.Push(value, true); err != nil {
					return err
				}
			}
		case core.PushDec:
			if value, ok := instruction.Parameter.(decimal.Decimal); ok {
				if err := runtime.Push(value, true); err != nil {
					return err
				}
			}
		case core.PushStr:
			if value, ok := instruction.Parameter.(string); ok {
				if err := runtime.Push(value, true); err != nil {
					return err
				}
			}
		case core.PushBool:
			if value, ok := instruction.Parameter.(bool); ok {
				if err := runtime.Push(value, true); err != nil {
					return err
				}
			}
		case core.Push:
			if err := runtime.Push(instruction.Parameter, true); err != nil {
				return err
			}
		default:
			return errors.New("invalid opcode")
		}
	}
	return nil
}