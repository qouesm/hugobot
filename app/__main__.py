import os

import hikari
import lightbulb

from app import BOT_TOKEN, GUILD_ID


def create_bot() -> lightbulb.BotApp:
    # Create the main bot instance with all intents.
    bot = lightbulb.BotApp(
        token=BOT_TOKEN,
        intents=hikari.Intents.ALL,
        default_enabled_guilds=GUILD_ID,
    )

    bot.load_extensions_from("./app/commands")

    return bot


if __name__ == "__main__":
    if os.name != "nt":
        # uvloop is only available on UNIX systems, but instead of
        # coding for the OS, we include this if statement to make life
        # easier.
        import uvloop

        uvloop.install()

    create_bot().run()
