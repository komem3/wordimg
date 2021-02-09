# Word Rand Img

This generate a random image from the entered words.

## CLI

### Install

```shell
go get github.com/komem3/word_rand_img/cmd/wordimg
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
      --size=      Font size.
  -w, --width=     Width of the generated image. (default: 512)
  -h, --height=    Height of the generated image. (default: 512)
  -c, --color=     Text color.

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
