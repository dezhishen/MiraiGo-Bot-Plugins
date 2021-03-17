# MiraiGo-Bot-Plugins
MiraiGo-Bot插件库
## 已实现

名称|描述
-|-
[一言](./pkg/plugins/hitokoto)|`.hitokoto`,调用一言接口,随机发一句骚话
[掷骰](./pkg/plugins/random)|`.r` 100内的整数
[天气](./pkg/plugins/weather)|`.weather 省份 城市 地区(可以不传)` 调用腾讯天气
[提醒](./pkg/plugins/tips)|`.tips 10:10 提示内容` 10:10时,@该发送群友+提示内容
[今日人品](./pkg/plugins/jrrp)|`.jrrp`展示今日人品<br> `.jrrp 7` 显示最多7日的历史人品
[日历](./pkg/plugins/calendar)|`.calendar` 展示今日的日历(阳历,周几,阳历节日,农历,农历节日)<br>`.calendar Y/N` 启用/禁用定时(早6点)发送
~~[微博监听](./pkg/plugins/weibolisten)~~|由于访问限制,暂时不可用