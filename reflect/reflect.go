package reflect

import (
	"errors"
	"reflect"
	"runtime"
	"strings"
)

type Utils interface {
	GetFuncName(fn any) string                                     //从函数实例获取函数名
	GetStructName(obj any) string                                  // 从实例获取结构名
	GetVarName(v any) string                                       // 获取变量名
	StructSet(obj any, nilField any, val any) error                // 给结构体设置field类型的值
	MatchReceiverMethods(receiver any, matchFn any) map[string]any // 匹配receiver的所有methods中与matchFn签名参数类似的方法
	GetFuncSignature(fn any) string                                // 获取函数签名信息
	InspectValue(v any) Value                                      // 检索Value的信息
	FuncEqual(fn1, fn2 any) bool                                   // 函数是否相等
	IsAssignableStruct(obj any) bool                               // 是否是可赋值的结构指针类型
	Indirect(a any) any                                            // 通过尽可能的解引用获取底层的类型
}

type Value struct {
	Name      string
	IsPointer bool
	Kind      string
	Items     []ValueItem
}

type ValueItem struct {
	Name  string
	Kind  string
	Value any
}

// IsAssignableStruct 是否是可赋值的结构指针类型
func IsAssignableStruct(obj any) bool {
	varType := reflect.TypeOf(obj)
	if varType.Kind() == reflect.Ptr && varType.Elem().Kind() == reflect.Struct {
		return true
	}
	return false
}

// GetFuncName 从函数实例获取函数名
func GetFuncName(fn any) string {
	tokens := strings.Split(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), ".")
	return strings.Split(tokens[len(tokens)-1], "-")[0]
}

// GetStructName 获取结构名
func GetStructName(obj any) string {
	if obj == nil {
		return ""
	}

	var st reflect.Type
	if t := reflect.TypeOf(obj); t.Kind() == reflect.Ptr {
		st = t.Elem()
	} else {
		st = t
	}
	if st.Kind() != reflect.Struct {
		return ""
	}
	return st.Name()
}

// GetVarName 获取变量名
func GetVarName(v any) string {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func FuncEqual(fn1, fn2 any) bool {
	v1 := reflect.ValueOf(&fn1).Elem()
	v2 := reflect.ValueOf(&fn2).Elem()
	return v1.Interface() == v2.Interface()
}

func InspectValue(v any) *Value {
	var isPointer bool
	var st reflect.Type
	var sv reflect.Value
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		isPointer = true
		st = t.Elem()
		sv = reflect.ValueOf(v).Elem()
	} else {
		st = t
		sv = reflect.ValueOf(v)
	}

	var items []ValueItem
	switch st.Kind() {
	case reflect.Struct:
		items = GetStructFields(st, sv)
	case reflect.Slice:
		items = GetSliceItems(sv)
	}

	return &Value{
		Name:      st.Name(),
		IsPointer: isPointer,
		Kind:      st.Kind().String(),
		Items:     items,
	}
}

func GetStructFields(st reflect.Type, sv reflect.Value) []ValueItem {
	fields := make([]ValueItem, 0)
	for i := 0; i < st.NumField(); i++ {
		switch v := sv.Field(i).Interface().(type) {
		default:
			fields = append(fields, ValueItem{
				Name:  st.Field(i).Name,
				Kind:  st.Field(i).Type.Kind().String(),
				Value: v,
			})
		}
	}
	return fields
}

func GetSliceItems(sv reflect.Value) []ValueItem {
	items := make([]ValueItem, 0)
	for i := 0; i < sv.Len(); i++ {
		switch v := sv.Index(i).Interface().(type) {
		default:
			items = append(items, ValueItem{
				Name:  "",
				Kind:  sv.Index(i).Type().Kind().String(),
				Value: v,
			})
		}
	}
	return items
}

// StructSet 通过反射给结构体obj的nilField类型匹配的field赋值val
func StructSet(obj any, nilField any, val any) error {
	if obj == nil {
		return errors.New("nil struct")
	}

	destv := reflect.ValueOf(obj)
	if destv.Kind() == reflect.Ptr && !destv.Elem().CanSet() {
		return errors.New("struct should be a pointer")
	} else {
		destv = destv.Elem()
	}

	if val == nil {
		return errors.New("cannot set to zero value")
	}

	numField := destv.NumField()
	for i := 0; i < numField; i++ {
		field := destv.Field(i)
		fieldType := field.Type().String()
		emptyFieldType := reflect.TypeOf(nilField).String()
		// 找到第一个匹配类型的field设置进去
		if fieldType == emptyFieldType || "*"+fieldType == emptyFieldType {
			if !field.CanSet() {
				return errors.New("struct field cannot set")
			}

			field.Set(reflect.ValueOf(val))
			return nil
		}
	}
	return errors.New("struct field not found")
}

// MatchReceiverMethods 匹配receiver的所有methods中与matchFn签名参数类似的方法
func MatchReceiverMethods(receiver any, matchFn any) map[string]any {
	if receiver == nil {
		return nil
	}

	st := reflect.TypeOf(receiver)
	sv := reflect.ValueOf(receiver)
	numMethod := sv.NumMethod()

	receivers := make(map[string]any)
	for i := 0; i < numMethod; i++ {
		vv := sv.Method(i)
		if vv.Type().ConvertibleTo(reflect.TypeOf(matchFn)) {
			receivers[st.Method(i).Name] = vv.Interface()
		}
	}
	return receivers
}

// GetFuncSignature 获取函数签名信息
func GetFuncSignature(fn any) string {
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		return ""
	}

	buf := strings.Builder{}
	buf.WriteString("func (")
	for i := 0; i < t.NumIn(); i++ {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(t.In(i).String())
	}
	buf.WriteString(")")
	if numOut := t.NumOut(); numOut > 0 {
		if numOut > 1 {
			buf.WriteString(" (")
		} else {
			buf.WriteString(" ")
		}
		for i := 0; i < t.NumOut(); i++ {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(t.Out(i).String())
		}
		if numOut > 1 {
			buf.WriteString(")")
		}
	}

	return buf.String()
}

// Indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func Indirect(a any) any {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}
