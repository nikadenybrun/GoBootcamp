package main

import (
	"errors"
	"fmt"
	"unsafe"
)

// getElement возвращает элемент массива по заданному индексу.
// Функция не использует прямой доступ к элементу по индексу, а вместо этого
// использует доступ к первому элементу и смещение указателя.
func getElement(arr []int, idx int) (int, error) {
	if len(arr) == 0 {
		return 0, errors.New("пустой срез")
	}
	if idx < 0 {
		return 0, errors.New("отрицательный индекс")
	}
	if idx >= len(arr) {
		return 0, errors.New("индекс выходит за пределы допустимого диапазона")
	}

	ptr := &arr[0]

	ptr = (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + uintptr(idx)*(unsafe.Sizeof(arr[0]))))

	return *ptr, nil
}

func main() {
	arr := []int{1, 2, 3, 4, 5}
	idx := 2

	element, err := getElement(arr, idx)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(element)
	}
}
