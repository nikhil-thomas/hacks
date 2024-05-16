package main

import (
	"context"
	"fmt"

	"github.com/operator-framework/operator-registry/alpha/action"
	"github.com/operator-framework/operator-registry/alpha/declcfg"
)

const (
	indexImage  string = "quay.io/openshift-pipeline/redhat-operator-index@sha256:93701f7481960e9021051e4da3782463842f7a7c7873fa07614fd767bc36553d"
	packageName string = "openshift-pipelines-operator-rh"
	channelName string = "pipelines-1.7"
)

func main() {
	//var write func(declcfg.DeclarativeConfig, io.Writer) error
	//	write = declcfg.WriteYAML
	render := action.Render{
		Refs: []string{indexImage},
	}
	cfg, err := render.Run(context.Background())
	if err != nil {
		panic(err)
	}
	var entries []declcfg.ChannelEntry
	for _, channel := range cfg.Channels {
		if channel.Package == packageName { //&&
			//channel.Name == channelName {
			entries = channel.Entries
		}
	}

	//if err := write(entries, os.Stdout); err != nil {
	//	log.Fatal(err)
	//}
	fmt.Println(entries)
}
