package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/spf13/cobra"
)

func buildAwardCmd(parentCmd *cobra.Command) {
	listCmd := &cobra.Command{
		Use:   "award",
		Short: "award a badge to a new npub",
	}

	parentCmd.AddCommand(listCmd)

	listCmd.Run = func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.PrintErr("please input a npub and a badge id")
			os.Exit(1)
		}

		badgeID := args[0]
		_, v, err := nip19.Decode(args[1])
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		awardeeNpub := v.(string)

		relays := []nostr.Relay{}
		for _, r := range *relayAddrsOpt {
			relay, err := nostr.RelayConnect(context.Background(), r)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			relays = append(relays, *relay)
		}

		_, nv, err := nip19.Decode(*nsecOpt)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		sk := nv.(string)

		pk, err := nostr.GetPublicKey(sk)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		filters := []nostr.Filter{{
			Kinds:   []int{nostr.KindBadgeAward},
			Authors: []string{pk},
			Tags: nostr.TagMap{
				"#a": []string{fmt.Sprintf("30009:%s:%s", pk, badgeID)},
			},
		}}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		e := &nostr.Event{
			PubKey:    pk,
			CreatedAt: nostr.Timestamp(time.Now().Unix()),
			Kind:      nostr.KindBadgeAward,
			Tags: nostr.Tags{
				{"a", fmt.Sprintf("30009:%s:%s", pk, badgeID)},
			},
		}
	find:
		for _, r := range relays {
			sub, err := r.Subscribe(ctx, filters)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			for ev := range sub.Events {
				e = ev

				break find
			}
		}

		e.Tags = append(e.Tags, nostr.Tag{"p", awardeeNpub})

		if err := e.Sign(sk); err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		cmd.Println("Badge Award Created successfully!")

		ej, err := e.MarshalJSON()
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		cmd.Println("Event:", string(ej))

		cmd.Println("Publishing to relays...")

		for _, url := range *relayAddrsOpt {
			relay, err := nostr.RelayConnect(ctx, url)
			if err != nil {
				cmd.Println(err)
				continue
			}

			if err := relay.Publish(ctx, *e); err != nil {
				cmd.Println(err)
				continue
			}

			cmd.Printf("published to %s\n", url)
		}
	}
}
