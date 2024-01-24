package testutil

import (
	"fmt"
	"io"
	"math/big"
	"strconv"
)

type Int struct {
	v *big.Int
}

func MakeInt(i int64) Int {
	return Int{big.NewInt(i)}
}

func MakeIntFromString(s string) Int {
	i, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic(fmt.Sprintf("invalid int: %s", s))
	}
	return Int{i}
}

func (n Int) IsZero() bool {
	return n.v == nil || (n.v.IsInt64() && n.v.Int64() == 0)
}

func (n Int) String() string {
	return n.v.String()
}

func (n Int) Format(s fmt.State, r rune) {
	n.v.Format(s, r)
}

func (n Int) Cmp(i Int) int {
	return n.v.Cmp(i.v)
}

func (n Int) Abs() Int {
	return Int{new(big.Int).Abs(n.v)}
}

func (n Int) Neg() Int {
	return Int{new(big.Int).Neg(n.v)}
}

func (n Int) Add(i Int) Int {
	return Int{new(big.Int).Add(n.v, i.v)}
}

func (n Int) Sub(i Int) Int {
	return Int{new(big.Int).Sub(n.v, i.v)}
}

func (n Int) Mul(i Int) Int {
	return Int{new(big.Int).Mul(n.v, i.v)}
}

func (n Int) Div(i Int) Int {
	return Int{new(big.Int).Div(n.v, i.v)}
}

func (n Int) DivFloat(i Int) Float {
	return n.Float().Quo(i.Float())
}

func (n Int) Int64() int64 {
	return n.v.Int64()
}

func (n Int) Float() Float {
	return Float{new(big.Float).SetPrec(100).SetInt(n.v)}
}

func (n Int) MarshalText() ([]byte, error) {
	return n.v.MarshalText()
}

func (n *Int) UnmarshalText(txt []byte) error {
	if len(txt) > 2 && txt[0] == '"' && txt[len(txt)-1] == '"' {
		txt = txt[1 : len(txt)-1]
	}
	i := big.NewInt(0)
	err := i.UnmarshalText(txt)
	if err != nil {
		return err
	}
	n.v = i
	return nil
}

type Float struct {
	v *big.Float
}

func MakeFloat(f float64) Float {
	return Float{new(big.Float).SetPrec(100).SetFloat64(f)}
}

func MakeFloatFromString(s string) Float {
	f, ok := new(big.Float).SetPrec(100).SetString(s)
	if !ok {
		panic(fmt.Sprintf("invalid float: %s", s))
	}
	return Float{f}
}

func (n Float) IsZero() bool {
	f, acc := n.v.Float64()
	return n.v == nil || (f == 0 && acc == big.Exact)
}

func (n Float) String() string {
	return n.v.String()
}

func (n Float) Format(s fmt.State, r rune) {
	if r == 's' {
		io.WriteString(s, n.String())
	} else {
		n.v.Format(s, r)
	}
}

func (n Float) Cmp(f Float) int {
	return n.v.Cmp(f.v)
}

func (n Float) Abs() Float {
	return Float{new(big.Float).Abs(n.v)}
}

func (n Float) Neg() Float {
	return Float{new(big.Float).Neg(n.v)}
}

func (n Float) Add(f Float) Float {
	return Float{new(big.Float).Add(n.v, f.v)}
}

func (n Float) Sub(f Float) Float {
	return Float{new(big.Float).Sub(n.v, f.v)}
}

func (n Float) Mul(f Float) Float {
	return Float{new(big.Float).Mul(n.v, f.v)}
}

func (n Float) Quo(f Float) Float {
	n.v.Float64()
	return Float{new(big.Float).Quo(n.v, f.v)}
}

func (n Float) Float64() float64 {
	v, _ := n.v.Float64()
	return v
}

func (b Float) Int() Int {
	i, _ := b.v.Int(nil)
	return Int{i}
}

func (n Float) MarshalText() ([]byte, error) {
	return n.v.MarshalText()
}

func (n *Float) UnmarshalText(txt []byte) error {
	if len(txt) > 2 && txt[0] == '"' && txt[len(txt)-1] == '"' {
		txt = txt[1 : len(txt)-1]
	}
	f := new(big.Float).SetPrec(100)
	err := f.UnmarshalText(txt)
	if err != nil {
		return err
	}
	n.v = f
	return nil
}

func FormatInt(i int64) string {
	return strconv.FormatInt(i, 10)
}

func FormatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
