# luna_http
lunahttp|一个抗指纹浏览器自动化测试框架|爬虫框架|anti-fingerprint


此项目为https://github.com/musiclover789/luna   golang项目的http 方式封装。

python版本 https://github.com/musiclover789/luna_python 也是基于lunahttp的再次封装。



此项目存在的意义在于、其他语言使用者，可以自行封装其他语言版本的luna框架。

因为仅通过http调用即可完成封装。



#### 如何使用：

1、熟悉golang。

​		如果您熟悉golang、建议您直接使用https://github.com/musiclover789/luna
    版本即可。

2、您不熟悉golang

​		您可以直接调用https://github.com/musiclover789/luna_http/tree/main/main
  包下面的可执行程序；

其中包括了、windows的exe、Mac inter芯片、arm芯片、linux等版本。



##### 原理：

​		当您通过命令行的方式调用这个可执行程序的时候、luna_http 会启动一个http的服务。我们拿mac arm版本举例：

​	当您命令行启动程序的时候 ，他会监听您指定的端口、并相应相关请求。

```
./mac_arm_auth 8899
```

执行结果:

```
xxx@honMacBook-Pro-2 main % ./mac_arm_auth 8899
Error: luna directory not found
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /test_browser             --> main.testBrowserPost.func1 (3 handlers)
[GIN-debug] POST   /close_browser            --> main.closeBrowser.func1 (3 handlers)
[GIN-debug] POST   /new_browser              --> main.newBrowser.func1 (3 handlers)
[GIN-debug] POST   /open_page                --> main.openPage.func1 (3 handlers)
[GIN-debug] POST   /close_page               --> main.closePage.func1 (3 handlers)
[GIN-debug] POST   /get_element_position_by_css_on_page --> main.getElementPositionByCssOnPage.func1 (3 handlers)
[GIN-debug] POST   /simulate_mouse_move_to_target --> main.simulateMouseMoveToTarget.func1 (3 handlers)
[GIN-debug] POST   /simulate_mouse_move_to_element --> main.simulateMouseMoveToElement.func1 (3 handlers)
[GIN-debug] POST   /simulate_mouse_move_on_page --> main.simulateMouseMoveOnPage.func1 (3 handlers)
[GIN-debug] POST   /simulate_mouse_click     --> main.simulateMouseClickOnPage.func1 (3 handlers)
[GIN-debug] POST   /simulate_keyboard_input_on_page --> main.simulateKeyboardInputOnPage.func1 (3 handlers)
[GIN-debug] POST   /simulate_scroll_to_element_by_selector --> main.simulateScrollToElementBySelector.func1 (3 handlers)
[GIN-debug] POST   /get_first_child_element_by_css --> main.getFirstChildElementByCss.func1 (3 handlers)
[GIN-debug] POST   /get_next_sibling_element_by_css --> main.getNextSiblingElementByCss.func1 (3 handlers)
[GIN-debug] POST   /get_previous_sibling_element_by_css --> main.getPreviousSiblingElementByCss.func1 (3 handlers)
[GIN-debug] POST   /get_last_child_element_by_css --> main.getLastChildElementByCss.func1 (3 handlers)
[GIN-debug] POST   /get_parent_element_by_css --> main.getParentElementByCss.func1 (3 handlers)
[GIN-debug] POST   /get_element_by_css       --> main.getElementByCss.func1 (3 handlers)
[GIN-debug] POST   /get_all_child_element_by_css --> main.getAllChildElementByCss.func1 (3 handlers)
[GIN-debug] POST   /simulate_drag            --> main.simulateDrag.func1 (3 handlers)
[GIN-debug] POST   /upload_files             --> main.uploadFiles.func1 (3 handlers)
[GIN-debug] POST   /run_js                   --> main.runJS.func1 (3 handlers)
[GIN-debug] POST   /run_js_sync              --> main.runJSSync.func1 (3 handlers)
[GIN-debug] POST   /get_cookie               --> main.getCookie.func1 (3 handlers)
[GIN-debug] POST   /set_cookie               --> main.setCookie.func1 (3 handlers)
[GIN-debug] GET    /get_html                 --> main.getHtml.func1 (3 handlers)
[GIN-debug] GET    /exit                     --> main.exit.func1 (3 handlers)
[GIN-debug] GET    /ready                    --> main.ready.func1 (3 handlers)
[GIN-debug] POST   /open_page_and_listen_network --> main.openPageAndListenNetwork.func1 (3 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :8899
```



#### 我如何封装？

​	您可以封装任何语言的框架、仅需要参考如下逻辑：

​	1、用您的语言 通过命令行方式启动对应平台的可执行程序，这个程序会启动一个http服务，并且监听您指定的端口，并且对您发出的http请求、来执行操作浏览器动作。

​	2、将相关http相应服务、封装成您的对应函数调用

 

#### 都有哪些接口目前？

您可以自行查看源码的lunahttp/luna_http/main/main.go部分

也可以参考下面示例列表:

```
接口：new_browser
方式:post
请求url示例
http://127.0.0.1:9876/new_browser
作用：启动一个新的浏览器
传入参数格式:application/json
json 示例:
{
        "chromiumPath": "path/to/chromium",        // chromium可执行文件的路径 必选<必须要有>
        "cachePath": "path/to/cache",        // 缓存文件的路径 可选
        "imgPath": "path/to/img",        // 截图文件的路径 可选
        "headless": true,        // 是否使用无头模式 可选
        "proxyStr": "",        // 代理服务器的地址 可选
        "windowSize": {        // 窗口大小 可选
                "width": 800,        // 窗口宽度
                "height": 600        // 窗口高度
        },
        "fingerprint": ["value1", "value2", "value3"]        // 指纹数组 可选
}
返回示例:
{
  "message": "very nice",
  "chrome_id": "your_chrome_id", //返回的浏览器id 可以用于后续操作此浏览器实例
  "status": 200
}
```



```
接口：close_browser
方式:post
请求url示例
http://127.0.0.1:9876/close_browser
作用：关闭指定浏览器实例
传入参数格式:application/x-www-form-urlencoded
form-data 示例:
chrome_id=your_chrome_id //你需要关闭的chrome id

返回示例:
{
  "message": "very nice",
  "status": 200
}
```






其他的自己参考https://github.com/musiclover789/luna_http/tree/main/main 里面的代码，自行封装即可。

太多了，就不一一列出了。

