![GitHub Repo stars](https://img.shields.io/github/stars/skanehira/remonade?style=social)
![GitHub](https://img.shields.io/github/license/skanehira/remonade)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/skanehira/remonade)
![GitHub all releases](https://img.shields.io/github/downloads/skanehira/remonade/total)
![GitHub CI Status](https://img.shields.io/github/workflow/status/skanehira/remonade/ci?label=CI)
![GitHub Release Status](https://img.shields.io/github/workflow/status/skanehira/remonade/Release?label=release)

# <img src="https://i.gyazo.com/85e13d8198dcb843ece467cad46350e7.png" width="30"/> remonade - Unofficial Nature Remo CLI
**UNDER DEVELOPMENT**
![](https://i.gyazo.com/e1e0e0e34c51b1bf1894bbd26a3f442b.png)

## Installation

```sh
$ go install github.com/skanehira/remonade@latest
```

## Usage
At first, you must generate token from [home.nature.global](https://home.nature.global).
And then, you can setup token to configuration file with run `remonade init`.

```sh
# setup your token
$ remonade init

# edit your config
$ remonade edit

# run
$ remonade
```

### Key maps

| Panel           | Key                | Description            |
|-----------------|--------------------|------------------------|
| Common          | `j`                | Move down              |
| Common          | `k`                | Move up                |
| Common          | `h`                | Move left              |
| Common          | `l`                | Move right             |
| Common          | `Ctrl+n`           | Next panel             |
| Common          | `Ctrl+p`           | Previous panel         |
| Appliances      | `u`                | Power on               |
| Appliances      | `d`                | Power off              |
| Appliances      | `o`                | Open settings          |
| AirCon Settings | `q`, `c`           | Close                  |
| AirCon Settings | `Ctrl+n`, `Ctrl+j` | Next item              |
| AirCon Settings | `Ctrl+p`, `Ctrl+k` | Previous item          |
| AirCon Settings | `Enter`, `j`, `k`  | Change value           |
| Light Settings  | `Enter`            | Apply button or signal |


## Author
skanehira
