package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

	// ssh namespace
	sshTools := &cobra.Command{
		Use:   "ssh",
		Short: "SSH tools",
	}
	rootCmd.AddCommand(sshTools)

	sshKnownHosts := &cobra.Command{
		Use:   "known-hosts",
		Short: "known_hosts",
	}
	sshTools.AddCommand(sshKnownHosts)

	sshKnownHostsList := &cobra.Command{
		Use:   "diff",
		Short: "Show diff between ~/.ssh/known_hosts and ~/.ssh/known_hosts.db",
		Run: func(cmd *cobra.Command, args []string) {
			trusted := sshKnownHostsRead(".db")
			active := sshKnownHostsRead("")

			// check diff
			temp := make([]string, 0)
			for _, a := range active {
				found := false
				for _, t := range trusted {
					if t == a {
						found = true
						break
					}
				}
				if !found {
					temp = append(temp, a)
				}
			}

			if len(temp) > 0 {
				fmt.Print(ctc.ForegroundYellow, "New hosts (", len(temp), ")", ctc.Reset, "\n")
				for _, s := range temp {
					fmt.Println(s)
				}
			} else {
				fmt.Print(ctc.ForegroundGreen, "No new hosts\n", ctc.Reset)
			}
		},
	}
	sshKnownHosts.AddCommand(sshKnownHostsList)

	sshKnownHostsReset := &cobra.Command{
		Use:   "reset",
		Short: "Replace ~/.ssh/known_hosts with ~/.ssh/known_hosts.db",
		Run: func(cmd *cobra.Command, args []string) {
			trusted := sshKnownHostsRead(".db")
			sshKnownHostsWrite("", trusted)
			fmt.Printf("%s%d hosts written%s\n", ctc.ForegroundGreen, len(trusted), ctc.Reset)
		},
	}
	sshKnownHosts.AddCommand(sshKnownHostsReset)

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

func sshKnownHostsRead(suffix string) []string {
	res := make([]string, 0)
	if f, err := os.Open(fmt.Sprint(os.Getenv("HOME"), "/.ssh/known_hosts", suffix)); err == nil {
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err == nil {
			for _, s := range strings.Split(string(b), "\n") {
				if strings.TrimSpace(s) != "" {
					res = append(res, s)
				}
			}
		}
	}

	return res
}

func sshKnownHostsWrite(suffix string, lines []string) {
	if f, err := os.OpenFile(fmt.Sprint(os.Getenv("HOME"), "/.ssh/known_hosts", suffix), os.O_WRONLY|os.O_APPEND, 0644); err == nil {
		defer f.Close()
		f.Truncate(0)

		for _, line := range lines {
			fmt.Fprintln(f, line)
		}
	}
}
