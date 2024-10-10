<p>
    <a href="README.md">English</a>
    | <a href="README_CN.md">中文</a>
</p>
<p align="center"><img src="https://github.com/AquaApps/AkuaX/blob/main/assets/horu_circle.png?raw=true" alt="1600" width="25%"/></p>
<p align="center">
    <strong>Horu</strong>
    <br>
    <p>一个非常简易的反向代理。</a>
    <br>
</p>
<br>


## Features

- 使用起来非常非常非常简单。
- 纯Go语言。
- 高性能。
- 支持http**3**!
- 自动启用 **Brotli**、 **zlib**、 **gzip** 压缩算法。

## 快速上手

<p align="center">
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="license"/>
</p>

1. 生成一个样例配置

```shell
./horu -demo
```

2. 编辑配置

```yaml
entry_list:
    - port: 8080
      maps:
        source1.com: https://destination1.com
        source2.com: https://destination2.com
    - port: 80
      maps:
        source3.com: https://destination3.com
        source4.com: https://destination4.com
key_file: path/to/your/server.key
cert_file: path/to/your/server.crt
```

3. 运行

```shell
./horu -c /etc/horu/config.yaml
```

## 注意

本程序**仅**支持**前台**运行，如果需要后台运行请考虑'screen'。


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