package main

import (
	"fmt"
	"os"
	"time"

	proto "gmsec/rpc/example"

	"github.com/xxjwxc/public/dev"

	"gmsec/internal/config"

	"github.com/gmsec/goplugins/api"

	"github.com/gin-gonic/gin"
	"github.com/gmsec/goplugins/plugin"
	"github.com/gmsec/micro"
	"github.com/xxjwxc/ginrpc"
	"github.com/xxjwxc/public/mydoc/myswagger"
	"github.com/xxjwxc/public/server"
)

// CallBack service call backe
func CallBack() {
	// swagger
	myswagger.SetHost("https://localhost:8080")
	myswagger.SetBasePath("gmsec")
	myswagger.SetSchemes(true, false)
	// -----end --

	// reg := registry.NewDNSNamingRegistry()
	// grpc 相关 初始化服务
	service := micro.NewService(
		micro.WithName("lp.srv.eg1"),
		// micro.WithRegisterTTL(time.Second*30),      //指定服务注册时间
		micro.WithRegisterInterval(time.Second*15), //让服务在指定时间内重新注册
		// micro.WithRegistryNaming(reg),
	)
	h := new(hello)
	proto.RegisterHelloServer(service.Server(), h) // 服务注册
	// ----------- end

	// gin restful 相关
	base := ginrpc.New(ginrpc.WithCtx(api.NewAPIFunc), ginrpc.WithDebug(dev.IsDev()))
	router := gin.Default()
	v1 := router.Group("/xxjwxc/api/v1")
	base.Register(v1, h) // 对象注册
	// ------ end

	plg, b := plugin.Run(plugin.WithMicro(service),
		plugin.WithGin(router),
		plugin.WithAddr(":82"))

	if b == nil {
		plg.Wait()
	}
	fmt.Println("done")
}

func main() {
	if config.GetIsDev() || len(os.Args) == 0 {
		CallBack()
	} else {
		server.On(config.GetServiceConfig()).Start(CallBack)
	}
}
