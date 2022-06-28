package geo

import (
	"fmt"
)

type QuadKeySystem struct {
	minZoom   int64
	maxZoom   int64
	quadCodes map[byte]rune
}

func NewQuadKeySystem(minZoom int64, maxZoom int64) *QuadKeySystem {
	return &QuadKeySystem{
		minZoom: minZoom,
		maxZoom: maxZoom,
		quadCodes: map[byte]rune{
			0: '0',
			1: '1',
			2: '2',
			3: '3',
		},
	}
}

func (q *QuadKeySystem) TileXYToQuadKey(tx, ty int64, zoom int64) QuadKey {
	quadKey := make([]rune, 0, MaxZoom)
	for i := zoom; i >= q.minZoom; i-- {
		var digit byte = 0
		var mask int64 = 1 << i
		if (tx & mask) != 0 {
			digit++
		}
		if (ty & mask) != 0 {
			digit++
			digit++
		}
		quadKey = append(quadKey, q.quadCodes[digit])
	}
	return quadKey
}

func (q *QuadKeySystem) QuadKeyToTileXY(quadKey QuadKey) (tileX, tileY int64, err error) {
	var tx int64 = 0
	var ty int64 = 0
	var z = len(quadKey)
	var mz = int(q.minZoom)
	for i := z; i > mz; i-- {
		var mask int64 = 1 << (i - 1)
		switch quadKey[z-i] {
		case '0':
		case '1':
			tx |= mask
		case '2':
			ty |= mask
		case '3':
			tx |= mask
			ty |= mask
		default:
			return 0, 0, fmt.Errorf("invalid quad key %v", quadKey)
		}
	}
	return tx, ty, nil
}

func (q *QuadKeySystem) QuadKeyRange(quadKey QuadKey) (min QuadKey, max QuadKey) {
	diff := int(q.maxZoom - int64(len(quadKey)))
	min = quadKey.Copy()
	max = quadKey.Copy()
	for i := 0; i <= diff; i++ {
		min = append(min, '0')
		max = append(max, '3')
	}
	return min, max
}

func (q *QuadKeySystem) Contains(point QuadKey, tile QuadKey) bool {
	val := point.Int64()
	min, max := q.QuadKeyRange(tile)
	minInt, maxInt := min.Int64(), max.Int64()
	if val >= minInt && val <= maxInt {
		return true
	}
	return false
}

func (q *QuadKeySystem) Base10Delta(len int64) int64 {
	diff := q.maxZoom - len
	return (1 << diff * 2) + 1
}

func (q *QuadKeySystem) ForEachZoom(qk QuadKey, cb func(x, y, z int64) error) error {
	n := len(qk)
	reduce := 0
	for {
		newLen := int64(n - reduce)
		if newLen <= q.minZoom {
			return nil
		}
		qk = qk[:newLen]
		tx, ty, err := q.QuadKeyToTileXY(qk)
		if err != nil {
			return err
		}
		err = cb(tx, ty, newLen)
		if err != nil {
			return err
		}
		reduce++
	}
}

func (q *QuadKeySystem) CreateMask(qk QuadKey, zoom int64, clusterLevel int64) QuadKey {
	var mask = qk.Copy()
	if clusterLevel <= 0 {
		clusterLevel = 1
	}
	diff := int(q.maxZoom - (zoom + clusterLevel))
	if diff < 0 {
		diff = 0
	}
	for i := 0; i < int(clusterLevel); i++ {
		mask = append(mask, '3')
	}

	for i := 0; i <= diff; i++ {
		mask = append(mask, '0')
	}
	return mask
}
