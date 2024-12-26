# nadge

a nostr badge client on cli written in golang for fun!

# how to install?

```
go install github.com/kehiy/nadge@latest
```

# how to use?

you can use `issue` command to issue new badge.

example:

```sh
nadge issue ./badge.json --nsec="nsec163yk7pd3gx3m59exje5c72tmsdeaatdesejjevumv9zw0e7z8wpqjll76r" --relays="wss://jellyfish.land,wss://nos.lol" 
```
> [example template file](/badge_issue_template.json)

you can use list command to get all badges issued by someone:

```sh
nadge list --npub="npub10q6ut93r6c7d3xxvea8nzuch5d80kevwrhf5ucw0tj7xkzjq765qd4test"
```

you can award a badge like this:

```sh
nadge award "*if#s0G" npub1h49w8en79xty6j2pwgnpm3znjhyf767jua6xgt3kvyn3w80ms86s2z9kay --nsec="nsec163yk7pd3gx3m59exje5c72tmsdeaatdesejjevumv9zw0e7z8wpqjll76r" --relays="wss://jellyfish.land,wss://nos.lol"
```

first argument is badge unique id, second one is awardee npub.


# license

this software is published under [mit license](./LICENSE)
