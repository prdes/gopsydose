# Welcome traveller!

I hope you find this interesting!

## What's the purpose of this project?

It's an attempt at making a Go module and application for enabling
a user to log their drug dosages and timings.
It uses the logged information,
along with fetched and stored locally information about a particular
substance, to make assumptions about where things might be going with
your experience.

The main purpose is harm reduction.

This project does NOT aim the endorsement of any psychoactive substance.
If it prevents a person from consuming more, then it's done it's job.

Showing logs to others as a bragging point makes you look like a fool.

## Source Information

You can find any info present in the "source" that's configured.

Currently the default source is [psychonautwiki](https://psychonautwiki.org).
This means the info present there, that's available through their API, will be
downloaded to your machine and stored in a local database. Later when you log
new dosages or want information about a substance, it will use the local
database, not the API, if it's already fetched the information.

You don't need to configure anything for the database or API by default.

Keep in mind, if you store data locally for too long, it might
be obsolete. You shouldn't clear your local db too quickly as well,
because this way you put more strain on the API servers for no reason and
you might lose data, which you might need in case you don't have an Internet
connection. It's also a lot slower to get data using the API compared to
the local database.

## Examples

If you want to log a dose:

`gopsydose -drug alcohol -route oral -dose 355 -units ml -perc 4.5`

or

`gopsydose -drug weed -route smoked -dose 100 -units "mg (THC)"`

Since both of these aren't consumed at once, there is a command to mark
when the dosing has ended.

`gopsydose -change-log -end-time now`

This will set when you finished your dose for the last log. You can use an unix timestamp like in the example below instead of `now`.

To change the start time of dose use: `gopsydose -change-log -start-time 1655443322`

Changing the start time also changes the "id" of a dose. This means, if you're looking for the last dose and you've changed the start time to an earlier moment, it will get pushed back in the list.

You can set the times for a specific id by using the `-for-id` command like so:

`gopsydose -change-log -end-time now -for-id 1655443322`

This works for both times.

If you're consuming something at once like
[LSD](https://en.wikipedia.org/wiki/Lysergic_acid_diethylamide) or
[Psilocybin mushrooms](https://en.wikipedia.org/wiki/Psilocybin_mushroom) or
anything else, there's no need for the
`-change-log -end-time` command. Just continue without doing it.

To see the newest dose only: `gopsydose -get-new-logs 1`

To see all dosages: `gopsydose -get-logs`

To search the dosages: `gopsydose -get-logs -search "whatever"`

To see the progress of your last dosage: `gopsydose -get-times`

You can combine `-get-logs` or `-get-times` with: `-for-id`

to get information for a specific ID.

You can see all IDs using: `gopsydose -get-logs`

If you want a log to be remembered and only set the dose for the next log:

`gopsydose -remember -drug weed -route smoked -dose 100 -units "mg (THC)"`

Remembers the config and then for the next dose: `gopsydose -dose 100`

Forgetting the last config: `gopsydose -forget`

If you're running Linux or another UNIX-like OS with GNU Watch,
you can do:

`watch -n 60 gopsydose -get-times`

This will run the command every minute and show you
the latest results.

There is a limit set in a config file about how many dosages you can do,
the default is 100, when the limit is reached it will not allow anymore,
but you can clean the logs like so: `gopsydose -clean-logs`

To clean dosages with a search: `gopsydose -clean-logs -search "whatever"`

This will clean only dosages which contain the "whatever" string.

No need to delete all logs, you can delete X number of the oldest logs
like so, for example to delete 3 of the oldest logs:

`gopsydose -clean-old-logs 3`

You can do the command like so: `gopsydose -clean-old-logs 1 -for-id 1655144869`

to remove a specific ID, works with `-clean-new-logs 1` as well.

After logging you can change the data of a log using `-change-log` it works for:

`-start-time` ; `-end-time` ; `-drug` ; `-dose` ; `-units` ; `-route`

So for example for changing the dose you would do:

`gopsydose -change-log -dose 123`

This will change the dose for the last log, to change for a specific log do:

`gopsydose -change-log -dose 123 -for-id 1655144869`

To see where your config files and database file are:

`gopsydose -get-dirs`

Checkout [this](#configs-explained) section for more info on configs!

If you're paranoid, to clean the whole database: `gopsydose -clean-db`

Also don't forget, if you're using sqlite, which is the default, you can always
do: `gopsydose -get-dirs`

Then you can delete the database (db) file itself manually.

## Security/Privacy

The issue is currently no files are encrypted and can't be 
until a proper implementation is done, also since
by default we're fetching drug information using the psychonautwiki API,
it would be wise not to spam their servers too much.
We store all information locally on first fetch and use only the local info
later for everything. This way even if the Internet goes down, the logger
can still be used properly and the API servers can relax.

For any more info, again: `gopsydose -help`

## Alias

You can use the `example.alias` file to setup an easier to use environment
for your system. The difference between this and the `-remember` command is
that, shell aliases are system wide and from the operating system. The
`-remember` command stores the "command" in the database, allowing even remote
use of those settings, which is OS-independent. Aliases are a lot more flexible
currently, so use them when you need them!

## Current status

The project is in a very early phase, making the very first steps
at being something remotely usable. If you feel like it, give it a try
and please if you have any ideas, concerns, etc. write in the Issues.

If you don't want to build from source, you can checkout the "Releases"
section. There should be a "snapshot" tag with archives attached to it.

Pull requests are very welcome as well and there currently aren't any
guidelines, just try sticking to the coding style and use `gofmt`!
It would be nice if single lines stay <= 120 characters.

Using this as a module now is a bad idea since the API will most
likely change a lot until the first stable tag.

## Installing Go

Check whether you already have Go, by typing `go version` in your
terminal. If information about Go shows up and you're on Go version 1.18+,
you're good to go. Otherwise remove the old installation and reinstall
with a newer version.

You need to download Go:
* If you're on Windows, you need to download Go
from the official website [here](https://go.dev/).
* If you're on Linux, use your package manager!

[Skip to GCC](#installing-gcc)

### Go Linux Installation

Works on OpenSUSE Tumbleweed. Leap hasn't been tested
and seems like it won't work there.

OpenSUSE: `sudo zypper install go`

Arch Linux: `sudo pacman -S go`

Fedora and Debian might be a bit out of date.

If for Debian the version is old, you need to use backports or sid.
Do this at your own risk, since this is not the normal way of installing
packages.

If for Fedora the version is old, you need to use a package from
a newer release, Rawhide or a third party repository.
The same warning applies as for Debian.

Fedora: `sudo dnf install golang`

Debian: `sudo apt install golang-go`

## Installing GCC

Check whether you already have GCC, by typing `gcc --version` in your
terminal. If information about GCC shows up and you're on GCC version 10+,
you're good to go. Otherwise remove the old installation and reinstall
with a newer version.

You need to download GCC as well:
* On Windows you can download
[tdm-gcc](https://jmeubank.github.io/tdm-gcc/download/).
* On Linux use your package manager!

[Skip to main installation](#installing-this-project)

### GCC Windows installation

The tdm-gcc GUI installer defaults are enough.

### GCC Linux installation

Same info from [Go section](#go-linux-installation)
about Leap, Debian and Fedora versions applies here.

OpenSUSE: `sudo zypper install gcc`

Arch Linux: `sudo pacman -S gcc`

Fedora: `sudo dnf install gcc`

Debian: `sudo apt install gcc-10` or `gcc-11`

## Installing this project

Then you need to run in a terminal:

`go install github.com/psybits/gopsydose@latest`

Afterwards you can do a quick test with: `gopsydose -help`

There are a lot of commands and currently it isn't very clear, but
a few examples and explanations will be left in this document.

## Viewing/Editing the database

If you wish to view/edit the database manually you can
use the [SQLite Browser](https://sqlitebrowser.org/dl/).

Get the database directory using: `gopsydose -get-dirs`

## Configs explained

### gpd-settings.toml

#### MaxLogsPerUser
How many logs a user can make. You can log as a different user using the
`-user` command.

#### UseAPI
The name of the API set in `gpd-sources.toml`. The API needs to have an
implementation, currently only psychonautwiki has one.

#### AutoFetch
Whether to fetch info from an API when logging. If set to false the
source table needs to be manually filled using other tools.

#### DBDir
This is relevant to sqlite.

The directory where the database file will be created. If changing this,
don't forget to check the old directory for any files left!

#### AutoRemove
If set to true, will remove the oldest log when adding a new one, if the
`MaxLogsPerUser` limit is reached, without telling the user.

#### DBDriver
Which database driver to use. Current options are "sqlite3" or "mysql".

#### MySQLAccess
Credentials to access a MySQL database.

Example: user:password@tcp(127.0.0.1:3306)/database

#### VerbosePrinting
If set to true, functions which print more verbose information, will print it.

#### PrefixPrinting
If set to true, all print function wrappers add the project's name and the name
of the code from which the printing is done.

This is mostly useful when importing the module into another project to
have a better distinction between different text information.
It can help with debugging a bit as well, but don't rely on it too much!

#### Timezone
This by default is "Local", which means the code tries to figure out the local
time zone using the operating system settings. If this fails, you can use
the other possible strings from below.

You can change it to something like "Europe/Paris" to change the time zone to
the one in that area.

All info about this string can be found here:
https://pkg.go.dev/time#LoadLocation

### gpd-sources.toml

```
[NameOfTheApi]
ApiUrl = 'api.whereitis.com'
```

An implementation needs to be present in the code for the name here.

