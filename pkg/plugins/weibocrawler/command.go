package weibocrawler

import "github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"

// Command 微博命令插件
type Command struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

// func init() {
// 	plugins.RegisterOnMessagePlugin(WeiboCrawlerCommand{})
// }
