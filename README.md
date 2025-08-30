# NapCatShellUpdater
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FSn0wo2%2FNapCatShellUpdater.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FSn0wo2%2FNapCatShellUpdater?ref=badge_shield)

###### 方便的更新 `Windows` 平台下的[` NapCat.Shell `](https://github.com/NapNeko/NapCatQQ)

---

将程序放在你的NapCat.Shell工作目录双击打开或使用命令

``` shell
./NapCatShellUpdater
-path=xxx(指定NapCat.Shell工作目录, 默认为当前目录)
-version=v0.0.0(指定更新版本, 默认留空自动获得最新版)
-download-url=xxx(指定下载url, 默认为: https://github.com/NapNeko/NapCatQQ/releases/download/%s/NapCat.Shell.zip)
-proxy=http://url:port(遵循proxy url规范, 默认留空)
-version=v0.0.0(指定更新版本, 默认留空自动获得最新版)
-debug=true(debug模式日志, 默认启用)
-exclude=xxx(忽略文件, 多个用逗号分隔, 默认为: config,logs,quickLoginExample.bat,update.bat(强制加入更新的压缩包和程序本身)
```

---

- [NapNeko/NapCatQQ](https://github.com/NapNeko/NapCatQQ)


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FSn0wo2%2FNapCatShellUpdater.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FSn0wo2%2FNapCatShellUpdater?ref=badge_large)