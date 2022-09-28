import hikari
import lightbulb
import miru
from lightbulb import commands


class TestComponents(miru.View):
    @miru.button(label="Rock", emoji=chr(129704), style=hikari.ButtonStyle.PRIMARY)
    async def rock_button(self, button: miru.Button, ctx: miru.Context) -> None:
        await ctx.respond("Paper!")

    @miru.button(label="Paper", emoji=chr(128220), style=hikari.ButtonStyle.PRIMARY)
    async def paper_button(self, button: miru.Button, ctx: miru.Context) -> None:
        await ctx.respond("Scissors!")

    @miru.button(label="Scissors", emoji=chr(9986), style=hikari.ButtonStyle.PRIMARY)
    async def scissors_button(self, button: miru.Button, ctx: miru.Context):
        await ctx.respond("Rock!")

    @miru.button(emoji=chr(9209), style=hikari.ButtonStyle.DANGER, row=2)
    async def stop_button(self, button: miru.Button, ctx: miru.Context):
        self.stop()  # Stop listening for interactions


@lightbulb.command('role', 'Commands to manage roles')
@lightbulb.implements(commands.SlashCommandGroup)
async def role(ctx: lightbulb.Context) -> None:
    await ctx.respond('called role (pass)')


@role.child
@lightbulb.command('init', 'Initialize a menu to a text chat')
@lightbulb.implements(commands.SlashSubCommand)
async def init(ctx: lightbulb.Context) -> None:
    await ctx.respond('called role init (pass)')


@role.child
@lightbulb.option('name', 'Name of the role menu to add', 'str')
@lightbulb.command('add', 'Add a role menu (up to 25 per text chat)')
@lightbulb.implements(commands.SlashSubCommand)
async def add(ctx: lightbulb.Context) -> None:
    name = ctx.options.name
    await ctx.respond(f'called role add (pass)\n`{name=}`')


@role.child
@lightbulb.option('name', 'Name of the role menu to delete', 'str')
@lightbulb.command('delete', 'Delete a role menu')
@lightbulb.implements(commands.SlashSubCommand)
async def delete(ctx: lightbulb.Context) -> None:
    name = ctx.options.name
    await ctx.respond(f'called role delete (pass)\n`{name=}`')


@role.child
@lightbulb.option('name', 'Name of the role menu to modify', 'str')
@lightbulb.command('edit', 'Edit a role menu')
@lightbulb.implements(commands.SlashSubCommand)
async def edit(ctx: lightbulb.Context) -> None:
    name = ctx.options.name
    await ctx.respond(f'called role edit (pass)\n`{name=}`')


@role.child
@lightbulb.command('test', 'testing embeds currently')
@lightbulb.implements(commands.SlashSubCommand)
async def test(ctx: lightbulb.Context) -> None:
    embed = (
        hikari.Embed(
            title='Title',
            description='Use the [menu|buttons] to pick a role type',
            color=hikari.Color(0xff00ff),
        )
        # \u200b: zero width space
        .add_field('\u200b', ':zero:\n:one:\n:two:\n:three:\n:four:', inline=True)
        # .add_field('\u200b', '\u200b', inline=True)
        .add_field('Role Type', 'one\ntwo\nthree\nfour\nfive', inline=True)
        # .set_thumbnail('https://i.imgur.com/EpuEOXC.jpg')
        # .set_image('https://i.imgur.com/EpuEOXC.jpg')
        # .set_footer('This is the footer')
    )
    await ctx.respond(embed)


plugin = lightbulb.Plugin('TestPlugin')


@role.child
@lightbulb.command('button', 'testing buttons')
@lightbulb.implements(commands.SlashSubCommand)
async def button(ctx: lightbulb.Context) -> None:
    # @plugin.listener(hikari.GuildReactionAddEvent)
    async def some_listener(event: hikari.GuildReactionAddEvent) -> None:
        print('now listening')
        await plugin.rest.create_message(
            channel=1023049089379741787,
            content='detected reaction add event',
        )
    plugin.listener(hikari.GuildReactionAddEvent, some_listener)
    await ctx.respond('now listening (theoretically)')

    # view = TestComponents(timeout=60)  # Create a new view
    # message = await event.message.respond("Rock Paper Scissors!", components=view.build())
    # view.start(message)  # Start listening for interactions
    # await view.wait() # Wait until the view times out or gets stopped
    # await event.message.respond("Thank you for playing!")


def load(bot: lightbulb.BotApp) -> None:
    bot.command(role)


def unload(bot: lightbulb.BotApp) -> None:
    bot.remove_command(bot.get_slash_command('role'))
