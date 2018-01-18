package main

import (
	"github.com/zish/simplessh"
	"log"
	"os"
	"io/ioutil"
	"encoding/json"
)

type input struct {
	Opts *simplessh.Opts
	Commands []*simplessh.Command
}

func main () {
	log.SetFlags(log.Lshortfile)
	log.SetFlags(log.Ltime)

	var (
		errs []error
		err error
		inputFile string
		inputJson []byte
		inputData *input
		sshClient *simplessh.Client
		output string
	)

	if len(os.Args) < 2 {
		log.Fatal("Please specify an input JSON file.")
	}

	inputFile = os.Args[1]

	if inputJson, err = ioutil.ReadFile(inputFile); err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(inputJson, &inputData); err != nil {
		log.Fatal(err)
	}

	if sshClient, err = simplessh.New(inputData.Opts); err != nil {
		log.Printf("Error creating new SSH Client: %e\n", err)
	} else {
		log.Printf("%s\n", "SSH Client created successfully")
	}

	log.Printf("%s\n", "SSH Client Start Connect To Host...")

	if err = sshClient.Connect(); err != nil {
		log.Printf("SSH Client Error Connecting: %q\n", err)

	} else {
		log.Printf("%s\n", "SSH Client Successfully Connected")
	}

	for _, command := range inputData.Commands {
		if output, errs = sshClient.RunCommand(command); len(errs) > 0 {
			log.Printf("command: ERROR(S): %q\n", errs)
		}
		log.Println("command output:\n"+output)
	}

	if err = sshClient.Disconnect(); err != nil {
		log.Printf("Error disconnecting from SSH Session: %q\n", err)
	}
}
