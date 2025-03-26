# 3d-rack-brackets

Learning project for 3d modelling using [OpenSCAD](https://openscad.org/) and [GhostSCAD](https://github.com/ljanyst/ghostscad/).

## Usage

```sh
> devbox shell
# To get the openscad code:
> devbox run render:prod
# or to get the stl:
> devbox run render:stl
```

No there's [output](./output/output.scad).

## Development setup

```sh
> devbox shell
> devbox run watch
```

Now you can edit the code and on every change the output file is re-rendered.

## Preview

Currently looks like this:

![current state of the 3d model](./output.png)
