package api

/*
* 使用 _ 导入路径的方式是为了触发 xxx 包的 init 方法。
* xxx 包的 init 方法会调用 router.Register，自动将用户路由注册到全局路由列表中。
 */

import (
	_ "github.com/MortalSC/IM-System/api/user"
)
