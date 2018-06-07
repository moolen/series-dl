# series-dl

Batch download series. inspired by [github.com/alexander-schoch/scripts](https://github.com/alexander-schoch/scripts/tree/master/series-stream)

## install

if you have a go toolchain installed: `go get -u github.com/moolen/series-dl`

You need to have `youtube-dl` installed. Having `phantomjs` is strongly recommended.

For archlinux users:

``` sh
$ yaourt -S phantomjs-bin youtube-dl
```

## build

```sh
$ make

```

## usage

``` sh
# search for a series slug
$ ./series-dl -search house
INFO[0000] house-of-cards-us
INFO[0000] mickey-mouse-clubhouse
INFO[0000] little-house-on-the-prairie
INFO[0000] the-real-housewives-of-potomac
INFO[0000] the-real-housewives-of-new-jersey
INFO[0000] the-real-housewives-of-orange-county
[...]

# download first two seasons
$ ./series-dl -series house-of-cards-us -season-start 1 -season-end 2 -concurrency 4
INFO[0000] found season: Season 1
INFO[0000] found season: Season 2
INFO[0000] found season: Season 3
INFO[0000] found season: Season 4
INFO[0000] found season: Season 5
INFO[0000] found season: Season 0
INFO[0000] found season: Season 6
INFO[0000] found episode: S1E1 Chapter 1
INFO[0000] found episode: S1E2 Chapter 2
INFO[0000] found episode: S1E3 Chapter 3
INFO[0000] found episode: S1E4 Chapter 4
INFO[0000] found episode: S1E5 Chapter 5
INFO[0000] found episode: S1E6 Chapter 6
INFO[0000] found episode: S1E7 Chapter 7
INFO[0000] found episode: S1E8 Chapter 8
INFO[0000] found episode: S1E9 Chapter 9
INFO[0000] found episode: S1E10 Chapter 10
INFO[0000] found episode: S1E11 Chapter 11
INFO[0000] found episode: S1E12 Chapter 12
INFO[0000] found episode: S1E13 Chapter 13
INFO[0000] found episode: S2E1 Chapter 14
INFO[0000] found episode: S2E2 Chapter 15
INFO[0000] found episode: S2E3 Chapter 16
[...]
INFO[0001] downloading S1E1 Chapter 1
INFO[0001] downloading S1E2 Chapter 2
INFO[0002] downloading S1E3 Chapter 3
INFO[0002] trying link [redacted]
INFO[0002] found watch button link: [redacted]
INFO[0002] downloading S1E4 Chapter 4
INFO[0002] trying link [redacted]
INFO[0003] found watch button link: [redacted]
INFO[0003] trying link [redacted]
INFO[0003] found watch button link: [redacted]
INFO[0003] trying link [redacted]
INFO[0003] found watch button link: [redacted]
INFO[0004] S1E1 [0.0%] [Unknown]
INFO[0004] S1E1 [0.0%] [Unknown]
INFO[0004] S1E1 [0.0%] [Unknown]
INFO[0004] S1E1 [0.0%] [12.32MiB/s]
INFO[0004] S1E1 [0.0%] [19.61MiB/s]
INFO[0005] S1E1 [0.0%] [1.44MiB/s]
INFO[0005] S1E1 [0.0%] [1.23MiB/s]
INFO[0005] S1E1 [0.1%] [1.85MiB/s]
INFO[0005] S1E1 [0.1%] [2.53MiB/s]
INFO[0005] S1E1 [0.2%] [3.57MiB/s]
```

## Known Issues

The naming of the downloaded files is a little bit off. This may be resolved using `youtube-dl --output <output-template>`.