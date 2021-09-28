package main

import (
	"bytes"
	"fmt"
	"sort"
	"text/template"

	"github.com/radovskyb/watcher"
)

func iterTemplate(arr []string, varMap map[string]interface{}) (r []string) {
	for _, el := range arr {
		r = append(r, execTemplate(el, varMap))
	}
	return
}

func execTemplate(tplStr string, varMap map[string]interface{}) string {
	tmpl := template.Must(
		template.New("new.tmpl").Parse(tplStr),
	)
	buf := &bytes.Buffer{}
	err := tmpl.Execute(buf, varMap)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func printAvailableVars() {
	vm := makeVarMap(watcher.Event{})
	var iterator []string
	for el := range vm {
		iterator = append(iterator, el)
	}
	sort.Strings(iterator)
	fmt.Printf("\nThe following vars are available:\n")
	for _, val := range iterator {
		fmt.Printf("    {{.%s}}\n", val)
	}
	fmt.Printf("\n")
}
