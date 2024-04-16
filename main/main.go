package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/musiclover789/luna/base_devtools/network"
	"github.com/musiclover789/luna/devtools"
	"github.com/musiclover789/luna/luna_utils"
	"github.com/musiclover789/luna/protocol"
	"github.com/tidwall/gjson"
	"luna_http/common/brower_map"
	"luna_http/tool"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	port := "9876"
	// 获取命令行参数
	args := os.Args

	// 第一个参数是程序名称，忽略它
	if len(args) > 1 {
		// 遍历打印其他参数
		for i, arg := range args[1:] {
			port = arg
			break
			fmt.Printf("参数 %d: %s\n", i+1, arg)
		}
	} else {
		fmt.Println("没有提供任何参数")
	}

	// 创建 Gin 的默认路由引擎
	router := gin.Default()
	//测试打开浏览器
	testBrowserPost(router)
	//打开浏览器

	//关闭浏览器
	closeBrowser(router)
	//测试打开浏览器
	newBrowser(router)
	//关闭浏览器

	//打开页面
	openPage(router)
	//关闭页面
	closePage(router)
	//操控、鼠标、键盘、dom
	//dom tree操作函数
	getElementPositionByCssOnPage(router)
	simulateMouseMoveToTarget(router)
	simulateMouseMoveToElement(router)
	simulateMouseMoveOnPage(router)
	simulateMouseClickOnPage(router)
	simulateKeyboardInputOnPage(router)

	simulateScrollToElementBySelector(router)
	//目前总共有几个页面
	getFirstChildElementByCss(router)
	getNextSiblingElementByCss(router)
	getPreviousSiblingElementByCss(router)
	getLastChildElementByCss(router)
	getParentElementByCss(router)
	getElementByCss(router)
	getAllChildElementByCss(router)

	simulateDrag(router)
	uploadFiles(router)
	//runjs
	runJS(router)

	runJSSync(router)
	//getcookie\set cookie
	getCookie(router)
	setCookie(router)
	//获取页面源代码
	getHtml(router)
	//退出
	exit(router)
	//
	ready(router)
	//
	openPageAndListenNetwork(router)
	// 启动 HTTP 服务器，监听端口 8080
	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}

/*
**
打开浏览器、
作为测试使用
*/
func testBrowserPost(router *gin.Engine) {
	router.POST("/test_browser", func(c *gin.Context) {
		// 获取表单数据
		chromiumPath := c.PostForm("chromium_path")
		url := c.PostForm("url")
		//-----------------大逻辑---------------end---------
		if len(chromiumPath) == 0 || len(url) == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "You forgot to pass the chromium_path and url parameters.",
				"status":  200,
			})
			return
		}
		//
		_, browserObj := devtools.NewBrowser(chromiumPath, &devtools.BrowserOptions{})
		browserObj.OpenPage(url)

		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"status":  200,
		})
	})
}

/*
**
初始化浏览器、
*/
func newBrowser(router *gin.Engine) {
	//luna_utils.KillProcess()
	router.POST("/new_browser", func(c *gin.Context) {
		// 获取请求body数据
		body, err := c.GetRawData()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		json := gjson.ParseBytes(body)

		// 获取表单数据
		chromiumPath := json.Get("chromiumPath").String()

		cachePath := json.Get("cachePath").String()
		img_path := json.Get("imgPath").String()
		headless := json.Get("headless").Bool()
		proxy_str := json.Get("proxyStr").String()
		var windowSizePtr *devtools.WindowSize

		if json.Get("windowSize").Exists() {
			windowSizePtr = &devtools.WindowSize{
				Width:  int(json.Get("windowSize").Get("width").Int()),
				Height: int(json.Get("windowSize").Get("height").Int()),
			}
		}
		fmt.Println()
		var fingerprintArray []string

		// 使用 ForEach() 函数遍历结果，并将每个值添加到字符串切片中
		json.Get("fingerprint").ForEach(func(key, value gjson.Result) bool {
			fingerprintArray = append(fingerprintArray, value.String())
			return true // 继续遍历
		})

		//-----------------大逻辑---------------end---------
		if len(chromiumPath) == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "You forgot to pass the chromium_path parameters.",
				"status":  200,
			})
			return
		}
		//
		chrome_id := tool.GenerateIncrementalID()
		var wgoutside sync.WaitGroup // 同步等待
		wgoutside.Add(1)
		var counter int32
		go func() {
			var wg sync.WaitGroup // 同步等待
			wg.Add(1)             // 增加等待的数量
			err, browserObj := devtools.NewBrowser(chromiumPath, &devtools.BrowserOptions{
				CachePath:   cachePath,
				ImgPath:     img_path,
				Headless:    headless,
				ProxyStr:    proxy_str,
				WindowSize:  windowSizePtr,
				Fingerprint: fingerprintArray,
			})
			if err != nil {
				atomic.AddInt32(&counter, 1)
				return
			}
			brower_map.Push(chrome_id, browserObj)
			wgoutside.Done()
			wg.Wait()
		}()
		wgoutside.Wait()
		value := atomic.LoadInt32(&counter)
		if value == 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"message":   "new_browser error",
				"chrome_id": chrome_id,
				"status":    200,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"message":   "very nice",
			"chrome_id": chrome_id,
			"status":    200,
		})
	})
}

/*
**
初始化浏览器、
*/
func openPage(router *gin.Engine) {
	router.POST("/open_page", func(c *gin.Context) {
		// 获取表单数据
		chromiumPath := c.PostForm("url")
		url := c.PostForm("url")
		chrome_id := c.PostForm("chrome_id")

		//-----------------大逻辑---------------end---------
		if len(chromiumPath) == 0 || len(url) == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "You forgot to pass the chromium_path and url parameters.",
				"status":  200,
			})
			return
		}
		page_id := tool.GenerateIncrementalID()
		browserObj := brower_map.Get(chrome_id).(*devtools.Browser)
		err, page := browserObj.OpenPageAndListen(url, func(devToolsConn *protocol.DevToolsConn) {

		})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "Unknown error.",
				"status":  403,
			})
			return
		}
		brower_map.Push(page_id, page)
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"page_id": page_id,
			"status":  200,
		})
	})
}

/*
**
初始化浏览器、
*/
func closePage(router *gin.Engine) {
	router.POST("/close_page", func(c *gin.Context) {
		// 获取表单数据
		page_id := c.PostForm("page_id")
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(page_id).(*devtools.Page)
		pageObj.Close()
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"status":  200,
		})
		return
	})
}

/*
**
初始化浏览器、
*/
func closeBrowser(router *gin.Engine) {
	router.POST("/close_browser", func(c *gin.Context) {
		// 获取表单数据
		chrome_id := c.PostForm("chrome_id")
		fmt.Println("chrome_id:", chrome_id)
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(chrome_id).(*devtools.Browser)
		pageObj.Close()
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"status":  200,
		})
		return
	})
}

func exit(router *gin.Engine) {
	router.GET("/exit", func(c *gin.Context) {
		os.Exit(0)
		return
	})
}

func ready(router *gin.Engine) {
	router.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"status":  200,
		})
		return
	})
}

func getHtml(router *gin.Engine) {
	router.GET("/get_html", func(c *gin.Context) {
		// 获取表单数据
		page_id := c.Query("page_id")
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(page_id).(*devtools.Page)

		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"data":    pageObj.GetHtml().Get("result.outerHTML").String(),
			"status":  200,
		})
		return
	})
}

// CookieResponse 结构体表示cookie响应数据
type CookieResponse struct {
	Cookies []struct {
		Name   string `json:"name"`
		Value  string `json:"value"`
		Domain string `json:"domain"`
	} `json:"cookies"`
}

func getCookie(router *gin.Engine) {
	router.POST("/get_cookie", func(c *gin.Context) {
		// 获取表单数据
		page_id := c.PostForm("page_id")
		url := c.PostForm("url")
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(page_id).(*devtools.Page)
		cookies, _ := network.GetCookies(pageObj.DevToolsConn, []string{url})
		var cookieResponse CookieResponse
		for _, result := range gjson.Parse(luna_utils.FormatJSONAsString(cookies)).Get("result.cookies").Array() {
			cookie := struct {
				Name   string `json:"name"`
				Value  string `json:"value"`
				Domain string `json:"domain"`
			}{
				Name:   result.Get("name").String(),
				Value:  result.Get("value").String(),
				Domain: result.Get("domain").String(),
			}
			cookieResponse.Cookies = append(cookieResponse.Cookies, cookie)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"data":    cookieResponse,
			"status":  200,
		})
		return
	})
}

func setCookie(router *gin.Engine) {
	router.POST("/set_cookie", func(c *gin.Context) {
		// 获取表单数据
		pageId := c.PostForm("page_id")
		key := c.PostForm("key")
		value := c.PostForm("value")
		domain := c.PostForm("url")
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(pageId).(*devtools.Page)
		network.SetCookie(pageObj.DevToolsConn, key, value, domain)
		network.SetCookieByURL(pageObj.DevToolsConn, key, value, domain)
		//有点麻烦，后面在写
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"status":  200,
		})
		return
	})
}

/***

//p1.GetElementPositionByXpathOnPage()
//p1.GetFirstChildElementByCss()
//p1.GetNextSiblingElementByCss()
//p1.GetPreviousSiblingElementByCss()
//p1.GetElementByCss()
//p1.GetLastChildElementByCss()
//p1.GetAllChildElementByCss()
//p1.GetParentElementByCss()
//操作
//p1.SimulateMouseMoveToTarget()
//p1.SimulateMouseMoveToElement()
//p1.SimulateMouseMoveOnPage()

//p1.SimulateDrag()

//p1.UploadFiles()

//p1.SimulateMouseClickOnPage()
//p1.SimulateKeyboardInputOnPage()

//p1.SimulateScrollToElementBySelector()

//p1.RunJS()
//p1.RunJSSync()

//p1.GetCurrentURL()

//p1.SetViewportSizeAndScale()
//p1.SetViewportSize()
*/

/*
*
runJS、
*/
func runJS(router *gin.Engine) {
	router.POST("/run_js", func(c *gin.Context) {
		// 获取表单数据
		pageId := c.PostForm("page_id")
		js := c.PostForm("js")
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(pageId).(*devtools.Page)
		pageObj.RunJS(js)
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"status":  200,
		})
		return
	})
}

/*
*
runJS、
*/
func runJSSync(router *gin.Engine) {
	router.POST("/run_js_sync", func(c *gin.Context) {
		// 获取表单数据
		pageId := c.PostForm("page_id")
		js := c.PostForm("js")
		time_out, err := strconv.Atoi(c.PostForm("time_out"))
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Timeout conversion failed",
				"status":  403,
			})
			return
		}
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(pageId).(*devtools.Page)
		err, result := pageObj.RunJSSync(js, time.Second*time.Duration(time_out))
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"data":    result.String(),
			"status":  200,
		})
		return
	})
}

func getElementPositionByCssOnPage(router *gin.Engine) {
	router.POST("/get_element_position_by_css_on_page", func(c *gin.Context) {
		// 获取表单数据
		pageId := c.PostForm("page_id")
		selector := c.PostForm("selector")
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(pageId).(*devtools.Page)
		err, x, y := pageObj.GetElementPositionByCssOnPage(selector)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"x":       x,
			"y":       y,
			"status":  200,
		})
		return
	})
}

/*
*
//p1.SimulateMouseClickOnPage()
*/
func simulateMouseClickOnPage(router *gin.Engine) {
	router.POST("/simulate_mouse_click", func(c *gin.Context) {
		// 获取表单数据
		pageId := c.PostForm("page_id")
		x := c.PostForm("x")
		y := c.PostForm("y")
		xf, err := strconv.ParseFloat(x, 64)
		if err != nil {
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{
					"message": err.Error(),
					"status":  403,
				})
				return
			}
		}
		yf, err := strconv.ParseFloat(y, 64)
		if err != nil {
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{
					"message": err.Error(),
					"status":  403,
				})
				return
			}
		}
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(pageId).(*devtools.Page)
		pageObj.SimulateMouseClickOnPage(xf, yf)
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"status":  200,
		})
		return
	})
}

/*
*
p1.SimulateMouseMoveToTarget()
*/
func simulateMouseMoveToTarget(router *gin.Engine) {
	router.POST("/simulate_mouse_move_to_target", func(c *gin.Context) {
		// 获取表单数据
		pageId := c.PostForm("page_id")
		x := c.PostForm("x")
		y := c.PostForm("y")
		xf, err := strconv.ParseFloat(x, 64)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		yf, err := strconv.ParseFloat(y, 64)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(pageId).(*devtools.Page)
		err = pageObj.SimulateMouseMoveToTarget(xf, yf)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"status":  200,
		})
		return
	})
}

func simulateMouseMoveToElement(router *gin.Engine) {
	router.POST("/simulate_mouse_move_to_element", func(c *gin.Context) {
		// 获取表单数据
		pageId := c.PostForm("page_id")
		selector := c.PostForm("selector")
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(pageId).(*devtools.Page)
		err, x, y := pageObj.SimulateMouseMoveToElement(selector)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"x":       x,
			"y":       y,
			"status":  200,
		})
		return
	})
}

func simulateMouseMoveOnPage(router *gin.Engine) {
	router.POST("/simulate_mouse_move_on_page", func(c *gin.Context) {
		// 获取表单数据
		pageId := c.PostForm("page_id")
		sx := c.PostForm("sx")
		sy := c.PostForm("sy")
		ex := c.PostForm("ex")
		ey := c.PostForm("ey")
		sxf, err := strconv.ParseFloat(sx, 64)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		syf, err := strconv.ParseFloat(sy, 64)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		exf, err := strconv.ParseFloat(ex, 64)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		eyf, err := strconv.ParseFloat(ey, 64)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(pageId).(*devtools.Page)
		pageObj.SimulateMouseMoveOnPage(sxf, syf, exf, eyf)
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"status":  200,
		})
		return
	})
}

/*
**
SimulateKeyboardInputOnPage
*/
func simulateKeyboardInputOnPage(router *gin.Engine) {
	router.POST("/simulate_keyboard_input_on_page", func(c *gin.Context) {
		// 获取表单数据
		pageId := c.PostForm("page_id")
		text := c.PostForm("text")
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(pageId).(*devtools.Page)
		pageObj.SimulateKeyboardInputOnPage(text)
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"status":  200,
		})
		return
	})
}

/*
*
SimulateScrollToElementBySelector
*/
func simulateScrollToElementBySelector(router *gin.Engine) {
	router.POST("/simulate_scroll_to_element_by_selector", func(c *gin.Context) {
		// 获取表单数据
		pageId := c.PostForm("page_id")
		selector := c.PostForm("selector")
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(pageId).(*devtools.Page)
		err := pageObj.SimulateScrollToElementBySelector(selector)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"status":  200,
		})
		return
	})
}

/*
**
GetFirstChildElementByCss
*/
func getFirstChildElementByCss(router *gin.Engine) {
	router.POST("/get_first_child_element_by_css", func(c *gin.Context) {
		// 获取表单数据
		pageId := c.PostForm("page_id")
		selector := c.PostForm("selector")
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(pageId).(*devtools.Page)
		err, node := pageObj.GetFirstChildElementByCss(selector)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"data":    node,
			"status":  200,
		})
		return
	})
}

/*
*
GetNextSiblingElementByCss
*/
func getNextSiblingElementByCss(router *gin.Engine) {
	router.POST("/get_next_sibling_element_by_css", func(c *gin.Context) {
		// 获取表单数据
		pageId := c.PostForm("page_id")
		selector := c.PostForm("selector")
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(pageId).(*devtools.Page)
		err, node := pageObj.GetNextSiblingElementByCss(selector)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"data":    node,
			"status":  200,
		})
		return
	})
}

/*
*
GetPreviousSiblingElementByCss
*/
func getPreviousSiblingElementByCss(router *gin.Engine) {
	router.POST("/get_previous_sibling_element_by_css", func(c *gin.Context) {
		// 获取表单数据
		pageId := c.PostForm("page_id")
		selector := c.PostForm("selector")
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(pageId).(*devtools.Page)
		err, node := pageObj.GetPreviousSiblingElementByCss(selector)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"data":    node,
			"status":  200,
		})
		return
	})
}

/*
*
GetLastChildElementByCss
*/
func getLastChildElementByCss(router *gin.Engine) {
	router.POST("/get_last_child_element_by_css", func(c *gin.Context) {
		// 获取表单数据
		pageId := c.PostForm("page_id")
		selector := c.PostForm("selector")
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(pageId).(*devtools.Page)
		err, node := pageObj.GetLastChildElementByCss(selector)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"data":    node,
			"status":  200,
		})
		return
	})
}

/*
*
GetParentElementByCss
*/
func getParentElementByCss(router *gin.Engine) {
	router.POST("/get_parent_element_by_css", func(c *gin.Context) {
		// 获取表单数据
		pageId := c.PostForm("page_id")
		selector := c.PostForm("selector")
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(pageId).(*devtools.Page)
		err, node := pageObj.GetParentElementByCss(selector)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"data":    node,
			"status":  200,
		})
		return
	})
}

/*
*
GetParentElementByCss
*/
func getElementByCss(router *gin.Engine) {
	router.POST("/get_element_by_css", func(c *gin.Context) {
		// 获取表单数据
		pageId := c.PostForm("page_id")
		selector := c.PostForm("selector")
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(pageId).(*devtools.Page)
		err, node := pageObj.GetElementByCss(selector)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"data":    node,
			"status":  200,
		})
		return
	})
}

/*
*
 */
func getAllChildElementByCss(router *gin.Engine) {
	router.POST("/get_all_child_element_by_css", func(c *gin.Context) {
		// 获取表单数据
		pageId := c.PostForm("page_id")
		selector := c.PostForm("selector")
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(pageId).(*devtools.Page)
		err, nodes := pageObj.GetAllChildElementByCss(selector)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"data":    nodes,
			"status":  200,
		})
		return
	})
}

/**

//p1.GetCurrentURL()

//p1.SetViewportSizeAndScale()
//p1.SetViewportSize()
*/

func uploadFiles(router *gin.Engine) {
	router.POST("/upload_files", func(c *gin.Context) {
		// 获取表单数据
		pageId := c.PostForm("page_id")
		selector := c.PostForm("selector")
		file := c.PostForm("file")
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(pageId).(*devtools.Page)
		pageObj.UploadFiles(selector, []string{file})
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"status":  200,
		})
		return
	})
}

/*
*
 */
func simulateDrag(router *gin.Engine) {
	router.POST("/simulate_drag", func(c *gin.Context) {
		// 获取表单数据
		pageId := c.PostForm("page_id")
		sx := c.PostForm("sx")
		sy := c.PostForm("sy")
		ex := c.PostForm("ex")
		ey := c.PostForm("ey")
		sxf, err := strconv.ParseFloat(sx, 64)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		syf, err := strconv.ParseFloat(sy, 64)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		exf, err := strconv.ParseFloat(ex, 64)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		eyf, err := strconv.ParseFloat(ey, 64)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
				"status":  403,
			})
			return
		}
		//-----------------大逻辑---------------end---------
		pageObj := brower_map.Get(pageId).(*devtools.Page)
		pageObj.SimulateDrag(sxf, syf, exf, eyf)
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"status":  200,
		})
		return
	})
}

func openPageAndListenNetwork(router *gin.Engine) {
	router.POST("/open_page_and_listen_network", func(c *gin.Context) {
		chrome_id := c.PostForm("chrome_id")
		url := c.PostForm("url")
		port := c.PostForm("port")
		//-----------------大逻辑---------------end---------
		page_id := tool.GenerateIncrementalID()
		browserObj := brower_map.Get(chrome_id).(*devtools.Browser)
		err, page := browserObj.OpenPageAndListen(url, func(devToolsConn *protocol.DevToolsConn) {
			network.EnableNetwork(devToolsConn)
			network.RequestResponseAsync(devToolsConn, func(requestId string, request, response map[string]interface{}) {
				sb := strings.Builder{}
				requestStr := luna_utils.FormatJSONAsString(request)
				responseStr := luna_utils.FormatJSONAsString(response)

				sb.WriteString("{")
				sb.WriteString("\"request_result\":")
				sb.WriteString(requestStr)
				sb.WriteString(",")
				body, err := network.GetResponseBody(devToolsConn, requestId, time.Minute)
				if err == nil {
					rBody := luna_utils.FormatJSONAsString(body)
					sb.WriteString("\"response_body_result\":")
					sb.WriteString(rBody)
					sb.WriteString(",")
				}
				sb.WriteString("\"response_result\":")
				sb.WriteString(responseStr)
				sb.WriteString("}")
				go func() {
					defer func() {
						if r := recover(); r != nil {
							fmt.Println("Recovered from panic:", r)
						}
					}()
					resp, err := http.Post("http://127.0.0.1"+port, "application/json", bytes.NewBuffer([]byte(sb.String())))
					if err != nil {
						fmt.Printf("Error sending POST request: %v\n", err)
						return
					}
					defer resp.Body.Close()
				}()
				/**
				将内容放入到缓存中;
				*/
			})
		})
		fmt.Println("我也忘记了是不是异步的....")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "Unknown error.",
				"status":  403,
			})
			return
		}
		brower_map.Push(page_id, page)
		c.JSON(http.StatusOK, gin.H{
			"message": "very nice",
			"page_id": page_id,
			"status":  200,
		})
	})
}
