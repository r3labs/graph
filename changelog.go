/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package graph

import "github.com/r3labs/diff"

func prefixChanges(prefix string, cl diff.Changelog) diff.Changelog {
	for i := 0; i < len(cl); i++ {
		cl[i].Path = append([]string{prefix}, cl[i].Path...)
	}

	return cl
}
