package main

import "github.com/r3labs/graph"

func basic() {
	instance := make(graph.GenericComponent)
	instance["_action"] = "create"
	instance["_component"] = "instances"
	instance["_component_id"] = "instances::test1"
	instance["_provider"] = "test"
	instance["name"] = "john"
	instance["size"] = "1024"

	instance2 := make(graph.GenericComponent)
	instance2["_action"] = "update"
	instance2["_component"] = "instances"
	instance2["_component_id"] = "instances::test2"
	instance2["_provider"] = "test"
	instance2["name"] = "snow"
	instance2["size"] = "1024"

	fromGraph := graph.New()

	fromGraph.ID = "id"
	fromGraph.Name = "name"
	fromGraph.Username = "john snow"

	println("Given I have a graph 'g' with 2 instances")
	_ = fromGraph.AddComponent(&instance)
	println("  And I have a graph 'g2' with 2 instances")
	_ = fromGraph.AddComponent(&instance2)

	toGraph := graph.New()
	toGraph.ID = "id"
	toGraph.Name = "name"
	toGraph.Username = "john snow"

	_ = toGraph.AddComponent(&instance)
	instance2["name"] = "lol"
	_ = toGraph.AddComponent(&instance2)

	println(" When I execute a diff between both")
	gx, err := fromGraph.Diff(toGraph)
	if err != nil {
		println(err.Error())
		return
	}

	instances := gx.GetComponents().ByType("instances")

	println(" Then resulting graph must have only one component of type instance")
	if len(instances) != 2 {
		println("[ ERROR ] : Expectation not met")
		return
	}

	println(" Then resulting graph must have only 4 edges")
	if len(gx.Edges) != 4 {
		println("[ ERROR ] : Expectation not met")
		return
	}
}

func main() {
	basic()
}
