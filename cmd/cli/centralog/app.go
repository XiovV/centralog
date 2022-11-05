package centralog

import (
	"bufio"
	"fmt"
	"github.com/XiovV/centralog-cli/pkg/prompt"
	"os"
	"text/tabwriter"
)

type App struct {
}

func (a *App) ListNodesCmd() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(w, "NAME\tCONTAINERS\tSTATUS\t")
	fmt.Fprintln(w, "node1\t2/3\tUP\t")
	fmt.Fprintln(w, "myNewNode\t0/5\tDOWN\t")

	w.Flush()
}

func (a *App) ListContainersCmd(node string) {
	if node != "node1" && node != "myNewNode" {
		fmt.Println(node, "does not exist")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTATUS")
	if node == "node1" {
		fmt.Fprintln(w, "instance1\tUP")
		fmt.Fprintln(w, "instance2\tUP")
		fmt.Fprintln(w, "instance3\tDOWN")
	} else {
		fmt.Fprintln(w, "instance1\tDOWN")
		fmt.Fprintln(w, "instance2\tDOWN")
		fmt.Fprintln(w, "instance3\tDOWN")
		fmt.Fprintln(w, "instance4\tDOWN")
		fmt.Fprintln(w, "instance5\tDOWN")
	}

	w.Flush()
}

func (a *App) ShowLogs() {

	w := bufio.NewWriterSize(os.Stdout, 2)

	w.WriteString("s1")
	w.WriteString("s2")
	w.WriteString("s3")
	w.WriteString("s4")
	w.WriteString("s5")
	w.WriteString("s6")
	w.WriteString("s")

	fmt.Println("")

	w.Flush()

	//
	//fmt.Fprintln(w, "testsdfsdfsdfsdfsdfdsf")
	//fmt.Fprintln(w, "testsdfsdfsdfsdfsdfdsf")
	//fmt.Fprintln(w, "testsdfsdfsdfsdfsdfdsf")
	//fmt.Fprintln(w, "testsdfsdfsdfsdfsdfdsf")

	//w.Flush()

	//fmt.Printf("node=%s container=%s message: %s\n", "node1")
}

func (a *App) AddNodeWithPrompt() {
	prompt.GetURL("Enter your node's URL")
	_, err := prompt.GetAPIKey("Enter your node's API key")
	if err != nil {
		fmt.Println(err)
	}

	name := prompt.GetString("Enter your node's custom name")

	fmt.Printf("Node %s added successfully\n", name)
}

func (a *App) AddNodeWithFlags(url, apiKey, name string) {
	fmt.Printf("Node %s added successfully\n", name)
}
