package tree

import "go/constant"

type Node interface {
	Version() int

	setParent(n Node)
	resetParent()
	incrementVersion()
}

var _ Node = &Slice{}
var _ Node = &Map{}
var _ Node = &Leaf{}

type base struct {
	version int
	parent  Node
}

func (b *base) Version() int {
	return b.version
}

func (b *base) setParent(n Node) {
	if b.parent != nil {
		panic("tree: node can not have multiple parents")
	}
	b.parent = n
}

func (b *base) resetParent() {
	if b.parent == nil {
		panic("tree: internal error")
	}
	b.parent = nil
}

func (b *base) incrementVersion() {
	b.version++
	if b.parent != nil {
		b.parent.incrementVersion()
	}
}

type Slice struct {
	*base
	value []Node
}

func MakeSlice(len, cap int) *Slice {
	return &Slice{
		base:  &base{version: 1},
		value: make([]Node, len, cap),
	}
}

func (s *Slice) Index(i int) Node {
	return s.value[i]
}

func (s *Slice) SetIndex(i int, value Node) {
	if prev := s.value[i]; prev != nil {
		prev.resetParent()
	}
	if value != nil {
		value.setParent(s)
	}
	s.value[i] = value
	s.incrementVersion()
}

func (s *Slice) Slice(i, j int) *Slice {
	return &Slice{
		base:  s.base,
		value: s.value[i:j],
	}
}

func (s *Slice) Slice3(i, j, k int) *Slice {
	return &Slice{
		base:  s.base,
		value: s.value[i:j:k],
	}
}

func (s *Slice) Append(values ...Node) *Slice {
	for _, v := range values {
		if v != nil {
			v.setParent(s)
		}
	}
	s2 := &Slice{
		base:  s.base,
		value: append(s.value, values...),
	}
	s.incrementVersion()
	return s2
}

func Move(dst, src *Slice) {
	copy(dst.value, src.value)
	for i, v := range src.value {
		if v != nil {
			v.resetParent()
			v.setParent(dst)
			src.value[i] = nil
		}
	}
	src.incrementVersion()
	dst.incrementVersion()
}

type Map struct {
	*base
	value map[string]Node
}

func MakeMap() *Map {
	return &Map{
		base:  &base{version: 1},
		value: make(map[string]Node),
	}
}

func (m *Map) MapIndex(key string) Node {
	return m.value[key]
}

func (m *Map) MapIndex2(key string) (Node, bool) {
	value, ok := m.value[key]
	return value, ok
}

func (m *Map) SetMapIndex(key string, value Node) {
	if prev := m.value[key]; prev != nil {
		prev.resetParent()
	}
	if value != nil {
		value.setParent(m)
	}
	m.value[key] = value
	m.incrementVersion()
}

func (m *Map) Delete(key string) {
	if prev := m.value[key]; prev != nil {
		prev.resetParent()
	}
	delete(m.value, key)
	m.incrementVersion()
}

type Leaf struct {
	*base
	value constant.Value
}

func MakeLeaf(value constant.Value) *Leaf {
	return &Leaf{
		base:  &base{version: 1},
		value: value,
	}
}

func (l *Leaf) Value() constant.Value {
	return l.value
}

func (l *Leaf) SetValue(value constant.Value) {
	l.value = value
	l.incrementVersion()
}

func MarshalJSON(tree Node) ([]byte, error) {
	panic("not implemented")
}

func UnmarshalJSON(b []byte) (Node, error) {
	panic("not implemented")
}
