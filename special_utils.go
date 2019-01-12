package utils

import (
	"bytes"
	"github.com/satori/go.uuid"
	"log"
	"runtime"
	"strings"
	"strconv"
	"net"
	"fmt"
	"os"
)

type Base struct {
	Page    string `json:"page"`
	Rows    string `json:"rows"`
	OrderBy string `json:"orderBy"`
}

func (b *Base) GetLimit() string {
	page := ParseInt(b.Page)
	rows := ParseInt(b.Rows)
	if (page == 0 || rows == 0) && page != -1 {
		page = 1
		rows = 10
	}
	if page > 0 && rows > 0 {
		return " limit " + strconv.Itoa((page-1)*rows) + ", " + strconv.Itoa(rows)
	}
	return ""
}

//create by cwj on 2017-08-24
func GenerateUuid() string {
	return strings.Replace(uuid.NewV4().String(), "-", EMPTY_STRING, -1)
}

//create by cwj on 2017-10-17
//check errArray
//if there are error, return it
func CheckError(errArray ...error) error {
	for _, e := range errArray {
		if e != nil {
			return e
		}
	}
	return nil
}

//create by cwj on 2017-10-17
//check string
//if there are Single quotation marks('),transfer it for prevent injection
func QuotationTransferred(s interface{}) string {
	str := ParseString(s)
	str = strings.Replace(str, "'", "''", -1)
	return str
}

func QuotationTransferredForLike(s interface{}) string {
	str := ParseString(s)
	str = strings.Replace(str, "'", "''", -1)
	str = strings.Replace(str, "%", "[%]", -1)
	str = strings.Replace(str, "_", "[_]", -1)
	return str
}

func TEE(isTrue bool, ele1 interface{}, ele2 interface{}) interface{}{
	if isTrue{
		return ele1
	} else {
		return ele2
	}
}

func GetLocalIp() string{
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		ULog.Println(err)
		return "127.0.0.1"
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

var ULog UtilLog

type UtilLog struct{
	*log.Logger
}

func (u *UtilLog)Println(con ...interface{})  {
	stat := fmt.Sprint(GoroutineId(), "-[UTILS]", con)
	log.Output(2, stat)
	u.Logger.Output(2, stat)
}

func init()  {
	ULog.Logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
}


func GoroutineId() string{
	b := make([]byte, 32)
	b = b[:runtime.Stack(b, false)]
	b = b[:bytes.IndexByte(b, '[')]
	return string(b)
}