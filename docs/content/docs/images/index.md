---
title: 'Images'
date: 2019-02-11T19:30:08+10:00
draft: false
weight: 5
---

HelmReleaser rewrite image references (repositories and tags) in the `values.yaml` at build time. This is designed to support a [semver tagging strategy](https://medium.com/@mccode/using-semantic-versioning-for-docker-image-tags-dfde8be06699) that matches the git release tag.

To identify the images, the helmreleaser.yaml configuration includes an `images` key. In this section, you can identify where your images are defined in the `values.yaml`, and define template strings that will be used to rewrite these at build time.

Given a `values.yaml` that identifies an image:

```yaml
# values.yaml

image: "docker.elastic.co/elasticsearch/elasticsearch"
imageTag: "latest"
```

This tag can be identified in a `helmreleaser.yaml` to be rewritten as a semver tag that matches the git tag like this:

```yaml
# .helmreleaser.yaml

images:
  - id: elasticsearch
      imageKey: image
      tagKey: imageTag
      imageTemplate: "docker.elastic.co/elasticsearch/elasticsearch"
      tagTemplate: {{ .Major }}.{{ .Minor }}.{{ .Patch }}
```

Note: this example shows all fields, but only fields that affect the outcome are required. The example above would still have the same effect without the `imageKey` and `imageTemplate` fields.

## Nested Keys

Sometimes images aren't defined as top-level keys in the `values.yaml`.

```yaml
# values.yaml

images:
  api:
    repository: myapp/api
    tag: latest
```

This is supported by specifying the full path in the `.helmreleaser.yaml`.

```yaml
# .helmreleaser.yaml

images:
  - id: api
      imageKey: images.api.repository
      tagKey: images.api.tag
      imageTemplate: myapp/api
      tagTemplate: {{ .Major }}.{{ .Minor }}.{{ .Patch }}
```
