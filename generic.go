/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package graph

// GenericComponent is a representation of a component backed by a map[string]interface{}
type GenericComponent map[string]interface{}

// GetID : returns the component's ID
func (gc *GenericComponent) GetID() string {
	return (*gc)["_component_id"].(string)
}

// GetName returns a components name
func (gc *GenericComponent) GetName() string {
	return ""
}

// GetProvider : returns the provider type
func (gc *GenericComponent) GetProvider() string {
	return (*gc)["_provider"].(string)
}

// GetProviderID returns a components provider id
func (gc *GenericComponent) GetProviderID() string {
	return ""
}

// GetType : returns the type of the component
func (gc *GenericComponent) GetType() string {
	return (*gc)["_component"].(string)
}

// GetState : returns the state of the component
func (gc *GenericComponent) GetState() string {
	return (*gc)["_state"].(string)
}

// SetState : sets the state of the component
func (gc *GenericComponent) SetState(state string) {
	(*gc)["_state"] = state
}

// GetAction : returns the action of the component
func (gc *GenericComponent) GetAction() string {
	return (*gc)["_action"].(string)
}

// SetAction : Sets the action of the component
func (gc *GenericComponent) SetAction(action string) {
	(*gc)["_action"] = action
}

// GetGroup : returns the components group
func (gc *GenericComponent) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (gc *GenericComponent) GetTags() map[string]string {
	return map[string]string{}
}

// GetTag returns a components tag
func (gc *GenericComponent) GetTag(string) string {
	return ""
}

// Dependencies : returns a list of component id's upon which the component depends
func (gc *GenericComponent) Dependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (gc *GenericComponent) Validate() error {
	return nil
}

// Diff : diff's the component against another component of the same type
func (gc *GenericComponent) Diff(v Component) bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (gc *GenericComponent) SetDefaultVariables() {}

// Rebuild : rebuilds the component's internal state, such as templated values
func (gc *GenericComponent) Rebuild(g *Graph) {}

// Update : updates the provider returned values of a component
func (gc *GenericComponent) Update(v Component) {}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (gc *GenericComponent) IsStateful() bool {
	return true
}
