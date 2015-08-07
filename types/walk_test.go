package types

import (
	"bytes"
	"testing"

	"github.com/attic-labs/noms/chunks"
	. "github.com/attic-labs/noms/dbg"
	"github.com/stretchr/testify/assert"
)

func TestWalkAll(t *testing.T) {
	assert := assert.New(t)
	cs := &chunks.MemoryStore{}

	write := func(v Value) Value {
		_, err := WriteValue(v, cs)
		Chk.NoError(err)
		return v
	}

	b := write(Bool(true))
	i := write(Int32(42))
	f := write(Float32(88.8))
	s := write(NewString("hi"))
	blv, _ := NewBlob(bytes.NewBuffer([]byte("hi")))
	bl := write(blv)
	l := write(NewList(b, i, f, s, bl))
	m := write(NewMap(b, i, f, s))
	se := write(NewSet(b, i, f, s, bl))
	l2 := write(NewList(l))

	tests := []struct {
		v        Value
		expected Set
	}{
		{b, NewSet(b)},
		{i, NewSet(i)},
		{f, NewSet(f)},
		{s, NewSet(s)},
		{bl, NewSet(bl)},
		{l, NewSet(l, b, i, f, s, bl)},
		{m, NewSet(m, b, i, f, s)},
		{se, NewSet(se, b, i, f, s, bl)},
		{l2, NewSet(l2, l, b, i, f, s, bl)},
	}

	for _, t := range tests {
		expected := t.expected
		All(t.v.Ref(), cs, func(f Future) {
			v, err := f.Deref(cs)
			Chk.NoError(err)
			assert.True(expected.Has(v))
			expected = expected.Remove(v)
		})
		assert.True(expected.Empty())
	}
}