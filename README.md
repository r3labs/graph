# graph

Graph is a library for handling graph functions.

## Core concepts

Graph : representation of a transition of a group of components to its updated version. A graph is a combination of components and edges that will become on a graph changes.

  - Component : represents an environment component like an instance, a sql server, anything can be a component as long as it exposes [this interface](component.go#L8)
  - Change : the processed component.
  - Edge : transition from a source to a destination

## Basic usage

In order to create a new graph you can use the default constructor:
```go
g := graph.New(
  ID:       "id",
  Name:     "name",
  UserID:   "uID",
  Username: "john snow"
)
```

Lets add an instance as a new component:
```go
instance := make(graph.GenericComponent)
instance["_action"] = "none"
instance["_component"] = "instances"
instance["_component_id"] = "instances::test1"
instance["_provider"] = "test"
instance["name"] = "john"
instance["size"] = "1024"

err = g.AddComponent(instance)
```

## Actions

In order to get the edges for the creation of this component, you may want to change the `_action` field, this actually accepts [different action](graph.go#L16). The previous example will look like:

```go
instance := make(graph.GenericComponent)
instance["_action"] = "create"
instance["_component"] = "instances"
instance["_component_id"] = "instances::test1"
instance["_provider"] = "test"
instance["name"] = "john"
instance["size"] = "1024"

err = g.AddComponent(instance)
```

inspecting graph, will result on a set of valid edges to be processed.

Note: **none** action will result on non processing the specific component, so it wont be processed.

## Diffs

Another interesting capability of graph library is the ability to calculate the diff between two graphs. This is accomplished with the `Diff` method and its variants.

You can find an example on how to diff graphs on the [basic example](examples/basic.go)


## Managing dependencies

In many scenarios you may want to process a component before another, this causes a dependency between both. Graph library is also capable to manage that with `Connect` function family.

An easy to understand example is a sql server and its databases, you can't create databases before the server itself.

```go
g := graph.New(
  ID:       "id",
  Name:     "name",
  UserID:   "uID",
  Username: "john snow"
)
server := make(graph.GenericComponent)
server["_component"] = "server"
server["_component_id"] = "server::test"
server["_provider"] = "test"
server["name"] = "mySQL"
err = g.AddComponent(server)

db := make(graph.GenericComponent)
db["_component"] = "db"
db["_component_id"] = "db::test"
db["_provider"] = "test"
db["name"] = "myDB"
err = g.AddComponent(db)

g.Connect(sql, db.GetID())
```

When you process this graph, you'll see that the sql server will be processed before the database.


## Build status

* master: [![CircleCI](https://circleci.com/gh/r3labs/graph/tree/master.svg?style=svg)](https://circleci.com/gh/r3labs/graph/tree/master)
* develop: [![CircleCI](https://circleci.com/gh/r3labs/graph/tree/develop.svg?style=svg)](https://circleci.com/gh/r3labs/graph/tree/develop)

## Installation

```
make deps
make install
```

or

```
go get github.com/r3labs/graph
```

## Running Tests

```
make test
```

## Contributing

Please read through our
[contributing guidelines](CONTRIBUTING.md).
Included are directions for opening issues, coding standards, and notes on
development.

Moreover, if your pull request contains patches or features, you must include
relevant unit tests.

## Versioning

For transparency into our release cycle and in striving to maintain backward
compatibility, this project is maintained under [the Semantic Versioning guidelines](http://semver.org/).

## Copyright and License

Code and documentation copyright since 2015 r3labs.io authors.

Code released under
[the Mozilla Public License Version 2.0](LICENSE).
