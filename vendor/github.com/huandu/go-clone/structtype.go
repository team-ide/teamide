// Copyright 2019 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package clone

import (
	"crypto/elliptic"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var (
	cachedStructTypes  sync.Map
	cachedPointerTypes sync.Map
)

func init() {
	// Some well-known scalar-like structs.
	MarkAsScalar(reflect.TypeOf(time.Time{}))
	MarkAsScalar(reflect.TypeOf(reflect.Value{}))

	// Special case for elliptic.Curve which is used by TLS ECC certificate.
	// Package crypto/tls uses elliptic.Curve as enum values
	// so that they should be treated as opaque pointers.
	//
	// As elliptic.Curve is an interface, it can be *elliptic.CurveParam or elliptic.p256Curve.
	p224 := elliptic.P224()
	MarkAsScalar(reflect.TypeOf(p224))
	MarkAsOpaquePointer(reflect.TypeOf(&elliptic.CurveParams{}))

	// Special case for reflect.Type (actually *reflect.rtype):
	// The *reflect.rtype should not be copied as it is immutable and
	// may point to a variable that actual type is not reflect.rtype,
	// e.g. *reflect.arrayType or *reflect.chanType.
	MarkAsOpaquePointer(reflect.TypeOf(reflect.TypeOf(0)))

	// Some well-known no-copy structs.
	//
	// Almost all structs defined in package "sync" and "sync/atomic" are set
	// except `sync.Once` which can be safely cloned with a correct done value.
	SetCustomFunc(reflect.TypeOf(sync.Mutex{}), emptyCloneFunc)
	SetCustomFunc(reflect.TypeOf(sync.RWMutex{}), emptyCloneFunc)
	SetCustomFunc(reflect.TypeOf(sync.WaitGroup{}), emptyCloneFunc)
	SetCustomFunc(reflect.TypeOf(sync.Cond{}), func(old, new reflect.Value) {
		// Copy the New func from old value.
		oldL := old.FieldByName("L")
		newL := noState.clone(oldL)
		new.FieldByName("L").Set(newL)
	})
	SetCustomFunc(reflect.TypeOf(sync.Pool{}), func(old, new reflect.Value) {
		// Copy the New func from old value.
		oldFn := old.FieldByName("New")
		newFn := noState.clone(oldFn)
		new.FieldByName("New").Set(newFn)
	})
	SetCustomFunc(reflect.TypeOf(sync.Map{}), func(old, new reflect.Value) {
		if !old.CanAddr() {
			return
		}

		// Clone all values inside sync.Map.
		oldMap := old.Addr().Interface().(*sync.Map)
		newMap := new.Addr().Interface().(*sync.Map)
		oldMap.Range(func(key, value interface{}) bool {
			k := Clone(key)
			v := Clone(value)
			newMap.Store(k, v)
			return true
		})
	})
	SetCustomFunc(reflect.TypeOf(atomic.Value{}), func(old, new reflect.Value) {
		if !old.CanAddr() {
			return
		}

		// Clone value inside atomic.Value.
		oldValue := old.Addr().Interface().(*atomic.Value)
		newValue := new.Addr().Interface().(*atomic.Value)
		v := Clone(oldValue.Load())
		newValue.Store(v)
	})
}

// MarkAsScalar marks t as a scalar type so that all clone methods will copy t by value.
// If t is not struct or pointer to struct, MarkAsScalar ignores t.
//
// In the most cases, it's not necessary to call it explicitly.
// If a struct type contains scalar type fields only, the struct will be marked as scalar automatically.
//
// Here is a list of types marked as scalar by default:
//     * time.Time
//     * reflect.Value
func MarkAsScalar(t reflect.Type) {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return
	}

	cachedStructTypes.Store(t, structType{})
}

// MarkAsOpaquePointer marks t as an opaque pointer so that all clone methods will copy t by value.
// If t is not a pointer, MarkAsOpaquePointer ignores t.
//
// Here is a list of types marked as opaque pointers by default:
//     * `elliptic.Curve`, which is `*elliptic.CurveParam` or `elliptic.p256Curve`;
//     * `reflect.Type`, which is `*reflect.rtype` defined in `runtime`.
func MarkAsOpaquePointer(t reflect.Type) {
	if t.Kind() != reflect.Ptr {
		return
	}

	cachedPointerTypes.Store(t, struct{}{})
}

// Func is a custom func to clone value from old to new.
// The new is a zero value
// which `new.CanSet()` and `new.CanAddr()` is guaranteed to be true.
//
// Func must update the new to return result.
type Func func(old, new reflect.Value)

// emptyCloneFunc is used to disable shadow copy.
// It's useful when cloning sync.Mutex as cloned value must be a zero value.
func emptyCloneFunc(old, new reflect.Value) {}

// SetCustomFunc sets a custom clone function for type t.
// If t is not struct or pointer to struct, SetCustomFunc ignores t.
func SetCustomFunc(t reflect.Type, fn Func) {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return
	}

	cachedStructTypes.Store(t, structType{
		fn: fn,
	})
}

type structType struct {
	PointerFields []structFieldType
	fn            Func
}

type structFieldType struct {
	Offset uintptr // The offset from the beginning of the struct.
	Index  int     // The index of the field.
}

// NewFrom creates a new value of src.Type() and shadow copies all content from src.
func (st *structType) Copy(src, nv reflect.Value) {
	dst := nv.Elem()

	if st.fn != nil {
		if !src.CanInterface() {
			src = forceClearROFlag(src)
		}

		st.fn(src, dst)
		return
	}

	ptr := unsafe.Pointer(nv.Pointer())
	shadowCopy(src, ptr)
}

func (st *structType) CanShadowCopy() bool {
	return len(st.PointerFields) == 0 && st.fn == nil
}

func loadStructType(t reflect.Type) (st structType) {
	if v, ok := cachedStructTypes.Load(t); ok {
		st = v.(structType)
		return
	}

	num := t.NumField()
	pointerFields := make([]structFieldType, 0, num)

	for i := 0; i < num; i++ {
		field := t.Field(i)
		ft := field.Type
		k := ft.Kind()

		if isScalar(k) {
			continue
		}

		switch k {
		case reflect.Array:
			if ft.Len() == 0 {
				continue
			}

			elem := ft.Elem()

			if isScalar(elem.Kind()) {
				continue
			}

			if elem.Kind() == reflect.Struct {
				if fst := loadStructType(elem); fst.CanShadowCopy() {
					continue
				}
			}
		case reflect.Struct:
			if fst := loadStructType(ft); fst.CanShadowCopy() {
				continue
			}
		}

		pointerFields = append(pointerFields, structFieldType{
			Offset: field.Offset,
			Index:  i,
		})
	}

	if len(pointerFields) == 0 {
		pointerFields = nil // Release memory ASAP.
	}

	st = structType{
		PointerFields: pointerFields,
	}
	cachedStructTypes.LoadOrStore(t, st)
	return
}

func isScalar(k reflect.Kind) bool {
	switch k {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String, reflect.Func,
		reflect.UnsafePointer,
		reflect.Invalid:
		return true
	}

	return false
}

func isOpaquePointer(t reflect.Type) (ok bool) {
	_, ok = cachedPointerTypes.Load(t)
	return
}

func copyScalarValue(src reflect.Value) reflect.Value {
	if src.CanInterface() {
		return src
	}

	// src is an unexported field value. Copy its value.
	switch src.Kind() {
	case reflect.Bool:
		return reflect.ValueOf(src.Bool())

	case reflect.Int:
		return reflect.ValueOf(int(src.Int()))
	case reflect.Int8:
		return reflect.ValueOf(int8(src.Int()))
	case reflect.Int16:
		return reflect.ValueOf(int16(src.Int()))
	case reflect.Int32:
		return reflect.ValueOf(int32(src.Int()))
	case reflect.Int64:
		return reflect.ValueOf(src.Int())

	case reflect.Uint:
		return reflect.ValueOf(uint(src.Uint()))
	case reflect.Uint8:
		return reflect.ValueOf(uint8(src.Uint()))
	case reflect.Uint16:
		return reflect.ValueOf(uint16(src.Uint()))
	case reflect.Uint32:
		return reflect.ValueOf(uint32(src.Uint()))
	case reflect.Uint64:
		return reflect.ValueOf(src.Uint())
	case reflect.Uintptr:
		return reflect.ValueOf(uintptr(src.Uint()))

	case reflect.Float32:
		return reflect.ValueOf(float32(src.Float()))
	case reflect.Float64:
		return reflect.ValueOf(src.Float())

	case reflect.Complex64:
		return reflect.ValueOf(complex64(src.Complex()))
	case reflect.Complex128:
		return reflect.ValueOf(src.Complex())

	case reflect.String:
		return reflect.ValueOf(src.String())
	case reflect.Func:
		t := src.Type()

		if src.IsNil() {
			return reflect.Zero(t)
		}

		// Don't use this trick unless we have no choice.
		return forceClearROFlag(src)
	case reflect.UnsafePointer:
		return reflect.ValueOf(unsafe.Pointer(src.Pointer()))
	}

	panic(fmt.Errorf("go-clone: <bug> impossible type `%v` when cloning private field", src.Type()))
}

var typeOfInterface = reflect.TypeOf((*interface{})(nil)).Elem()

// forceClearROFlag clears all RO flags in v to make v accessible.
// It's a hack based on the fact that InterfaceData is always available on RO data.
// This hack can be broken in any Go version.
// Don't use it unless we have no choice, e.g. copying func in some edge cases.
func forceClearROFlag(v reflect.Value) reflect.Value {
	var i interface{}
	indirect := 0

	// Save flagAddr.
	for v.CanAddr() {
		v = v.Addr()
		indirect++
	}

	v = v.Convert(typeOfInterface)
	nv := reflect.ValueOf(&i)
	*(*interfaceData)(unsafe.Pointer(nv.Pointer())) = parseReflectValue(v)
	cleared := nv.Elem().Elem()

	for indirect > 0 {
		cleared = cleared.Elem()
		indirect--
	}

	return cleared
}
