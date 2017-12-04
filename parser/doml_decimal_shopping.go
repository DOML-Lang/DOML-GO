// +build useShoppingDecimal

/*
 * Copyright (c) 2017. Braedon Wooding
 * Created under LICENSE, see the file LICENSE for information
 */

package parser

import (
	"doml/internal/core"

	"github.com/shopspring/decimal"
)

func ParseDecimal(text string) (value interface{}, opcode byte, err error) {
	value, err = decimal.NewFromString(text[1:])
	opcode = core.PushDec
	return
}
