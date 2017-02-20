/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package graph

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type testComponent struct {
	Name   string `json:"name"`
	State  string `json:"state"`
	Action string `json:"action"`
}

func (tv *testComponent) GetID() string {
	return tv.Name
}

func (tv *testComponent) GetName() string {
	return tv.Name
}

func (tv *testComponent) GetProvider() string {
	return "test"
}

func (tv *testComponent) GetProviderID() string {
	return "test"
}

func (tv *testComponent) GetType() string {
	return "test"
}

func (tv *testComponent) GetState() string {
	return tv.State
}

func (tv *testComponent) SetState(state string) {
	tv.State = state
}

func (tv *testComponent) GetAction() string {
	return tv.Action
}

func (tv *testComponent) SetAction(action string) {
	tv.Action = action
}

func (tv *testComponent) GetGroup() string {
	return "test"
}

func (tv *testComponent) GetTags() map[string]string {
	return map[string]string{}
}

func (tv *testComponent) GetTag(string) string {
	return ""
}

func (tv *testComponent) Dependencies() []string {
	return []string{}
}

func (tv *testComponent) Validate() error {
	return nil
}

func (tv *testComponent) Diff(v Component) bool {
	return true
}

func (tv *testComponent) SetDefaultVariables() {}

func (tv *testComponent) Rebuild(g *Graph) {}

func (tv *testComponent) Update(v Component) {}

func (tv *testComponent) IsStateful() bool {
	return true
}

func TestGraph(t *testing.T) {
	Convey("Given a new graph", t, func() {
		g := New()

		Convey("When adding a new component", func() {
			g.AddComponent(&testComponent{Name: "test"})
			Convey("It should be stored on the graph", func() {
				So(len(g.Components), ShouldEqual, 1)
				So(g.Components[0].GetID(), ShouldEqual, "test")
			})
		})

		Convey("When adding a duplicate component", func() {
			g.AddComponent(&testComponent{Name: "test"})
			Convey("It should not be stored on the graph", func() {
				So(len(g.Components), ShouldEqual, 1)
			})
		})

		Convey("When connecting two verticies", func() {
			g.AddComponent(&testComponent{Name: "test1"})
			g.AddComponent(&testComponent{Name: "test2"})
			err := g.Connect("test1", "test2")
			Convey("It should create an edge between the two verticies", func() {
				So(err, ShouldBeNil)
				So(len(g.Components), ShouldEqual, 2)
				So(len(g.Edges), ShouldEqual, 1)
				So(g.Edges[0].Source, ShouldEqual, "test1")
				So(g.Edges[0].Destination, ShouldEqual, "test2")
			})
		})

		Convey("When connecting two verticies mutually", func() {
			g.AddComponent(&testComponent{Name: "test1"})
			g.AddComponent(&testComponent{Name: "test2"})
			err := g.ConnectMutually("test1", "test2")
			Convey("It should create an edge between the two verticies", func() {
				So(err, ShouldBeNil)
				So(len(g.Components), ShouldEqual, 2)
				So(len(g.Edges), ShouldEqual, 2)
				So(g.Edges[0].Source, ShouldEqual, "test1")
				So(g.Edges[0].Destination, ShouldEqual, "test2")
				So(g.Edges[1].Source, ShouldEqual, "test2")
				So(g.Edges[1].Destination, ShouldEqual, "test1")
			})
		})

		Convey("When connecting two verticies that don't exist", func() {
			err := g.ConnectMutually("fake1", "fake2")
			Convey("It should error", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When getting the origin verticies of a component", func() {
			g.AddComponent(&testComponent{Name: "test1"})
			g.AddComponent(&testComponent{Name: "test2"})
			g.AddComponent(&testComponent{Name: "test3"})
			erra := g.Connect("test2", "test1")
			errb := g.Connect("test3", "test1")
			origins := g.Origins("test1")
			Convey("It should return the correct verticies", func() {
				So(erra, ShouldBeNil)
				So(errb, ShouldBeNil)
				So(len(*origins), ShouldEqual, 2)
				So((*origins)[0].GetID(), ShouldEqual, "test2")
				So((*origins)[1].GetID(), ShouldEqual, "test3")
			})
		})

		Convey("When getting the neighbouring verticies of a component", func() {
			g.AddComponent(&testComponent{Name: "test1"})
			g.AddComponent(&testComponent{Name: "test2"})
			g.AddComponent(&testComponent{Name: "test3"})
			erra := g.Connect("test1", "test2")
			errb := g.Connect("test1", "test3")
			neighbours := g.Neighbours("test1")
			Convey("It should return the correct verticies", func() {
				So(erra, ShouldBeNil)
				So(errb, ShouldBeNil)
				So(len(*neighbours), ShouldEqual, 2)
				So((*neighbours)[0].GetID(), ShouldEqual, "test2")
				So((*neighbours)[1].GetID(), ShouldEqual, "test3")
			})
		})

		Convey("When disconnecting a component from the graph", func() {
			g.AddComponent(&testComponent{Name: "test1"})
			g.AddComponent(&testComponent{Name: "test2"})
			g.AddComponent(&testComponent{Name: "test3"})
			g.AddComponent(&testComponent{Name: "test4"})
			g.AddComponent(&testComponent{Name: "test5"})
			erra := g.Connect("test1", "test2")
			errb := g.Connect("test2", "test3")
			errc := g.Connect("test2", "test4")
			errd := g.Connect("test4", "test5")

			err := g.DisconnectComponent("test2")
			Convey("It should disconnect matching edges and reconnect them", func() {
				So(erra, ShouldBeNil)
				So(errb, ShouldBeNil)
				So(errc, ShouldBeNil)
				So(errd, ShouldBeNil)
				So(err, ShouldBeNil)
				So(len(g.Edges), ShouldEqual, 3)
				So(g.Edges[0].Source, ShouldEqual, "test4")
				So(g.Edges[0].Destination, ShouldEqual, "test5")
				So(g.Edges[1].Source, ShouldEqual, "test1")
				So(g.Edges[1].Destination, ShouldEqual, "test4")
				So(g.Edges[2].Source, ShouldEqual, "test1")
				So(g.Edges[2].Destination, ShouldEqual, "test3")
			})
		})

		Convey("When getting a component by name", func() {
			g.AddComponent(&testComponent{Name: "test1"})
			g.AddComponent(&testComponent{Name: "test2"})
			component := g.Component("test2")
			Convey("It should return the correct component", func() {
				So(component, ShouldNotBeNil)
				So(component.GetID(), ShouldEqual, "test2")
			})
		})

		Convey("When testing for an existing component", func() {
			g.AddComponent(&testComponent{Name: "test1"})
			exists := g.HasComponent("test1")
			Convey("It should return true", func() {
				So(exists, ShouldBeTrue)
			})
		})

		Convey("When testing for a non-existent component", func() {
			exists := g.HasComponent("test1")
			Convey("It should return false", func() {
				So(exists, ShouldBeFalse)
			})
		})
	})
}
