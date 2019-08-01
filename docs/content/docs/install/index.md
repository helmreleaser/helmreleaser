---
title: 'Install'
date: 2019-02-11T19:27:37+10:00
weight: 2
---

To install HelmReleaser on your workstation, use one of the packaging tools below, or download and compile from source.

## Homebrew (MacOS)

```shell
$ brew install helmreleaser/tag/helmreleaser
```

## Snapcraft

```shell
$ sudop snap install --classic helmreleaser
```

## Compiling from source

```shell
git clone https://github.com/helmreleaser/helmrelease
cd helmreleaser
make
./helmreleaser --versionA
```
