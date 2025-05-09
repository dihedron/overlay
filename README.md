# overlay - A simple tool to write arbitrary text to an existing image

Overlay is a simple tool to write arbitrary text, in arbitrary user-provided TTF fonts onto existing images. 
It supports several image formats (GIT, BMP, PNG and JPEG) and can be used as a filter in a shell pipeline to apply multiple text sections to the same image incrementally.

## How to build

In order to install the application, you need to have Golang 1.23+ and `make` installed.

To build the application, simply run:

```bash
$> make
```

This will build the `linux/amd64` version by default; in order to cross compile for other platforms, specify the `$GOOS/$GOARCH` combination like so:

```bash
$> make windows/amd64
```

The built artifact will be stored in the `dist/${GOOS}/${GOARCH}` subdirectory (e.g. `dist/linux/amd64`).

## How to package

The application can be packaged for Debian/Ubuntu, Red Hat/Fedora and Alpine linux by using nFPM. The tool can be installed by running 

```bash
$> make setup-tools
```

and then employed by running it like so:

```bash
$> make deb
```

The `rpm` and `apk` targets product Red Hat/Fedora and Alpine linux install packages, respectively.

## How to get help  on the build process

The `Makefile` is self-documented and can print information about its targets by running:

```
$> make help
```

## How to run

Check the `Makefile`'s `test` and `test-pipeline` targets to check how to run the command against a file or in a pipeline (respectively):

```bash
$> 	cat input.jpg | overlay --point=600,100 --size=72 --font=Economica-Regular.ttf --color=#FFFFFF --format=jpg --text="HALLO, WORLD..." | overlay --point=700,160 --size=48 --font=Economica-Regular.ttf --color=#00FF0033  --text="... from me!" --output=output.jpg
```

The application can be used to convert the image from one format to the other; when the input image is on the filesystem, the extension is used to automatically detect the image format; when the image is piped into the command, it is necessary to specify the (input) format via the `--format` flag. 


## Licenses

The tests make use of the `Economica` font ( Copyright (c) 2012, Vicente Lamonaca).

