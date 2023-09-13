# -*- coding: utf-8 -*-
"""
easyCDFT.py
Briefly describe the functionality and purpose of the file.

This is a Main function file!

This file is part of EasyCDFT.
EasyCDFT is a Python script that automate the use of Multiwfn

@author:
Kimariyb (kimariyb@163.com)

@license:
Licensed under the MIT License.
For details, see the LICENSE file.

@Data:
2023-09-13
"""
import configparser
import subprocess
import os

from datetime import datetime

# 获取当前文件被修改的最后一次时间
time_last = os.path.getmtime(os.path.abspath(__file__))
# 全局的静态变量
__version__ = "v1.0.0"
__developer__ = "Kimariyb, Ryan Hsiun"
__address__ = "XiaMen University, School of Electronic Science and Engineering"
__website__ = "https://github.com/kimariyb/easy-cdft"
__release__ = str(datetime.fromtimestamp(time_last).strftime("%b-%d-%Y"))


def welcome():
    """
    主页面信息
    """

    print(f"@Name: EasyCDFT")
    print(f"@Version: {__version__}, @Release date: {__release__}")
    print(f"@Developer: {__developer__}")
    print(f"@Address: {__address__}")
    print(f"@Website: {__website__}")

    # 获取当前日期日期和时间
    now_time = datetime.now().strftime("%b-%d-%Y, %H:%M:%S")
    # 输出版权信息和问候语
    print(f"(Copyright (C) 2023 Kimariyb. Currently timeline: {now_time})\n")


def read_config(config_path: str) -> configparser.ConfigParser():
    """
    读取当前文件夹下的 config.ini 文件

    Args:
        config_path(str): config.ini 文件的路径

    Returns:
        config(configparser.ConfigParser()): 创建 ini config 对象
    """
    # 获取 ini config 对象
    config = configparser.ConfigParser()
    # 设置注释前缀
    config.comment_prefixes = ';'
    # 添加一个虚拟的 section 头部
    config.read_string('[DEFAULT]\n' + open(config_path, encoding='utf-8').read())

    return config


def get_files_byType(input_type: str) -> list[str]:
    """
    得到符合 config 中 inputType 属性配置的文件类型的所有文件

    Args:
        input_type(str): 输入文件的类型

    Returns:
        results(list[str]): 返回符合要求的所有文件
    """
    # 得到当前文件夹路径
    current_path = os.getcwd()

    # 存储符合要求的文件
    results = []

    # 遍历当前文件夹下的所有文件
    for file_name in os.listdir(current_path):
        # 判断文件名的后缀是否与输入的类型匹配
        if file_name.endswith(input_type):
            # 如果匹配，则将文件名添加到结果列表中
            results.append(file_name)

    # 返回符合要求的所有文件
    return results


def process_file(file: str, multiwfn_path: str, commands: list[str]):
    """
    自动化执行 Multiwfn 程序

    Args:
        file(str): 需要处理的文件，可以为 fchk、gjf、mol 等文件
        multiwfn_path(str): Multiwfn 程序的路径
        commands(list[str]): 需要执行的命令

    Returns:
        None
    """
    print(f"Documents being processed: {file}")
    # 创建一个命令对象，执行 Multiwfn 程序
    cmd = [multiwfn_path, file]

    # 启动命令并获取标准输入和输出管道
    proc = subprocess.Popen(cmd, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE,
                            universal_newlines=True)

    # 向命令的标准输入写入指令
    for command in commands:
        proc.stdin.write(command + "\n")

    # 关闭标准输入管道
    proc.stdin.close()

    # 等待命令执行完成
    proc.wait()

    # 在每次执行完Multiwfn后会在当前文件夹下生成一个CDFT.txt文件
    # 将该文件名修改为 ${name}-CDFT.txt 文件，${name} 为一个变量，代表 file 文件的名字
    output_file = os.path.splitext(file)[0] + "-CDFT.txt"
    # 重命名 CDFT.txt 文件
    os.rename("CDFT.txt", output_file)
    print("Renamed CDFT.txt to", output_file)


def main():
    # 显示 hello 页面
    welcome()
    # config.ini 文件路径，必须在运行程序的当前路径下，名字为 config.ini
    config_path = os.path.join(os.getcwd(), 'config.ini')
    # 读取 config.ini 文件
    config = read_config(config_path)
    # 读取 inputType 和 multiwfnPath 两个配置
    input_type = str(config.get('DEFAULT', 'inputType')).strip('\"')
    multiwfn_path = str(config.get('DEFAULT', 'multiwfnPath')).strip('\"')
    # 同时打印 inputType 和 multiwfnPath 两个配置
    print(f"inputType: {input_type}")
    print(f"multiwfnPath: {multiwfn_path}\n")
    # 根据 inputType 所配置的文件类型，得到当前文件下所有符合 inputType 类型的文件
    files = get_files_byType(input_type)
    print(f"Files with inputType {input_type}:")
    # 打印每一个文件的名字
    for file in files:
        print(file)
    # 读取 mission 配置
    mission = int(config.get('DEFAULT', 'mission'))
    # 读取 calcLevel 配置
    calc_level = str(config.get('DEFAULT', 'calcLevel')).strip('\"')
    # 读取 ChargeSpin1、2、3 配置
    chargeSpin1 = str(config.get('DEFAULT', 'chargeSpin1')).strip('\"')
    chargeSpin2 = str(config.get('DEFAULT', 'chargeSpin2')).strip('\"')
    chargeSpin3 = str(config.get('DEFAULT', 'chargeSpin3')).strip('\"')

    # 根据 config.ini 设置储存命令
    if mission == 0:
        # 如果 mission 为 0 则使用计算指数命令
        commands = ["22", "1", calc_level, chargeSpin1, chargeSpin2, chargeSpin3, "y", "2", "q"]
    elif mission == 1:
        print("Warning: This feature is not yet developed")
        return
    elif mission == 2:
        print("Warning: This feature is not yet developed")
        return
    else:
        # 如果没有符合的命令，则抛出异常
        print("Error: Invalid mission value")
        return
    # 批量执行 Multiwfn
    for file in files:
        process_file(file, multiwfn_path, commands)

    print()
    print("The mission has been successfully completed!")


if __name__ == '__main__':
    main()
