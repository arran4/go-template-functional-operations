package funtemplates

import "errors"

var (
	ErrInputFuncMustTake0or1Arguments  = errors.New("expected second parameter function to take 0 or 1 parameters")
	ErrExpectedFirstParameterToBeSlice = errors.New("expected first parameter to be an slice")
	ErrExpected2ndArgumentToBeFunction = errors.New("expected second parameter to be a function")
	ErrExpectedSecondReturnToBeError   = errors.New("expected second return type to be assignable to error")
	ErrExpected1Or2ReturnTypes         = errors.New("expected return with 1 or 2 arguments of types (any, error?)")
	ErrExpectedFirstReturnToBeBool     = errors.New("expected first return type to be assignable to bool")
)
