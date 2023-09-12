# EasyCDFT

## 简介

EasyCDFT 是厦门大学电子科学系研究生 Kimariyb 开发的一个 Go 语言程序。EasyCDFT 可以批量执行 Multiwfn 并使用 Multiwfn 生成各种概率密度泛函理论（CDFT）定义的各种量。包括局部亲核指数、局部亲电指数、全局亲核指数、全局亲电指数等等。

## 安装

EasyCDFT 十分易安装，下载地址为：https://github.com/kimariyb/easy-cdft/releases/download/v1.0.0/EasyCDFT.zip

## 使用

将下载好的 EasyCDFT.zip 解压后，可以修改其中的 `config.ini`。`config.ini` 是 EasyCDFT 的配置文件，请根据各自电脑的情况，配置 `multiwfnPath` 属性，其他的属性可以根据自己的需求修改。

```ini
; 用于配置需要在 Multiwfn 中做批处理的输入文件
; 需要含有体系的结构信息
; 可以为 fchk、mol、mol2、xyz、pdb、gjf 等文件
inputType = "fchk"
; Multiwfn 可执行文件所在的路径，例如 /home/kimariyb/Multiwfn/Multiwfn.exe
; 同时需要修改 Multiwfn 的 settings.ini，并且配置 gopath
multiwfnPath = "/home/kimariyb/Multiwfn/Multiwfn"
; 需要使用 Multiwfn 批量计算概念密度泛函理论中定义的各种量的类型
; 0. 计算各种指数，包括全局亲核、亲电；局部亲核、亲电指数
; 1. 计算福井函数和双描述符 （暂未开发）
; 2. 计算 ωcubic （暂未开发）
mission = 0
; 考察 CDFT 定义的量时使用的计算级别。必须为 Gaussian 能读懂的关键词
calcLevel = "B3LYP/6-31G*"
; 计算 N、N+1 和 N-1 态用的电荷和自旋多重度。
ChargeSpin1 = "0 1"
ChargeSpin2 = "-1 2"
ChargeSpin3 = "1 2"
; 其他功能，暂未开发
```

配置完成后，如果在 Windows 系统下，可以点击 `easy-cdft.exe` 运行；如果在 Linux 系统下，可以在 Bash 终端中输入 `./easy-cdft` 运行。

## 注意

无论使不使用 EasyCDFT，要想在 Multiwfn 中做 CDFT 分析时运行 Gaussian 必须配置 Multiwfn 的 settings.ini 中的 `gopath` 属性！
