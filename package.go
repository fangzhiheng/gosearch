package gosearch

// Package 代指一个Go模块
type Package struct {
	// 包名称
	Name string

	// 导入路径
	ImportPath string

	// 主页
	HomeSite string

	// 概要信息
	Synopsis string

	// 许可证
	License string
}
