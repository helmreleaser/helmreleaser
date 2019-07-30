---
title: 'Configuration'
date: 2019-02-11T19:30:08+10:00
draft: false
weight: 4
---

## .helmreleaser.yaml

HelmReleaser reads configuration from a `.helmreleaser.yaml` file that contains the following top level keys:

- [Images](#images)

### Images {#images}
Images define the Docker images that are included in the Helm chart. Specifying the `imageTemplate` will update the image name in the `values.yaml` to the rendered output of the template.

```yaml
# .helmreleasder.yaml

images:
  - id: api-server
      imageKey: image
      tagKey: imageTag
      imageTemplate: my-org/api-server
      tagTemplate: "{{ .Major }}.{{ .Minor }}.{{ .Patch }}"
```

The value of the `id` must be unique in the `images` array.

`imageKey` is the fully qualified path to the `image` key in the `values.yaml`. For example, if your image is defined as:

```yaml
# values.yaml

apiserver:
  image: "apiserver"
  tag: "1.29.3"
```

Then you can reference this in a `helmreleaser.yaml` as:

```yaml
# .helmreleaser.yaml

images:
  - id: apiserver
      imageKey: apiserver.image
      tagKey: apiserver.tag
      imageTemplate: my-org/api-server
      tagTemplate: "{{ .Major }}.{{ .Minor }}.{{ .Patch }}"
```
