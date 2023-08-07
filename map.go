package funtemplates

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrInputFuncMustTake0or1Arguments     = errors.New("expected second parameter function to take 0 or 1 parameters")
	ErrExpectedFirstParameterToBeSlice    = errors.New("expected first parameter to be an slice")
	ErrExpected2ndArgumentToBeFunction    = errors.New("expected second parameter to be a function")
	ErrExpectedSecondReturnToBeError      = errors.New("expected second return type to be assignable to error")
	ErrExpectedSecondArgumentToBeFunction = errors.New("expected second return of function f to take 1 or 2 parameters instead")
	ErrExpected2ReturnTypes               = errors.New("expected return with 1 or 2 arguments of types (any, error?)")
)

func MapTemplateFunc(slice any, f any) (any, error) {
	av := reflect.ValueOf(slice)
	if av.Kind() != reflect.Slice {
		return slice, fmt.Errorf("%w not %s", ErrExpectedFirstParameterToBeSlice, av.Kind())
	}
	fv := reflect.ValueOf(f)
	if fv.Kind() != reflect.Func {
		return slice, ErrExpected2ndArgumentToBeFunction
	}
	var fvfpt reflect.Type
	switch fv.Type().NumIn() {
	case 0:
	case 1:
		fvfpt = fv.Type().In(0)
	default:
		return slice, ErrInputFuncMustTake0or1Arguments
	}
	var fvfrt reflect.Type
	switch fv.Type().NumOut() {
	case 1:
		fvfrt = fv.Type().Out(0)
	case 2:
		fvsrt := fv.Type().Out(1)
		if !fvsrt.AssignableTo(reflect.TypeOf(error(nil))) {
			return slice, fmt.Errorf("%w instead got: %s", ErrExpectedSecondReturnToBeError, fvsrt)
		}
	default:
		return slice, fmt.Errorf("%w got: %d", ErrExpectedSecondArgumentToBeFunction, fv.Type().NumOut())
	}
	l := av.Len()
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
			return slice, ErrInputFuncMustTake0or1Arguments
		}
		if r == nil {
			continue
		}
		if len(r) != 1 && len(r) != 2 {
			return slice, fmt.Errorf("f execution number %d returned: %d results %w", i, len(r), ErrExpected2ReturnTypes)
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
