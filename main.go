package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	nsecOpt       *string
	relayAddrsOpt *[]string
)

func main() {
	rootCmd := &cobra.Command{
		Use:               "nadge",
		Short:             "a nostr badge client.",
		Version:           "nadge v0.1.0",
		CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
	}

	relayAddrsOpt = rootCmd.PersistentFlags().StringSlice("relays", []string{}, "relay websocket address")
	nsecOpt = rootCmd.PersistentFlags().String("nsec", "", "nostr secret key")

	buildListCmd(rootCmd)
	buildIssueCmd(rootCmd)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("%v", err)
	}
}
