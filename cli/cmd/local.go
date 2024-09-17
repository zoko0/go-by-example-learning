package cmd

import (
  "log"
  "os/exec"
  "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(localCmd)
}

var localCmd = &cobra.Command{
  Use:   "local",
  Short: "Run go-by-example server locally",
  Long:  `Runs go-by-example server locally.
  By default runs with MySql local default setup. More info in project's README`,
  Run: func(cmd *cobra.Command, args []string) {
    command := exec.Command("go","run","main.go", "root", "new-password", "&>./app.log", "&")

    if err := command.Start(); err != nil {
      log.Println(err.Error())
      return
    }
  },
}
