package funtemplates

import (
	"fmt"
	"reflect"
)

func FilterTemplateFunc(slice any, f any) (any, error) {
	av := reflect.ValueOf(slice)
	if av.Kind() != reflect.Slice && av.Kind() != reflect.Invalid {
		return nil, fmt.Errorf("%w not %s", ErrExpectedFirstParameterToBeSlice, av.Kind())
	}
	fv := reflect.ValueOf(f)
	if fv.Kind() != reflect.Func {
		return nil, ErrExpected2ndArgumentToBeFunction
	}

	numIn := fv.Type().NumIn()
	if numIn > 1 {
		return nil, ErrInputFuncMustTake0or1Arguments
	}

	var elemType reflect.Type
	checkInside := false

	if numIn == 1 {
		elemType = fv.Type().In(0)
		if av.Kind() == reflect.Slice {
			sliceElemType := av.Type().Elem()
			// If the slice contains interfaces, we cannot statically guarantee that
			// the dynamic values inside are assignable to the function argument type.
			// We must check inside the loop to provide a friendly error instead of panicking,
			// or to allow valid assignments if the types align (e.g. interface{} -> interface{}).
			if sliceElemType.Kind() == reflect.Interface {
				checkInside = true
			} else if !sliceElemType.AssignableTo(elemType) {
				return nil, fmt.Errorf("item 0 not assignable to: %s", elemType)
			}
		}
	} else {
		// NumIn == 0
		if av.Kind() == reflect.Slice {
			elemType = av.Type().Elem()
		} else {
			// av is Invalid (nil)
			elemType = reflect.TypeOf((*interface{})(nil)).Elem()
		}
	}

	switch fv.Type().NumOut() {
	case 1:
		fvfrt := fv.Type().Out(0)
		if !fvfrt.AssignableTo(reflect.TypeOf(true)) {
			return nil, fmt.Errorf("%w instead got: %s", ErrExpectedFirstReturnToBeBool, fvfrt)
		}
	case 2:
		fvsrt := fv.Type().Out(1)
		errorType := reflect.TypeOf((*error)(nil)).Elem()
		if !fvsrt.AssignableTo(errorType) {
			return nil, fmt.Errorf("%w instead got: %s", ErrExpectedSecondReturnToBeError, fvsrt)
		}
		fvfrt := fv.Type().Out(0)
		if !fvfrt.AssignableTo(reflect.TypeOf(true)) {
			return nil, fmt.Errorf("%w instead got: %s", ErrExpectedFirstReturnToBeBool, fvfrt)
		}
	default:
		return nil, fmt.Errorf("%w got: %d", ErrExpected1Or2ReturnTypes, fv.Type().NumOut())
	}

	l := 0
	if av.Kind() != reflect.Invalid && !av.IsNil() {
		l = av.Len()
	}

	// Create result slice with correct type and capacity
	nra := reflect.MakeSlice(reflect.SliceOf(elemType), 0, l)

	// Pre-allocate args
	var args []reflect.Value
	if numIn == 1 {
		args = make([]reflect.Value, 1)
	}

	for i := 0; i < l; i++ {
		var r []reflect.Value
		if numIn == 1 {
			ev := av.Index(i)
			if checkInside {
				if !ev.Type().AssignableTo(elemType) {
					return nil, fmt.Errorf("item %d not assignable to: %s", i, elemType)
				}
			}
			args[0] = ev
			r = fv.Call(args)
		} else {
			r = fv.Call(nil)
		}

		if len(r) == 2 && !r[1].IsNil() {
			return nil, fmt.Errorf("f execution number %d returned: %w", i, r[1].Interface().(error))
		}

		if r[0].Bool() {
			nra = reflect.Append(nra, av.Index(i))
		}
	}

	return nra.Interface(), nil
}
