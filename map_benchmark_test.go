package funtemplates

import (
	"testing"
)

func BenchmarkMapTemplateFunc(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}
	inc := func(i int) int {
		return i + 1
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := MapTemplateFunc(data, inc)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMapTemplateFunc_InterfaceReturn(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}
	// Function returning interface{}, forcing dynamic type checking in the original implementation
	inc := func(i int) interface{} {
		return i + 1
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := MapTemplateFunc(data, inc)
		if err != nil {
			b.Fatal(err)
		}
	}
}
