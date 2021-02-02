package main

import (
	"crypto/rand"
	"encoding/base64"
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

	// namespaces

	randomGen := &cobra.Command{
		Use:   "rand",
		Short: "Random related generators",
	}
	rootCmd.AddCommand(randomGen)

	golang := &cobra.Command{
		Use:   "go",
		Short: "Go language related generators",
	}
	rootCmd.AddCommand(golang)

	// go namespace

	golangMain := &cobra.Command{
		Use:   "main",
		Short: "Generate initial main()",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(golangMainTpl)
		},
	}
	golang.AddCommand(golangMain)

	// rand namespace

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
	randomGen.AddCommand(uuid4)

	randomBytes := &cobra.Command{
		Use:   "bytes",
		Short: "Generates bytes encoded with base64",
		Run: func(cmd *cobra.Command, args []string) {
			size, err := cmd.Flags().GetUint32("size")
			if err != nil {
				logError(err)
			}

			b := make([]byte, size)
			_, err = rand.Read(b)
			if err != nil {
				logError(err)
			}
			encoded := base64.StdEncoding.EncodeToString(b)
			fmt.Println(encoded)
		},
	}
	randomBytes.Flags().Uint32("size", 64, "specify bytes length")
	randomGen.AddCommand(randomBytes)

	//

	if err := rootCmd.Execute(); err != nil {
		logError(err)
		os.Exit(1)
	}
}

const golangMainTpl = `package main

import "fmt"

func main() {
	fmt.Println("Hello")
}`
