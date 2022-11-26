package ignorepattern

import barpkg "general/y"

type Graph int32 // want Graph:"^Graph_GRAPH_UNSPECIFIED,Graph_GRAPH_LINE,Graph_GRAPH_PIE$"

const (
	Graph_GRAPH_UNSPECIFIED Graph = 0
	Graph_GRAPH_LINE        Graph = 1
	Graph_GRAPH_PIE         Graph = 2
)

func _a() {
	var e Graph
	switch e { // want "^missing cases in switch of type ignorepattern.Graph: ignorepattern.Graph_GRAPH_PIE$"
	case Graph_GRAPH_LINE:
	}

	_ = map[Graph]int{ // want "^missing keys in map of key type ignorepattern.Graph: ignorepattern.Graph_GRAPH_PIE$"
		Graph_GRAPH_LINE: 1,
	}
}

func _b() {
	var p barpkg.Phylum
	switch p { // want "^missing cases in switch of type bar.Phylum: bar.Mollusca$"
	case barpkg.Chordata:
	}

	_ = map[barpkg.Phylum]int{ // want "^missing keys in map of key type bar.Phylum: bar.Mollusca$"
		barpkg.Chordata: 1,
	}
}
