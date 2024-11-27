package common_b2c

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hyzx-go/common-b2c/config"
	innerLog "github.com/hyzx-go/common-b2c/log"
	"github.com/hyzx-go/common-b2c/utils"
	"log"
	"time"
)

type ApplicationService interface {
	Run([]func(r *gin.RouterGroup))
}

func newDefaultApplication(startTime time.Time) ApplicationService {
	service := newMicroService(startTime)
	return service
}
func newMicroService(startTime time.Time) *Service {
	service := &Service{startTime: startTime}
	service.parser = config.GetParser()

	sys, err := service.parser.GetSystemConf()
	if err != nil {
		panic(err)
	}

	log.Printf(fmt.Sprintf("Load micro service env:%s,serviceName:%s,local:%s,timeZone:%s", service.parser.GetEnv(), sys.ServiceName, sys.Local, sys.TimeZone))

	return service
}

func (s *Service) Run(routerModules []func(r *gin.RouterGroup)) {
	// 创建 Gin 实例
	r := gin.New()

	// 注册模块路由
	group := r.Group("", innerLog.RequestLogger(), innerLog.GinRecovery())
	for _, module := range routerModules {
		module(group)
	}

	sysConf, err := s.parser.GetSystemConf()
	if err != nil {
		log.Fatalf("Failed to start server get sys conf: %v", err)
	}
	// 启动服务
	log.Printf("Starting server on port %s...", sysConf.ServePort)
	if err := r.Run(":" + sysConf.ServePort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

type Service struct {
	startTime time.Time
	parser    config.Parser
}

func NewsStartService(routers []func(r *gin.RouterGroup), applications ...ApplicationService) *Starter {
	starter := &Starter{startTime: time.Now(), routers: routers}
	if len(applications) > 0 {
		starter.application = applications[0]
	}
	return starter
}

var banner = "//                          _ooOoo_                               //\n//                         o8888888o                              //\n//                         88\" . \"88                              //\n//                         (| ^_^ |)                              //\n//                         O\\  =  /O                              //\n//                      ____/`---'\\____                           //\n//                    .'  \\\\|     |//  `.                         //\n//                   /  \\\\|||  :  |||//  \\                        //\n//                  /  _||||| -:- |||||-  \\                       //\n//                  |   | \\\\\\  -  /// |   |                       //\n//                  | \\_|  ''\\---/''  |   |                       //\n//                  \\  .-\\__  `-`  ___/-. /                       //\n//                ___`. .'  /--.--\\  `. . ___                     //\n//              .\"\" '<  `.___\\_<|>_/___.'  >'\"\".                  //\n//            | | :  `- \\`.;`\\ _ /`;.`/ - ` : | |                 //\n//            \\  \\ `-.   \\_ __\\ /__ _/   .-` /  /                 //\n//      ========`-.____`-.___\\_____/___.-`____.-'========         //\n//                           `=---='                              //\n//      ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^        //\n//             佛祖保佑          永无BUG         永不修改         //"

var version = fmt.Sprintf("Welcome to the HYZX common framework. current version %s", "1.0.0")

type Starter struct {
	startTime   time.Time
	configOpts  func() []config.Option
	application ApplicationService
	routers     []func(r *gin.RouterGroup)
}

var (
	beforeInitializeConfigs []func() error
	afterInitializeConfigs  []func(p config.Parser) error
)

func (a *Starter) Start() {

	defer func() {
		if err := recover(); err != nil {
			config.Red("Stater application failed, please check configs \n", err)
		}
	}()

	// Logger Output defines the standard output of the print functions. By default, os.Stdout
	config.Cyan(banner)
	config.Blue(version)

	service := a.Init()

	// Run service
	service.Run(a.routers)
}

func (a *Starter) Init() ApplicationService {

	// Initialise config
	config.NewParserManager().
		BeforeInitializeConfigs(beforeInitializeConfigs).AfterInitializeConfigs(afterInitializeConfigs).Initialize()

	service := a.application
	if a.application == nil {
		service = newDefaultApplication(a.startTime)
	}

	return service
}

// MainTest used this method start your test
func MainTest(opts []config.Option, routers []func(r *gin.RouterGroup)) {

	path, err := utils.LookUpFilePath("config.yaml", 8)
	if err != nil {
		return
	}

	opts = append(opts, config.SetConfigFilePath(path))
	config.NewParserManager(opts...).Initialize()

	service := newDefaultApplication(time.Now())
	service.Run(routers)
}
