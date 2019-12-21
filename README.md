# Golang Mod version

Simple program to receive latest version (for go.mod) of passed package.

## Install

```
go get -u github.com/im-kulikov/modver
```  

## Usage

*Help:*
```
Simple program to receive latest version (for go.mod) of passed package.
Version v0.0.5 (2019-12-21T13:42:08).
modver [global options] {repo-path or url}
      --verbose         verbose
  -h, --help            show help
  -v, --version         show version
      --commit          display latest commit version (for example v0.0.0-<hash>-<date>)
      --branch string   use passed branch to receive version (for remote repos only) (default "master")
```

*Usage:*
```
$ modver /path/to/local/repo
go.mod version: v0.2.0

$ modver --commit /path/to/local/repo
go.mod version: v0.0.0-20191113131239-3f7fc0db5b05

$ modver /path/to/local/repo/without/tags
go.mod version: v0.0.0-20191113131239-3f7fc0db5b05

$ modver /path/to/local/repo/without/references
reference not found

$ modver /path/to/not/repo
repository does not exist

$ modver github.com/im-kulikov/helium
Enumerating objects: 553, done.
Counting objects: 100% (553/553), done.
Compressing objects: 100% (275/275), done.
Total 553 (delta 272), reused 511 (delta 247), pack-reused 0
go.mod version: v0.12.2

$ modver git@github.com:im-kulikov/helium.git
Enumerating objects: 553, done.
Counting objects: 100% (553/553), done.
Compressing objects: 100% (275/275), done.
Total 553 (delta 272), reused 511 (delta 247), pack-reused 0
go.mod version: v0.12.2

$ modver --branch dev.cc  github.com/golang/go
Enumerating objects: 80356, done.
Counting objects: 100% (80356/80356), done.
Compressing objects: 100% (38577/38577), done.
Total 80356 (delta 57743), reused 57903 (delta 40144), pack-reused 0
go.mod version: go1.12.13

$ modver --branch dev.cc --commit  github.com/golang/go
Enumerating objects: 80356, done.
Counting objects: 100% (80356/80356), done.
Compressing objects: 100% (38577/38577), done.
Total 80356 (delta 57743), reused 57903 (delta 40144), pack-reused 0
go.mod version: v0.0.0-20150223214927-a91c2e0d2d19
```
