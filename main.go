package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tacherasasi/goofer/cli"
	"github.com/tacherasasi/goofer/logger"
)

func main() {
	if len(os.Args) > 1 {
		args := os.Args[1:]
		logger.Debug.Printf("invoking command %+v", args)

		switch args[0] {
		case "prefetch":
			// just run goofer -v to trigger the download
			if err := cli.Run([]string{"-v"}, true); err != nil {
				panic(err)
			}
			os.Exit(0)
			return
		case "init":
			// override default init flags
			args = append(args, "--generator-provider", "go run github.com/tacherasasi/goofer")
			if err := cli.Run(args, true); err != nil {
				panic(err)
			}
			os.Exit(0)
			return
		case "generate":
			// just run goofer -v to trigger the download
			if len(args) == 1 {
				args = append(args, "--generator-provider", "go run github.com/tacherasasi/goofer")
			}

			// goofer CLI
			if err := cli.Run(args, true); err != nil {
				log.Fatalf("error running goofer CLI: %s", err)
			}

			return
		}

		// goofer CLI
		if err := cli.Run(args, true); err != nil {
			panic(err)
		}

		return
	}

	// running the goofer generator

	logger.Debug.Printf("invoking goofer")

	// if this wasn't actually invoked by the goofer generator, print a warning and exit
	if os.Getenv("GOOFER_GENERATOR_INVOCATION") == "" {
		logger.Info.Printf("This command is only meant to be invoked internally. Please run the following instead:")
		logger.Info.Printf("`go run github.com/tacherasasi/goofer <command>`")
		logger.Info.Printf("e.g.")
		logger.Info.Printf("`go run github.com/tacherasasi/goofer generate`")
		os.Exit(1)
	}

	// exit when signal triggers
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		os.Exit(1)
	}()

	if err := invokeGoofer(); err != nil {
		log.Printf("error occurred when invoking goofer: %s", err)
		os.Exit(1)
	}

	logger.Debug.Printf("success")
}

func invokeGoofer() error {
	// TODO: implement invokeGoofer
	return nil
}
