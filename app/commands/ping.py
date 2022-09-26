import lightbulb
from lightbulb import commands


@lightbulb.option('a', "Use 'a' and 'b'", bool, default=False)
@lightbulb.command('ping', 'Roll one or more dice.')
@lightbulb.implements(commands.SlashCommand)
async def ping(ctx: lightbulb.context.Context) -> None:
    a = ctx.options.a
    print(f'{a=}')

    if a:
        await ctx.respond('B')
    else:
        await ctx.respond('pong')


def load(bot: lightbulb.BotApp) -> None:
    bot.command(ping)


def unload(bot: lightbulb.BotApp) -> None:
    bot.remove_command(bot.get_slash_command('ping'))
