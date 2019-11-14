# Golang Mod version

Simple program to receive latest version (for go.mod) of passed package.

## Install

```
go get -u github.com/im-kulikov/modver
```  

## Usage

```
$ modver /path/to/local/repo
go.mod version: v0.2.0

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
```
