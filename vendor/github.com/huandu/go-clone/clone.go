// Copyright 2019 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

// Package clone provides functions to deep clone any Go data.
// It also provides a wrapper to protect a pointer from any unexpected mutation.
package clone

import (
	"fmt"
	"reflect"
	"unsafe"
)

// Clone recursively deep clone v to a new value.
// It assumes that there is no pointer cycle in v,
// e.g. v has a pointer points to v itself.
// If there is a pointer cycle, use Slowly instead.
//
// Clone allocates memory and deeply copies values inside v in depth-first sequence.
// There are a few special rules for following types.
//
//     * Scalar types: all number-like types are copied by value.
//     * func: Copied by value as func is an opaque pointer at runtime.
//     * string: Copied by value as string is immutable by design.
//     * unsafe.Pointer: Copied by value as we don't know what's in it.
//     * chan: A new empty chan is created as we cannot read data inside the old chan.
//
// Unlike many other packages, Clone is able to clone unexported fields of any struct.
// Use this feature wisely.
func Clone(v interface{}) interface{} {
	if v == nil {
		return v
	}

	val := reflect.ValueOf(v)
	cloned := noState.clone(val)
	return cloned.Interface()
}

// Slowly recursively deep clone v to a new value.
// It marks all cloned values internally, thus it can clone v with cycle pointer.
//
// Slowly works exactly the same as Clone. See Clone doc for more details.
func Slowly(v interface{}) interface{} {
	if v == nil {
		return v
	}

	val := reflect.ValueOf(v)
	state := &cloneState{
		visited: visitMap{},
		invalid: invalidPointers{},
	}
	cloned := state.clone(val)
	state.fix(cloned)
	return cloned.Interface()
}

type cloneState struct {
	visited visitMap
	invalid invalidPointers
}

var noState *cloneState

type visit struct {
	p     uintptr
	extra int
	t     reflect.Type
}

type visitMap map[visit]reflect.Value
type invalidPointers map[visit]reflect.Value

func (state *cloneState) clone(v reflect.Value) reflect.Value {
	if isScalar(v.Kind()) {
		return copyScalarValue(v)
	}

	switch v.Kind() {
	case reflect.Array:
		return state.cloneArray(v)
	case reflect.Chan:
		return reflect.MakeChan(v.Type(), v.Cap())
	case reflect.Interface:
		return state.cloneInterface(v)
	case reflect.Map:
		return state.cloneMap(v)
	case reflect.Ptr:
		return state.clonePtr(v)
	case reflect.Slice:
		return state.cloneSlice(v)
	case reflect.Struct:
		return state.cloneStruct(v)
	default:
		panic(fmt.Errorf("go-clone: <bug> unsupported type `%v`", v.Type()))
	}
}

func (state *cloneState) cloneArray(v reflect.Value) reflect.Value {
	dst := reflect.New(v.Type())
	state.copyArray(v, dst)
	return dst.Elem()
}

func (state *cloneState) copyArray(src, nv reflect.Value) {
	p := unsafe.Pointer(nv.Pointer()) // dst must be a Ptr.
	dst := nv.Elem()
	num := src.Len()

	if isScalar(src.Type().Elem().Kind()) {
		shadowCopy(src, p)
		return
	}

	for i := 0; i < num; i++ {
		dst.Index(i).Set(state.clone(src.Index(i)))
	}
}

func (state *cloneState) cloneInterface(v reflect.Value) reflect.Value {
	if v.IsNil() {
		return reflect.Zero(v.Type())
	}

	t := v.Type()
	elem := v.Elem()
	return state.clone(elem).Convert(elem.Type()).Convert(t)
}

func (state *cloneState) cloneMap(v reflect.Value) reflect.Value {
	if v.IsNil() {
		return reflect.Zero(v.Type())
	}

	t := v.Type()

	if state != nil {
		vst := visit{
			p: v.Pointer(),
			t: t,
		}

		if val, ok := state.visited[vst]; ok {
			return val
		}
	}

	nv := reflect.MakeMap(t)

	if state != nil {
		vst := visit{
			p: v.Pointer(),
			t: t,
		}
		state.visited[vst] = nv
	}

	for iter := mapIter(v); iter.Next(); {
		key := state.clone(iter.Key())
		value := state.clone(iter.Value())
		nv.SetMapIndex(key, value)
	}

	return nv
}

func (state *cloneState) clonePtr(v reflect.Value) reflect.Value {
	if v.IsNil() {
		return reflect.Zero(v.Type())
	}

	t := v.Type()

	if isOpaquePointer(t) {
		if v.CanInterface() {
			return v
		}

		ptr := reflect.New(t)
		p := unsafe.Pointer(ptr.Pointer())
		shadowCopy(v, p)
		return ptr.Elem()
	}

	if state != nil {
		vst := visit{
			p: v.Pointer(),
			t: t,
		}

		if val, ok := state.visited[vst]; ok {
			return val
		}
	}

	src := v.Elem()
	elemType := src.Type()
	elemKind := src.Kind()
	nv := reflect.New(elemType)

	if state != nil {
		vst := visit{
			p: v.Pointer(),
			t: t,
		}
		state.visited[vst] = nv
	}

	switch elemKind {
	case reflect.Struct:
		state.copyStruct(src, nv)
	case reflect.Array:
		state.copyArray(src, nv)
	default:
		nv.Elem().Set(state.clone(src))
	}

	// If this pointer is the address of a struct field and it's a cycle pointer,
	// it may be updated.
	if state != nil {
		vst := visit{
			p: v.Pointer(),
			t: t,
		}
		nv = state.visited[vst]
	}

	return nv
}

func (state *cloneState) cloneSlice(v reflect.Value) reflect.Value {
	if v.IsNil() {
		return reflect.Zero(v.Type())
	}

	t := v.Type()
	num := v.Len()

	if state != nil {
		vst := visit{
			p:     v.Pointer(),
			extra: num,
			t:     t,
		}

		if val, ok := state.visited[vst]; ok {
			return val
		}
	}

	c := v.Cap()
	nv := reflect.MakeSlice(t, num, c)

	if state != nil {
		vst := visit{
			p:     v.Pointer(),
			extra: num,
			t:     t,
		}
		state.visited[vst] = nv
	}

	// For scalar slice, copy underlying values directly.
	if isScalar(t.Elem().Kind()) {
		src := unsafe.Pointer(v.Pointer())
		dst := unsafe.Pointer(nv.Pointer())
		sz := int(t.Elem().Size())
		l := num * sz
		cc := c * sz
		copy((*[maxByteSize]byte)(dst)[:l:cc], (*[maxByteSize]byte)(src)[:l:cc])
	} else {
		for i := 0; i < num; i++ {
			nv.Index(i).Set(state.clone(v.Index(i)))
		}
	}

	return nv
}

func (state *cloneState) cloneStruct(v reflect.Value) reflect.Value {
	t := v.Type()
	nv := reflect.New(t)
	state.copyStruct(v, nv)
	return nv.Elem()
}

func (state *cloneState) copyStruct(src, nv reflect.Value) {
	t := src.Type()
	st := loadStructType(t)
	ptr := unsafe.Pointer(nv.Pointer())

	st.Copy(src, nv)

	// If the struct type is a scalar type, a.k.a type without any pointer,
	// there is no need to iterate over fields.
	if len(st.PointerFields) == 0 {
		return
	}

	for _, pf := range st.PointerFields {
		i := int(pf.Index)
		p := unsafe.Pointer(uintptr(ptr) + pf.Offset)
		field := src.Field(i)

		// This field can be referenced by a pointer or interface inside itself.
		// Put the pointer to this field to visited to avoid any error.
		//
		// See https://github.com/huandu/go-clone/issues/3.
		if state != nil && field.CanAddr() {
			ft := field.Type()
			fp := field.Addr().Pointer()
			vst := visit{
				p: fp,
				t: reflect.PtrTo(ft),
			}
			nv := reflect.NewAt(ft, p)

			// The address of this field was visited, so fp must be a cycle pointer.
			// As this field is not fully cloned, the val stored in visited[visit] must be wrong.
			// It must be replaced by nv which will be the right value (it's incomplete right now).
			//
			// Unfortunately, if the val was used by previous clone routines,
			// there is no easy way to fix wrong values - all pointers must be traversed and fixed.
			if val, ok := state.visited[vst]; ok {
				state.invalid[visit{
					p: val.Pointer(),
					t: vst.t,
				}] = nv
			}

			state.visited[vst] = nv
		}

		v := state.clone(field)
		shadowCopy(v, p)
	}
}

var typeOfString = reflect.TypeOf("")

func shadowCopy(src reflect.Value, p unsafe.Pointer) {
	switch src.Kind() {
	case reflect.Bool:
		*(*bool)(p) = src.Bool()
	case reflect.Int:
		*(*int)(p) = int(src.Int())
	case reflect.Int8:
		*(*int8)(p) = int8(src.Int())
	case reflect.Int16:
		*(*int16)(p) = int16(src.Int())
	case reflect.Int32:
		*(*int32)(p) = int32(src.Int())
	case reflect.Int64:
		*(*int64)(p) = src.Int()
	case reflect.Uint:
		*(*uint)(p) = uint(src.Uint())
	case reflect.Uint8:
		*(*uint8)(p) = uint8(src.Uint())
	case reflect.Uint16:
		*(*uint16)(p) = uint16(src.Uint())
	case reflect.Uint32:
		*(*uint32)(p) = uint32(src.Uint())
	case reflect.Uint64:
		*(*uint64)(p) = src.Uint()
	case reflect.Uintptr:
		*(*uintptr)(p) = uintptr(src.Uint())
	case reflect.Float32:
		*(*float32)(p) = float32(src.Float())
	case reflect.Float64:
		*(*float64)(p) = src.Float()
	case reflect.Complex64:
		*(*complex64)(p) = complex64(src.Complex())
	case reflect.Complex128:
		*(*complex128)(p) = src.Complex()

	case reflect.Array:
		t := src.Type()

		if src.CanAddr() {
			srcPtr := unsafe.Pointer(src.UnsafeAddr())
			sz := t.Size()
			copy((*[maxByteSize]byte)(p)[:sz:sz], (*[maxByteSize]byte)(srcPtr)[:sz:sz])
			return
		}

		val := reflect.NewAt(t, p).Elem()

		if src.CanInterface() {
			val.Set(src)
			return
		}

		sz := t.Elem().Size()
		num := src.Len()

		for i := 0; i < num; i++ {
			elemPtr := unsafe.Pointer(uintptr(p) + uintptr(i)*sz)
			shadowCopy(src.Index(i), elemPtr)
		}
	case reflect.Chan:
		*((*uintptr)(p)) = src.Pointer()
	case reflect.Func:
		t := src.Type()
		src = copyScalarValue(src)
		val := reflect.NewAt(t, p).Elem()
		val.Set(src)
	case reflect.Interface:
		*((*interfaceData)(p)) = parseReflectValue(src)
	case reflect.Map:
		*((*uintptr)(p)) = src.Pointer()
	case reflect.Ptr:
		*((*uintptr)(p)) = src.Pointer()
	case reflect.Slice:
		*(*sliceHeader)(p) = sliceHeader{
			Data: src.Pointer(),
			Len:  src.Len(),
			Cap:  src.Cap(),
		}
	case reflect.String:
		s := src.String()
		val := reflect.NewAt(typeOfString, p).Elem()
		val.SetString(s)
	case reflect.Struct:
		t := src.Type()
		val := reflect.NewAt(t, p).Elem()

		if src.CanInterface() {
			val.Set(src)
			return
		}

		num := t.NumField()

		for i := 0; i < num; i++ {
			field := t.Field(i)
			fieldPtr := unsafe.Pointer(uintptr(p) + field.Offset)
			shadowCopy(src.Field(i), fieldPtr)
		}
	case reflect.UnsafePointer:
		// There is no way to copy unsafe.Pointer value.
		*((*uintptr)(p)) = src.Pointer()

	default:
		panic(fmt.Errorf("go-clone: <bug> impossible type `%v` when cloning private field", src.Type()))
	}
}

// fix tranverses v to update all pointer values in state.invalid.
func (state *cloneState) fix(v reflect.Value) {
	if state == nil || len(state.invalid) == 0 {
		return
	}

	fix := &fixState{
		fixed:   fixMap{},
		invalid: state.invalid,
	}
	fix.fix(v)
}

type fixState struct {
	fixed   fixMap
	invalid invalidPointers
}

type fixMap map[visit]struct{}

func (fix *fixState) fix(v reflect.Value) (copied reflect.Value, changed int) {
	if isScalar(v.Kind()) {
		return
	}

	switch v.Kind() {
	case reflect.Array:
		return fix.fixArray(v)
	case reflect.Chan:
		// Do nothing.
		return
	case reflect.Interface:
		return fix.fixInterface(v)
	case reflect.Map:
		return fix.fixMap(v)
	case reflect.Ptr:
		return fix.fixPtr(v)
	case reflect.Slice:
		return fix.fixSlice(v)
	case reflect.Struct:
		return fix.fixStruct(v)
	default:
		panic(fmt.Errorf("go-clone: <bug> unsupported type `%v`", v.Type()))
	}
}

func (fix *fixState) fixArray(v reflect.Value) (copied reflect.Value, changed int) {
	t := v.Type()
	et := t.Elem()
	kind := et.Kind()

	if isScalar(kind) {
		return
	}

	l := v.Len()

	for i := 0; i < l; i++ {
		elem := v.Index(i)

		if kind == reflect.Ptr {
			vst := visit{
				p: elem.Pointer(),
				t: et,
			}

			if nv, ok := fix.invalid[vst]; ok {
				// If elem cannot be set, v must be copied to make it settable.
				// Don't do it unless there is no other choices.
				if !elem.CanSet() {
					copied = reflect.New(t).Elem()
					shadowCopy(v, unsafe.Pointer(copied.Addr().Pointer()))
					_, changed = fix.fixArray(copied)
					return
				}

				elem.Set(nv)
				changed++
				continue
			}
		}

		fixed, c := fix.fix(elem)
		changed += c

		if fixed.IsValid() {
			// If elem cannot be set, v must be copied to make it settable.
			// Don't do it unless there is no other choices.
			if !elem.CanSet() {
				copied = reflect.New(t).Elem()
				shadowCopy(v, unsafe.Pointer(copied.Addr().Pointer()))
				_, changed = fix.fixArray(copied)
				return
			}

			elem.Set(fixed)
		}
	}

	return
}

func (fix *fixState) fixInterface(v reflect.Value) (copied reflect.Value, changed int) {
	if v.IsNil() {
		return
	}

	elem := v.Elem()
	t := elem.Type()
	kind := elem.Kind()

	if kind == reflect.Ptr {
		vst := visit{
			p: elem.Pointer(),
			t: t,
		}

		if nv, ok := fix.invalid[vst]; ok {
			copied = nv.Convert(v.Type())
			changed++
			return
		}
	}

	copied, changed = fix.fix(elem)

	if copied.IsValid() {
		copied = copied.Convert(v.Type())
	}

	return
}

func (fix *fixState) fixMap(v reflect.Value) (copied reflect.Value, changed int) {
	if v.IsNil() {
		return
	}

	t := v.Type()
	vst := visit{
		p: v.Pointer(),
		t: t,
	}

	if _, ok := fix.fixed[vst]; ok {
		return
	}

	fix.fixed[vst] = struct{}{}

	kt := t.Key()
	et := t.Elem()
	keyKind := kt.Kind()
	elemKind := et.Kind()

	if isScalar(keyKind) && isScalar(elemKind) {
		return
	}

	invalidKeys := map[reflect.Value][2]reflect.Value{}

	for iter := mapIter(v); iter.Next(); {
		key := iter.Key()
		elem := iter.Value()
		var fixed reflect.Value
		c := 0

		if elemKind == reflect.Ptr {
			vst := visit{
				p: elem.Pointer(),
				t: et,
			}

			if nv, ok := fix.invalid[vst]; ok {
				fixed = nv
				c++
			} else {
				fixed, c = fix.fixPtr(elem)
			}
		} else {
			fixed, c = fix.fix(elem)
		}

		changed += c
		c = 0

		if fixed.IsValid() {
			v = forceSetMapIndex(v, key, fixed)
			elem = fixed
			fixed = reflect.Value{}
		}

		if keyKind == reflect.Ptr {
			vst := visit{
				p: key.Pointer(),
				t: kt,
			}

			if nv, ok := fix.invalid[vst]; ok {
				fixed = nv
				c++
			} else {
				fixed, c = fix.fixPtr(key)
			}
		} else {
			fixed, c = fix.fix(key)
		}

		changed += c

		// Key cannot be changed immediately inside map range iteration.
		// Do it later.
		if fixed.IsValid() {
			invalidKeys[key] = [2]reflect.Value{fixed, elem}
		}
	}

	for key, kv := range invalidKeys {
		v = forceSetMapIndex(v, key, reflect.Value{})
		v = forceSetMapIndex(v, kv[0], kv[1])
	}

	return
}

func forceSetMapIndex(v, key, elem reflect.Value) (nv reflect.Value) {
	nv = v

	if !v.CanInterface() {
		nv = forceClearROFlag(v)
	}

	if !key.CanInterface() {
		key = forceClearROFlag(key)
	}

	if elem.IsValid() && !elem.CanInterface() {
		elem = forceClearROFlag(elem)
	}

	nv.SetMapIndex(key, elem)
	return
}

func (fix *fixState) fixPtr(v reflect.Value) (copied reflect.Value, changed int) {
	if v.IsNil() {
		return
	}

	vst := visit{
		p: v.Pointer(),
		t: v.Type(),
	}

	if _, ok := fix.invalid[vst]; ok {
		panic(fmt.Errorf("go-clone: <bug> invalid pointers must have been fixed in other methods"))
	}

	if _, ok := fix.fixed[vst]; ok {
		return
	}

	fix.fixed[vst] = struct{}{}

	elem := v.Elem()
	_, changed = fix.fix(elem)
	return
}

func (fix *fixState) fixSlice(v reflect.Value) (copied reflect.Value, changed int) {
	if v.IsNil() {
		return
	}

	t := v.Type()
	et := t.Elem()
	kind := et.Kind()

	if isScalar(kind) {
		return
	}

	l := v.Len()
	p := unsafe.Pointer(v.Pointer())
	vst := visit{
		p:     uintptr(p),
		extra: l,
		t:     t,
	}

	if _, ok := fix.fixed[vst]; ok {
		return
	}

	fix.fixed[vst] = struct{}{}

	for i := 0; i < l; i++ {
		elem := v.Index(i)
		var fixed reflect.Value
		c := 0

		if kind == reflect.Ptr {
			vst := visit{
				p: elem.Pointer(),
				t: et,
			}

			if nv, ok := fix.invalid[vst]; ok {
				fixed = nv
			} else {
				fixed, c = fix.fixPtr(elem)
			}
		} else {
			fixed, c = fix.fix(elem)
		}

		changed += c

		if fixed.IsValid() {
			sz := et.Size()
			elemPtr := unsafe.Pointer(uintptr(p) + sz*uintptr(i))
			shadowCopy(fixed, elemPtr)
		}
	}

	return
}

func (fix *fixState) fixStruct(v reflect.Value) (copied reflect.Value, changed int) {
	t := v.Type()
	st := loadStructType(t)

	if len(st.PointerFields) == 0 {
		return
	}

	for _, pf := range st.PointerFields {
		i := int(pf.Index)
		field := v.Field(i)

		ft := field.Type()

		if ft.Kind() == reflect.Ptr {
			vst := visit{
				p: field.Pointer(),
				t: ft,
			}

			if nv, ok := fix.invalid[vst]; ok {
				// If v is not addressable, a new struct must be allocated.
				// Don't do it unless there is no other choices.
				if !v.CanAddr() {
					copied = reflect.New(t).Elem()
					shadowCopy(v, unsafe.Pointer(copied.Addr().Pointer()))
					_, changed = fix.fixStruct(copied)
					return
				}

				ptr := unsafe.Pointer(v.Addr().Pointer())
				p := unsafe.Pointer(uintptr(ptr) + pf.Offset)
				shadowCopy(nv, p)
				continue
			}
		}

		fixed, c := fix.fix(field)
		changed += c

		if fixed.IsValid() {
			// If v is not addressable, a new struct must be allocated.
			// Don't do it unless there is no other choices.
			if !v.CanAddr() {
				copied = reflect.New(t).Elem()
				shadowCopy(v, unsafe.Pointer(copied.Addr().Pointer()))
				_, changed = fix.fixStruct(copied)
				return
			}

			ptr := unsafe.Pointer(v.Addr().Pointer())
			p := unsafe.Pointer(uintptr(ptr) + pf.Offset)
			shadowCopy(fixed, p)
		}
	}

	return
}
