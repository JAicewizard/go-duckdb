package duckdb

/*
#include <duckdb.h>
*/
import "C"

import (
	"reflect"
	"runtime"
	"time"
	"unsafe"
)

type structEntry struct {
	TypeInfo
	name string
}

// StructEntry is an interface to provide STRUCT entry information.
type StructEntry interface {
	// Info returns a STRUCT entry's type information.
	Info() TypeInfo
	// Name returns a STRUCT entry's name.
	Name() string
}

// NewStructEntry returns a STRUCT entry.
// info contains information about the entry's type, and name holds the entry's name.
func NewStructEntry(info TypeInfo, name string) (StructEntry, error) {
	if name == "" {
		return nil, getError(errAPI, errEmptyName)
	}

	return &structEntry{
		TypeInfo: info,
		name:     name,
	}, nil
}

// Info returns a STRUCT entry's type information.
func (entry *structEntry) Info() TypeInfo {
	return entry.TypeInfo
}

// Name returns a STRUCT entry's name.
func (entry *structEntry) Name() string {
	return entry.name
}

type baseTypeInfo struct {
	Type
	structEntries []StructEntry
	decimalWidth  uint8
	decimalScale  uint8
}

type vectorTypeInfo struct {
	baseTypeInfo
	dict map[string]uint32
}

type typeInfo struct {
	baseTypeInfo
	childTypes []TypeInfo
	enumNames  []string
}

// TypeInfo is an interface for a DuckDB type.
type TypeInfo interface {
	logicalType() C.duckdb_logical_type
}

// NewTypeInfo returns type information for DuckDB's primitive types.
// It returns the TypeInfo, if the Type parameter is a valid primitive type.
// Else, it returns nil, and an error.
// Valid types are:
// TYPE_[BOOLEAN, TINYINT, SMALLINT, INTEGER, BIGINT, UTINYINT, USMALLINT, UINTEGER,
// UBIGINT, FLOAT, DOUBLE, TIMESTAMP, DATE, TIME, INTERVAL, HUGEINT, VARCHAR, BLOB,
// TIMESTAMP_S, TIMESTAMP_MS, TIMESTAMP_NS, UUID, TIMESTAMP_TZ, ANY].
func NewTypeInfo(t Type) (TypeInfo, error) {
	name, inMap := unsupportedTypeToStringMap[t]
	if inMap && t != TYPE_ANY {
		return nil, getError(errAPI, unsupportedTypeError(name))
	}

	switch t {
	case TYPE_DECIMAL:
		return nil, getError(errAPI, tryOtherFuncError(funcName(NewDecimalInfo)))
	case TYPE_ENUM:
		return nil, getError(errAPI, tryOtherFuncError(funcName(NewEnumInfo)))
	case TYPE_LIST:
		return nil, getError(errAPI, tryOtherFuncError(funcName(NewListInfo)))
	case TYPE_STRUCT:
		return nil, getError(errAPI, tryOtherFuncError(funcName(NewStructInfo)))
	case TYPE_MAP:
		return nil, getError(errAPI, tryOtherFuncError(funcName(NewMapInfo)))
	}

	return &typeInfo{
		baseTypeInfo: baseTypeInfo{Type: t},
	}, nil
}

// NewDecimalInfo returns DECIMAL type information.
// Its input parameters are the width and scale of the DECIMAL type.
func NewDecimalInfo(width uint8, scale uint8) (TypeInfo, error) {
	if width < 1 || width > MAX_DECIMAL_WIDTH {
		return nil, getError(errAPI, errInvalidDecimalWidth)
	}
	if scale > width {
		return nil, getError(errAPI, errInvalidDecimalScale)
	}

	return &typeInfo{
		baseTypeInfo: baseTypeInfo{
			Type:         TYPE_DECIMAL,
			decimalWidth: width,
			decimalScale: scale,
		},
	}, nil
}

// NewEnumInfo returns ENUM type information.
// Its input parameters are the dictionary values.
func NewEnumInfo(first string, others ...string) (TypeInfo, error) {
	// Check for duplicate names.
	m := map[string]bool{}
	m[first] = true
	for _, name := range others {
		_, inMap := m[name]
		if inMap {
			return nil, getError(errAPI, duplicateNameError(name))
		}
		m[name] = true
	}

	info := &typeInfo{
		baseTypeInfo: baseTypeInfo{
			Type: TYPE_ENUM,
		},
		enumNames: make([]string, 0),
	}

	info.enumNames = append(info.enumNames, first)
	info.enumNames = append(info.enumNames, others...)
	return info, nil
}

// NewListInfo returns LIST type information.
// childInfo contains the type information of the LIST's elements.
func NewListInfo(childInfo TypeInfo) (TypeInfo, error) {
	if childInfo == nil {
		return nil, getError(errAPI, interfaceIsNilError("childInfo"))
	}

	info := &typeInfo{
		baseTypeInfo: baseTypeInfo{Type: TYPE_LIST},
		childTypes:   make([]TypeInfo, 1),
	}
	info.childTypes[0] = childInfo
	return info, nil
}

// NewStructInfo returns STRUCT type information.
// Its input parameters are the STRUCT entries.
func NewStructInfo(firstEntry StructEntry, others ...StructEntry) (TypeInfo, error) {
	if firstEntry == nil {
		return nil, getError(errAPI, interfaceIsNilError("firstEntry"))
	}
	if firstEntry.Info() == nil {
		return nil, getError(errAPI, interfaceIsNilError("firstEntry.Info()"))
	}
	for i, entry := range others {
		if entry == nil {
			return nil, getError(errAPI, addIndexToError(interfaceIsNilError("entry"), i))
		}
		if entry.Info() == nil {
			return nil, getError(errAPI, addIndexToError(interfaceIsNilError("entry.Info()"), i))
		}
	}

	// Check for duplicate names.
	m := map[string]bool{}
	m[firstEntry.Name()] = true
	for _, entry := range others {
		name := entry.Name()
		_, inMap := m[name]
		if inMap {
			return nil, getError(errAPI, duplicateNameError(name))
		}
		m[name] = true
	}

	info := &typeInfo{
		baseTypeInfo: baseTypeInfo{
			Type:          TYPE_STRUCT,
			structEntries: make([]StructEntry, 0),
		},
	}
	info.structEntries = append(info.structEntries, firstEntry)
	info.structEntries = append(info.structEntries, others...)
	return info, nil
}

// NewMapInfo returns MAP type information.
// keyInfo contains the type information of the MAP keys.
// valueInfo contains the type information of the MAP values.
func NewMapInfo(keyInfo TypeInfo, valueInfo TypeInfo) (TypeInfo, error) {
	if keyInfo == nil {
		return nil, getError(errAPI, interfaceIsNilError("keyInfo"))
	}
	if valueInfo == nil {
		return nil, getError(errAPI, interfaceIsNilError("valueInfo"))
	}

	info := &typeInfo{
		baseTypeInfo: baseTypeInfo{Type: TYPE_MAP},
		childTypes:   make([]TypeInfo, 2),
	}
	info.childTypes[0] = keyInfo
	info.childTypes[1] = valueInfo
	return info, nil
}

func (info *typeInfo) logicalType() C.duckdb_logical_type {
	switch info.Type {
	case TYPE_BOOLEAN, TYPE_TINYINT, TYPE_SMALLINT, TYPE_INTEGER, TYPE_BIGINT, TYPE_UTINYINT, TYPE_USMALLINT,
		TYPE_UINTEGER, TYPE_UBIGINT, TYPE_FLOAT, TYPE_DOUBLE, TYPE_TIMESTAMP, TYPE_TIMESTAMP_S, TYPE_TIMESTAMP_MS,
		TYPE_TIMESTAMP_NS, TYPE_TIMESTAMP_TZ, TYPE_DATE, TYPE_TIME, TYPE_INTERVAL, TYPE_HUGEINT, TYPE_VARCHAR,
		TYPE_BLOB, TYPE_UUID, TYPE_ANY:
		return C.duckdb_create_logical_type(C.duckdb_type(info.Type))

	case TYPE_DECIMAL:
		return C.duckdb_create_decimal_type(C.uint8_t(info.decimalWidth), C.uint8_t(info.decimalScale))
	case TYPE_ENUM:
		return info.logicalEnumType()
	case TYPE_LIST:
		return info.logicalListType()
	case TYPE_STRUCT:
		return info.logicalStructType()
	case TYPE_MAP:
		return info.logicalMapType()
	}
	return nil
}

func (info *typeInfo) logicalEnumType() C.duckdb_logical_type {
	count := len(info.enumNames)
	size := C.size_t(unsafe.Sizeof((*C.char)(nil)))
	names := (*[1 << 31]*C.char)(C.malloc(C.size_t(count) * size))

	for i, name := range info.enumNames {
		(*names)[i] = C.CString(name)
	}
	cNames := (**C.char)(unsafe.Pointer(names))
	logicalType := C.duckdb_create_enum_type(cNames, C.idx_t(count))

	for i := 0; i < count; i++ {
		C.duckdb_free(unsafe.Pointer((*names)[i]))
	}
	C.duckdb_free(unsafe.Pointer(names))
	return logicalType
}

func (info *typeInfo) logicalListType() C.duckdb_logical_type {
	child := info.childTypes[0].logicalType()
	logicalType := C.duckdb_create_list_type(child)
	C.duckdb_destroy_logical_type(&child)
	return logicalType
}

func (info *typeInfo) logicalStructType() C.duckdb_logical_type {
	count := len(info.structEntries)
	size := C.size_t(unsafe.Sizeof(C.duckdb_logical_type(nil)))
	types := (*[1 << 31]C.duckdb_logical_type)(C.malloc(C.size_t(count) * size))

	size = C.size_t(unsafe.Sizeof((*C.char)(nil)))
	names := (*[1 << 31]*C.char)(C.malloc(C.size_t(count) * size))

	for i, entry := range info.structEntries {
		(*types)[i] = entry.Info().logicalType()
		(*names)[i] = C.CString(entry.Name())
	}

	cTypes := (*C.duckdb_logical_type)(unsafe.Pointer(types))
	cNames := (**C.char)(unsafe.Pointer(names))
	logicalType := C.duckdb_create_struct_type(cTypes, cNames, C.idx_t(count))

	for i := 0; i < count; i++ {
		C.duckdb_destroy_logical_type(&types[i])
		C.duckdb_free(unsafe.Pointer((*names)[i]))
	}
	C.duckdb_free(unsafe.Pointer(types))
	C.duckdb_free(unsafe.Pointer(names))
	return logicalType
}

func (info *typeInfo) logicalMapType() C.duckdb_logical_type {
	key := info.childTypes[0].logicalType()
	value := info.childTypes[1].logicalType()
	logicalType := C.duckdb_create_map_type(key, value)

	C.duckdb_destroy_logical_type(&key)
	C.duckdb_destroy_logical_type(&value)
	return logicalType
}

// NewDuckdbType creates a new Type for T. All valid T are guaranteeed to have a valid
// representation in duckdb, thus no error is returned.
func NewDuckdbType[T SaveTypes]() TypeInfo {
	t, _ := tryGetDuckdbType[T]()
	return t
}

// TryNewDuckdbType creates a new Type for T. Since not all valid T are guaranteeed to have
// a valid representation in duckdb, an error may be returned.
func TryNewDuckdbType[T any]() (TypeInfo, error) {
	return tryGetDuckdbType[T]()
}

// TryNewDuckdbTypeFrom Value creates a new Type for the type of v. Since not all valid T
// guaranteeed to have a valid representation in duckdb, an error may be returned.
func TryNewDuckdbTypeFromValue(v any) (TypeInfo, error) {
	return tryGetDuckdbTypeFromValue(reflect.TypeOf(v))
}

func tryGetDuckdbType[T any]() (TypeInfo, error) {
	var v T
	return tryGetDuckdbTypeFromValue(reflect.TypeOf(v))
}

func canConvertToDuckdb(rt reflect.Type) bool {
	switch rt {
	case reflect.TypeOf(time.Time{}),
		reflect.TypeOf(UUID{}),
		reflect.TypeOf([]byte{}):
		return true
	}
	switch rt.Kind() {
	// Invalid types
	case reflect.Chan, reflect.Func,
		reflect.UnsafePointer, reflect.Uintptr,
		reflect.Int, reflect.Uint,
		reflect.Complex64, reflect.Complex128:
		return false
	// Valid types
	case reflect.Bool,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return true
	case reflect.Struct:
		for i := 0; i < rt.NumField(); i++ {
			if rt.Field(i).IsExported() && !canConvertToDuckdb(rt.Field(i).Type) {
				return false
			}
		}
		return true
	case reflect.Array, reflect.Slice:
		elemt := rt.Elem()
		return canConvertToDuckdb(elemt)
	case reflect.Map:
		keyt := rt.Key()
		elemt := rt.Elem()
		return canConvertToDuckdb(elemt) && canConvertToDuckdb(keyt)
	case reflect.Pointer, reflect.Interface:
		return canConvertToDuckdb(rt.Elem())
	// This case should never be reached
	default:
		return false
	}
}

func tryGetDuckdbTypeFromValue(rt reflect.Type) (TypeInfo, error) {
	switch rt {
	case reflect.TypeOf(time.Time{}):
		return NewTypeInfo(Type(TYPE_TIMESTAMP_NS))
	case reflect.TypeOf(UUID{}):
		return NewTypeInfo(Type(TYPE_UHUGEINT))
	case reflect.TypeOf([]byte{}):
		return NewTypeInfo(Type(TYPE_UUID))
	}
	switch rt.Kind() {
	// Invalid types
	case reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Int, reflect.Uint, reflect.Uintptr, reflect.Complex64, reflect.Complex128:
		return nil, unsupportedTypeError(rt.String())
	// Valid types
	case reflect.Bool:
		return NewTypeInfo(Type(TYPE_BOOLEAN))
	case reflect.Int8:
		return NewTypeInfo(Type(TYPE_TINYINT))
	case reflect.Int16:
		return NewTypeInfo(Type(TYPE_SMALLINT))
	case reflect.Int32:
		return NewTypeInfo(Type(TYPE_INTEGER))
	case reflect.Int64:
		return NewTypeInfo(Type(TYPE_BIGINT))
	case reflect.Uint8:
		return NewTypeInfo(Type(TYPE_UTINYINT))
	case reflect.Uint16:
		return NewTypeInfo(Type(TYPE_USMALLINT))
	case reflect.Uint32:
		return NewTypeInfo(Type(TYPE_UINTEGER))
	case reflect.Uint64:
		return NewTypeInfo(Type(TYPE_UBIGINT))
	case reflect.Float32:
		return NewTypeInfo(Type(TYPE_FLOAT))
	case reflect.Float64:
		return NewTypeInfo(Type(TYPE_DOUBLE))
	case reflect.String:
		return NewTypeInfo(Type(TYPE_VARCHAR))
	case reflect.Struct:
		fields := make([]StructEntry, 0, rt.NumField())
		for i := 0; i < rt.NumField(); i++ {
			if rt.Field(i).IsExported() {
				field := rt.Field(i)
				typ, err := tryGetDuckdbTypeFromValue(field.Type)
				if err != nil {
					return nil, err
				}

				entry, err := NewStructEntry(typ, field.Name)
				fields = append(fields, entry)
			}
		}
		return NewStructInfo(fields[0], fields[1:]...)
	case reflect.Slice:
		elemt := rt.Elem()
		t, err := tryGetDuckdbTypeFromValue(elemt)
		if err != nil {
			return nil, err
		}
		return NewListInfo(t)
	case reflect.Map:
		keyt := rt.Key()
		kt, err := tryGetDuckdbTypeFromValue(keyt)
		if err != nil {
			return nil, err
		}
		elemt := rt.Elem()
		vt, err := tryGetDuckdbTypeFromValue(elemt)
		if err != nil {
			return nil, err
		}
		return NewMapInfo(kt, vt)
	case reflect.Pointer, reflect.Interface:
		return tryGetDuckdbTypeFromValue(rt.Elem())
	// This case should never be reached
	default:
		return nil, unsupportedTypeError(rt.String())
	}
}

func funcName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
