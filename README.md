# MiraiGo-Bot-Plugins
MiraiGo-Bot插件库


## 已实现
名称|描述
-|-
[运势插件](./pkg/plugins/todayFortune)|`.tf` 随机运势,素材等来源[`https://github.com/FloatTech/Plugin-Fortune`](https://github.com/FloatTech/Plugin-Fortune),为了避免流量占用过多,素材复制至-> https://github.com/dezhishen/raw/tree/master/fortune
[彩虹屁](./pkg/plugins/caihongpi)|`.chp` 随机发送一条彩虹屁
[日历](./pkg/plugins/calendar)|`.calendar` 展示今日的日历(阳历,周几,阳历节日,农历,农历节日)<br>`.calendar Y/N` 启用/禁用定时(早6点)发送
[毒鸡汤](./pkg/plugins/dujitang)|`.djt` 随机发送一条毒鸡汤
[古风小姐姐](./pkg/plugins/haimage)|`.hapic` 随机发送一张古风小姐姐图
[一言](./pkg/plugins/hitokoto)|`.hitokoto`,调用一言接口,随机发一句骚话
[今日人品](./pkg/plugins/jrrp)|`.jrrp`展示今日人品<br> `.jrrp 7` 显示最多7日的历史人品
[lpl赛事](./pkg/plugins/lpl)|`.lpl`查看最近的赛事
[menhera图片](./pkg/plugins/mc)|`.mc`随机一张menhera酱,她真可爱.jpg
[pixiv](./pkg/plugins/pixiv)|`.pixiv`,随机一张 懂的都懂
[掷骰](./pkg/plugins/random)|`.r` 100内的整数
[苏联笑话](./pkg/plugins/sovietjokes)|`.sj`随机一条人类政治精华
[猫猫图](./pkg/plugins/thecat)|`.thecat`/`.cat` 随机发送一条猫猫图
[狗狗图](./pkg/plugins/thedog)|`.thedog`/`.dog` 随机发送一条狗狗图
~~[舔狗语录](./pkg/plugins/tiangou)~~|`.tg` 随机发送一条舔狗语录(重复率较高),已停用
[提醒](./pkg/plugins/tips)|`.tips 10:10 提示内容` 10:10时,@该发送群友+提示内容
[天气](./pkg/plugins/weather)|`.weather 所在地` 调用 https://github.com/schachmat/wego
~~[微博监听](./pkg/plugins/weibolisten)~~|由于访问限制,暂时不可用
[B站粉丝数查询](./pkg/plugins/bilifan)|`.bilifan UID` 发送该UID的粉丝数量
[翻译插件](./pkg/plugins/translate)|`.tr test` 翻译文本<br> `.tr -f zh 三点多,先喝茶 -t yue`
[表情包](./pkg/plugins/facesave)|`.face-save -n/--name xxx` 紧跟着发送一张图片,以后发送 xxx(图片名称),bot会发出该图片

## 启动方式
### 宿主机方式
1.在[releases](https://github.com/dezhiShen/MiraiGo-Bot-Plugins/releases)中下载对应的包

2.执行
### docker
- 1.选择对应的镜像
- 2.第一次执行时,需要生成设备信息,验证设备合法性等
  - 2.1.`docker run -it -v ${数据目录}:/data dezhishen/miraigo-bot:${version}`
  - 2.2.按照提示,输入账户/密码,验证设备等
  - 2.3.关闭
- 3.`docker run -d --restart=always -v ${数据目录}:/data dezhishen/miraigo-bot:${version}`

### 二次开发
参考本项目启动方式 [cmd/main](./cmd/main.go)
```
package main

import (
  //引入你要使用的插件
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/calendar"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/haimage"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/hitokoto"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/jrrp"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/lpl"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/mc"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/pixiv"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/random"

	// _ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/tiangou"

	// _ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/segmentation"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/thecat"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/thedog"

	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/tips"

	// _ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/vader"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/caihongpi"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/dujitang"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/weather"

	// _ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/weibolisten"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/sovietjokes"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/server"
)

func main() {
  //启动
	server.Start()
}

```
