/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package graph

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type testVertex struct {
	Name   string `json:"name"`
	State  string `json:"state"`
	Action string `json:"action"`
}

func (tv *testVertex) GetID() string {
	return tv.Name
}

func (tv *testVertex) GetProvider() string {
	return "test"
}

func (tv *testVertex) GetType() string {
	return "test"
}

func (tv *testVertex) GetState() string {
	return tv.State
}

func (tv *testVertex) SetState(state string) {
	tv.State = state
}

func (tv *testVertex) GetAction() string {
	return tv.Action
}

func (tv *testVertex) SetAction(action string) {
	tv.Action = action
}

func (tv *testVertex) GetGroup() string {
	return "test"
}

func (tv *testVertex) Dependencies() []string {
	return []string{}
}

func (tv *testVertex) Diff(i interface{}) {}

func (tv *testVertex) Rebuild(i interface{}) {}

func (tv *testVertex) Update(i interface{}) bool {
	return true
}

func (tv *testVertex) IsStateful() bool {
	return true
}

func TestGraph(t *testing.T) {
	Convey("Given a new graph", t, func() {
		g := New()

		Convey("When adding a new vertex", func() {
			g.AddVertex(&testVertex{Name: "test"})
			Convey("It should be stored on the graph", func() {
				So(len(g.Vertices), ShouldEqual, 1)
				So(g.Vertices[0].GetID(), ShouldEqual, "test")
			})
		})

		Convey("When adding a duplicate vertex", func() {
			g.AddVertex(&testVertex{Name: "test"})
			Convey("It should not be stored on the graph", func() {
				So(len(g.Vertices), ShouldEqual, 1)
			})
		})

		Convey("When connecting two verticies", func() {
			g.AddVertex(&testVertex{Name: "test1"})
			g.AddVertex(&testVertex{Name: "test2"})
			err := g.Connect("test1", "test2")
			Convey("It should create an edge between the two verticies", func() {
				So(err, ShouldBeNil)
				So(len(g.Vertices), ShouldEqual, 2)
				So(len(g.Edges), ShouldEqual, 1)
				So(g.Edges[0].Source, ShouldEqual, "test1")
				So(g.Edges[0].Destination, ShouldEqual, "test2")
			})
		})

		Convey("When connecting two verticies mutually", func() {
			g.AddVertex(&testVertex{Name: "test1"})
			g.AddVertex(&testVertex{Name: "test2"})
			err := g.ConnectMutually("test1", "test2")
			Convey("It should create an edge between the two verticies", func() {
				So(err, ShouldBeNil)
				So(len(g.Vertices), ShouldEqual, 2)
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

		Convey("When getting the origin verticies of a vertex", func() {
			g.AddVertex(&testVertex{Name: "test1"})
			g.AddVertex(&testVertex{Name: "test2"})
			g.AddVertex(&testVertex{Name: "test3"})
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

		Convey("When getting the neighbouring verticies of a vertex", func() {
			g.AddVertex(&testVertex{Name: "test1"})
			g.AddVertex(&testVertex{Name: "test2"})
			g.AddVertex(&testVertex{Name: "test3"})
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

		Convey("When disconnecting a vertex from the graph", func() {
			g.AddVertex(&testVertex{Name: "test1"})
			g.AddVertex(&testVertex{Name: "test2"})
			g.AddVertex(&testVertex{Name: "test3"})
			g.AddVertex(&testVertex{Name: "test4"})
			g.AddVertex(&testVertex{Name: "test5"})
			erra := g.Connect("test1", "test2")
			errb := g.Connect("test2", "test3")
			errc := g.Connect("test2", "test4")
			errd := g.Connect("test4", "test5")

			err := g.RemoveVertex("test2")
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

		Convey("When getting a vertex by name", func() {
			g.AddVertex(&testVertex{Name: "test1"})
			g.AddVertex(&testVertex{Name: "test2"})
			vertex := g.Vertex("test2")
			Convey("It should return the correct vertex", func() {
				So(vertex, ShouldNotBeNil)
				So(vertex.GetID(), ShouldEqual, "test2")
			})
		})

		Convey("When testing for an existing vertex", func() {
			g.AddVertex(&testVertex{Name: "test1"})
			exists := g.HasVertex("test1")
			Convey("It should return true", func() {
				So(exists, ShouldBeTrue)
			})
		})

		Convey("When testing for a non-existent vertex", func() {
			exists := g.HasVertex("test1")
			Convey("It should return false", func() {
				So(exists, ShouldBeFalse)
			})
		})
	})
}
