package format

import (
	"strconv"
)

// Byte is just one byte.
const Byte uint64 = 1

// binary
const (
	_ = Byte << (10 * iota)
	KibiByte
	MebiByte
	GibiByte
	TebiByte
	PebiByte
	ExbiByte
)

// decimal
const (
	KiloByte = 1000 * Byte
	MegaByte = 1000 * KiloByte
	GigaByte = 1000 * MegaByte
	TeraByte = 1000 * GigaByte
	PetaByte = 1000 * TeraByte
	ExaByte  = 1000 * PetaByte
)

type unit struct {
	size  uint64
	short string
	long  string
}

var byteUnit = unit{Byte, "B", "Byte"}

var binunits = []unit{
	byteUnit,
	{KibiByte, "KiB", "Kibibyte"},
	{MebiByte, "MiB", "Mibibyte"},
	{GibiByte, "GiB", "Gibibyte"},
	{TebiByte, "TiB", "Tebibyte"},
	{PebiByte, "PiB", "Pebibyte"},
	{ExbiByte, "EiB", "Exbibyte"},
}

var decunits = []unit{
	byteUnit,
	{KiloByte, "KB", "Kilobyte"},
	{MegaByte, "MB", "Megabyte"},
	{GigaByte, "GB", "Gigabyte"},
	{TeraByte, "TB", "Terabyte"},
	{PetaByte, "PB", "Petabyte"},
	{ExaByte, "EB", "Exabyte"},
}

// FmtBytesUnit returns the string representation of bytes size with unit.
func FmtBytesUnit(bytes uint64, unit uint64, longUnit bool) string {
	units := append(binunits, decunits[1:]...)
	for i := len(units) - 1; i >= 0; i-- {
		if unit == units[i].size {
			return fmtBytes(bytes, &units[i], longUnit)
		}
	}
	return ""
}

// FmtBytes returns the string representation of bytes size with auto unit.
func FmtBytes(bytes uint64, longUnit, useDecimal bool) string {
	var u *unit
	units := binunits
	if useDecimal {
		units = decunits
	}

	for i := len(units) - 1; i >= 0; i-- {
		if bytes >= units[i].size {
			u = &units[i]
			break
		}
	}

	return fmtBytes(bytes, u, longUnit)
}

func fmtBytes(bytes uint64, u *unit, long bool) string {
	if u == nil {
		u = &byteUnit
	}
	var res = make([]byte, 0, 32)
	res = strconv.AppendFloat(res, float64(bytes)/float64(u.size), 'f', 3, 64)
	res = append(res, ' ')
	if long {
		res = append(res, u.long...)
	} else {
		res = append(res, u.short...)
	}

	return string(res)
}
