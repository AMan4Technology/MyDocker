package datastructure

func NewBitmap(capacity int) Bitmap {
    return Bitmap{
        capacity: (capacity/8 + 1) * 8,
        data:     make([]byte, capacity/8+1),
    }
}

type Bitmap struct {
    capacity, length int
    data             []byte
}

func (b Bitmap) Full() bool {
    return b.length == b.capacity
}

func (b Bitmap) Empty() bool {
    return b.length == 0
}

func (b Bitmap) Cap() int {
    return b.capacity
}

func (b Bitmap) Len() int {
    return b.length
}

func (b Bitmap) Have(val int) bool {
    if val < 0 || val >= b.capacity {
        return false
    }
    var (
        byteIndex = val / 8
        bitIndex  = uint(val % 8)
    )
    return 1<<bitIndex&b.data[byteIndex] != 0
}

func (b *Bitmap) Save(val int) {
    if val < 0 {
        return
    }
    for ; val >= b.capacity; b.capacity += 8 {
        b.data = append(b.data, 0)
    }
    b.update(val, true)
}

func (b *Bitmap) Remove(val int) {
    if val < 0 || val >= b.capacity {
        return
    }
    b.update(val, false)
}

func (b *Bitmap) update(val int, save bool) {
    var (
        byteIndex = val / 8
        bitIndex  = uint(val % 8)
        operator  = uint8(1 << bitIndex)
    )
    if save {
        if operator&b.data[byteIndex] == 0 {
            b.data[byteIndex] |= operator
            b.length++
        }
    } else {
        if operator&b.data[byteIndex] != 0 {
            b.data[byteIndex] &= ^operator
            b.length--
        }
    }
}
