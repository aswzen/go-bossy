package main

import "github.com/gin-gonic/gin"
import "net/http"
import "net"
import "log"
import "syscall"
import "time"
import "bytes"
import "io/ioutil"
import "fmt"
import "os"
import "strings"
import "encoding/hex"
import "io"
import "encoding/json"
import "unicode/utf16"
import "github.com/pbnjay/memory"
import "crypto/md5"
import "crypto/rand"
// import "github.com/bobziuchkovski/digest"
// import "github.com/tidwall/gjson"


type ServerDetail struct {
    Id string `json:"id"`
    Name string `json:"name"`
    Icon string `json:"icon"`
    Address string `json:"address"`
    Port string `json:"port"`
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


	router.GET("/server-list", func(c *gin.Context) {
		file, err := os.Open("config.json")
		if err != nil {
	      panic(err)
	    }
		defer file.Close()
	  	byteValue, _ := ioutil.ReadAll(file)
	    var result map[string]interface{}
	    json.Unmarshal([]byte(byteValue), &result)


	    listSlice, ok := result["server"].([]interface{})
	    if !ok {
      		panic(err)
	    }

	    for _, v2 := range listSlice { 
    		fmt.Println(v2)
    		nestedMap, _ := v2.(map[string]interface{})

    		// fmt.Println(nestedMap["address"])
    		// fmt.Println(url)
    		nestedMap["deploymentName"] = "TEST"
		 	var serverDet = getServerVersion(nestedMap["api"].(string), nestedMap["username"].(string),nestedMap["password"].(string));
		    var result map[string]interface{}
		    json.Unmarshal([]byte(serverDet), &result)

		    var sap = []string{}
    		deployments, _ := result["deployment"].(map[string]interface{})
		    for _, v3 := range deployments { 
		 		var lastDeploy = getLastDeployed(nestedMap["api"].(string)+"/deployment/"+v3.(string)+"?operation=attribute&name=enabled", nestedMap["username"].(string),nestedMap["password"].(string));
		 		sap = append(lastDeploy, _);
		    }

    		// nestedMap["serverVersion"] = serverDet
    		nestedMap["lastDeploy"] = sap
    		nestedMap["deployment"] = result["deployment"]
    		nestedMap["version"] = nil
    		if result["product-name"] != nil {
	    		nestedMap["version"] = result["product-name"].(string)+" - "+result["product-version"].(string)
	    	}
    		// fmt.Println(nestedMap["address"])
		}

		c.JSON(200, result["server"])
	})

	// router.GET("/server-list", func(c *gin.Context) {
	// 	file, err := os.Open("config.json")
	// 	if err != nil {
	//       panic(err)
	//     }
	// 	defer file.Close()
	//   	byteValue, _ := ioutil.ReadAll(file)
	//     var result map[string]interface{}
	//     json.Unmarshal([]byte(byteValue), &result)

	//     for k, v := range result { 
	// 	    fmt.Printf("key[%s] value[%s]\n", k, v)
	// 	    listSlice, ok := v.([]interface{})
	// 	    if !ok {
	//       		panic(err)
	// 	    }

	// 	    for _, v2 := range listSlice {
 //        		fmt.Println(v2)
	// 	    }
	// 	}

	// 	c.JSON(200, result["server"])
	// })

	// router.GET("/server-list", func(c *gin.Context) {
	// 	file, err := os.Open("config.json")
	// 	if err != nil {
	//       panic(err)
	//     }
	// 	defer file.Close()
	//   	buf := new(bytes.Buffer)
	//     buf.ReadFrom(file)
	//     contents := buf.String()
	// 	c.String(200, contents)
	// })

	router.GET("/dummy", func(c *gin.Context) {
		file, err := os.Open("config.json")
		if err != nil {
	      panic(err)
	    }
		defer file.Close()

		print(err);

	  	// buf := new(bytes.Buffer)
	    // buf.ReadFrom(file)
	    // contents := buf.String()
		
	  	byteValue, _ := ioutil.ReadAll(file)

	    var result map[string]interface{}
	    json.Unmarshal([]byte(byteValue), &result)
    	fmt.Println("Results")
    	fmt.Println(result["server"])

		// var serverDetails ServerDetails
		// json.Unmarshal(buf, &serverDetails)

		// for i := 0; i < len(serverDetails); i++ {
		//     fmt.Println("User Type: " + serverDetails.ServerDetails[i].Name)
		// }
		//var result map[string]interface{}
		//json.Unmarshal([]byte(file), &result)
		
	 	// ServerDetails := []ServerDetail{{Id: "ravi", Name: "travel", Address: "golang"},{Id: "ab", Name: "cd", Address: "ef"}}
	 	//fmt.Println(result)
	 	//fmt.Println(err)
	 	//k1 := gjson.Get(json, result)

		// c.JSON(200, file)
		c.String(200, "test")
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
	router.Run("localhost:8080")
}

func getCnonce() string {
    b := make([]byte, 8)
    io.ReadFull(rand.Reader, b)
    return fmt.Sprintf("%x", b)[:16]
}

func getMD5(text string) string {
    hasher := md5.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}

func getDigestAuthrization(digestParts map[string]string) string {
    d := digestParts
    ha1 := getMD5(d["username"] + ":" + d["realm"] + ":" + d["password"])
    ha2 := getMD5(d["method"] + ":" + d["uri"])
    nonceCount := 00000001
    cnonce := getCnonce()
    response := getMD5(fmt.Sprintf("%s:%s:%v:%s:%s:%s", ha1, d["nonce"], nonceCount, cnonce, d["qop"], ha2))
    authorization := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", cnonce="%s", nc="%v", qop="%s", response="%s"`,
        d["username"], d["realm"], d["nonce"], d["uri"], cnonce, nonceCount, d["qop"], response)
    return authorization
}

func digestParts(resp *http.Response) map[string]string {
    result := map[string]string{}
    if len(resp.Header["Www-Authenticate"]) > 0 {
        wantedHeaders := []string{"nonce", "realm", "qop"}
        responseHeaders := strings.Split(resp.Header["Www-Authenticate"][0], ",")
        for _, r := range responseHeaders {
            for _, w := range wantedHeaders {
                if strings.Contains(r, w) {
                    if(strings.Contains(r, `"`)){
                    	var a = strings.Split(r, `"`)
                    	result[w] = a[1];
                    } else {
                    	var b = strings.Split(r, `=`)
                    	result[w] = b[1];
                    }
					//fmt.Println("Fin >> ", w, result[w] )
                }
            }
        }
    }
    return result
}

func getAuth(url string, username string, password string ) (version string){

	os.Setenv("HTTP_PROXY", "http://10.64.100.223:8080")
	os.Setenv("HTTPS_PROXY", "http://10.64.100.223:8080")

    var query = []byte("")
    req, err := http.NewRequest("GET", url, bytes.NewBuffer(query))
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("URL >> ", url)

    fmt.Println("----------------Challenge Resp----------")
    fmt.Println("RESP >> ",resp)

  	digestParts := digestParts(resp)


    // digestParts := map[string]string{}
    digestParts["uri"] = "/management"
    digestParts["domain"] = "/management"
    digestParts["method"] = "GET"
    digestParts["username"] = username
    digestParts["password"] = password

    fmt.Println("----------------Making Digest----------")
    fmt.Println("Digest >> ",digestParts)

    var digData = getDigestAuthrization(digestParts)
    fmt.Println("Digest Compiled >> ",digData)
    return digData;
}

func getServerVersion(url string, username string, password string ) (version string){
	
	os.Setenv("HTTP_PROXY", "http://10.64.100.223:8080")
	os.Setenv("HTTPS_PROXY", "http://10.64.100.223:8080")

    client := &http.Client{}	
    var query = []byte("")
    req, err := http.NewRequest("GET", url, bytes.NewBuffer(query))
    req.Header.Set("Authorization", getAuth(url, username, password))
    req.Header.Set("Content-Type", "application/json")

    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    fmt.Println("----------------After Challenge----------")
    fmt.Println(resp)

    fmt.Println("StatusCode >> ", resp.StatusCode )
    fmt.Println("Status >> ", resp.Status)
    fmt.Println("Header >> ", resp.Header)
    fmt.Println("-----------------------------\n\n\n")
    body, _ := ioutil.ReadAll(resp.Body)
    //fmt.Println("response Body:", string(body))
    return string(body);
}

func getLastDeployed(url string, username string, password string ) (version string){
	
	os.Setenv("HTTP_PROXY", "http://10.64.100.223:8080")
	os.Setenv("HTTPS_PROXY", "http://10.64.100.223:8080")

    client := &http.Client{}
    var query = []byte("")
    req, err := http.NewRequest("GET", url, bytes.NewBuffer(query))
    req.Header.Set("Authorization", getAuth(url, username, password))
    req.Header.Set("Content-Type", "application/json")

    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    fmt.Println("----------------After Challenge----------")
    fmt.Println(resp)

    fmt.Println("StatusCode >> ", resp.StatusCode )
    fmt.Println("Status >> ", resp.Status)
    fmt.Println("Header >> ", resp.Header)
    fmt.Println("-----------------------------\n\n\n")
    body, _ := ioutil.ReadAll(resp.Body)
    //fmt.Println("response Body:", string(body))
    return string(body);
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