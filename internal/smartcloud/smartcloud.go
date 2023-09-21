package smartcloud

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"smartcloud/internal/config"
	"smartcloud/internal/core"
	"smartcloud/internal/database"
	"smartcloud/internal/validate"

	"golang.org/x/net/webdav"
)

var fs *webdav.Handler
var validator *validate.ValidatorService
var service *core.WebDav

func Run() {

	log.Println("smartcloud is running")
	config := config.Read()
	fs = &webdav.Handler{
		FileSystem: webdav.Dir(config.Rootdir),
		LockSystem: webdav.NewMemLS(),
	}

	log.Print(config.DbConfig.Url)
	database.Init_database(config.DbConfig)
	validator = validate.NewValidatorService(config.Auth)
	service = core.NewWebDav(config.Permission, fs)

	port := fmt.Sprintf(":%s", config.Server.Port)

	http.HandleFunc(config.Rootpath, handler)
	if config.Server.Protocol == "https" {
		// 启用客户端证书验证
		serverCert := config.Server.Tls.Cert // 服务器证书文件路径
		serverKey := config.Server.Tls.Key   // 服务器私钥文件路径

		// 加载服务器证书和私钥
		cert, err := tls.LoadX509KeyPair(serverCert, serverKey)
		if err != nil {
			log.Fatal(err)
		}

		// 创建自定义的 TLS 配置
		tlsConfig := &tls.Config{
			RootCAs:      validate.Roots,
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.VerifyClientCertIfGiven, // 要求并验证客户端证书
		}

		// 创建带有自定义 TLS 配置的 HTTP 服务器
		server := &http.Server{
			Addr:      port,
			Handler:   nil, // 使用默认路由处理程序
			TLSConfig: tlsConfig,
		}

		// 启动服务器并监听连接
		log.Printf("Server listening on port %s...", port)
		log.Fatal(server.ListenAndServeTLS("", ""))
	} else {
		// 创建带有自定义 TLS 配置的 HTTP 服务器
		server := &http.Server{
			Addr:    port,
			Handler: nil, // 使用默认路由处理程序
		}

		// 启动服务器并监听连接
		log.Printf("Server listening on port %s...", port)
		log.Fatal(server.ListenAndServe())
	}

}

func handler(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		// No Authorization header found, return a 401 Unauthorized response
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// if !validator.Validate(r) {
	// 	log.Println("validate failed")
	// 	return
	// }

	service.DoService(w, r)

	fs.ServeHTTP(w, r)
	// 处理客户端请求
	// fmt.Fprintln(w, "Hello, TLS client!")
}
