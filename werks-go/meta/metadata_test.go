package meta_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/egustafson/werks/werks-go/meta"
)

func TestNewMetadata_empty(t *testing.T) {

	empty_md := meta.NewMetadata()
	assert.Zero(t, empty_md.Length())
	assert.Zero(t, empty_md.Count("any-key"))
	assert.False(t, empty_md.Has("any-key"))
	assert.False(t, empty_md.HasExact("any-key", "any-value"))
	assert.Zero(t, len(empty_md.All()))
	values, ok := empty_md.Find("any-key")
	assert.Zero(t, len(values))
	assert.False(t, ok)
	empty_md.RemoveAll("any-key")
	empty_md.RemoveExact("any-key", "any-value")
	assert.Zero(t, empty_md.Length())
}

func TestMetadata_set_remove(t *testing.T) {
	md := meta.NewMetadata()

	md.Set("k", "value")
	assert.Equal(t, 1, md.Length())
	assert.Equal(t, 1, md.Count("k"))
	assert.True(t, md.Has("k"))
	assert.True(t, md.HasExact("k", "value"))
	assert.Equal(t, 1, len(md.All()))

	values, ok := md.Find("k")
	assert.True(t, ok)
	assert.Equal(t, 1, len(values))
	assert.Equal(t, "value", values[0])

	md.Append("k", "value2")
	assert.Equal(t, 2, md.Length())
	assert.Equal(t, 2, md.Count("k"))
	assert.True(t, md.HasExact("k", "value"))
	assert.True(t, md.HasExact("k", "value2"))
	assert.Equal(t, 2, len(md.All()))

	values, ok = md.Find("k")
	assert.True(t, ok)
	assert.Equal(t, 2, len(values))
	assert.Equal(t, "value", values[0])
	assert.Equal(t, "value2", values[1])

	md.RemoveExact("k", "value")
	assert.Equal(t, 1, md.Length())
	assert.Equal(t, 1, md.Count("k"))
	assert.True(t, md.Has("k"))
	assert.False(t, md.HasExact("k", "value"))
	assert.True(t, md.HasExact("k", "value2"))
	assert.Equal(t, 1, len(md.All()))

	values, ok = md.Find("k")
	assert.True(t, ok)
	assert.Equal(t, 1, len(values))
	assert.Equal(t, "value2", values[0])

	md.Append("k2", "v2")
	assert.True(t, md.HasExact("k2", "v2"))

	md.RemoveAll("k")
	md.RemoveAll("k2")
	assert.Zero(t, md.Length())
}

func TestMetadata_WithMeta(t *testing.T) {

	bootstrap := meta.Meta{K: "boot-key", V: "boot-value"}
	md := meta.NewMetadata(meta.WithMeta(bootstrap))

	assert.Equal(t, 1, md.Length())
	assert.Equal(t, 1, md.Count("boot-key"))
	assert.True(t, md.HasExact("boot-key", "boot-value"))
}

func TestMetadata_WithMetaList(t *testing.T) {

	bootstrap := []meta.Meta{
		{K: "boot-k1", V: "boot-v1"},
		{K: "boot-k2", V: "boot-v2"},
	}
	md := meta.NewMetadata(meta.WithMetaList(bootstrap))

	assert.Equal(t, 2, md.Length())
	assert.True(t, md.HasExact("boot-k1", "boot-v1"))
	assert.True(t, md.HasExact("boot-k2", "boot-v2"))
}

func TestMetadata_WithMetadata(t *testing.T) {

	bootstrap := meta.NewMetadata(meta.WithMeta(meta.Meta{"boot-key", "boot-value"}))
	md := meta.NewMetadata(meta.WithMetadata(bootstrap))

	assert.Equal(t, 1, md.Length())
	assert.Equal(t, 1, md.Count("boot-key"))
	assert.True(t, md.HasExact("boot-key", "boot-value"))
}
