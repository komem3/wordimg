# Word Img

This generate a random image from the entered words.

## Http

### Endpoint

Endpoint to generate image.

`https://wordimg-otho5yxlgq-an.a.run.app`

### Query

The size of the image is fixed at 512 * 512.

- `text` Text to write to image.
- `size` Font size. Supports j${line} as a special format. If you specify 'j1', the font size will fit on one line.
- `color` Text color.
- `align` Word postion. Default is 'left'. 'left' or 'center'.

## CLI

### Install

```shell
go get github.com/komem3/wordimg/cmd/wordimg
```

### Usage

#### Help

```shell
$ wordimg -h
Usage:
  wordimg [OPTIONS]

Application Options:
  -m, --message=   Message to write to image. required.
  -i, --imagePath= Path of the image to write. Default is 'unix_time.png
  -f, --font=      Path to font file. Only support ttf.
      --size=      Font size. Supports j${line} as a special format. If you
                   specify 'j1', the font size will fit on one line.
  -a, --align=     Word postion. 'left' or 'center'
  -w, --width=     Width of the generated image. (default: 512)
  -h, --height=    Height of the generated image. (default: 512)
  -c, --color=     Text color. Format is 'R:G:B'. Example: 255:255:0(yellow)

Help Options:
  -h, --help       Show this help message
```

#### Simple case

```shell
$ wordimg -m "Hello World"
wrote: 1612796180.png
```

#### Specify font
```shell
$ wordimg -m "こんにちは世界" -f ./SawarabiGothic-Regular.ttf
wrote: 1612796255.png
```

## License

MIT

## Author
komem3
