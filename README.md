# NCMconverter
[![PkgGoDev](https://pkg.go.dev/badge/github.com/closetool/NCMconverter)](https://pkg.go.dev/github.com/closetool/NCMconverter)  

NCMconverter将ncm文件转换为mp3或者flac文件

实现参考了[yoki123/ncmdump][1]，重构了代码，并且添加了多线程支持

## 使用

* `NCMconverter [options] <files/dirs>`

* `--output value, -o value  指定输出目录，默认为原音频文件夹
   --tag, -t                 是否使用给转换后的文件添加meta信息（有bug，这个参数没有用）
   --deepth value, -d value  文件目录寻找的最大深度，默认为0，无视目录
   --thread value, -n value  线程数
   --help, -h                help
   --version, -v             version`

---
[1]:https://github.com/yoki123/ncmdump

