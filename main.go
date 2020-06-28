package main

import (
	"errors"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const connectString = "root:@/database_name?charset=utf8&parseTime=True&loc=Local"

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/pictures", "./pictures")
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("user", store))

	inits := []func() error{
		dbInit_users,
		dbInit_linux,
		dbInit_server,
		dbInit_router,
	}
	for _, f := range inits {
		err := f()
		if err != nil {
			log.Fatal(err)
		}
	}

	//ログインページ
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})

	//サインアップの処理
	r.POST("/signup", func(c *gin.Context) {
		name := c.PostForm("name")
		password := c.PostForm("password")
		err := dbSignup(name, password)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("dbSignup失敗"))
			return
		}
		c.HTML(http.StatusOK, "signup.html", gin.H{"name": name, "password": password})
	})

	//ログインの処理
	r.POST("/login", func(c *gin.Context) {
		name := c.PostForm("name")
		password := c.PostForm("password")
		session := sessions.Default(c)
		_, err := dblogin(name, password)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("ログインしてない"))
			return
		}
		session.Set("user_name", name)
		session.Save()
		c.Redirect(302, "/smoty")
	})

	//SMOTYトップページ
	r.GET("/smoty", func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("user_name") == nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("ログインしてない"))
			return
		}
		c.HTML(200, "smoty.html", gin.H{"user_name": session.Get("user_name")})
	})

	//Linux問題のページ
	r.GET("/smoty/linux", func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("user_name") == nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("ログインしてない"))
			return
		}
		linux, err := linuxGetAll()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		sort.Slice(linux, func(i, j int) bool {
			return linux[i].ID < linux[j].ID
		})
		c.HTML(200, "linux.html", gin.H{"user_name": session.Get("user_name"), "linux": linux})
	})

	//Linux問題の正解か不正化の判断
	r.POST("/smoty/linux/check/:id", func(c *gin.Context) {
		session := sessions.Default(c)
		name := session.Get("user_name")
		a := c.PostForm("anser")
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		linux, anser, err := check_linux(id, a)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.HTML(http.StatusOK, "linuxCheck.html", gin.H{"user_name": name, "linux": linux, "anser": anser, "a": a})
	})

	//Server問題のページ
	r.GET("/smoty/server", func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("user_name") == nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("ログインしてない"))
			return
		}
		server, err := serverGetAll()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		sort.Slice(server, func(i, j int) bool {
			return server[i].ID < server[j].ID
		})
		c.HTML(200, "server.html", gin.H{"user_name": session.Get("user_name"), "server": server})
	})

	//Server問題の正解か不正化の判断
	r.POST("/smoty/server/check/:id", func(c *gin.Context) {
		session := sessions.Default(c)
		name := session.Get("user_name")
		a := c.PostForm("anser")
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		server, anser, err := check_linux(id, a)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.HTML(http.StatusOK, "serverCheck.html", gin.H{"user_name": name, "server": server, "anser": anser, "a": a})
	})

	//router問題のページ
	r.GET("/smoty/router", func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("user_name") == nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("ログインしてない"))
			return
		}
		router, err := routerGetAll()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		sort.Slice(router, func(i, j int) bool {
			return router[i].ID < router[j].ID
		})
		c.HTML(200, "router.html", gin.H{"user_name": session.Get("user_name"), "router": router})
	})

	//router問題の正解か不正化の判断
	r.POST("/smoty/router/check/:id", func(c *gin.Context) {
		session := sessions.Default(c)
		name := session.Get("user_name")
		a := c.PostForm("anser")
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		router, anser, err := check_router(id, a)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.HTML(http.StatusOK, "routerCheck.html", gin.H{"user_name": name, "router": router, "anser": anser, "a": a})
	})

	//ログアウト処理
	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("user_name") == nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("ログインしてない"))
			return
		}
		session.Clear()
		session.Save()
		c.HTML(200, "logout.html", gin.H{})
	})

	//管理者ページ
	r.GET("/root", func(c *gin.Context) {
		c.HTML(http.StatusOK, "root.html", gin.H{})
	})

	//Linuxページの編集
	r.GET("/root/linux", func(c *gin.Context) {
		linux, err := linuxGetAll()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.HTML(http.StatusOK, "rootLinux.html", gin.H{
			"linux": linux,
		})
	})

	//問題の新規追加処理
	r.POST("/root/linux/new", func(c *gin.Context) {
		question := c.PostForm("question")
		anser := c.PostForm("anser")
		hint := c.PostForm("hint")
		err := linuxInsert(question, anser, hint)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Redirect(302, "/root/linux")
	})

	//問題の修正
	r.GET("/root/linux/detail/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		linux, err := linuxGetOne(id)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.HTML(200, "rootLinuxDetail.html", gin.H{"linux": linux})
	})

	//問題の削除
	r.GET("/root/linux/deleteCheck/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		linux, err := linuxGetOne(id)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.HTML(200, "rootLinuxDelete.html", gin.H{"linux": linux})
	})

	r.POST("/root/linux/delete/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		err = linuxDelete(id)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Redirect(302, "/root/linux")
	})

	//Linuxの問題の更新処理
	r.POST("/root/linux/update/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		question := c.PostForm("question")
		hint := c.PostForm("hint")
		anser := c.PostForm("anser")
		err = linuxUpdate(id, question, hint, anser)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Redirect(302, "/root/linux")
	})

	//serverページの編集
	r.GET("/root/server", func(c *gin.Context) {
		server, err := serverGetAll()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.HTML(http.StatusOK, "rootServer.html", gin.H{
			"server": server,
		})
	})

	//問題の新規追加処理
	r.POST("/root/server/new", func(c *gin.Context) {
		question := c.PostForm("question")
		anser := c.PostForm("anser")
		hint := c.PostForm("hint")
		err := serverInsert(question, anser, hint)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Redirect(302, "/root/server")
	})

	//問題の修正
	r.GET("/root/server/detail/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		server, err := serverGetOne(id)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.HTML(200, "rootServerDetail.html", gin.H{"server": server})
	})

	//問題の削除
	r.GET("/root/server/deleteCheck/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		server, err := serverGetOne(id)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.HTML(200, "rootServerDelete.html", gin.H{"server": server})
	})

	r.POST("/root/server/delete/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		err = serverDelete(id)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Redirect(302, "/root/server")
	})

	//Linuxの問題の更新処理
	r.POST("/root/server/update/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		question := c.PostForm("question")
		anser := c.PostForm("anser")
		hint := c.PostForm("hint")
		err = serverUpdate(id, question, hint, anser)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Redirect(302, "/root/server")
	})

	//routerページの編集
	r.GET("/root/router", func(c *gin.Context) {
		router, err := routerGetAll()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.HTML(http.StatusOK, "rootRouter.html", gin.H{
			"router": router,
		})
	})

	//問題の新規追加処理
	r.POST("/root/router/new", func(c *gin.Context) {
		question := c.PostForm("question")
		anser := c.PostForm("anser")
		hint := c.PostForm("hint")
		err := routerInsert(question, anser, hint)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Redirect(302, "/root/router")
	})

	//問題の修正
	r.GET("/root/router/detail/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		router, err := routerGetOne(id)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.HTML(200, "rootRouterDetail.html", gin.H{"router": router})
	})

	//問題の削除
	r.GET("/root/router/deleteCheck/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		router, err := routerGetOne(id)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.HTML(200, "rootRouterDelete.html", gin.H{"router": router})
	})

	r.POST("/root/router/delete/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		err = routerDelete(id)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Redirect(302, "/root/router")
	})

	//Linuxの問題の更新処理
	r.POST("/root/router/update/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		question := c.PostForm("question")
		anser := c.PostForm("anser")
		hint := c.PostForm("hint")
		err = routerUpdate(id, question, hint, anser)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Redirect(302, "/root/router")
	})

	//URLにないページにアクセスが来たとき
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "error.html", nil)
	})

	r.Run(":80")
}
