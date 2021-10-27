# discordcmd 

A opinionated wrapper around [discordgo](https://github.com/bwmarrin/discordgo) aimed at managing and providing slash commands as first class citizens. 

The core commands struct, registers the commands added into it into all the servers that the bot is a member of, and un-registers them on app close. The allows the user to simply write a command object.s

