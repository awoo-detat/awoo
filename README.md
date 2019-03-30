# Werewolf
## Or: Awoooo!

## What this is
A client/server Werewolf app. It simulates the experience of playing with cards, but instead of cards you use things that cost $600+ each.

You start the server (this repo/file), then everyone goes to the website it creates on their phone, tablet, computer, etc. From there you play, with each person submitting votes, seeing their role, etc on their own device!

## How to use it
Right now we have no releases. Sorry. That's coming in [#26](https://github.com/awoo-detat/awoo/issues/26).

Run the executable (`awoo.exe` on Windows, `awoo` on Mac/Linux). It will create a log file named `werewolf.log` in the directory you ran it from.

People can go to the IP of the computer it was run from, port 42300: ie `10.0.1.4:42300`

Play!

## Why 42300?
`4 23 0 0`

`4 W 0 0`

`A W O O`

## How to build it (if you want to change the code)
The simplest, though most time consuming, option is to run:

```
go get github.com/jteeuwen/go-bindata/...
go get github.com/elazarl/go-bindata-assetfs/...
npm i -g elm@0.18
make
```

That will run the testsuite, rebuild the assets and compile the code into the `./awoo` binary.

If you're only editing the go code and want to test your changes:

	go run bindata.go main.go

If you've made any changes to the static files (html, js, css), before you compile:

	make assets
