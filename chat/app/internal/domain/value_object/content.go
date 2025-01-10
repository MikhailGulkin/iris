package value_object

type Content struct {
	Text  *string
	Type  ContentType
	Bytes []byte
}
