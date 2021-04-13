# Content Files Directory

The contents of this directory contains raw plugins documentation files and partials that can
be used with the `packer-sdc renderdocs` command to auto generate full docs for one or more components.

| Directory contents | Description |
-----
|partials/ | Shared documents that can be included in one or doc files |
|docs/[builders|datasources|provivisoners|post-processers]/\*.mdx |  Generic documentation files for the respective components |
|docs/README.md | Default README to include as GitHub index for docs directory|


### Auto Generating Docs
Make the necessary updates the content-files; once ready use packer-sdc to generate the update docs

```
packer-sdc renderdocs -src content-files/docs -partials content-files/partials -dst docs/
```

