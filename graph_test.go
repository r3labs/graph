/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package graph

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type testComponent struct {
	Name       string   `json:"name"`
	State      string   `json:"state"`
	Action     string   `json:"action"`
	Deps       []string `json:"deps"`
	Sequential []string `json:"sequential"`
	TestVal    int      `json:"test_val"`
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
	return tv.Deps
}

func (tv *testComponent) Validate() error {
	return nil
}

func (tv *testComponent) Diff(v Component) bool {
	return tv.TestVal != v.(*testComponent).TestVal
}

func (tv *testComponent) SequentialDependencies() []string {
	return tv.Sequential
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

	Convey("Given an existing graph", t, func() {
		ng := New()
		_ = ng.AddComponent(&testComponent{Name: "1", TestVal: 1})
		_ = ng.AddComponent(&testComponent{Name: "2", Deps: []string{"1"}, TestVal: 1})
		_ = ng.AddComponent(&testComponent{Name: "3", Deps: []string{"1"}, TestVal: 1})
		_ = ng.AddComponent(&testComponent{Name: "4", Deps: []string{"2", "3"}, TestVal: 1})

		Convey("That has no verticies", func() {
			eg := New()
			Convey("When diffing the new populated graph", func() {
				g, err := ng.Diff(eg)
				Convey("It should return the correct changes", func() {
					So(err, ShouldBeNil)
					So(len(g.Changes), ShouldEqual, 4)
					So(g.Changes[0].GetID(), ShouldEqual, "1")
					So(g.Changes[0].GetAction(), ShouldEqual, ACTIONCREATE)
					So(g.Changes[1].GetID(), ShouldEqual, "2")
					So(g.Changes[1].GetAction(), ShouldEqual, ACTIONCREATE)
					So(g.Changes[2].GetID(), ShouldEqual, "3")
					So(g.Changes[2].GetAction(), ShouldEqual, ACTIONCREATE)
					So(g.Changes[3].GetID(), ShouldEqual, "4")
					So(g.Changes[3].GetAction(), ShouldEqual, ACTIONCREATE)
				})
				Convey("It should return the correct edges", func() {
					So(len(g.Edges), ShouldEqual, 6)
					So(g.Edges[0].Source, ShouldEqual, "1")
					So(g.Edges[0].Destination, ShouldEqual, "2")
					So(g.Edges[1].Source, ShouldEqual, "1")
					So(g.Edges[1].Destination, ShouldEqual, "3")
					So(g.Edges[2].Source, ShouldEqual, "2")
					So(g.Edges[2].Destination, ShouldEqual, "4")
					So(g.Edges[3].Source, ShouldEqual, "3")
					So(g.Edges[3].Destination, ShouldEqual, "4")
					So(g.Edges[4].Source, ShouldEqual, "start")
					So(g.Edges[4].Destination, ShouldEqual, "1")
					So(g.Edges[4].Source, ShouldEqual, "start")
					So(g.Edges[4].Destination, ShouldEqual, "1")
					So(g.Edges[5].Source, ShouldEqual, "4")
					So(g.Edges[5].Destination, ShouldEqual, "end")
				})
			})
		})

		Convey("That has no verticies with sequential components", func() {
			sng := New()
			_ = sng.AddComponent(&testComponent{Name: "1", TestVal: 1})
			_ = sng.AddComponent(&testComponent{Name: "2", Deps: []string{"1"}, Sequential: []string{"1"}, TestVal: 1})
			_ = sng.AddComponent(&testComponent{Name: "3", Deps: []string{"1"}, Sequential: []string{"1"}, TestVal: 1})
			_ = sng.AddComponent(&testComponent{Name: "4", Deps: []string{"2", "3"}, TestVal: 1})

			eg := New()
			Convey("When diffing the new populated graph", func() {
				g, err := sng.Diff(eg)
				Convey("It should return the correct changes", func() {
					So(err, ShouldBeNil)
					So(len(g.Changes), ShouldEqual, 4)
					So(g.Changes[0].GetID(), ShouldEqual, "1")
					So(g.Changes[0].GetAction(), ShouldEqual, ACTIONCREATE)
					So(g.Changes[1].GetID(), ShouldEqual, "2")
					So(g.Changes[1].GetAction(), ShouldEqual, ACTIONCREATE)
					So(g.Changes[2].GetID(), ShouldEqual, "3")
					So(g.Changes[2].GetAction(), ShouldEqual, ACTIONCREATE)
					So(g.Changes[3].GetID(), ShouldEqual, "4")
					So(g.Changes[3].GetAction(), ShouldEqual, ACTIONCREATE)
				})
				Convey("It should return the correct edges", func() {
					So(len(g.Edges), ShouldEqual, 6)
					So(g.Edges[0].Source, ShouldEqual, "1")
					So(g.Edges[0].Destination, ShouldEqual, "2")
					So(g.Edges[1].Source, ShouldEqual, "2")
					So(g.Edges[1].Destination, ShouldEqual, "3")
					So(g.Edges[2].Source, ShouldEqual, "2")
					So(g.Edges[2].Destination, ShouldEqual, "4")
					So(g.Edges[3].Source, ShouldEqual, "3")
					So(g.Edges[3].Destination, ShouldEqual, "4")
					So(g.Edges[4].Source, ShouldEqual, "start")
					So(g.Edges[4].Destination, ShouldEqual, "1")
					So(g.Edges[4].Source, ShouldEqual, "start")
					So(g.Edges[4].Destination, ShouldEqual, "1")
					So(g.Edges[5].Source, ShouldEqual, "4")
					So(g.Edges[5].Destination, ShouldEqual, "end")
				})
			})
		})

		Convey("That is missing vertex '3'", func() {
			eg := New()
			_ = eg.AddComponent(&testComponent{Name: "1", TestVal: 1})
			_ = eg.AddComponent(&testComponent{Name: "2", Deps: []string{"1"}, TestVal: 1})
			_ = eg.AddComponent(&testComponent{Name: "4", Deps: []string{"2"}, TestVal: 1})
			Convey("When diffing the new populated graph", func() {
				g, err := ng.Diff(eg)
				Convey("It should mark vertex '3' for creation", func() {
					So(err, ShouldBeNil)
					So(len(g.Changes), ShouldEqual, 1)
					So(g.Changes[0].GetID(), ShouldEqual, "3")
					So(g.Changes[0].GetAction(), ShouldEqual, ACTIONCREATE)
				})
				Convey("It should return the correct edges", func() {
					So(len(g.Edges), ShouldEqual, 2)
					So(g.Edges[0].Source, ShouldEqual, "start")
					So(g.Edges[0].Destination, ShouldEqual, "3")
					So(g.Edges[1].Source, ShouldEqual, "3")
					So(g.Edges[1].Destination, ShouldEqual, "end")
				})
			})
		})

		Convey("That has an additional vertex '5'", func() {
			eg := New()
			_ = eg.AddComponent(&testComponent{Name: "1", TestVal: 1})
			_ = eg.AddComponent(&testComponent{Name: "2", Deps: []string{"1"}, TestVal: 1})
			_ = eg.AddComponent(&testComponent{Name: "3", Deps: []string{"1"}, TestVal: 1})
			_ = eg.AddComponent(&testComponent{Name: "4", Deps: []string{"2"}, TestVal: 1})
			_ = eg.AddComponent(&testComponent{Name: "5", Deps: []string{"1"}, TestVal: 1})
			Convey("When diffing the new populated graph", func() {
				g, err := ng.Diff(eg)
				Convey("It should mark vertex '3' for creation", func() {
					So(err, ShouldBeNil)
					So(len(g.Changes), ShouldEqual, 1)
					So(g.Changes[0].GetID(), ShouldEqual, "5")
					So(g.Changes[0].GetAction(), ShouldEqual, ACTIONDELETE)
				})
				Convey("It should return the correct edges", func() {
					So(len(g.Edges), ShouldEqual, 2)
					So(g.Edges[0].Source, ShouldEqual, "start")
					So(g.Edges[0].Destination, ShouldEqual, "5")
					So(g.Edges[1].Source, ShouldEqual, "5")
					So(g.Edges[1].Destination, ShouldEqual, "end")
				})
			})
		})

		Convey("That has changes on its verticies", func() {
			eg := New()
			_ = eg.AddComponent(&testComponent{Name: "1", TestVal: 1})
			_ = eg.AddComponent(&testComponent{Name: "2", Deps: []string{"1"}, TestVal: 2})
			_ = eg.AddComponent(&testComponent{Name: "3", Deps: []string{"1"}, TestVal: 2})
			_ = eg.AddComponent(&testComponent{Name: "4", Deps: []string{"2"}, TestVal: 1})
			Convey("When diffing the new populated graph", func() {
				g, err := ng.Diff(eg)
				Convey("It should mark vertex '2' and '3' for update", func() {
					So(err, ShouldBeNil)
					So(len(g.Changes), ShouldEqual, 2)
					So(g.Changes[0].GetID(), ShouldEqual, "2")
					So(g.Changes[0].GetAction(), ShouldEqual, ACTIONUPDATE)
					So(g.Changes[1].GetID(), ShouldEqual, "3")
					So(g.Changes[1].GetAction(), ShouldEqual, ACTIONUPDATE)
				})
				Convey("It should return edges connected sequentially", func() {
					So(len(g.Edges), ShouldEqual, 3)
					So(g.Edges[0].Source, ShouldEqual, "start")
					So(g.Edges[0].Destination, ShouldEqual, "2")
					So(g.Edges[1].Source, ShouldEqual, "2")
					So(g.Edges[1].Destination, ShouldEqual, "3")
					So(g.Edges[2].Source, ShouldEqual, "3")
					So(g.Edges[2].Destination, ShouldEqual, "end")
				})
			})
		})

		Convey("That has sequential dependencies", func() {
			eg := New()
			_ = eg.AddComponent(&testComponent{Name: "1", TestVal: 1})
			_ = eg.AddComponent(&testComponent{Name: "2", Deps: []string{"1"}, Sequential: []string{"1"}, TestVal: 2})
			_ = eg.AddComponent(&testComponent{Name: "3", Deps: []string{"1"}, Sequential: []string{"1"}, TestVal: 2})
			_ = eg.AddComponent(&testComponent{Name: "4", Deps: []string{"2"}, TestVal: 1})
			Convey("When diffing the new populated graph", func() {
				g, err := ng.Diff(eg)
				Convey("It should mark vertex '2' and '3' for update", func() {
					So(err, ShouldBeNil)
					So(len(g.Changes), ShouldEqual, 2)
					So(g.Changes[0].GetID(), ShouldEqual, "2")
					So(g.Changes[0].GetAction(), ShouldEqual, ACTIONUPDATE)
					So(g.Changes[1].GetID(), ShouldEqual, "3")
					So(g.Changes[1].GetAction(), ShouldEqual, ACTIONUPDATE)
				})
				Convey("It should return edges connected sequentially", func() {
					So(len(g.Edges), ShouldEqual, 3)
					So(g.Edges[0].Source, ShouldEqual, "start")
					So(g.Edges[0].Destination, ShouldEqual, "2")
					So(g.Edges[1].Source, ShouldEqual, "2")
					So(g.Edges[1].Destination, ShouldEqual, "3")
					So(g.Edges[2].Source, ShouldEqual, "3")
					So(g.Edges[2].Destination, ShouldEqual, "end")
				})
			})
		})

		Convey("That is missing vertex all verticies and has sequential dependencies", func() {
			sng := New()
			_ = sng.AddComponent(&testComponent{Name: "1", TestVal: 1})
			_ = sng.AddComponent(&testComponent{Name: "2", Deps: []string{"1"}, Sequential: []string{"1"}, TestVal: 1})
			_ = sng.AddComponent(&testComponent{Name: "3", Deps: []string{"1"}, Sequential: []string{"1"}, TestVal: 1})
			_ = sng.AddComponent(&testComponent{Name: "4", Deps: []string{"2", "3"}, TestVal: 1})

			eg := New()
			Convey("When diffing the new populated graph", func() {
				g, err := eg.Diff(sng)
				Convey("It should mark all verticies for deletion", func() {
					So(err, ShouldBeNil)
					So(len(g.Changes), ShouldEqual, 4)
					So(g.Changes[0].GetID(), ShouldEqual, "1")
					So(g.Changes[0].GetAction(), ShouldEqual, ACTIONDELETE)
					So(g.Changes[1].GetID(), ShouldEqual, "2")
					So(g.Changes[1].GetAction(), ShouldEqual, ACTIONDELETE)
					So(g.Changes[2].GetID(), ShouldEqual, "3")
					So(g.Changes[2].GetAction(), ShouldEqual, ACTIONDELETE)
					So(g.Changes[3].GetID(), ShouldEqual, "4")
					So(g.Changes[3].GetAction(), ShouldEqual, ACTIONDELETE)
				})
				Convey("It should return the correct edges", func() {
					So(len(g.Edges), ShouldEqual, 6)
					So(g.Edges[0].Source, ShouldEqual, "2")
					So(g.Edges[0].Destination, ShouldEqual, "1")
					So(g.Edges[1].Source, ShouldEqual, "3")
					So(g.Edges[1].Destination, ShouldEqual, "1")
					So(g.Edges[2].Source, ShouldEqual, "4")
					So(g.Edges[2].Destination, ShouldEqual, "2")
					So(g.Edges[3].Source, ShouldEqual, "2")
					So(g.Edges[3].Destination, ShouldEqual, "3")
					So(g.Edges[4].Source, ShouldEqual, "1")
					So(g.Edges[4].Destination, ShouldEqual, "end")
					So(g.Edges[5].Source, ShouldEqual, "start")
					So(g.Edges[5].Destination, ShouldEqual, "4")
				})
			})
		})
	})
}
