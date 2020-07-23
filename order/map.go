package order

import (
	"bytes"
	"encoding/json"
	"fmt"

	goyaml "gopkg.in/yaml.v3"
)

type MapSlice []MapItem

type MapItem struct {
	Key, Val interface{}
}

func (m MapSlice) Merge(in MapSlice) MapSlice {
	t := make(map[interface{}]struct{}, len(in))
	for _, item := range in {
		t[item.Key] = struct{}{}
	}
	var out MapSlice
	for _, item := range m {
		if _, ok := t[item.Key]; !ok {
			out = append(out, item)
		}
	}
	for _, item := range in {
		out = append(out, item)
	}
	return out
}

var MapSliceEscapeHTML = false

func (m MapSlice) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.Write([]byte{'{'})
	for i, item := range m {
		ibuf := &bytes.Buffer{}
		enc := json.NewEncoder(ibuf)
		enc.SetEscapeHTML(MapSliceEscapeHTML)
		if err := enc.Encode(&item.Val); err != nil {
			return nil, err
		}
		buf.WriteString(fmt.Sprintf("%q:", fmt.Sprint(item.Key)))
		buf.Write(ibuf.Bytes())
		if i < len(m)-1 {
			buf.Write([]byte{','})
		}
	}
	buf.Write([]byte{'}'})
	return buf.Bytes(), nil
}

func (m MapSlice) MarshalYAML() (interface{}, error) {
	var node goyaml.Node
	if err := node.Encode(map[string]interface{}{}); err != nil {
		return nil, err
	}
	node.Style = 0
	node.Content = make([]*goyaml.Node, 0, len(m)*2)
	for _, item := range m {
		var knode, vnode goyaml.Node
		if err := knode.Encode(item.Key); err != nil {
			return nil, err
		}
		if err := vnode.Encode(item.Val); err != nil {
			return nil, err
		}
		node.Content = append(node.Content, &knode, &vnode)
	}
	return &node, nil
}
