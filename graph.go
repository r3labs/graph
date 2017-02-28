/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package graph

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// Graph ...
type Graph struct {
	ID         string      `json:"id"`
	Action     string      `json:"action"`
	Components []Component `json:"components"`
	Changes    []Component `json:"changes,omitempty"`
	Edges      []Edge      `json:"edges,omitempty"`
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

// ComponentAll returns a component from either changes or components given the name matches
func (g *Graph) ComponentAll(component string) Component {
	for i, v := range g.Changes {
		if v.GetID() == component {
			return g.Changes[i]
		}
	}
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

// DeleteComponent deletes a component from the graph
func (g *Graph) DeleteComponent(component Component) {
	for i := len(g.Components) - 1; i >= 0; i-- {
		if g.Components[i].GetID() == component.GetID() {
			g.Components = append(g.Components[:i], g.Components[i+1:]...)
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

// connect is the internal method for connecting two verticies, it provides less checks than publicly exposed methods
func (g *Graph) connect(source, destination string) {
	if g.Connected(source, destination) != true {
		g.Edges = append(g.Edges, Edge{Source: source, Destination: destination, Length: 1})
	}
}

// Connect adds a dependency between two vertices
func (g *Graph) Connect(source, destination string) error {
	if !g.HasComponent(source) || !g.HasComponent(destination) {
		return errors.New("Could not connect Component, does not exist")
	}

	g.connect(source, destination)

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

// ConnectSequential adds a dependency between two vertices. If the source has more than 1 neighbouring vertex, the destination vertex will be connected to that.
func (g *Graph) ConnectSequential(source, destination string) error {
	if !g.HasComponent(source) {
		source = "start"
	}

	if !g.HasComponent(destination) {
		return errors.New("Could not connect Component, does not exist")
	}

	c := g.Component(destination)
	gc := g.Neighbours(source).GetComponentGroup(c.GetGroup())

	for gc != nil {
		source = gc.GetID()
		gc = g.Neighbours(source).GetComponentGroup(c.GetGroup())
	}

	g.connect(source, destination)

	return nil
}

// Connected returns true if two components are connected
func (g *Graph) Connected(source, destination string) bool {
	for _, edge := range g.Edges {
		if edge.Source == source && edge.Destination == destination {
			return true
		}
	}

	return false
}

// GetComponents returns a component group that can be filtered
func (g *Graph) GetComponents() ComponentGroup {
	return g.Components
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

// Diff : diff two graphs, new, modified or removed components will be moved to Changes, and components will be
func (g *Graph) Diff(og *Graph) (*Graph, error) {
	// new temporary graph
	ng := New()

	for _, c := range g.Components {
		if c.GetAction() == "none" {
			continue
		}

		oc := og.Component(c.GetID())
		if oc != nil {
			if c.Diff(oc) {
				c.SetAction("update")
				c.SetState("waiting")
				ng.AddComponent(c)
			}
		} else {
			if c.GetAction() != "find" {
				c.SetAction("create")
			}
			c.SetState("waiting")
			ng.AddComponent(c)
		}
	}

	for _, oc := range og.Components {
		c := g.Component(oc.GetID())
		if c == nil {
			oc.SetAction("delete")
			oc.SetState("waiting")
			ng.AddComponent(oc)
		}
	}

	ng.SetDiffDependencies()

	// Set changes + edges on original graph
	og.Changes = ng.Components
	og.Edges = ng.Edges

	return og, nil
}

// Graphviz outputs the graph in graphviz format
func (g *Graph) Graphviz() string {
	var output []string

	output = append(output, "digraph G {")

	for _, edge := range g.Edges {
		dest := g.ComponentAll(edge.Destination)
		if dest != nil {
			output = append(output, fmt.Sprintf("  \"%s\" -> \"%s\" [label=\"%s\"]", edge.Source, edge.Destination, dest.GetAction()))
		} else {
			output = append(output, fmt.Sprintf("  \"%s\" -> \"%s\"", edge.Source, edge.Destination))
		}
	}

	output = append(output, "}")

	return strings.Join(output, "\n")
}

// SetDiffDependencies rebuilds the graph's dependencies when diffing
func (g *Graph) SetDiffDependencies() {
	// This needs improvement

	g.Edges = make([]Edge, 0)

	for _, c := range g.Components {
		for _, dep := range c.Dependencies() {
			switch c.GetAction() {
			case "delete":
				if c.IsStateful() {
					g.Connect(c.GetID(), dep)
				}
			case "update":
				g.ConnectSequential(dep, c.GetID())
			case "create", "find":
				g.Connect(dep, c.GetID())
			}
		}
	}

	g.SetStartFinish()
}

// SetStartFinish sets a start and finish point
func (g *Graph) SetStartFinish() {
	for _, c := range g.Components {
		o := g.Origins(c.GetID())
		n := g.Neighbours(c.GetID())

		if len(*o) < 1 {
			g.connect("start", c.GetID())
		}

		if len(*n) < 1 {
			g.connect(c.GetID(), "end")
		}
	}
}

// ToJSON serialises the graph as json
func (g *Graph) ToJSON() ([]byte, error) {
	return json.Marshal(g)
}

// Load loads a graph from json
func (g *Graph) Load(gg map[string]interface{}) error {
	components := gg["components"].([]interface{})
	changes := gg["changes"].([]interface{})

	for i := 0; i < len(components); i++ {
		c := components[i].(map[string]interface{})
		components[i] = MapGenericComponent(c)
	}

	for i := 0; i < len(changes); i++ {
		c := changes[i].(map[string]interface{})
		changes[i] = MapGenericComponent(c)
	}

	return mapstructure.Decode(gg, g)
}

// func (g *Graph) DepthFirstSearch()

// func (g *Graph) CycleDetection()
