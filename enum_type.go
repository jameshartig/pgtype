package pgtype

import "fmt"

// EnumType represents a enum type. While it implements Value, this is only in service of its type conversion duties
// when registered as a data type in a ConnType. It should not be used directly as a Value.
type EnumType struct {
	value string
	valid bool

	typeName   string            // PostgreSQL type name
	members    []string          // enum members
	membersMap map[string]string // map to quickly lookup member and reuse string instead of allocating
}

// NewEnumType initializes a new EnumType. It retains a read-only reference to members. members must not be changed.
func NewEnumType(typeName string, members []string) *EnumType {
	et := &EnumType{typeName: typeName, members: members}
	et.membersMap = make(map[string]string, len(members))
	for _, m := range members {
		et.membersMap[m] = m
	}
	return et
}

func (et *EnumType) NewTypeValue() Value {
	return &EnumType{
		value: et.value,
		valid: et.valid,

		typeName:   et.typeName,
		members:    et.members,
		membersMap: et.membersMap,
	}
}

func (et *EnumType) TypeName() string {
	return et.typeName
}

func (et *EnumType) Members() []string {
	return et.members
}

// Set assigns src to dst. Set purposely does not check that src is a member. This allows continued error free
// operation in the event the PostgreSQL enum type is modified during a connection.
func (dst *EnumType) Set(src interface{}) error {
	if src == nil {
		dst.valid = false
		return nil
	}

	if value, ok := src.(interface{ Get() interface{} }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {
	case string:
		dst.value = value
		dst.valid = true
	case *string:
		if value == nil {
			dst.valid = false
		} else {
			dst.value = *value
			dst.valid = true
		}
	case []byte:
		if value == nil {
			dst.valid = false
		} else {
			dst.value = string(value)
			dst.valid = true
		}
	default:
		if originalSrc, ok := underlyingStringType(src); ok {
			return dst.Set(originalSrc)
		}
		return fmt.Errorf("cannot convert %v to enum %s", value, dst.typeName)
	}

	return nil
}

func (dst EnumType) Get() interface{} {
	if !dst.valid {
		return nil
	}
	return dst.value
}

func (src *EnumType) AssignTo(dst interface{}) error {
	if !src.valid {
		return NullAssignTo(dst)
	}

	switch v := dst.(type) {
	case *string:
		*v = src.value
		return nil
	case *[]byte:
		*v = make([]byte, len(src.value))
		copy(*v, src.value)
		return nil
	default:
		if nextDst, retry := GetAssignToDstType(dst); retry {
			return src.AssignTo(nextDst)
		}
		return fmt.Errorf("unable to assign to %T", dst)
	}
}

func (EnumType) PreferredResultFormat() int16 {
	return TextFormatCode
}

func (dst *EnumType) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		dst.valid = false
		return nil
	}

	// Lookup the string in membersMap to avoid an allocation.
	if s, found := dst.membersMap[string(src)]; found {
		dst.value = s
	} else {
		// If an enum type is modified after the initial connection it is possible to receive an unexpected value.
		// Gracefully handle this situation. Purposely NOT modifying members and membersMap to allow for sharing members
		// and membersMap between connections.
		dst.value = string(src)
	}
	dst.valid = true

	return nil
}

func (dst *EnumType) DecodeBinary(ci *ConnInfo, src []byte) error {
	return dst.DecodeText(ci, src)
}

func (EnumType) PreferredParamFormat() int16 {
	return TextFormatCode
}

func (src EnumType) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	if !src.valid {
		return nil, nil
	}

	return append(buf, src.value...), nil
}

func (src EnumType) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	return src.EncodeText(ci, buf)
}
