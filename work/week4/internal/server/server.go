package server

import "github.com/google/wire"

// http和grpc实例的创建和配置,或者直接直接独立进程方式启动
var ProviderSet = wire.NewSet(NewCronServer)
