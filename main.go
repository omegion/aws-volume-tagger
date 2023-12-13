package main

import (
	"github.com/omegion/aws-volume-tagger/cmd"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	root := &cobra.Command{
		Use:          "volume-tagger",
		Short:        "AWS Volume Tagger.",
		Long:         "AWS Volume Tagger for EBS Volumes.",
		SilenceUsage: true,
	}

	root.AddCommand(cmd.TagCommand())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
