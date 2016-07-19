/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package graph

// Edge ...
type Edge struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Event       string `json:"event"`
	Length      int    `json:"-"`
}
