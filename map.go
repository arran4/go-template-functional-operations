package funtemplates

import (
	"fmt"
	"reflect"
)

func MapTemplateFunc(slice any, f any) (any, error) {
	av := reflect.ValueOf(slice)
	if av.Kind() != reflect.Slice && av.Kind() != reflect.Invalid {
		return nil, fmt.Errorf("%w not %s", ErrExpectedFirstParameterToBeSlice, av.Kind())
	}
	fv := reflect.ValueOf(f)
	if fv.Kind() != reflect.Func {
		return nil, ErrExpected2ndArgumentToBeFunction
	}
	var fvfpt reflect.Type
	switch fv.Type().NumIn() {
	case 0:
	case 1:
		fvfpt = fv.Type().In(0)
	default:
		return nil, ErrInputFuncMustTake0or1Arguments
	}
	var fvfrt reflect.Type
	switch fv.Type().NumOut() {
	case 1:
		fvfrt = fv.Type().Out(0)
	case 2:
		fvsrt := fv.Type().Out(1)
		if !fvsrt.AssignableTo(reflect.TypeOf(error(nil))) {
			return nil, fmt.Errorf("%w instead got: %s", ErrExpectedSecondReturnToBeError, fvsrt)
		}
	default:
		return nil, fmt.Errorf("%w got: %d", ErrExpected1Or2ReturnTypes, fv.Type().NumOut())
	}
	l := 0
	if av.Kind() != reflect.Invalid && !av.IsNil() {
		l = av.Len()
	}
	ra := make([]reflect.Value, l)
	var newType reflect.Type = fvfrt
	toan := reflect.TypeOf(any(nil))
	for i := 0; i < l; i++ {
		var r []reflect.Value
		switch fv.Type().NumIn() {
		case 1:
			if fvfpt != nil {
				ev := av.Index(i)
				if !ev.Type().AssignableTo(fvfpt) {
					return nil, fmt.Errorf("item %d not assignable to: %s", i, fvfpt)
				}
				r = fv.Call([]reflect.Value{ev})
				break
			}
			fallthrough
		case 0:
			r = fv.Call([]reflect.Value{})
		default:
			return nil, ErrInputFuncMustTake0or1Arguments
		}
		if r == nil {
			continue
		}
		if len(r) != 1 && len(r) != 2 {
			return nil, fmt.Errorf("f execution number %d returned: %d results %w", i, len(r), ErrExpected1Or2ReturnTypes)
		}
		if len(r) == 2 && !r[1].IsNil() {
			return nil, fmt.Errorf("f execution number %d returned: %w", i, r[1].Interface().(error))
		}
		rt := reflect.TypeOf(r[0])
		if newType == nil {
			newType = rt
		} else if rt != newType && !rt.AssignableTo(newType) {
			rt = toan
		}
		ra[i] = r[0]
	}
	if newType == nil {
		return ra, nil
	}
	nra := reflect.MakeSlice(reflect.SliceOf(newType), l, l)
	for i, e := range ra {
		nra.Index(i).Set(e)
	}
	return nra.Interface(), nil
}
