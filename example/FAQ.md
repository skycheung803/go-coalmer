# 
[go-rod Manual](https://go-rod.github.io/i18n/zh-CN/#/)
# 参考
[使用 Golang Rod 解析浏览器中动态渲染的内容](https://soulteary.com/2022/12/15/rsscan-use-golang-rod-to-parse-the-content-dynamically-rendered-in-the-browser-part-4.html)

* 操作系统里本身就安装了 Chrome，那么可以使用 --remote-debugging-port=9222 --headless 参数启动一个可以被 Rod 使用的 Headless 浏览器容器环境。以 macOS 为例，完整命令如下：（其他系统需要调整路径）
$ /Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome --remote-debugging-port=9222 --headless
* 如果你需要浏览器访问的地址需要代理服务器或者堡垒机中转，那么你还可以在配置中添加 --proxy-server 参数：
$  /Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome --remote-debugging-port=9222 --headless --proxy-server=10.11.12.90:8001



# ubuntu install Dynamic Library
apt install  libatk1.0-0 libatk-bridge2.0-0 libcups2 libxcomposite1  libxdamage1 libxfixes3 libxrandr2 libgbm1 libpango-1.0-0 libcairo2 libasound2t64