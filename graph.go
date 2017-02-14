/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package graph

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Graph ...
type Graph struct {
	Components        []Component `json:"components"`
	ChangedComponents []Component `json:"changes"`
	Edges             []Edge      `json:"edges"`
}

// New returns a new graph
func New() *Graph {
	return &Graph{
		Components: make([]Component, 0),
		Edges:      make([]Edge, 0),
	}
}

// Component returns a component given the name matches
func (g *Graph) Component(component string) Component {
	for i, v := range g.Components {
		if v.GetID() == component {
			return g.Components[i]
		}
	}
	return nil
}

// HasComponent finds if the specified component exists
func (g *Graph) HasComponent(componentID string) bool {
	for _, v := range g.Components {
		if v.GetID() == componentID {
			return true
		}
	}
	return false
}

// AddComponent adds a component to the graphs vertices if it does not already exist
func (g *Graph) AddComponent(component Component) error {
	if g.HasComponent(component.GetID()) {
		return errors.New("Component already exists: " + component.GetID())
	}
	g.Components = append(g.Components, component)

	return nil
}

// UpdateComponent updates the graph
func (g *Graph) UpdateComponent(component Component) {
	for i := 0; i < len(g.Components); i++ {
		if g.Components[i].GetID() == component.GetID() {
			g.Components[i] = component
			return
		}
	}
}

// DisconnectComponent removes a component from the graph. It will connect any neighbour/origin components together
func (g *Graph) DisconnectComponent(name string) error {
	origins := g.Origins(name)

	for i := len(g.Edges) - 1; i >= 0; i-- {
		// Remove any edges that connect to the disconnected component
		if g.Edges[i].Destination == name {
			g.Edges = append(g.Edges[:i], g.Edges[i+1:]...)
		}

		// Remove any neighbouring connections and reconnect them to origins
		if g.Edges[i].Source == name {
			for _, ov := range *origins {
				err := g.Connect(ov.GetID(), g.Edges[i].Destination)
				if err != nil {
					return err
				}
			}
			g.Edges = append(g.Edges[:i], g.Edges[i+1:]...)
		}
	}

	return nil
}

// Connect adds a dependency between two vertices
func (g *Graph) Connect(source, destination string) error {
	if !g.HasComponent(source) || !g.HasComponent(destination) {
		return errors.New("Could not connect Component, does not exist")
	}

	g.Edges = append(g.Edges, Edge{Source: source, Destination: destination, Length: 1})

	return nil
}

// ConnectMutually connects two vertices to eachother
func (g *Graph) ConnectMutually(source, destination string) error {
	err := g.Connect(source, destination)
	if err != nil {
		return err
	}
	return g.Connect(destination, source)
}

// ComponentByProviderID : returns a single component by matching provider id
func (g *Graph) ComponentByProviderID(id string) Component {
	for _, component := range g.Components {
		if component.GetProviderID() == id {
			return component
		}
	}

	return nil
}

// Neighbours returns all depencencies of a component
func (g *Graph) Neighbours(component string) *Neighbours {
	var n Neighbours

	for _, edge := range g.Edges {
		if edge.Source == component {
			n = append(n, g.Component(edge.Destination))
		}
	}

	return n.Unique()
}

// Origins returns all source components of a component
func (g *Graph) Origins(component string) *Neighbours {
	var n Neighbours

	for _, edge := range g.Edges {
		if edge.Destination == component {
			n = append(n, g.Component(edge.Source))
		}
	}

	return n.Unique()
}

// LengthBetween returns the length between two edges
func (g *Graph) LengthBetween(source, destination string) int {
	for _, e := range g.Edges {
		if e.Source == source && e.Destination == destination {
			return e.Length
		}
	}
	return -1
}

// ToJSON serialises the graph as json
func (g *Graph) ToJSON() ([]byte, error) {
	return json.Marshal(g)
}

// FromJSON loads the graph from json
func (g *Graph) FromJSON(data []byte) error {
	return json.Unmarshal(data, g)
}

// LoadEdges loads a graphs edges
func (g *Graph) LoadEdges(edges []Edge) {
	g.Edges = edges
}

// Graphviz outputs the graph in graphviz format
func (g *Graph) Graphviz() string {
	var output []string

	output = append(output, "digraph G {")

	for _, edge := range g.Edges {
		output = append(output, fmt.Sprintf("  \"%s\" -> \"%s\"", edge.Source, edge.Destination))
	}

	output = append(output, "}")

	return strings.Join(output, "\n")
}

// Diff two graphs
func (g *Graph) Diff(og *Graph) (*Graph, error) {
	//for _, ov := g.Components {

	//}

	return nil, nil
}

//func (g *Graph) DepthFirstSearch()
