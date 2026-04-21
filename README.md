# Roku-Cell - A Sixel Image Viewer

`rkcel` shows images on Terminals which support **Sixel**.

## Build

```
make
```

## Usage

```
Usage: ./rkcel [OPTIONS] PATH

PATH: Filename of image file
      (BMP, GIF, JPEG, PNG, TIFF, WebP)
OPTIONS:
  -n N: Use N colors (N: max 255)
  -d: Disable dithering
  -m: Disable median cut
  -c: Run calibration
  -f: Disable fitting
  -sb: Approximate bilinear scaling
  -sn: Nearest neighbor scaling
  -cover: Cover fitting
  -wait: Wait enter key
  -raw: Output raw Sixel
```
