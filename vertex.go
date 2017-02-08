/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package graph

// Vertex : representation of a component
type Vertex interface {
	GetID() string           // returns the ID of the component
	GetProvider() string     // returns the components provider
	GetType() string         // returns the type of component
	GetState() string        // returns the state of the component. i.e. waiting, running, completed, errored
	SetState(string)         // sets the state of the component
	GetAction() string       // returns the action of the component, i.e. create, update, delete, get
	SetAction(string)        // sets the action of the component
	GetGroup() string        // returns the components group name
	Diff(interface{})        // should return changelog
	Update(interface{}) bool // updates the values stored on the component
	Rebuild(interface{})     // rebuilds the internal state of the component, a component set is passed in
	IsStateful() bool        // returns if the component is stateful. This is important to work out if a component can be skipped when deleting its dependencies (pruning).
}
