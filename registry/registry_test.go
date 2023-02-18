package registry

import "testing"

func TestRegistry(t *testing.T) {
	registry := GetConnectionRegistry()
	connection := CreateRegistryItem("test", newMockConn([]byte("test")))
	err := registry.Add(connection)
	if err != nil {
		t.Errorf("Error adding connection to registry: %s", err)
	}
	_, err = registry.Get("test")
	if err != nil {
		t.Errorf("Error getting connection from registry: %s", err)
	}
	err = registry.Del(connection)
	if err != nil {
		t.Errorf("Error deleting connection from registry: %s", err)
	}
	_, err = registry.Get("test")
	if err == nil {
		t.Errorf("Error getting connection from registry: %s", err)
	}
}
