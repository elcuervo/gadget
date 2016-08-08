package main

import (
	"testing"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/network"
	"github.com/stretchr/testify/assert"
)

func fakeFinder() func(string) (types.ContainerJSON, error) {
	return func(id string) (types.ContainerJSON, error) {
		return fakeContainer(id), nil
	}

}

func fakeContainer(id string) types.ContainerJSON {
	networks := make(map[string]*network.EndpointSettings)

	return types.ContainerJSON{
		Config: &container.Config{
			Image: "image",
		},

		NetworkSettings: &types.NetworkSettings{
			Networks: networks,
		},

		ContainerJSONBase: &types.ContainerJSONBase{
			ID:    id,
			Image: "image",
			State: &types.ContainerState{Status: "running"},
			Name:  "name",
		},
	}

}

func TestContainerLookup(t *testing.T) {
	lookup := &ContainerLookup{Finder: fakeFinder()}
	name := "68487e1d320a"
	c, _ := lookup.FindContainer(name)

	assert.Equal(t, name, c.Info.ID)
}
