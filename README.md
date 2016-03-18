##tcpdemo是什么?

golang书写的示例    
阐述tcp client 与 tcp server 通信。

其中client与server间的通信协议格式采用TLV(Type,Length,Value)。

client传送json格式的value,server收到消息后将json数据解析到结构体内。  
tcpdemo采用seelog存储log。

keywords:  
* json解析  
* socket 编解码  
* tcp demo


##感谢

感谢网上已有的资料，从这里面学到了有用的知识，列位有需求的也拿来做个参考：

* [从零开始写Socket Server](http://studygolang.com/articles/4998)
* [Go语言TCP Socket编程](http://www.tuicool.com/articles/EbAFzei)


##关于作者
* 昵称：强尼老三
* 邮箱：whereshallyoube@163.com
* 英文名：Johnny-three
