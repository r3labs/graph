/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package graph

// Neighbours represents a collection of dependent vertices
type Neighbours []Component

// Exists checks if the group contains the component
func (n *Neighbours) Exists(component string) bool {
	for _, v := range *n {
		if v.GetID() == component {
			return true
		}
	}
	return false
}

// Unique returns a new collection of Vertices that are unique
func (n *Neighbours) Unique() *Neighbours {
	var un Neighbours

	for _, v := range *n {
		if v == nil {
			continue
		}
		if !un.Exists(v.GetID()) {
			un = append(un, v)
		}
	}

	return &un
}

// GetComponentGroup returns the first item in a components group
func (n *Neighbours) GetComponentGroup(group string) Component {
	for _, v := range *n {
		if v.GetGroup() != "" && v.GetGroup() == group {
			return v
		}
	}

	return nil
}

// GetSequentialDependency retuens the sequentially dependent component based on
// a given id if given
func (n *Neighbours) GetSequentialDependency(id string) Component {
	if id == "" {
		return nil
	}

	for _, v := range *n {
		for _, d := range v.SequentialDependencies() {
			if d == id {
				return v
			}
		}
	}

	return nil
}
