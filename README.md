# gmc
Memcached like server in golang

## Supported protocol functions

#### Memcached protocol
| Function |      Status |
|----------|:-------------:|
| get |  finished   |
| set |  finished   |
| add |  finished   |
| replace |  finished   |
| delete |  finished   |
| gat |  finished   |
| gats |  finished   |
| flush_all | finished |
| stats |  not started   |
| cas |  not started   |
| incr  |  finished   |
| decr   |  finished   |
| append | not started |
| prepend | not started |

#### Custom functions
| Function |      Status | What it does |
|----------|:-------------:|----------:|
| has |  finished   | returns if a key exist without all the payload |

## Building
#### Source
Install the latest version of golang
```
> git clone git@github.com:BlizzTrack/gmc.git
> cd gmc
> go build -o gmc ./cmd/main.go
```

#### Go Get
```
> go get -u github.com/BlizzTrack/gmc/cmd/gmc
```

## Binary builds
-- Coming Soon --

## Who use this
[BlizzTrack](https://www.blizztrack.com)
[Badge Directory](https://www.badgedirectory.com)
