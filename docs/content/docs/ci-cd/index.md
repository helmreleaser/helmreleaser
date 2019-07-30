---
title: 'CI / CD Workflows'
date: 2019-02-11T19:27:37+10:00
draft: false
weight: 3
---


### CircleCI

```yaml
# .circleci/config.yml

version: 2
jobs:
  release:
    docker:
      - image: circleci/golang:1.10
    steps:
      - checkout
      - run: curl -sL https://git.io/helmreleaser | bash
workflows:
  version: 2
  release:
    jobs:
      - release:
          filters:
            branches:
              ignore: /.*/
            tags:
              only: /v[0-9]+(\.[0-9]+)*(-.*)*/
```
