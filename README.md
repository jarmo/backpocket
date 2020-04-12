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
* stores articles on a local disk (see [an example](examples/successful.html));
* images are stored as base64 data sources;
* supports all kind of formats - when not html, then will be stored AS IS;
* stored articles can be opened just with a browser or other tools depending on the format;
* since articles are stored in their native format, then can do whatever with them - perform full-text search, convert to text etc;
* articles which fail to be stored for some reason will still save a reference to them storing the original link (see [an example](examples/failed.html)).


## Installation

Download latest binary from [releases](https://github.com/jarmo/backpocket/releases), extract it and add it to somewhere in your **PATH**. That's it.

*Of course, you're free to compile your own version of binary to be 100% sure that it has not been tampered with, since this is an open-source project after all.*


## Usage

Using backpocket is simple too:

```
$ backpocket ARTICLE_URL
```

To open archived article immediately from command-line:

```
$ open `backpocket ARTICLE_URL`
```

To open any archived article later:

```
$ open examples/successful.html
```

Running backpocket will try to create a readable version of the article specified at `ARTICLE_URL` and stores it on a local disk.
When readable version cannot be created (for example, `ARTICLE_URL` points to an image) then that URL will be downloaded and stored AS IS.
You can configure storage dir by editing backpocket configuration file `config.json`, which is stored in a location specified by [XDG Base Directory standard](https://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html)).


## What about full text search?!

Easy!

```
$ grep -iR "trump" examples
```


## What about showing only text instead of HTML?

Just use html2text together with other UNIX tools:

```
$ cat examples/successful.html | sed 's/<img.*//g' | html2text | less
```


## Importing articles from Pocket

It is pretty easy and straightforward to import articles from Pocket too:

1. Install backpocket
2. Install [pup](https://github.com/EricChiang/pup) and [jq](https://stedolan.github.io/jq/);
3. [Export](https://getpocket.com/export) from Pocket to HTML file; 
4. Run the following command to import all archived articles where `ril_export.html` is the file Pocket exported:
```
cat ril_export.html | pup 'ul:last-of-type a json{}' | jq -r '.[] | "\(.href) \(.time_added)"' | tail -r | xargs -P4 -L1 backpocket
mkdir ~/backpocket/archive && mv ~/backpocket/*.* ~/backpocket/archive
```
5. Run the following command to import all non-archived articles from Pocket:
```
cat ril_export.html | pup 'ul:first-of-type a json{}' | jq -r '.[] | "\(.href) \(.time_added)"' | tail -r | xargs -P4 -L1 backpocket
```

This all might take some time depending on the count of articles, size of
articles, speed of your internet connection and so on. Also, please not that
there's a high probability that many of the articles saved to the Pocket in the
past do not exist anymore in the Internet will be saved as failed articles - that's the problem with Internet and that's exactly the reason why archiving interesting/important things on your local disk makes sense in case you ever
want to return to any content.


## Backup, syncing, mobile support etc.

Backpocket does not offer anything else except storing articles for future use.


## Windows support

Since backpocket is written using Go and it has a wonderful cross-compile
support then it also works without any problems on Windows. However, if you
need to use any UNIX tools described in this README then use [WSL](https://docs.microsoft.com/en-us/windows/wsl/install-win10) or some
other Linux subsystem on Windows.
