package geo

import (
	"fmt"
	"strconv"
)

type QuadKey []rune

func NewQuadKeyFromInt64(val int64) QuadKey {
	return QuadKey(strconv.FormatInt(val, 4))
}

func (q QuadKey) String() string {
	return string(q)
}

func (q QuadKey) Int64() int64 {
	val, _ := strconv.ParseInt(string(q), 4, 64)
	return val
}

func (q QuadKey) Bits() string {
	val, _ := strconv.ParseInt(string(q), 4, 64)
	return Bits(val)
}

func (q QuadKey) Copy() QuadKey {
	qk := make([]rune, len(q))
	copy(qk, q)
	return qk
}

func (q QuadKey) View() string {
	return fmt.Sprintf("%s:%20d %s:%s %s:%s",
		Colorize(colorYellow, "base-10"), q.Int64(),
		Colorize(colorYellow, "base-4"), q,
		Colorize(colorYellow, "base-2"), q.Bits())
}
