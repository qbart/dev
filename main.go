package main

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/wzshiming/ctc"
)

func logError(err error) {
	fmt.Print(ctc.ForegroundRed, err, ctc.Reset, "\n")
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "dev",
		Short: "dev is a collection of utilities that might come handy when developing a software",
	}

	uuid4 := &cobra.Command{
		Use:   "uuid4",
		Short: "Generates UUID v4",
		Run: func(cmd *cobra.Command, args []string) {
			uuid, err := uuid.NewUUID()
			if err != nil {
				logError(err)
				os.Exit(1)
			}

			fmt.Println(uuid.String())

		},
	}
	rootCmd.AddCommand(uuid4)

	if err := rootCmd.Execute(); err != nil {
		logError(err)
		os.Exit(1)
	}
}
