# preload_dllhijack 插件编写说明

可以查看文件夹下 preload_dllhijack_*.go 样例，需要实现一个接口 SubPluginPreloadDllHijackIface

`GetPluginName() string` 返回一个插件名，需要在项目根目录的 data 文件夹下创建一个和插件名相同的文件夹，然后放置白exe程序

如果白exe程序不止导入了一个非操作系统的dll，则除了你要劫持的dll之外，其余的非操作系统dll需要放置到上面创建的文件夹中，并且实现方法 `GetExtraFileList() ([]string)`，可参见示例 [preload_dllhijack_steam.go](./preload_dllhijack_steam.go)
