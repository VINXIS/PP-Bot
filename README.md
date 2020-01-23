# PP-Bot
A bot for using osu-tools on discord

## Installation
1. [Install golang](https://golang.org/doc/install). 
2. Clone the repository using `git clone --recurse-submodules https://github.com/VINXIS/PP-Bot.git`.
3. Install the dependencies with `go get`.
4. Duplicate `config.example.json` and call it `config.json`. Fill it in.
5. Invite the bot to your server by replacing `PUT_CLIENT_ID_HERE` in the URL below with the discord application's client ID obtained here [here](https://discordapp.com/developers/applications). https://discordapp.com/api/oauth2/authorize?client_id=PUT_CLIENT_ID_HERE&permissions=536870912&scope=bot.
7. Run the program by running `go build -o bot bot.go` and then `./bot` in your instance / computer.