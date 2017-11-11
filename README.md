Waifu Bot
==========

A not-serious bot that welcomes new members to a server and compliments users when told to, 
based on a file containing a set of compliments to use.

### Setup
- Uses Golang 1.9.2
- Uses [discordgo v0.17.0](https://github.com/bwmarrin/discordgo)
- [Create a bot and client token associated with it](https://github.com/reactiflux/discord-irc/wiki/Creating-a-discord-bot-&-getting-a-token)
- Create a file called `discord_token.json` based off of the template and place your client token in there
- Supply the bot with a list of compliments called `compliments.txt`, where the compliments are on separate lines.

### Usage

Bot listens on the event that a new user joins a server (guild)

To interact with the bot, commands have a prefix of `!waifu`

`!waifu compliment me` will make the bot give a random compliment to the message sender.

`!waifu compliment @someone` will make the bot give a random compliment to the specified user. This assumes that the user exists
