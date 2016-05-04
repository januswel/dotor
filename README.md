# dotor v0.1.5

[![Circle CI](https://circleci.com/gh/januswel/dotor/tree/master.svg?style=shield)](https://circleci.com/gh/:user/:repo/tree/master)

dotfile setuper

## synopsis

- creates symbolic links to dotfiles in your home directory.

## install

Just get a newest binary for your os and arch and "dotor.yml" from [release page](https://github.com/januswel/dotor/releases).

## usage

### dotor.yml

See your "dotor.yml". If you need to add some dotfiles, write settings for them to your "dotor.yml".

### your dotfiles

Fetch your dotfiles out of your sources.

- from github
- from your other machines

```sh
# from github
mkdir -p /path/to/your/dotfiles
cd /path/to/your/dotfiles
git clone git@github.com:youraccount/dotfiles.git .
```

### dotor

Run dotor with some arguments.

```sh
# creates symlinks in keeping with a setting file
dotor_<youros>_<yourarch> dotor.yml /path/to/your/dotfiles
```

## future works

- dry run
- (may) creates symlinks for well-known dotfiles without setting files
- (may) deletes symbolic links made by dotor from your home directory.
