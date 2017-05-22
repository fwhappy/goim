package msgpack

import (
	"bytes"

	"github.com/fwhappy/util"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

// Marshal JsonMap to []byte
func Marshal(data util.JsonMap) []byte {
	buf := new(bytes.Buffer)
	enc := msgpack.NewEncoder(buf)
	_ = enc.Encode(data)
	return buf.Bytes()
}

// Unmarshal []byte to JsonMap
func Unmarshal(p []byte) (util.JsonMap, error) {
	buf := new(bytes.Buffer)
	buf.Write(p)
	dec := msgpack.NewDecoder(buf)

	dec.DecodeMapFunc = func(d *msgpack.Decoder) (interface{}, error) {
		n, err := d.DecodeMapLen()
		if err != nil {
			return nil, err
		}

		m := make(map[string]interface{}, n)
		for i := 0; i < n; i++ {
			mk, err := d.DecodeString()
			if err != nil {
				return nil, err
			}

			mv, err := d.DecodeInterface()
			if err != nil {
				return nil, err
			}

			m[mk] = mv
		}
		return m, nil
	}

	out, err := dec.DecodeInterface()
	if err != nil {
		return nil, err
	}

	return util.JsonMap(out.(map[string]interface{})), nil
}
