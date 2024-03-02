<p align="center">
  <img src="assets/tget.png" />
<br>
<a href="http://makeapullrequest.com"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg"></a>
<a href="#Linux"><img src="https://img.shields.io/badge/os-linux-brightgreen">
<a href="#MacOS"><img src="https://img.shields.io/badge/os-mac-brightgreen">
<a href="#Android"><img src="https://img.shields.io/badge/os-android-brightgreen">
<a href="#Windows"><img src="https://img.shields.io/badge/os-windows-yellowgreen">
<a href="#iOS"><img src="https://img.shields.io/badge/os-ios-yellow">
<a href="#Steam-deck"><img src="https://img.shields.io/badge/os-steamdeck-yellow">
<br>
<a href="https://www.buymeacoffee.com/sweetbabyalaska"><img src="https://img.shields.io/badge/Buy%20Me%20a%20Coffee-ffdd00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=black"></a>
<a href="https://github.com/sweetbbak"><img src="https://img.shields.io/badge/creator-sweet-green"></a>
<br>
</p>

<p align="center">
<a href="#golang"><img src="https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white">
<a href="go"><img src="https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black">
<a href="linux"><img src="https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white">
<a href="bsd"><img src="https://img.shields.io/badge/-OpenBSD-%23FCC771?style=for-the-badge&logo=openbsd&logoColor=black">
<a href="mac"><img src="https://img.shields.io/badge/mac%20os-000000?style=for-the-badge&logo=macos&logoColor=F0F0F0">
</p>

<h3 align="center">
  wget but for torrents
</h3>

tget, simple torrent downloading cli

![example of toru in progress](assets/search.png)

## Table of Contents

## Install

<details closed>
  <summary>Install Go</summary>
  <a href="https://go.dev/doc/install">Install go</a>
  This project requires go 1.21.7 or higher.
</details>

```sh
go install github.com/sweetbbak/tget@latest
```

```sh
nix profile install github:sweetbbak/tget
```

<details closed>
  <summary>Build from source</summary>

```sh
git clone https://github.com/sweetbbak/tget.git && cd toru
go build -o tget .
```

you can also use the justfile

```sh
git clone https://github.com/sweetbbak/tget.git && cd tget
just
```

### Building for different platforms and architectures

Run to find your target architecture and platform:

```sh
go tool dist list
```

then use the environment variables `GOOS` and `GOARCH` before using
the build command.

Example:

```sh
GOOS=linux GOARCH=arm64 go build -o toru ./cmd/toru
```

</details>
## Examples

Search for an anime:

```sh
tget --torrent "magnet:..."
# or
tget "magnet:..."
```

## Support

Consider creating a PR, taking up a minor issue on the TODO list, leaving an issue to help improve functionality or buy
me a coffee!

![moe-visitor-counter](https://count.getloli.com/get/@sweetbbak?theme=asoul)
