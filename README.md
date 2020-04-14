# Backpocket

Backpocket is a command line utility for storing a reading list
of articles from the Internet to your local disk for the future.

It's an alternative to [Pocket](https://getpocket.com/) offering all the
required features without handing over all your private reading materials to
any 3rd party. And it's totally free, offering even full-text search via other
tools.

Backpocket is based on the algorithm created by Mozilla used to create Firefox
[Readability.js](https://github.com/mozilla/readability) functionality.

![successful.png](examples/successful.png)

## Features

Backpocket aims to be as simple as possible while offering a lot of features:

* command-line tool for very easy usage;
* stores articles on a local disk (see [an example](https://raw.githack.com/jarmo/backpocket/master/examples/successful.html));
* images are stored as base64 data sources;
* supports all kind of formats - when not html, then will be stored AS IS;
* stored articles can be opened just with a browser or other tools depending on the format;
* since articles are stored in their native format, then can do whatever with them - perform full-text search, convert to text etc;
* articles which fail to be stored for some reason will still save a reference to them storing the original link (see [an example](https://raw.githack.com/jarmo/backpocket/master/examples/failed.html)).


## Installation

Download latest binary from [releases](https://github.com/jarmo/backpocket/releases), extract it and add it to somewhere in your **PATH**. That's it.

*Of course, you're free to compile your own version of binary to be 100% sure that it has not been tampered with, since this is an open-source project after all.*


## Usage

Using backpocket is simple too:

```sh
$ backpocket 'ARTICLE_URL'
```

To open archived article immediately from command-line:

```sh
$ open `backpocket 'ARTICLE_URL'`
```

To open any archived article later:

```sh
$ open examples/successful.html
```

Running backpocket will try to create a readable version of the article specified at `ARTICLE_URL` and stores it on a local disk.
When readable version cannot be created (for example, `ARTICLE_URL` points to an image) then that URL will be downloaded and stored AS IS.
You can configure storage dir by editing backpocket configuration file `config.json`, which is stored in a location specified by [XDG Base Directory standard](https://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html)).


## What about full text search?!

Easy!

```sh
$ grep -niR 'trump' $(backpocket path)
```


## Aliases for reading, archiving and searching

Easiest way to read oldest article would be to create alias for command line:

```sh
alias bp-read='open `ls -Adp $(backpocket path)/* | grep -v "/$" | head -1`'
```

Also alias for archival makes sense:

```sh
alias bp-archive='mkdir -p `backpocket path`/archive && mv -v `ls -Adp $(backpocket path)/* | grep -v "/$" | head -1` `backpocket path`/archive'
```

And create a function for search:

```sh
function bp-search() { grep -noiRE ".{0,70}$1.{0,70}" `backpocket path` }
```

And why not create an alias for the backpocket:

```sh
alias bp=backpocket
```

And now just use these aliases:

```sh
$ bp 'ARTICLE_URL'
$ bp-read
$ bp-archive
$ bp-search 'TERM'
```


## What about showing only text instead of HTML?

Just use html2text together with other UNIX tools:

```sh
$ cat examples/successful.html | sed 's/<img.*//g' | html2text | less
```


## Importing articles from Pocket

It is pretty easy and straightforward to import articles from Pocket too:

1. Install backpocket
2. Install [pup](https://github.com/EricChiang/pup) and [jq](https://stedolan.github.io/jq/);
3. [Export](https://getpocket.com/export) from Pocket to HTML file; 
4. Run the following command to import all archived articles where `ril_export.html` is the file Pocket exported:
```sh
cat ril_export.html | pup 'ul:last-of-type a json{}' | jq -r '.[] | "\(.href) \(.time_added)"' | tail -r | xargs -P4 -L1 backpocket
mkdir -p `backpocket path`/archive && mv $(ls -Adp `backpocket path`/* | grep -v "/$") `backpocket path`/archive
```
5. Run the following command to import all non-archived articles from Pocket:
```sh
cat ril_export.html | pup 'ul:first-of-type a json{}' | jq -r '.[] | "\(.href) \(.time_added)"' | tail -r | xargs -P4 -L1 backpocket
```

This all might take some time depending on the count of articles, size of
articles, speed of your internet connection and so on. Also, please not that
there's a high probability that many of the articles saved to the Pocket in the
past do not exist anymore in the Internet will be saved as failed articles - that's the problem with Internet and that's exactly the reason why archiving interesting/important things on your local disk makes sense in case you ever
want to return to any content.


## Backup, syncing, mobile support etc.

Backpocket does not offer anything else except storing articles for future use.
Use any 3rd party tools for backup, syncing or showing articles on mobile (for
example, serve articles using a static web server).


## Windows support

Since backpocket is written using Go and it has a wonderful cross-compile
support then it also works without any problems on Windows. However, if you
need to use any UNIX tools described in this README then use [WSL](https://docs.microsoft.com/en-us/windows/wsl/install-win10) or some
other Linux subsystem on Windows.

There are some tricks which need to be done under WSL Ubuntu so that aliases written above would work.

For reading:

```sh
alias bp-read='cmd.exe /c start $(wslpath -aw `ls -Adp $(backpocket path)/* | grep -v "/$" | head -1`)'
```

For archival:

```sh
alias bp-archive='mkdir -p $(backpocket path)/archive && mv -v `ls $(backpocket path)/* | grep -v "/$" | head -1` `backpocket path`/archive'
```
