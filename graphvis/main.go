package main

import (
	"bytes"
	"flag"
	"fmt"
	"devgo/graphvis/config"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

type Node struct {
	Attrs map[string]string
}
type Edge struct {
	Attrs map[string]string
}
type Graph struct {
	Title string
	Node  map[string]map[string]*Node
	Edge  map[string]map[string]*Edge
	Attrs map[string]string
	group []string
}

func NewGraph(name string) *Graph {
	return &Graph{
		Title: name,
		Node:  make(map[string]map[string]*Node),
		Edge:  make(map[string]map[string]*Edge),
		Attrs: map[string]string{},
		group: []string{},
	}
}

type StorageSpec struct {
	StorageId     int    `json:"storageid,omitempy"`
	MgmtEndpoint  string `json:"mgmtendpoint,omitempy"`
	StorageVendor string `json:"storagevendor,omitempy"`
	Username      string `json:"username,omitempy"`
	Password      string `json:"password,omitempy"`
}

func (g *Graph) AddNode(group string, name string, N *Node) error {
	if _, exists := g.Node[group]; !exists {
		g.Node[group] = make(map[string]*Node)
		g.Node[group][name] = N
	} else {
		g.Node[group][name] = N
	}
	return nil
}
func (g *Graph) AddEdge(from string, to string, E *Edge) error {
	if _, exists := g.Edge[from]; !exists {
		g.Edge[from] = make(map[string]*Edge)
		g.Edge[from][to] = E
	} else {
		g.Edge[from][to] = E
	}
	return nil
}

func (g *Graph) WriteDot(w io.Writer) error {
	log.Printf("writeing dot...")
	tab := string("\t")
	fmt.Fprintf(w, "digraph %s {\n", g.Title)
	fmt.Fprintf(w, "%s%s\n", tab, g.GetLines(g.Attrs))
	for key, val := range g.Node {
		fmt.Fprintf(w, "%ssubgraph %s {\n", tab, key)
		fmt.Fprintf(w, "%sgraph [rank=same];\n", tab)
		//fmt.Fprintf(w, "%node [fontsize=10];\n", tab)
		for k, v := range val {
			fmt.Fprintf(w, "%s\"%s\" [ %s ]; \n", tab, k, g.GetString(v.Attrs))
		}
		fmt.Fprintf(w, "%s}\n", tab)
	}
	//str := strings.Join(g.group, "->")
	//fmt.Fprintf(w, "%s\n", str)
	for k1, v1 := range g.Edge {
		for k2, v2 := range v1 {
			fmt.Fprintf(w, "%s\"%s\"->\"%s\" [ %s ];\n", tab, k1, k2, g.GetString(v2.Attrs))
		}
	}
	fmt.Fprintf(w, "}\n")
	return nil
}

func (g *Graph) GetList(a map[string]string) []string {
	str := []string{}
	for k, atr := range a {
		str = append(str, fmt.Sprintf("%s=%q", k, atr))
	}
	return str
}
func (g *Graph) GetString(a map[string]string) string { //[괄호안의 속성]
	return strings.Join(g.GetList(a), " ")
}
func (g *Graph) GetLines(a map[string]string) string {
	str := strings.Join(g.GetList(a), "; ")
	return fmt.Sprintf("%s; ", str)
}

func (g *Graph) CreateHpe3par(stg *StorageSpec) error {

	// d := &hpe3par.Driver{}
	// d.Setup(stg)
	// g.group = []string{"PORT", "VLAN", "IP", "VHOST", "VOL", "CPG", "RAID", "DISK"}

	// portlist, _ := d.Client.GetPortsSys()
	// vhostlist, _ := d.Client.GetVhostSys()
	// vollist, _ := d.Client.GetVolumeSys()
	// vlunlist, _ := d.Client.GetVLunSys()
	// cpglist, _ := d.Client.GetCpgSys()
	// disklist, _ := d.Client.GetPhysicalDisk()
	// //Group
	return nil
}

func (g *Graph) CreateNetapp(stg *StorageSpec) error {

	// d := &netapp.Driver{}
	// d.Setup(stg)
	// g.group = []string{"PORT", "VLAN", "LIF", "SVM", "VOL", "AGGR", "DISK"}

	// portlist, _ := d.Client.GetPortSys()
	// vlanlist, _ := d.Client.GetVlanSys()
	// liflist, _ := d.Client.GetLifSys()
	// svmlist, _ := d.Client.GetSvmSys()
	// vollist, _ := d.Client.GetVolumeSys()
	// aggrlist, _ := d.Client.GetAggregateSys()
	// disklist, _ := d.Client.GetDiskSys()
	return nil
}

// DotToImage generates a SVG using the 'dot' utility, returning the filepath
func DotToImage(outfname string, format string, dot []byte) (string, error) {
	var dotExe string
	if dotExe == "" {
		dot, err := exec.LookPath("dot")
		if err != nil {
			log.Fatalln("unable to find program 'dot', please install it or check your PATH")
		}
		dotExe = dot
	}

	var img string
	//outfname = "go-callvis_export"
	img = fmt.Sprintf("%s.%s", outfname, format)

	cmd := exec.Command(dotExe, fmt.Sprintf("-T%s", format), "-o", img)
	cmd.Stdin = bytes.NewReader(dot)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return img, nil
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	env := flag.String("env", "config", "Environment")
	flag.Parse()
	config.LoadConfigFile(*env)
}

func main() {
	//targets := []string{"netapp", "hpe3par", "ceph"}
	targets := []string{"netapp", "hpe3par", "ceph"}
	storagelist := []StorageSpec{}
	for idx, target := range targets {
		storagelist = append(storagelist, StorageSpec{
			StorageId:     idx,
			MgmtEndpoint:  config.GetString(fmt.Sprintf("storage.%s.endpoint", target)),
			StorageVendor: config.GetString(fmt.Sprintf("storage.%s.name", target)),
			Username:      config.GetString(fmt.Sprintf("storage.%s.username", target)),
			Password:      config.GetString(fmt.Sprintf("storage.%s.passwd", target)),
		})
	}
	for _, stg := range storagelist {
		g := NewGraph(stg.StorageVendor)
		switch stg.StorageVendor {
		case "netapp":
			g.Attrs = map[string]string{
				"labeljust": "l",
				"fontname":  "Arial",
				"fontsize":  "10",
				"rankdir":   "LR",
				"penwidth":  "0.5",
				"nodesep":   "0.35",
				"minlen":    "5",
			}
			g.CreateNetapp(&stg)

		case "hpe3par":
			g.Attrs = map[string]string{
				"labeljust": "l",
				"fontname":  "Arial",
				"fontsize":  "10",
				"rankdir":   "LR",
				"penwidth":  "0.5",
				"nodesep":   "0.35",
				"minlen":    "5",
			}
			g.CreateHpe3par(&stg)
		}

		var buf bytes.Buffer
		g.WriteDot(&buf)
		//WriteDot(os.Stdout, g)
		fmt.Printf("%s", buf.String())
		filename := fmt.Sprintf("%s_%d", stg.StorageVendor, stg.StorageId)
		writeErr := ioutil.WriteFile(filename+".gv", []byte(buf.String()), 0755)
		if writeErr != nil {
			log.Fatalf("%v\n", writeErr)
		}
		DotToImage(filename, "svg", []byte(buf.String()))
	}
}
