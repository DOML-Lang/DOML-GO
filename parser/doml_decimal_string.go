// +build !useShoppingDecimal

/*
 * Copyright (c) 2017. Braedon Wooding
 * Created under LICENSE, see the file LICENSE for information
 */

package parser

import (
	"doml/internal/core"
)

func ParseDecimal(text string) (value interface{}, opcode byte, err error) {
	value = text
	opcode = core.PushStr
	return
}