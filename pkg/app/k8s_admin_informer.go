package app

import (
	"github.com/gin-gonic/gin"
	"k8s-admin-informer/pkg/handler"
)

type App struct {
	// router is the app router engine
	engine          *gin.Engine
	baseHandler     *handler.Handler
	workloadHandler *handler.WorkloadHandler
	rscHandler      *handler.ResourceHandler
}

func NewK8sAdminInformerApp() *App {
	baseHandler, err := handler.NewHandler()
	if err != nil {
		panic(err)
	}
	return &App{
		engine:          gin.Default(),
		baseHandler:     baseHandler,
		workloadHandler: handler.NewWorkloadHandler(baseHandler),
		rscHandler:      handler.NewResourceHandler(baseHandler),
	}
}

func (a *App) registerRoute() {
	// 查询工作负载后面的pod和event
	a.engine.POST("/informer/v1/getWorkloadInstance", a.workloadHandler.GetWorkloadInstance)
	// 检查当前请求资源是否超过部门配额
	a.engine.POST("/informer/v1/resource/dept/checkLimit", a.rscHandler.ComputeDeptResourceQuotaLimit)
	// 获取节点资源
	a.engine.GET("/informer/v1/resource/node", a.rscHandler.NodeResources)
	// 获取部门资源
	a.engine.GET("/informer/v1/resource/dept", a.rscHandler.DeptResources)
}

func (a *App) Run() error {
	// 注册路由
	a.registerRoute()

	// 启动各个informer
	if err := a.baseHandler.Start(); err != nil {
		return err
	}

	// 运行server
	err := a.engine.Run(":8080")
	if err != nil {
		return err
	}
	return nil
}
