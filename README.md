<p>
    <a href="README.md">English</a>
    | <a href="README_CN.md">中文</a>
</p>
<p align="center"><img src="https://github.com/AquaApps/AkuaX/blob/main/assets/horu_circle.png?raw=true" alt="1600" width="25%"/></p>
<p align="center">
    <strong>Horu</strong>
    <br>
    <p>A simple reverse proxy.</a>
    <br>
</p>
<br>


## Features

- Verrrryyyy easy for use.
- Pure golang.
- High performance.
- Automatically select **Brotli**, **zlib**, and **gzip** compression algorithms.

## Usage

<p align="center">
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="license"/>
</p>

1. Generate a demo config.

```shell
./horu -demo
```

2. Edit the config.

```yaml
point_lists:
    - port: 8080
      points:
        https://source.com:443: https://destination1.com
        https://source.cn: http://destination2.com
        http://source.cat: https://destination3.com
    - port: 80
      points:
        https://source2.com: https://localhost:8088
        https://source2.com: https://destination4.com
key_file: path/to/your/server.key
cert_file: path/to/your/server.crt
```

3. Run

```shell
./horu -c /path/to/config.yaml
```

## Notice

This program **only** supports running in the **foreground**. If you need to run it in the background, please use 'screen'.

## License

```
MIT License

Copyright (c) 2023 Aqua/A-kua

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```