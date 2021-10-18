package main

import "github.com/gin-gonic/gin"
import "net/http"
import "net"
import "log"
import "syscall"
import "time"
import "fmt"
import "os"
import "encoding/json"
import "unicode/utf16"
import "github.com/pbnjay/memory"
import "github.com/tidwall/gjson"


type ServerDetail struct {
    Id string `json:"id"`
    Name string `json:"name"`
    Address string `json:"address"`
}

func main() {
    start := time.Now()
	router := gin.Default()

	router.LoadHTMLGlob("web/templates/*")
	router.Static("/assets","web/public") 
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "ServoGo",
		})
	})

	router.GET("/dummy", func(c *gin.Context) {
		file, err := os.Open("config.json")

		var result map[string]interface{}
		json.Unmarshal([]byte(file), &result)
		
	 	// ServerDetails := []ServerDetail{{Id: "ravi", Name: "travel", Address: "golang"},{Id: "ab", Name: "cd", Address: "ef"}}
	 	fmt.Println(result)
	 	fmt.Println(err)
	 	k1 := gjson.Get(json, result)

		c.JSON(200, file)
	})
	router.GET("/ping", func(c *gin.Context) {

		elapsed := time.Since(start)

    	processingStart := time.Now()

		c.JSON(200, gin.H{
			"ipAddress": GetOutboundIP(),
			"userName": ComputerName(),
			"memory": DeviceMemory(),
			"agentName": AgentName(c.Request),
			"serverRuntime": elapsed.String(),
			"responseTime": time.Since(processingStart).String(),
		})
	})
	router.Run(":8080")
}

func DeviceMemory() (mem uint64){
	return memory.TotalMemory()
}

func AgentName(r *http.Request) (name string) {
 	ua := r.Header.Get("User-Agent")
 	return ua
}

func ComputerName() (name string) {
    var n uint32 = syscall.MAX_COMPUTERNAME_LENGTH + 1
    b := make([]uint16, n)
    e := syscall.GetComputerName(&b[0], &n)
    if e != nil {
        return ""
    }
    return string(utf16.Decode(b[0:n]))
}

func GetOutboundIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP
}