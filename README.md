# `deta-drive` - A command line companion for Deta Drive

`deta-drive` is a command line companion for Deta Drive. It allows you to manage your Deta Drive files from the command line.

## Installation

```bash
go install github.com/pomdtr/deta-drive@latest
```

## Usage

The cli supports most of the classic unix commands you are used to. Use `deta-drive help` to get a list of all commands.

## Examples

### Upload a file

```bash
deta-drive cp ./file.txt deta://my-drive/file.txt
```

### Download a file

```bash
deta-drive cp deta://my-drive/file.txt ./file.txt
```

### List files

```bash
deta-drive ls deta://my-drive/
```

### Print file content

```bash
deta-drive cat deta://my-drive/file.txt
```

### Remove a file

```bash
deta-drive rm deta://my-drive/file.txt
```
