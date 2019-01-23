package convert

// TODO: []byte -> io.ReadCloser / io.Writer
type Encoding interface {
	String() string
	Encode(input interface{}) ([]byte, error)
	Decode(input []byte) (interface{}, error)
}
