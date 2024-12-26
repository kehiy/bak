package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip13"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/spf13/cobra"
)

func buildListCmd(parentCmd *cobra.Command) {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list all badges issued by an npub",
	}

	parentCmd.AddCommand(listCmd)

	npubOpt := listCmd.Flags().StringP("npub", "p", "", "target npub")

	listCmd.Run = func(cmd *cobra.Command, _ []string) {

		relays := []nostr.Relay{}
		for _, r := range *relayAddrsOpt {
			relay, err := nostr.RelayConnect(context.Background(), r)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			relays = append(relays, *relay)
		}

		prefix, npub, err := nip19.Decode(*npubOpt)
		if err != nil || prefix != "npub" {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		pub := npub.(string)
		filters := []nostr.Filter{{
			Kinds:   []int{nostr.KindBadgeDefinition},
			Authors: []string{pub},
		}}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		resp := make(map[string]nostr.Event)
		for _, r := range relays {
			sub, err := r.Subscribe(ctx, filters)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			for ev := range sub.Events {
				resp[ev.ID] = *ev
			}
		}

		for _, e := range resp {
			cmd.Println("ID:", fmt.Sprintf("https://njump.me/%s", e.ID))
			cmd.Println("Unique Name:", e.Tags.GetD())
			cmd.Println("Name:", e.Tags.GetFirst([]string{"name"}).Value())
			cmd.Println("Description:", e.Tags.GetFirst([]string{"description"}).Value())
			cmd.Println("Image:", e.Tags.GetFirst([]string{"image"}).Value())
			cmd.Println("PoW Rarity:", nip13.Difficulty(e.ID))
			cmd.Println("Created At:", e.CreatedAt.Time().Format("2006-January-02"))
			cmd.Println("Age:", time.Since(e.CreatedAt.Time()))
			cmd.Println()
		}

		cmd.Println("Total Badges issued:", len(resp))
	}
}
