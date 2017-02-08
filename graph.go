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
	Vertices []Vertex `json:"verticies"`
	Edges    []Edge   `json:"edges"`
}

// New returns a new graph
func New() *Graph {
	return &Graph{
		Vertices: make([]Vertex, 0),
		Edges:    make([]Edge, 0),
	}
}

// AddVertex adds a vertex to the graphs vertices if it does not already exist
func (g *Graph) AddVertex(vertex Vertex) error {
	if g.HasVertex(vertex.GetID()) {
		return errors.New("Vertex already exists")
	}
	g.Vertices = append(g.Vertices, vertex)

	return nil
}

// HasVertex finds if the specified vertex exists
func (g *Graph) HasVertex(vertex string) bool {
	for _, v := range g.Vertices {
		if v.GetID() == vertex {
			return true
		}
	}
	return false
}

// Vertex returns a vertex given the name matches
func (g *Graph) Vertex(vertex string) Vertex {
	for i, v := range g.Vertices {
		if v.GetID() == vertex {
			return g.Vertices[i]
		}
	}
	return nil
}

// UpdateVertex updates the graph
func (g *Graph) UpdateVertex(vertex Vertex) {
	for i := 0; i < len(g.Vertices); i++ {
		if g.Vertices[i].GetID() == vertex.GetID() {
			g.Vertices[i] = vertex
			return
		}
	}
}

// Connect adds a dependency between two vertices
func (g *Graph) Connect(source, destination string) error {
	if !g.HasVertex(source) || !g.HasVertex(destination) {
		return errors.New("Could not connect vertex, does not exist")
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

// RemoveVertex removes a vertex from the graph. It will connect any neighbour/origin verticies together
func (g *Graph) RemoveVertex(name string) error {
	origins := g.Origins(name)

	for i := len(g.Edges) - 1; i >= 0; i-- {
		// Remove any edges that connect to the disconnected vertex
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

// Neighbours returns all depencencies of a vertex
func (g *Graph) Neighbours(vertex string) *Neighbours {
	var n Neighbours

	for _, edge := range g.Edges {
		if edge.Source == vertex {
			n = append(n, g.Vertex(edge.Destination))
		}
	}

	return n.Unique()
}

// Origins returns all source verticies of a vertex
func (g *Graph) Origins(vertex string) *Neighbours {
	var n Neighbours

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
func (g *Graph) Graphviz(start Vertex) string {
	var output []string

	output = append(output, "digraph G {")

	for _, edge := range g.Edges {
		output = append(output, fmt.Sprintf("  \"%s\" -> \"%s\"", edge.Source, edge.Destination))
	}

	output = append(output, "}")

	return strings.Join(output, "\n")
}

// Diff two graphs
func (g *Graph) Diff(og *Graph) error {
	//for _, ov := g.Vertices {

	//}

	return nil
}

//func (g *Graph) DepthFirstSearch()
