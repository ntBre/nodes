package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Node struct {
	Mom        string   `json:"Mom"`
	Type       string   `json:"ntype"`
	State      string   `json:"state"`
	Pcpus      int      `json:"pcpus"`
	Jobs       []string `json:"jobs"`
	Avail      Resource `json:"resources_available"`
	Assign     Resource `json:"resources_assigned"`
	Resv       string   `json:"resv_enable"`
	Sharing    string   `json:"sharing"`
	LastChange int      `json:"last_state_change_time"`
	LastUse    int      `json:"last_used_time"`
}

type PBS struct {
	Time    int             `json:"timestamp"`
	Version string          `json:"pbs_version"`
	Server  string          `json:"pbs_server"`
	Nodes   map[string]Node `json:"nodes"`
}

type Resource struct {
	Arch   string `json:"arch"`
	Host   string `json:"host"`
	Mem    string `json:"mem"`
	Cpus   int    `json:"ncpus"`
	Queues string `json:"Qlist"`
	Vnode  string `json:"vnode"`
}

// TODO percentages seem nice for graphing, but I think a flat number
// option might be convenient for humans

// TODO run the pbsnodes command here:
// pbsnodes -aF json
func main() {
	// this can be replaced with function to call pbsnodes
	f, err := os.Open("pbs.json")
	if err != nil {
		panic(err)
	}
	byts, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	pbs := new(PBS)
	err = json.Unmarshal(byts, pbs)
	if err != nil {
		panic(err)
	}
	// below is probably a function, pass it node list
	fmt.Printf("%5s%8s%8s\n", "Node", "%Mem", "%CPU")
	fmt.Println("---------------------")
	names := make([]string, 0)
	for _, node := range pbs.Nodes {
		names = append(names,
			strings.TrimRight(node.Mom, "cm.cluster"))
	}
	sort.Strings(names)
	for _, n := range names {
		node := pbs.Nodes[n]
		avail, _ := strconv.ParseFloat(
			strings.TrimRight(node.Avail.Mem, "kb"),
			64,
		)
		assign, _ := strconv.ParseFloat(
			strings.TrimRight(node.Assign.Mem, "kb"),
			64,
		)
		fmt.Printf("%5s%8.2f", n, 100*assign/avail)
		fmt.Printf("%8.2f", 100*float64(node.Assign.Cpus)/
			float64(node.Avail.Cpus))
		fmt.Printf("%8s", node.Avail.Queues)
		fmt.Print("\n")
	}
}
