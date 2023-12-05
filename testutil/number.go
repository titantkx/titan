package testutil

import (
	"fmt"
	"math/big"
	"strconv"
)

type Int int64

func (n Int) Int64() int64 {
	return int64(n)
}

func (n *Int) UnmarshalText(txt []byte) error {
	if len(txt) > 2 && txt[0] == '"' && txt[len(txt)-1] == '"' {
		txt = txt[1 : len(txt)-1]
	}
	i, err := strconv.ParseInt(string(txt), 10, 64)
	if err != nil {
		return err
	}
	*n = Int(i)
	return nil
}

type Float float64

func (n Float) Float64() float64 {
	return float64(n)
}

func (n *Float) UnmarshalText(txt []byte) error {
	if len(txt) > 2 && txt[0] == '"' && txt[len(txt)-1] == '"' {
		txt = txt[1 : len(txt)-1]
	}
	f, err := strconv.ParseFloat(string(txt), 10)
	if err != nil {
		return err
	}
	*n = Float(f)
	return nil
}

type BigInt struct {
	v *big.Int
}

func MakeBigInt(i int64) BigInt {
	return BigInt{big.NewInt(i)}
}

func MakeBigIntFromString(s string) BigInt {
	i, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic(fmt.Sprintf("invalid int: %s", s))
	}
	return BigInt{i}
}

func (n BigInt) IsZero() bool {
	return n.v == nil || (n.v.IsInt64() && n.v.Int64() == 0)
}

func (n BigInt) String() string {
	return n.v.String()
}

func (n BigInt) Format(s fmt.State, r rune) {
	n.v.Format(s, r)
}

func (n BigInt) Cmp(i BigInt) int {
	return n.v.Cmp(i.v)
}

func (n BigInt) Neg() BigInt {
	return BigInt{new(big.Int).Neg(n.v)}
}

func (n BigInt) Add(i BigInt) BigInt {
	return BigInt{new(big.Int).Add(n.v, i.v)}
}

func (n BigInt) Sub(i BigInt) BigInt {
	return BigInt{new(big.Int).Sub(n.v, i.v)}
}

func (n BigInt) Mul(i BigInt) BigInt {
	return BigInt{new(big.Int).Mul(n.v, i.v)}
}

func (n BigInt) Div(i BigInt) BigInt {
	return BigInt{new(big.Int).Div(n.v, i.v)}
}

func (n BigInt) BigFloat() BigFloat {
	return BigFloat{new(big.Float).SetPrec(100).SetInt(n.v)}
}

func (n BigInt) MarshalText() ([]byte, error) {
	return n.v.MarshalText()
}

func (n *BigInt) UnmarshalText(txt []byte) error {
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

type BigFloat struct {
	v *big.Float
}

func MakeBigFloat(f float64) BigFloat {
	return BigFloat{new(big.Float).SetPrec(100).SetFloat64(f)}
}

func MakeBigFloatFromString(s string) BigFloat {
	f, ok := new(big.Float).SetPrec(100).SetString(s)
	if !ok {
		panic(fmt.Sprintf("invalid float: %s", s))
	}
	return BigFloat{f}
}

func (n BigFloat) IsZero() bool {
	f, acc := n.v.Float64()
	return n.v == nil || (f == 0 && acc == big.Exact)
}

func (n BigFloat) String() string {
	return n.String()
}

func (n BigFloat) Format(s fmt.State, r rune) {
	n.v.Format(s, r)
}

func (n BigFloat) Cmp(f BigFloat) int {
	return n.v.Cmp(f.v)
}

func (n BigFloat) Neg() BigFloat {
	return BigFloat{new(big.Float).Neg(n.v)}
}

func (n BigFloat) Add(f BigFloat) BigFloat {
	return BigFloat{new(big.Float).Add(n.v, f.v)}
}

func (n BigFloat) Sub(f BigFloat) BigFloat {
	return BigFloat{new(big.Float).Sub(n.v, f.v)}
}

func (n BigFloat) Mul(f BigFloat) BigFloat {
	return BigFloat{new(big.Float).Mul(n.v, f.v)}
}

func (b BigFloat) BigInt() BigInt {
	i, _ := b.v.Int(nil)
	return BigInt{i}
}

func (n BigFloat) MarshalText() ([]byte, error) {
	return n.v.MarshalText()
}

func (n *BigFloat) UnmarshalText(txt []byte) error {
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
