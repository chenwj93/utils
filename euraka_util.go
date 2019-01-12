package utils

import (
	"encoding/json"
	"errors"
	"github.com/chenwj93/eureka"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Applications = make(map[string]eureka.Application)
var client *eureka.Client

var (
	hostName, app, ip string
	port              int
	instanceMap       map[string]string
	machines          []string
)

func StartEureka(App string, Port int, insMap map[string]string, Machines []string, serviceId ...string) {
	once := sync.Once{}
	for i := 0; i < len(Machines); i++ {
		if !strings.HasPrefix(Machines[i], "http://") {
			Machines[i] = "http://" + Machines[i]
		}
	}
	if len(serviceId) != 0 && serviceId[0] != "" {
		hostName = serviceId[0]
	} else {
		hostName = GetLocalIp()
	}

	ip = hostName
	app, port, instanceMap, machines = App, Port, insMap, Machines
	once.Do(eurekaService)
}

func ServiceCall(app, operate, method string, param, head map[string]interface{}, result interface{}) error {
	application, e := GetApp(app)
	if e != nil {
		ULog.Println(e)
		if a, ok := Applications[app]; ok {
			application = &a
		} else {
			return e
		}
	}
	if len(application.Instances) == 0 {
		return errors.New("目标模块地址未发现")
	}
	url := "http://" + application.Instances[0].IpAddr + ":" + strconv.Itoa(application.Instances[0].Port.Port) + "/" + operate
	ULog.Println(url)
	m := make(map[string]interface{})
	switch method {
	case "GET", "get", "Get":
		m, e = GetDataByHttpGet(url, param, head)
	case "POST", "post", "Post":
		m, e = GetDataByHttpPost(url, param, head)
	}
	if e != nil {
		return e
	}
	e = ParseStruct(m, result)
	return e
}

func GetServiceAddress(app string) (string, error) {
	application, e := GetApp(app)
	if e != nil {
		ULog.Println(e)
		if a, ok := Applications[app]; ok {
			application = &a
		} else {
			return "", e
		}
	}
	if len(application.Instances) == 0 {
		return "", errors.New("目标模块地址未发现")
	}
	url := "http://" + application.Instances[0].IpAddr + ":" + strconv.Itoa(application.Instances[0].Port.Port)
	return url, nil
}

func eurekaService() {

	var e error
	//client, e = eureka.NewClientFromFile("./eureka_util/config.json")
	//if nil != e{
	//	log.Println(e)
	//	return
	//}

	client = eureka.NewClient(machines)

	//instance := eureka.NewInstanceInfo("test.com", "myapp", "192.168.1.107", 8003, 30, false) //Create a new instance to register
	instance := eureka.NewInstanceInfo(hostName, app, ip, port, 30, false) //Create a new instance to register
	instance.Metadata = &eureka.MetaData{
		Map: instanceMap,
	}
	//instance.Metadata.Map["foo"] = "bar" //add metadata for example
	e = client.RegisterInstance(app, instance) // Register new instance in your eureka(s)
	if e != nil {
		panic("注册错误：" + e.Error())
	}

	go appsInit()
	go heartBeat(instance.App, instance.InstanceId)
}

func GetApp(app string) (application *eureka.Application, e error) {
	if client == nil {
		return nil, errors.New("未注册eureka")
	}
	application, e = client.GetApplication(app)
	if e != nil {
		ULog.Println("get application error:", e)
	}
	return
}

func GetIns(app string, hostname string) (application *eureka.InstanceInfo, e error) {
	if client == nil {
		return nil, errors.New("未注册eureka")
	}
	application, e = client.GetInstance(app, hostname)
	if e != nil {
		ULog.Println("get instance error:", e)
	}
	return
}

func appsInit() {
	tick := time.Tick(10 * time.Second)
	ifSuccess := false
	for !ifSuccess {
		apps, e := client.GetApplications() // Retrieves all applications from eureka server(s)
		if nil == e {
			for _, ele := range apps.Applications {
				Applications[ele.Name] = ele
			}
			print, _ := json.Marshal(Applications)
			ULog.Println(string(print))
			ifSuccess = true
		} else {
			ULog.Println("get apps error:", e)
		}
		<-tick
	}
	return
}

func heartBeat(app string, instanceId string) {
	tick := time.Tick(30 * time.Second)
	for {
		<-tick
		//fmt.Println("\nheart...")
		client.SendHeartbeat(app, instanceId) // say to eureka that your app is alive (here you must send heartbeat before 30 sec)
	}
}
