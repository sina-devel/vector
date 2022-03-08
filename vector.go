package vector

import (
	"unsafe"
)

// TODO: After coming the type parameters, it is removed
const intSize uintptr = unsafe.Sizeof(int(0))

type Vector struct {
	array    unsafe.Pointer
	length   int
	capacity int
}

// New returns a Vector with zero length and capacity.
func New() Vector {
	return Vector{
		array:    nil,
		length:   0,
		capacity: 0,
	}
}

// Make returns a Vector with the given length and capacity.
//
// It panics if length < 0 or capacity < length.
func Make(length int, capacity int) Vector {
	if length < 0 || capacity < length {
		panic("vector: out of range")
	}

	return Vector{
		array:    mallocgc(intSize*uintptr(capacity), nil, true),
		length:   length,
		capacity: capacity,
	}
}

// At returns the value of an element at the given index.
//
// It panics if i >= length or i < 0.
func (v *Vector) At(i int) int {
	if i < 0 || i >= v.length {
		panic("vector: out of range")
	}
	return *(*int)(unsafe.Add(v.array, intSize*uintptr(i)))
}

// Set sets the given value to the given index.
//
// It panics if i >= length or i < 0.
func (v *Vector) Set(i int, value int) {
	if i < 0 || i >= v.length {
		panic("vector: out of range")
	}
	*(*int)(unsafe.Add(v.array, intSize*uintptr(i))) = value
}

// Append appends the given values to the end of the vector.
//
// If it does not have enough capacity, a new infrastructure array is allocated and the vector grows.
func (v *Vector) Append(values ...int) {
	if len(values) > (v.capacity - v.length) {
		idlcap := len(values) + v.length // 12

		newCap := v.capacity // 10
		capx2 := v.capacity * 2 // 20
		if idlcap > capx2 { // 12 > 20
			newCap = idlcap
		} else {
			if v.capacity < 1024 { // 10 < 1024
				newCap = capx2 // 20
			} else {
				for 0 < newCap && newCap < idlcap {
					newCap += newCap / 4
				}

				if newCap <= 0 {
					newCap = idlcap
				}
			}
		}

		newArray := mallocgc(intSize*uintptr(newCap), nil, true)
		memmove(newArray, v.array, intSize*uintptr(v.length))
		v.array = newArray
		v.capacity = newCap
	}

	under := (*Vector)(unsafe.Pointer(&values)).array
	memmove(unsafe.Add(v.array, intSize*uintptr(v.length)), under, intSize*uintptr((len(values))))

	v.length += len(values)
}

// Len returns the length of the vector.
func (v *Vector) Len() int {
	return v.length
}

// Cap returns the capacity of the vector.
func (v *Vector) Cap() int {
	return v.capacity
}

//go:linkname mallocgc runtime.mallocgc
func mallocgc(size uintptr, typ interface{}, needzero bool) unsafe.Pointer

//go:linkname memmove runtime.memmove
func memmove(to, from unsafe.Pointer, n uintptr)
