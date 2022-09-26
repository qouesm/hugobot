__version__ = "0.0.1"

with open("./.env/bot_token") as f:
    BOT_TOKEN = f.read().strip()

# from os import environ
# BOT_TOKEN = environ['BOT']

GUILD_ID = 1023048878301384804
