package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip13"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/spf13/cobra"
)

type Template struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	PoW         int    `json:"pow"`
}

func buildIssueCmd(parentCmd *cobra.Command) {
	listCmd := &cobra.Command{
		Use:   "issue",
		Short: "issue a new badge",
	}

	parentCmd.AddCommand(listCmd)

	listCmd.Run = func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.PrintErr("please input badge information file!")
			os.Exit(1)
		}

		file, err := os.ReadFile(args[0])
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		info := new(Template)
		if err := json.Unmarshal(file, info); err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		_, v, err := nip19.Decode(*nsecOpt)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		sk := v.(string)

		pk, err := nostr.GetPublicKey(sk)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		e := nostr.Event{
			PubKey:    pk,
			CreatedAt: nostr.Timestamp(time.Now().Unix()),
			Kind:      nostr.KindBadgeDefinition,
			Content:   "",
			Tags: nostr.Tags{
				{"d", info.ID},
				{"name", info.Name},
				{"description", info.Description},
				{"image", info.Image},
			},
		}

		tag, err := nip13.DoWork(context.Background(), e, info.PoW)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		e.Tags = append(e.Tags, tag)

		if err = e.Sign(sk); err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		cmd.Println("Badge Created successfully!")

		ej, err := e.MarshalJSON()
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		cmd.Println("Event:", string(ej))

		cmd.Println("Publishing to relays...")

		ctx := context.Background()
		for _, url := range *relayAddrsOpt {
			relay, err := nostr.RelayConnect(ctx, url)
			if err != nil {
				cmd.Println(err)
				continue
			}
			if err := relay.Publish(ctx, e); err != nil {
				cmd.Println(err)
				continue
			}

			cmd.Printf("published to %s\n", url)
		}
	}
}
