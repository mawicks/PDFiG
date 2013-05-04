package pdf

type StreamFilter interface {
	Encode([]byte) []byte
	Decode([]byte) (decoded []byte,ok bool)
	Name() string
}