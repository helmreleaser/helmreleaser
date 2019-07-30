---
title: 'Snapshots'
date: 2019-02-11T19:27:37+10:00
weight: 7
---

To make it easier to test, HelmReleaser supports **snapshot** builds. A snapshot is similar to a tagged release, but snapshots don't require a git tag to be present and they don't publish to a Helm reposistory by default.

To create a snapshot release, run HelmReleaser with the `--snapshot` flag. Snapshot release overrides are read from the `snapshot` section of the `helmreleaser.yaml`.

```yaml
# .helmreleaser.yaml

snapshot:
  images:
    id: imageid
      imageTemplate: localhost:32000/my-image
      tagTemplate: latest
```

In the example above, HelmReleaser will replace the image name and template for the image with id `imageid` with a locally-built image reference.

[Learn more about image references](/docs/images)

The idea behind HelmReleaser's snapshots is for local testing, builds or to validate a CI pipeline.
