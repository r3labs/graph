/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package graph

import "errors"

// Graph ...
type Graph struct {
	Vertices []Vertex
	Edges    []Edge `json:"arcs"`
}

// New returns a new graph
func New() *Graph {
	return &Graph{
		Vertices: make([]Vertex, 0),
		Edges:    make([]Edge, 0),
	}
}

// AddVertex adds a vertex to the graphs vertices if it does not already exist
func (g *Graph) AddVertex(vertex Vertex) {
	for _, v := range g.Vertices {
		if v.Name() == vertex.Name() {
			return
		}
	}
	g.Vertices = append(g.Vertices, vertex)
}

// HasVertex finds if the specified vertex exists
func (g *Graph) HasVertex(name string) bool {
	for _, v := range g.Vertices {
		if v.Name() == name {
			return true
		}
	}
	return false
}

// Vertex retuns the matching vertex by name
func (g *Graph) Vertex(name string) Vertex {
	for _, v := range g.Vertices {
		if v.Name() == name {
			return v
		}
	}
	return nil
}

// Connect adds a dependency between two vertices
func (g *Graph) Connect(source, destination, event string) error {
	if !g.HasVertex(source) || !g.HasVertex(destination) {
		return errors.New("Could not connect vertex, does not exist")
	}

	g.Edges = append(g.Edges, Edge{Source: source, Destination: destination, Event: event, Length: 1})

	return nil
}

// ConnectMutually connects two vertices to eachother
func (g *Graph) ConnectMutually(source, destination, event string) error {
	err := g.Connect(source, destination, event)
	if err != nil {
		return err
	}
	return g.Connect(destination, source, event)
}

// DisconnectVertex removes a vertex from the graph. It will connect any neighbour/origin verticies together
func (g *Graph) DisconnectVertex(name string) error {
	origins := g.Origins(name)

	for i := len(g.Edges) - 1; i >= 0; i-- {
		// Remove any edges that connect to the disconnected vertex
		if g.Edges[i].Destination == name {
			g.Edges = append(g.Edges[:i], g.Edges[i+1:]...)
		}

		// Remove any neighbouring connections and reconnect them to origins
		if g.Edges[i].Source == name {
			for _, ov := range *origins {
				err := g.Connect(ov.Name(), g.Edges[i].Destination, g.Edges[i].Event)
				if err != nil {
					return err
				}
			}
			g.Edges = append(g.Edges[:i], g.Edges[i+1:]...)
		}
	}

	return nil
}

// Neighbours returns all depencencies of a vertex
func (g *Graph) Neighbours(vertex string) *Neighbours {
	n := Neighbours{}

	for _, edge := range g.Edges {
		if edge.Source == vertex {
			n = append(n, g.Vertex(edge.Destination))
		}
	}

	return n.Unique()
}

// Origins returns all source verticies of a vertex
func (g *Graph) Origins(vertex string) *Neighbours {
	n := Neighbours{}

	for _, edge := range g.Edges {
		if edge.Destination == vertex {
			n = append(n, g.Vertex(edge.Source))
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
