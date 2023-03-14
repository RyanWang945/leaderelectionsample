package main

import (
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "test",
	Short: "test for leader election",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runController()
	},
}

func runController() error {
	glog.Info("controller is running...")
	return nil
}
