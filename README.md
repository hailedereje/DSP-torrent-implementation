# Torrent client

A basic bittorrent client written in Go. To anyone who wants to try to make their own, take a look at this awesome blog post to get started: https://blog.jse.li/posts/torrent/. This is a great project to get an idea of how http, binary operations and channels work in Go.

## Usage

If you have the project locally, you can use

```sh
go run . <path/to/torrent/file> <path/to/store/downloaded/file>

eg. go run test-torrent\deb.iso.torrent download\
```

## TODO

- Add support for magnet links
- Multiple file support
- Seeding
- Store intermediate downloaded parts to disk instead of memory
