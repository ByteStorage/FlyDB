package cmd

import (
	"github.com/desertbit/grumble"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestServer_StartServer(t *testing.T) {
	c := grumble.Context{Args: map[string]*grumble.ArgMapItem{"name": {"server1", false}}}
	err := startServer(&c)
	assert.Nil(t, err)
	assert.Nil(t, db)

	c = grumble.Context{Args: nil}
	err = startServer(&c)
	assert.Nil(t, err)
	assert.NotNil(t, db)
}

func TestServer_StopServer(t *testing.T) {
	c := grumble.Context{Args: map[string]*grumble.ArgMapItem{"name": {"server1", false}}}
	err := stopServer(&c)
	assert.Nil(t, err)
	assert.Nil(t, db)

	c = grumble.Context{Args: nil}
	err = stopServer(&c)
	assert.NotNil(t, err)
	assert.Nil(t, db)

	err = startServer(&c)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = stopServer(&c)
	assert.Nil(t, err)
	assert.Nil(t, db)

}

func TestServer_CleanServer(t *testing.T) {
	c := grumble.Context{Args: map[string]*grumble.ArgMapItem{"name": {"server1", false}}}
	err := cleanServer(&c)
	assert.Nil(t, err)
	assert.Nil(t, db)

	c = grumble.Context{Args: nil}
	err = cleanServer(&c)
	assert.NotNil(t, err)
	assert.Nil(t, db)

	c = grumble.Context{Args: nil}
	err = startServer(&c)
	assert.Nil(t, err)
	assert.NotNil(t, db)
	err = db.Put([]byte("key"), []byte("value"))
	assert.NotEmpty(t, db.GetListKeys())

	err = cleanServer(&c)
	assert.Nil(t, err)
	assert.Empty(t, db.GetListKeys())

}
