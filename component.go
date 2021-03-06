/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package graph

import "github.com/r3labs/diff"

// Component : representation of a component
type Component interface {
	GetID() string                          // returns the ID of the component
	GetName() string                        // returns the name of the component
	GetProvider() string                    // returns the components provider
	GetProviderID() string                  // returns the components provider specific ID
	GetType() string                        // returns the type of component
	GetState() string                       // returns the state of the component. i.e. waiting, running, completed, errored
	SetState(string)                        // sets the state of the component
	GetAction() string                      // returns the action of the component, i.e. create, update, delete, get
	SetAction(string)                       // sets the action of the component
	GetGroup() string                       // returns the components group name
	GetTags() map[string]string             // returns the tags associated with the component
	GetTag(string) string                   // returns the tag associated with the component
	Diff(Component) (diff.Changelog, error) // should return changelog
	Update(Component)                       // updates the values stored on the component
	Rebuild(*Graph)                         // rebuilds the internal state of the component, a component set is passed in
	Validate() error                        // validates the component's values
	Dependencies() []string                 // returns a collection of parent component id's
	SetDefaultVariables()                   // sets the default variables for a component
	IsStateful() bool                       // returns if the component is stateful. This is important to work out if a component can be skipped when deleting its dependencies (pruning).
	SequentialDependencies() []string       // returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
}
