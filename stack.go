/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package graph

// NodeStack stores a collection of verticies
type NodeStack []Component

// Append a verticies onto the stack
func (n NodeStack) Append(i []Component) {
	n = append(n, i...)
}

// Prepend a verticies onto the stack
func (n NodeStack) Prepend(i []Component) {
	n = append(i, n...)
}

// Pop a component from the stack
func (n NodeStack) Pop() Component {
	var x Component
	x, n = n[len(n)-1], n[:len(n)-1]
	return x
}

// Empty returns true if there are no more verticies left
func (n NodeStack) Empty() bool {
	return len(n) < 1
}
