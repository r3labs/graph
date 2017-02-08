/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package graph

// Neighbours represents a collection of dependent vertices
type Neighbours []Vertex

// Exists checks if the group contains the vertex
func (n *Neighbours) Exists(vertex string) bool {
	for _, v := range *n {
		if v.GetID() == vertex {
			return true
		}
	}
	return false
}

// Unique returns a new collection of Vertices that are unique
func (n *Neighbours) Unique() *Neighbours {
	var un Neighbours

	for _, v := range *n {
		if !un.Exists(v.GetID()) {
			un = append(un, v)
		}
	}

	return &un
}
