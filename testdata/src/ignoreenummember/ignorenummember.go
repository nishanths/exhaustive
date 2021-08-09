// want package:"^Exchange:Exchange_EXCHANGE_UNSPECIFIED,Exchange_EXCHANGE_BITMEX,Exchange_EXCHANGE_BINANCE$"

package ignoreenummember

import barpkg "general/y"

type Exchange int32

const (
	Exchange_EXCHANGE_UNSPECIFIED Exchange = 0
	Exchange_EXCHANGE_BITMEX      Exchange = 1
	Exchange_EXCHANGE_BINANCE     Exchange = 2
)

func _a() {
	var e Exchange
	switch e { // want "^missing cases in switch of type Exchange: Exchange_EXCHANGE_BINANCE$"
	case Exchange_EXCHANGE_BITMEX:
	}
}

func _b() {
	var p barpkg.Phylum
	switch p { // want "^missing cases in switch of type bar.Phylum: Mollusca$"
	case barpkg.Chordata:
	}
}
