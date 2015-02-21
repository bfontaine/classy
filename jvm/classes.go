package jvm

// JClass is a Java class parsed from a .class file.
type JClass struct {
	majorVersion u2
	minorVersion u2
	accessFlags  u2

	constants  []JConst
	fields     []JField
	interfaces []u2

	classIndex      u2
	superClassIndex u2
}

// init the .constants field of a JClass with the given size
func (cls *JClass) initConstantPool(size u2) {
	cls.constants = make([]JConst, size)
}

// init the .fields field of a JClass with the given size
func (cls *JClass) initFields(size u2) {
	cls.fields = make([]JField, size)
}

// init the .interfaces field of a JClass with the given size
func (cls *JClass) initInterfaces(size u2) {
	cls.interfaces = make([]u2, size)
}

// add a constant to the pool.
// u2 is its index, tag its tag, and data its associated data
func (cls *JClass) addConstant(index u2, tag u1, data []byte) {
	cls.constants[index] = JConst{tag, data}
}

// add a field to the fields slice.
func (cls *JClass) addField(index u2, field JField) {
	cls.fields[index] = field
}

// return the constant at the given index in the pool, following references if
// there are some.
func (cls *JClass) resolveConstantIndex(index u2) JConst {
	cst := cls.constants[index]

	switch cst.tag {
	case TAG_CLASS_REF:
		fallthrough
	case TAG_STRING_REF:
		fallthrough
	case TAG_METHOD_TYPE:
		return cls.resolveConstantIndex(bytesToU2(cst.value))
	default:
		return cst
	}
}

// HasAccessFlag tests if the class has the given access flag.
func (cls *JClass) HasAccessFlag(flag u2) bool {
	return cls.accessFlags&flag == flag
}

// ClassName returns the class name as a string.
func (cls *JClass) ClassName() string {
	return cls.resolveConstantIndex(cls.classIndex).valueAsString()
}

// Constants returns all the constants of this class.
// Note that the first index (0) should be ignored.
func (cls *JClass) Constants() []JConst {
	return cls.constants
}

// JavaVersion returns the Java Version used by a class.
func (cls *JClass) JavaVersion() string {
	// we could also use the minor version here
	switch cls.majorVersion {
	case 45:
		return "JDK 1.1"
	case 46:
		return "JDK 1.2"
	case 47:
		return "JDK 1.3"
	case 48:
		return "JDK 1.4"
	case 49:
		return "J2SE 5.0"
	case 50:
		return "J2SE 6.0"
	case 51:
		return "J2SE 7"
	case 52:
		return "J2SE 8"
	}
	return "Unknown version"
}
