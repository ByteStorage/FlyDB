package store

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/datastore"
	"github.com/hashicorp/raft"
)

// DataStoreFactory is a function type that creates a new instance of a raft.LogStore.
// It takes a configuration map as input and returns the created LogStore or an error if the creation fails.
type DataStoreFactory func(conf config.Config) (raft.LogStore, error)

// datastoreFactories is a map that associates a string key (name) with a DataStoreFactory.
// It will be used to store and retrieve different DataStoreFactory implementations.
var datastoreFactories = make(map[string]DataStoreFactory)

// Init initializes the datastoreFactories by registering different DataStoreFactory implementations.
// It calls the Register function to associate each implementation with a unique name.
// Currently, "memory" and "bolt" are registered as the names for the corresponding factories.
// The function returns an error if any registration fails, but in this implementation, the error is ignored.
func Init() error {
	// Register the "memory" DataStoreFactory implementation with the name "memory"
	_ = Register("memory", datastore.NewLogInMemStorage)
	// Register the "bolt" DataStoreFactory implementation with the name "boltdb"
	_ = Register("boltdb", datastore.NewLogBoltDbStorage)
	// Register the "flydb" DataStoreFactory implementation with the name "flydb"
	_ = Register("flydb", datastore.NewLogFlyDbStorage)
	// Register the "rocksdb" DataStoreFactory implementation with the name "rocksdb"
	_ = Register("rocksdb", datastore.NewLogRocksDbStorage)

	return nil
}

// newRaftLog is a function that returns a new instance of raft.LogStore.
// It initializes the datastoreFactories by calling the Init function.
// Then, it retrieves the "memory" DataStoreFactory from the datastoreFactories map using the "memory" string key.
// Finally, it creates a new LogStore using the retrieved factory and an empty configuration map.
// The created LogStore and an error (if any) are returned.
func newRaftLog(conf config.Config) (raft.LogStore, error) {
	_ = Init()

	// Get the "memory" DataStoreFactory from the map
	return getDataStore(conf)
}

// Register is a function that registers a DataStoreFactory implementation with a given name.
// It takes the name string and the factory function as input.
// The function checks if the factory is nil and returns an error if it is.
// Then, it checks if the name is already registered in the datastoreFactories map.
// If the name is already registered, it returns an error.
// Otherwise, it adds the factory to the map with the name as the key and returns nil.
func Register(name string, factory DataStoreFactory) error {
	if factory == nil {
		return fmt.Errorf("datastore factory %s does not exist", name)
	}

	// Check if the name is already registered
	_, registered := datastoreFactories[name]
	if registered {
		return fmt.Errorf(`datastore factory %s already registered`, name)
	}
	// Add the factory to the datastoreFactories map
	datastoreFactories[name] = factory
	return nil
}

// getDataStore is a function that retrieves a LogStore implementation from the datastoreFactories map.
// It takes the name of the datastore and a configuration map as input.
// The function first checks if the requested datastore exists in the map.
// If the datastore is not found, it returns an error.
// Otherwise, it retrieves the corresponding DataStoreFactory from the map.
// Finally, it creates a new LogStore using the factory and the configuration map and returns it
// along with an error (if any).
func getDataStore(datastore config.Config) (raft.LogStore, error) {
	// Get the DataStoreFactory for the requested datastore
	dsFactory, ok := datastoreFactories[datastore.LogDataStorage]
	if !ok {
		return nil, fmt.Errorf("datastore not valid")
	}
	// Create a new LogStore using the DataStoreFactory and the configuration map
	return dsFactory(datastore)
}
