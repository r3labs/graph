/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package graph

// ComponentGroup holds a collection of components that can be filtered
type ComponentGroup []Component

// ByProviderID ...
func (cg ComponentGroup) ByProviderID(id string) Component {
	for _, component := range cg {
		if component.GetProviderID() == id {
			return component
		}
	}

	return nil
}

// ByType ...
func (cg ComponentGroup) ByType(ctype string) ComponentGroup {
	var c []Component

	for _, component := range cg {
		if component.GetType() == ctype {
			c = append(c, component)
		}
	}

	return c
}

// ByGroup ...
func (cg ComponentGroup) ByGroup(tag, group string) ComponentGroup {
	var c []Component

	for _, component := range cg {
		if component.GetTag(tag) == group {
			c = append(c, component)
		}
	}

	return c
}

// ByName ...
func (cg ComponentGroup) ByName(name string) ComponentGroup {
	var c []Component

	for _, component := range cg {
		if component.GetName() == name {
			c = append(c, component)
		}
	}

	return c
}

// NameValues ...
func (cg ComponentGroup) NameValues() []string {
	var names []string

	for _, component := range cg {
		names = append(names, component.GetName())
	}

	return names
}

// TagValues ...
func (cg ComponentGroup) TagValues(tag string) []string {
	var tv []string

	for _, component := range cg {
		tv = appendUnique(tv, component.GetTag(tag))
	}

	return tv
}

func appendUnique(s []string, item string) []string {
	for _, i := range s {
		if i == item {
			return s
		}
	}

	return append(s, item)
}
