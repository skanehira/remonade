![GitHub Repo stars](https://img.shields.io/github/stars/skanehira/remonade?style=social)
![GitHub](https://img.shields.io/github/license/skanehira/remonade)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/skanehira/remonade)
![GitHub all releases](https://img.shields.io/github/downloads/skanehira/remonade/total)
![GitHub CI Status](https://img.shields.io/github/workflow/status/skanehira/remonade/ci?label=CI)
![GitHub Release Status](https://img.shields.io/github/workflow/status/skanehira/remonade/Release?label=release)

# remonade - UNDER DEVELOPMENT
Unofficial Nature Remo CLI.

![](https://i.gyazo.com/5c2c9c5979368be9ee89a0521166104b.png)

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

| Panel      | Key | Description |
|------------|-----|-------------|
| Common     | `j` | move down   |
| Common     | `k` | move up     |
| Common     | `h` | move left   |
| Common     | `l` | move right  |
| Appliances | `u` | Power on    |
| Appliances | `d` | Power off   |


## Author
skanehira
