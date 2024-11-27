package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func Distinct[T comparable](arrays []T) []T {
	if len(arrays) == 0 {
		return arrays
	}

	var dm = make(map[T]bool)
	for _, item := range arrays {
		dm[item] = true
	}

	var res []T
	for k := range dm {
		res = append(res, k)
	}

	return res
}

func Exist[T comparable](arrays []T, key T) bool {

	var dm = make(map[T]struct{})
	for _, item := range arrays {
		dm[item] = struct{}{}
	}

	_, ok := dm[key]
	return ok
}
func Array2Map[T comparable](arrays []T) map[T]struct{} {

	var dm = make(map[T]struct{})
	for _, item := range arrays {
		dm[item] = struct{}{}
	}
	return dm
}

// CompareSlices 使用泛型比较任意类型的两个数组是否相同;utils.CompareSlices(x, y)
func CompareSlices[T comparable](a, b []T) bool {
	// 长度不相等，直接返回 false
	if len(a) != len(b) {
		return false
	}

	// 比较每个元素
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// ExtractField 提取结构体数组中指定字段的值，并返回字段值数组
func ExtractField(slice interface{}, fieldName string) ([]interface{}, error) {
	// 获取 slice 的反射值
	v := reflect.ValueOf(slice)

	// 检查是否是 slice 类型
	if v.Kind() != reflect.Slice {
		return nil, errors.New("input is not a slice")
	}

	// 如果传入的切片为空，返回空数组
	if v.Len() == 0 {
		return []interface{}{}, nil
	}

	// 获取切片元素的类型并检查是否为结构体
	elemType := v.Index(0).Type()
	if elemType.Kind() != reflect.Struct {
		return nil, errors.New("elements of the slice are not structs")
	}

	// 检查结构体是否包含指定的字段
	_, exists := elemType.FieldByName(fieldName)
	if !exists {
		return nil, errors.New("field " + fieldName + " does not exist in the struct")
	}

	// 用于存储提取的字段值
	result := make([]interface{}, v.Len())

	// 遍历切片，提取指定字段的值
	for i := 0; i < v.Len(); i++ {
		structVal := v.Index(i)                      // 结构体值
		fieldVal := structVal.FieldByName(fieldName) // 指定字段的值

		// 确保字段值是可导出的
		if !fieldVal.CanInterface() {
			return nil, errors.New("field " + fieldName + " is not accessible")
		}

		result[i] = fieldVal.Interface() // 获取字段值并存入结果数组
	}

	return result, nil
}

// GetFieldValueByTag 从任意结构体中根据指定的 tag 获取字段值
func GetFieldValueByTag(input any, tagName, tagValue string) (any, error) {
	// 确保输入是一个结构体或结构体的指针
	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, errors.New("input must be a struct or a pointer to a struct")
	}

	// 遍历结构体字段
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		// 检查字段的指定 tag 值
		if tag := field.Tag.Get(tagName); tag == tagValue {
			return v.Field(i).Interface(), nil
		}
	}

	return nil, fmt.Errorf("no field found with tag '%s' and value '%s'", tagName, tagValue)
}

// ExtractFieldValues 从数组/切片中提取指定字段的值，并以逗号分隔的字符串返回
func ExtractFieldValues(slice any, tagName, tagValue string) (string, error) {
	// 确保输入是切片或数组
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return "", errors.New("input must be a slice or an array")
	}

	// 遍历数组/切片中的每个元素
	var result []string
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		if elem.Kind() != reflect.Struct {
			return "", errors.New("elements must be structs or pointers to structs")
		}

		// 遍历结构体字段，匹配 tag
		t := elem.Type()
		for j := 0; j < elem.NumField(); j++ {
			field := t.Field(j)
			if tag := field.Tag.Get(tagName); tag == tagValue {
				result = append(result, fmt.Sprintf("%v", elem.Field(j).Interface()))
				break
			}
		}
	}

	return strings.Join(result, ","), nil
}
