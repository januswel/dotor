# dotor

dotfile setuper

## synopsis

- creates symbolic links to dotfiles in your home directory.

## usage

```sh
# clone your dotfiles
mkdir -p /path/to/your/dotfiles
cd /path/to/your/dotfiles
git clone git@github.com:youraccount/dotfiles.git .

# creates symlinks in keeping with a setting file
dotor dotorrc.yml /path/to/your/dotfiles
```

## future works

- dry run
- (may) creates symlinks for well-known dotfiles without setting files
- (may) deletes symbolic links made by dotor from your home directory.
