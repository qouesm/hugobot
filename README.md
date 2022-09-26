# hugobot

## A Discord bot for the New Paltz Computer Science server 

hugobot is designed to aid in management of the NPCS server

- Role management
- Mass role removal
- Notification of new member nickname changes

hugobot is written with the [hikari](https://github.com/hikari-py/hikari) library plus [hikari-lightbulb](https://github.com/tandemdude/hikari-lightbulb) as a command handler.

## Quickstart

### Requirements

- [python](https://www.python.org/)>=3.10

### Installation

- A bot token will need to be provides in `/.env/bot_token`. Create this file

- Install the python dependancies from `requirements.txt`:
```sh
pip install -r requirements.txt
```

- Start the bot:
```sh
python -OOm app
```

### Optional for developers

If you would like to restart the bot when file changes are detected, [nodemon](https://nodemon.io/) has been configured to do so. This requires [npm](https://nodejs.org/en/) to be installed.

- Install node packages:
```sh
npm i
```

- To run normally:
```sh
npm run start
```

- To reload on file changes:
```sh
npm run dev
```

