/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package graph

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type testVertex struct {
	name    string
	handler string
}

func (tv *testVertex) Name() string {
	return tv.name
}

func (tv *testVertex) Handler(data interface{}) (string, error) {
	return "", nil
}

func TestGraph(t *testing.T) {
	Convey("Given a new graph", t, func() {
		g := New()
		Convey("When adding a new vertex", func() {
			g.AddVertex(&testVertex{name: "test"})
			Convey("It should be stored on the graph", func() {
				So(len(g.Vertices), ShouldEqual, 1)
				So(g.Vertices[0].Name(), ShouldEqual, "test")
			})
		})

		Convey("When adding a duplicate vertex", func() {
			g.AddVertex(&testVertex{name: "test"})
			Convey("It should not be stored on the graph", func() {
				So(len(g.Vertices), ShouldEqual, 1)
			})
		})

		Convey("When connecting two verticies", func() {
			g.AddVertex(&testVertex{name: "test1"})
			g.AddVertex(&testVertex{name: "test2"})
			err := g.Connect("test1", "test2", "test-event")
			Convey("It should create an edge between the two verticies", func() {
				So(err, ShouldBeNil)
				So(len(g.Vertices), ShouldEqual, 2)
				So(len(g.Edges), ShouldEqual, 1)
				So(g.Edges[0].Source, ShouldEqual, "test1")
				So(g.Edges[0].Destination, ShouldEqual, "test2")
				So(g.Edges[0].Event, ShouldEqual, "test-event")
			})
		})

		Convey("When connecting two verticies mutually", func() {
			g.AddVertex(&testVertex{name: "test1"})
			g.AddVertex(&testVertex{name: "test2"})
			err := g.ConnectMutually("test1", "test2", "test-event")
			Convey("It should create an edge between the two verticies", func() {
				So(err, ShouldBeNil)
				So(len(g.Vertices), ShouldEqual, 2)
				So(len(g.Edges), ShouldEqual, 2)
				So(g.Edges[0].Source, ShouldEqual, "test1")
				So(g.Edges[0].Destination, ShouldEqual, "test2")
				So(g.Edges[0].Event, ShouldEqual, "test-event")
				So(g.Edges[1].Source, ShouldEqual, "test2")
				So(g.Edges[1].Destination, ShouldEqual, "test1")
				So(g.Edges[1].Event, ShouldEqual, "test-event")
			})
		})

		Convey("When connecting two verticies that don't exist", func() {
			err := g.ConnectMutually("fake1", "fake2", "fake-event")
			Convey("It should error", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When getting the origin verticies of a vertex", func() {
			g.AddVertex(&testVertex{name: "test1"})
			g.AddVertex(&testVertex{name: "test2"})
			g.AddVertex(&testVertex{name: "test3"})
			erra := g.Connect("test2", "test1", "test-event")
			errb := g.Connect("test3", "test1", "test-event")
			origins := g.Origins("test1")
			Convey("It should return the correct verticies", func() {
				So(erra, ShouldBeNil)
				So(errb, ShouldBeNil)
				So(len(*origins), ShouldEqual, 2)
				So((*origins)[0], ShouldEqual, "test2")
				So((*origins)[1], ShouldEqual, "test3")
			})
		})

		Convey("When getting the neighbouring verticies of a vertex", func() {
			g.AddVertex(&testVertex{name: "test1"})
			g.AddVertex(&testVertex{name: "test2"})
			g.AddVertex(&testVertex{name: "test3"})
			erra := g.Connect("test1", "test2", "test-event")
			errb := g.Connect("test1", "test3", "test-event")
			neighbours := g.Neighbours("test1")
			Convey("It should return the correct verticies", func() {
				So(erra, ShouldBeNil)
				So(errb, ShouldBeNil)
				So(len(*neighbours), ShouldEqual, 2)
				So((*neighbours)[0], ShouldEqual, "test2")
				So((*neighbours)[1], ShouldEqual, "test3")
			})
		})

		Convey("When getting a vertex by name", func() {
			g.AddVertex(&testVertex{name: "test1"})
			g.AddVertex(&testVertex{name: "test2"})
			vertex := g.Vertex("test2")
			Convey("It should return the correct vertex", func() {
				So(vertex, ShouldNotBeNil)
				So(vertex.Name(), ShouldEqual, "test2")
			})
		})

		Convey("When testing for an existing vertex", func() {
			g.AddVertex(&testVertex{name: "test1"})
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
