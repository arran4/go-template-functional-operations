package funtemplates

import (
	"fmt"
	"reflect"
)

func FindTemplateFunc(slice any, f any) (any, error) {
	i, err := FindIndexTemplateFunc(slice, f)
	if err != nil {
		return nil, err
	}
	if i == -1 {
		return nil, nil
	}
	av := reflect.ValueOf(slice)
	return av.Index(i).Interface(), nil
}

func FindIndexTemplateFunc(slice any, f any) (int, error) {
	av := reflect.ValueOf(slice)
	if av.Kind() != reflect.Slice && av.Kind() != reflect.Invalid {
		return -1, fmt.Errorf("%w not %s", ErrExpectedFirstParameterToBeSlice, av.Kind())
	}
	fv := reflect.ValueOf(f)
	if fv.Kind() != reflect.Func {
		return -1, ErrExpected2ndArgumentToBeFunction
	}
	var fvfpt reflect.Type
	switch fv.Type().NumIn() {
	case 0:
	case 1:
		fvfpt = fv.Type().In(0)
	default:
		return -1, ErrInputFuncMustTake0or1Arguments
	}
	switch fv.Type().NumOut() {
	case 1:
		fvfrt := fv.Type().Out(0)
		if !fvfrt.AssignableTo(reflect.TypeOf(true)) {
			return -1, fmt.Errorf("%w instead got: %s", ErrExpectedFirstReturnToBeBool, fvfrt)
		}
	case 2:
		fvsrt := fv.Type().Out(1)
		if !fvsrt.AssignableTo(reflect.TypeOf(error(nil))) {
			return -1, fmt.Errorf("%w instead got: %s", ErrExpectedSecondReturnToBeError, fvsrt)
		}
	default:
		return -1, fmt.Errorf("%w got: %d", ErrExpected1Or2ReturnTypes, fv.Type().NumOut())
	}
	l := 0
	if av.Kind() != reflect.Invalid && !av.IsNil() {
		l = av.Len()
	}
	switch fv.Type().NumIn() {
	case 1:
		if fvfpt != nil {
			for i := 0; i < l; i++ {
				ev := av.Index(i)
				if ev.Kind() == reflect.Interface && !ev.IsNil() {
					ev = ev.Elem()
				}
				if !ev.Type().AssignableTo(fvfpt) {
					return -1, fmt.Errorf("item %d not assignable to: %s", i, fvfpt)
				}
				r := fv.Call([]reflect.Value{ev})
				if r == nil {
					continue
				}
				if len(r) != 1 && len(r) != 2 {
					return -1, fmt.Errorf("f execution number %d returned: %d results %w", i, len(r), ErrExpected1Or2ReturnTypes)
				}
				if len(r) == 2 && !r[1].IsNil() {
					return -1, fmt.Errorf("f execution number %d returned: %w", i, r[1].Interface().(error))
				}
				if b1, b2 := r[0].Interface().(bool); b1 && b2 {
					return i, nil
				}
			}
			return -1, nil
		}
		fallthrough
	case 0:
		for i := 0; i < l; i++ {
			r := fv.Call([]reflect.Value{})
			if r == nil {
				continue
			}
			if len(r) != 1 && len(r) != 2 {
				return -1, fmt.Errorf("f execution number %d returned: %d results %w", i, len(r), ErrExpected1Or2ReturnTypes)
			}
			if len(r) == 2 && !r[1].IsNil() {
				return -1, fmt.Errorf("f execution number %d returned: %w", i, r[1].Interface().(error))
			}
			if b1, b2 := r[0].Interface().(bool); b1 && b2 {
				return i, nil
			}
		}
		return -1, nil
	default:
		return -1, ErrInputFuncMustTake0or1Arguments
	}
}
