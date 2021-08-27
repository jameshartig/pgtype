// Code generated by erb. DO NOT EDIT.

package pgtype

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"reflect"
	"time"

	"github.com/jackc/pgio"
)

type DateArray struct {
	Elements   []Date
	Dimensions []ArrayDimension
	Valid      bool
}

func (dst *DateArray) Set(src interface{}) error {
	// untyped nil and typed nil interfaces are different
	if src == nil {
		*dst = DateArray{}
		return nil
	}

	if value, ok := src.(interface{ Get() interface{} }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	// Attempt to match to select common types:
	switch value := src.(type) {

	case []time.Time:
		if value == nil {
			*dst = DateArray{}
		} else if len(value) == 0 {
			*dst = DateArray{Valid: true}
		} else {
			elements := make([]Date, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = DateArray{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Valid:      true,
			}
		}

	case []*time.Time:
		if value == nil {
			*dst = DateArray{}
		} else if len(value) == 0 {
			*dst = DateArray{Valid: true}
		} else {
			elements := make([]Date, len(value))
			for i := range value {
				if err := elements[i].Set(value[i]); err != nil {
					return err
				}
			}
			*dst = DateArray{
				Elements:   elements,
				Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Valid:      true,
			}
		}

	case []Date:
		if value == nil {
			*dst = DateArray{}
		} else if len(value) == 0 {
			*dst = DateArray{Valid: true}
		} else {
			*dst = DateArray{
				Elements:   value,
				Dimensions: []ArrayDimension{{Length: int32(len(value)), LowerBound: 1}},
				Valid:      true,
			}
		}
	default:
		// Fallback to reflection if an optimised match was not found.
		// The reflection is necessary for arrays and multidimensional slices,
		// but it comes with a 20-50% performance penalty for large arrays/slices
		reflectedValue := reflect.ValueOf(src)
		if !reflectedValue.IsValid() || reflectedValue.IsZero() {
			*dst = DateArray{}
			return nil
		}

		dimensions, elementsLength, ok := findDimensionsFromValue(reflectedValue, nil, 0)
		if !ok {
			return fmt.Errorf("cannot find dimensions of %v for DateArray", src)
		}
		if elementsLength == 0 {
			*dst = DateArray{Valid: true}
			return nil
		}
		if len(dimensions) == 0 {
			if originalSrc, ok := underlyingSliceType(src); ok {
				return dst.Set(originalSrc)
			}
			return fmt.Errorf("cannot convert %v to DateArray", src)
		}

		*dst = DateArray{
			Elements:   make([]Date, elementsLength),
			Dimensions: dimensions,
			Valid:      true,
		}
		elementCount, err := dst.setRecursive(reflectedValue, 0, 0)
		if err != nil {
			// Maybe the target was one dimension too far, try again:
			if len(dst.Dimensions) > 1 {
				dst.Dimensions = dst.Dimensions[:len(dst.Dimensions)-1]
				elementsLength = 0
				for _, dim := range dst.Dimensions {
					if elementsLength == 0 {
						elementsLength = int(dim.Length)
					} else {
						elementsLength *= int(dim.Length)
					}
				}
				dst.Elements = make([]Date, elementsLength)
				elementCount, err = dst.setRecursive(reflectedValue, 0, 0)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		if elementCount != len(dst.Elements) {
			return fmt.Errorf("cannot convert %v to DateArray, expected %d dst.Elements, but got %d instead", src, len(dst.Elements), elementCount)
		}
	}

	return nil
}

func (dst *DateArray) setRecursive(value reflect.Value, index, dimension int) (int, error) {
	switch value.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		if len(dst.Dimensions) == dimension {
			break
		}

		valueLen := value.Len()
		if int32(valueLen) != dst.Dimensions[dimension].Length {
			return 0, fmt.Errorf("multidimensional arrays must have array expressions with matching dimensions")
		}
		for i := 0; i < valueLen; i++ {
			var err error
			index, err = dst.setRecursive(value.Index(i), index, dimension+1)
			if err != nil {
				return 0, err
			}
		}

		return index, nil
	}
	if !value.CanInterface() {
		return 0, fmt.Errorf("cannot convert all values to DateArray")
	}
	if err := dst.Elements[index].Set(value.Interface()); err != nil {
		return 0, fmt.Errorf("%v in DateArray", err)
	}
	index++

	return index, nil
}

func (dst DateArray) Get() interface{} {
	if !dst.Valid {
		return nil
	}
	return dst
}

func (src *DateArray) AssignTo(dst interface{}) error {
	if !src.Valid {
		return NullAssignTo(dst)
	}

	if len(src.Dimensions) <= 1 {
		// Attempt to match to select common types:
		switch v := dst.(type) {

		case *[]time.Time:
			*v = make([]time.Time, len(src.Elements))
			for i := range src.Elements {
				if err := src.Elements[i].AssignTo(&((*v)[i])); err != nil {
					return err
				}
			}
			return nil

		case *[]*time.Time:
			*v = make([]*time.Time, len(src.Elements))
			for i := range src.Elements {
				if err := src.Elements[i].AssignTo(&((*v)[i])); err != nil {
					return err
				}
			}
			return nil

		}
	}

	// Try to convert to something AssignTo can use directly.
	if nextDst, retry := GetAssignToDstType(dst); retry {
		return src.AssignTo(nextDst)
	}

	// Fallback to reflection if an optimised match was not found.
	// The reflection is necessary for arrays and multidimensional slices,
	// but it comes with a 20-50% performance penalty for large arrays/slices
	value := reflect.ValueOf(dst)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return fmt.Errorf("cannot assign %T to %T", src, dst)
	}

	if len(src.Elements) == 0 {
		if value.Kind() == reflect.Slice {
			value.Set(reflect.MakeSlice(value.Type(), 0, 0))
			return nil
		}
	}

	elementCount, err := src.assignToRecursive(value, 0, 0)
	if err != nil {
		return err
	}
	if elementCount != len(src.Elements) {
		return fmt.Errorf("cannot assign %v, needed to assign %d elements, but only assigned %d", dst, len(src.Elements), elementCount)
	}

	return nil
}

func (src *DateArray) assignToRecursive(value reflect.Value, index, dimension int) (int, error) {
	switch kind := value.Kind(); kind {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		if len(src.Dimensions) == dimension {
			break
		}

		length := int(src.Dimensions[dimension].Length)
		if reflect.Array == kind {
			typ := value.Type()
			if typ.Len() != length {
				return 0, fmt.Errorf("expected size %d array, but %s has size %d array", length, typ, typ.Len())
			}
			value.Set(reflect.New(typ).Elem())
		} else {
			value.Set(reflect.MakeSlice(value.Type(), length, length))
		}

		var err error
		for i := 0; i < length; i++ {
			index, err = src.assignToRecursive(value.Index(i), index, dimension+1)
			if err != nil {
				return 0, err
			}
		}

		return index, nil
	}
	if len(src.Dimensions) != dimension {
		return 0, fmt.Errorf("incorrect dimensions, expected %d, found %d", len(src.Dimensions), dimension)
	}
	if !value.CanAddr() {
		return 0, fmt.Errorf("cannot assign all values from DateArray")
	}
	addr := value.Addr()
	if !addr.CanInterface() {
		return 0, fmt.Errorf("cannot assign all values from DateArray")
	}
	if err := src.Elements[index].AssignTo(addr.Interface()); err != nil {
		return 0, err
	}
	index++
	return index, nil
}

func (dst *DateArray) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = DateArray{}
		return nil
	}

	uta, err := ParseUntypedTextArray(string(src))
	if err != nil {
		return err
	}

	var elements []Date

	if len(uta.Elements) > 0 {
		elements = make([]Date, len(uta.Elements))

		for i, s := range uta.Elements {
			var elem Date
			var elemSrc []byte
			if s != "NULL" || uta.Quoted[i] {
				elemSrc = []byte(s)
			}
			err = elem.DecodeText(ci, elemSrc)
			if err != nil {
				return err
			}

			elements[i] = elem
		}
	}

	*dst = DateArray{Elements: elements, Dimensions: uta.Dimensions, Valid: true}

	return nil
}

func (dst *DateArray) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = DateArray{}
		return nil
	}

	var arrayHeader ArrayHeader
	rp, err := arrayHeader.DecodeBinary(ci, src)
	if err != nil {
		return err
	}

	if len(arrayHeader.Dimensions) == 0 {
		*dst = DateArray{Dimensions: arrayHeader.Dimensions, Valid: true}
		return nil
	}

	elementCount := arrayHeader.Dimensions[0].Length
	for _, d := range arrayHeader.Dimensions[1:] {
		elementCount *= d.Length
	}

	elements := make([]Date, elementCount)

	for i := range elements {
		elemLen := int(int32(binary.BigEndian.Uint32(src[rp:])))
		rp += 4
		var elemSrc []byte
		if elemLen >= 0 {
			elemSrc = src[rp : rp+elemLen]
			rp += elemLen
		}
		err = elements[i].DecodeBinary(ci, elemSrc)
		if err != nil {
			return err
		}
	}

	*dst = DateArray{Elements: elements, Dimensions: arrayHeader.Dimensions, Valid: true}
	return nil
}

func (src DateArray) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	if !src.Valid {
		return nil, nil
	}

	if len(src.Dimensions) == 0 {
		return append(buf, '{', '}'), nil
	}

	buf = EncodeTextArrayDimensions(buf, src.Dimensions)

	// dimElemCounts is the multiples of elements that each array lies on. For
	// example, a single dimension array of length 4 would have a dimElemCounts of
	// [4]. A multi-dimensional array of lengths [3,5,2] would have a
	// dimElemCounts of [30,10,2]. This is used to simplify when to render a '{'
	// or '}'.
	dimElemCounts := make([]int, len(src.Dimensions))
	dimElemCounts[len(src.Dimensions)-1] = int(src.Dimensions[len(src.Dimensions)-1].Length)
	for i := len(src.Dimensions) - 2; i > -1; i-- {
		dimElemCounts[i] = int(src.Dimensions[i].Length) * dimElemCounts[i+1]
	}

	inElemBuf := make([]byte, 0, 32)
	for i, elem := range src.Elements {
		if i > 0 {
			buf = append(buf, ',')
		}

		for _, dec := range dimElemCounts {
			if i%dec == 0 {
				buf = append(buf, '{')
			}
		}

		elemBuf, err := elem.EncodeText(ci, inElemBuf)
		if err != nil {
			return nil, err
		}
		if elemBuf == nil {
			buf = append(buf, `NULL`...)
		} else {
			buf = append(buf, QuoteArrayElementIfNeeded(string(elemBuf))...)
		}

		for _, dec := range dimElemCounts {
			if (i+1)%dec == 0 {
				buf = append(buf, '}')
			}
		}
	}

	return buf, nil
}

func (src DateArray) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	if !src.Valid {
		return nil, nil
	}

	arrayHeader := ArrayHeader{
		Dimensions: src.Dimensions,
	}

	if dt, ok := ci.DataTypeForName("date"); ok {
		arrayHeader.ElementOID = int32(dt.OID)
	} else {
		return nil, fmt.Errorf("unable to find oid for type name %v", "date")
	}

	for i := range src.Elements {
		if !src.Elements[i].Valid {
			arrayHeader.ContainsNull = true
			break
		}
	}

	buf = arrayHeader.EncodeBinary(ci, buf)

	for i := range src.Elements {
		sp := len(buf)
		buf = pgio.AppendInt32(buf, -1)

		elemBuf, err := src.Elements[i].EncodeBinary(ci, buf)
		if err != nil {
			return nil, err
		}
		if elemBuf != nil {
			buf = elemBuf
			pgio.SetInt32(buf[sp:], int32(len(buf[sp:])-4))
		}
	}

	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *DateArray) Scan(src interface{}) error {
	if src == nil {
		return dst.DecodeText(nil, nil)
	}

	switch src := src.(type) {
	case string:
		return dst.DecodeText(nil, []byte(src))
	case []byte:
		srcCopy := make([]byte, len(src))
		copy(srcCopy, src)
		return dst.DecodeText(nil, srcCopy)
	}

	return fmt.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src DateArray) Value() (driver.Value, error) {
	buf, err := src.EncodeText(nil, nil)
	if err != nil {
		return nil, err
	}
	if buf == nil {
		return nil, nil
	}

	return string(buf), nil
}
