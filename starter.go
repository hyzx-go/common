package common_b2c

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hyzx-go/common-b2c/config"
	"log"
	"time"
)

type ApplicationService interface {
	Run(registerRoutes func(r *gin.Engine))
}

func (s *Service) Run(registerRoutes func(r *gin.Engine)) {
	// 创建 Gin 实例
	r := gin.New()

	// 注册路由
	if registerRoutes != nil {
		registerRoutes(r)
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

var banner = "//                          _ooOoo_                               //\n//                         o8888888o                              //\n//                         88\" . \"88                              //\n//                         (| ^_^ |)                              //\n//                         O\\  =  /O                              //\n//                      ____/`---'\\____                           //\n//                    .'  \\\\|     |//  `.                         //\n//                   /  \\\\|||  :  |||//  \\                        //\n//                  /  _||||| -:- |||||-  \\                       //\n//                  |   | \\\\\\  -  /// |   |                       //\n//                  | \\_|  ''\\---/''  |   |                       //\n//                  \\  .-\\__  `-`  ___/-. /                       //\n//                ___`. .'  /--.--\\  `. . ___                     //\n//              .\"\" '<  `.___\\_<|>_/___.'  >'\"\".                  //\n//            | | :  `- \\`.;`\\ _ /`;.`/ - ` : | |                 //\n//            \\  \\ `-.   \\_ __\\ /__ _/   .-` /  /                 //\n//      ========`-.____`-.___\\_____/___.-`____.-'========         //\n//                           `=---='                              //\n//      ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^        //\n//             佛祖保佑          永无BUG         永不修改         //"

var version = fmt.Sprintf("Welcome to the HYZX common framework. current version %s", "1.0.0")

type Starter struct {
	startTime   time.Time
	configOpts  func() []config.Option
	application ApplicationService
}

var (
	beforeInitializeConfigs []func() error
	afterInitializeConfigs  []func(p config.Parser) error
)

func (a *Starter) Start(registerRoutes func(r *gin.Engine)) {

	defer func() {
		if err := recover(); err != nil {
			config.Red("Stater application failed, please check configs \n", err)
		}
	}()

	// Logger Output defines the standard output of the print functions. By default, os.Stdout
	config.Cyan(banner)
	config.Blue(version)
	if a.configOpts == nil {
		config.Red("Cannot find server config, please check configs \n")
		panic("please setting configs ")
	}

	// Initialise configs and new application service
	service := a.Init()

	// Run service
	service.Run(registerRoutes)
}

func (a *Starter) Init() ApplicationService {

	// Initialise config
	config.NewParserManager(a.configOpts()...).
		BeforeInitializeConfigs(beforeInitializeConfigs).AfterInitializeConfigs(afterInitializeConfigs).Initialize()

	service := a.application
	if a.application == nil {
	}

	return service
}
