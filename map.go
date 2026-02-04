package funtemplates

import (
	"fmt"
	"reflect"
)

var (
	errorType = reflect.TypeOf((*error)(nil)).Elem()
	anyType   = reflect.TypeOf((*any)(nil)).Elem()
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

	fvType := fv.Type()
	numIn := fvType.NumIn()
	if numIn != 0 && numIn != 1 {
		return nil, ErrInputFuncMustTake0or1Arguments
	}

	numOut := fvType.NumOut()
	if numOut != 1 && numOut != 2 {
		return nil, fmt.Errorf("%w got: %d", ErrExpected1Or2ReturnTypes, numOut)
	}

	var fvfpt reflect.Type
	if numIn == 1 {
		fvfpt = fvType.In(0)
	}

	var fvfrt reflect.Type
	if numOut == 1 {
		fvfrt = fvType.Out(0)
	} else {
		// numOut == 2
		fvsrt := fvType.Out(1)
		if !fvsrt.AssignableTo(errorType) && !fvsrt.Implements(errorType) {
			return nil, fmt.Errorf("%w instead got: %s", ErrExpectedSecondReturnToBeError, fvsrt)
		}
	}

	l := 0
	if av.Kind() != reflect.Invalid && !av.IsNil() {
		l = av.Len()
	}

	// Optimization: Fast path for known single return type
	if numOut == 1 {
		// Pre-allocate result slice
		nra := reflect.MakeSlice(reflect.SliceOf(fvfrt), l, l)

		if numIn == 1 {
			// Reuse argument slice to avoid allocation per iteration
			args := make([]reflect.Value, 1)
			for i := 0; i < l; i++ {
				ev := av.Index(i)
				if fvfpt != nil && !ev.Type().AssignableTo(fvfpt) {
					return nil, fmt.Errorf("item %d not assignable to: %s", i, fvfpt)
				}
				args[0] = ev
				r := fv.Call(args)
				// Direct assignment avoiding intermediate reflection overhead
				nra.Index(i).Set(r[0])
			}
		} else {
			// numIn == 0
			args := []reflect.Value{}
			for i := 0; i < l; i++ {
				r := fv.Call(args)
				nra.Index(i).Set(r[0])
			}
		}
		return nra.Interface(), nil
	}

	// Slow path: Dynamic return type or error handling (numOut == 2)
	ra := make([]reflect.Value, l)
	var newType reflect.Type // Initially nil

	// Optimization: Lift loop invariants
	var args []reflect.Value
	if numIn == 0 {
		args = []reflect.Value{}
	} else {
		args = make([]reflect.Value, 1)
	}

	for i := 0; i < l; i++ {
		var r []reflect.Value
		if numIn == 1 {
			ev := av.Index(i)
			if fvfpt != nil && !ev.Type().AssignableTo(fvfpt) {
				return nil, fmt.Errorf("item %d not assignable to: %s", i, fvfpt)
			}
			args[0] = ev
			r = fv.Call(args)
		} else {
			r = fv.Call(args)
		}

		// Check for error return (2nd value)
		if len(r) == 2 && !r[1].IsNil() {
			return nil, fmt.Errorf("f execution number %d returned: %w", i, r[1].Interface().(error))
		}

		rt := r[0].Type()
		if newType == nil {
			newType = rt
		} else if rt != newType && !rt.AssignableTo(newType) {
			// Fallback to []interface{} if types are incompatible
			newType = anyType
		}
		ra[i] = r[0]
	}

	if newType == nil {
		newType = anyType
	}
	nra := reflect.MakeSlice(reflect.SliceOf(newType), l, l)
	for i, e := range ra {
		nra.Index(i).Set(e)
	}
	return nra.Interface(), nil
}
