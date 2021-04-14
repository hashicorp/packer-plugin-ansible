# Content Files Directory

The contents of this directory contains raw plugin documentation files and partials that can
be used with the `packer-sdc renderdocs` command to auto generate full docs for one or more components.

| Directory contents | Description |
|:------------- |:-------------|
|partials/\*\*/\*.mdx | Shared documents that can be included in one or more doc files |
|docs/component/\*.mdx |  Generic documentation files for the respective components (i.e builders, datasources, provisioners, post-processors)|
|docs/README.md | Default README to include as GitHub index for docs directory|


### Auto Generating Docs
Make the necessary updates to content-files; once ready use packer-sdc to generate the update docs

```
packer-sdc renderdocs -src content-files/docs -partials content-files/partials -dst docs/
```

