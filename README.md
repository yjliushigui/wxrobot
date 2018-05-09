# 微信机器人协议详情


		时间戳：timestamp
		时间戳取反:^timestamp
		Uuid:Uuid
		Skey:Skey
		sid:sid
		Uin:uin
		回调地址:redirectURL
		请求地址头:httpHeader
		请求必须包体头:BaseRequest
		通行证:passticket


## 1、登录二维码(获取UUID)

		微信登录二维码是使用微信的登录接口生成的uuid继而生成相应的二维码废话不多说，直接上接口


### （1）获取uuid


<table>
	<tr>
	<td>URL</td>
	<td><code>https://login.wx.qq.com/jslogin?appid=wx782c26e4c19acffb&fun=new& lang=zh_CN&_=时间戳</code></td>
	</tr>
	<tr>
	<td>Method</td>
	<td><code>GET</code></td>
	</tr>
</table>


### (2)生成二维码


	https://login.weixin.qq.com/l/+uuid
	根据语言不同去生成二维码


## 2、登录


### (1)扫码登录

<table>
	<tr>
	<td>URL</td>
	<td><code>https://login.wx.qq.com/cgi-bin/mmwebwx-bin/login?loginicon=true&uuid=Uuid&tip=0&r=^timestamp&_=timestamp</code></td>
	</tr>
	<tr>
	<td>Method</td>
	<td><code>GET</code></td>
	</tr>
</table>


		扫码成功返回code为201并返回相应用户的头像的<code>base64</code>编码，登录成功返回信息code为200并且返回回调地址<code>redirectURL</code>，回调地址根据不同的微信号申请时间分为两个版本，以<code>wx.qq.com</code>打头为2012年以前的微信，称为微信第一版本，之后的称为微信第二版本，以<code>wx2.qq.com</code>,获取到的地址头为后面所有请求的地址头，根据微信版本选择自己的请求头，其余code皆为失败,这里将获取到的请求地址头定义为<code>httpHeader</code>



### (2)获取回调信息(XML)

<table>
	<tr>
	<td>URL</td>
	<td><code>redirectURL</code></td>
	</tr>
	<tr>
	<td>Method</td>
	<td><code>GET</code></td>
	</tr>
</table>


		解析返回的XML,<code>Ret</code>不为0失败，失败原因(1)参数不对值:1;(2)登录超时值:1101;(3):非法操作:120x--等待2-4个小时以后重新尝试。后面的流程失败成功都以<code>Ret</code>进行判断，以下参数带<code>*</code>的都是重点参数，直接影响返回的数据成功或者失败
				
		*BaseRequest{Uin:"", Sid:"", Skey:"",DeviceID:"16位随机数"}
		*passticket:""
		在这一步还需要将获取的cookie存储下来

## 3、初始化


		这一部分开始url后面的参数需要进行urlencode


<table>
	<tr>
	<td>URL</td>
	<td><code>httpHeader+webwxinit?pass_ticket=&skey=&r=^timestamp</code></td>
	</tr>
	<tr>
	<td>Method</td>
	<td><code>POST</code></td>
	</tr>
	<tr>
	<td>Param</td>
	<td><code>{Baserequest:""}</code></td>
	</tr>
	<tr>
	<td>Content-type</td>
	<td><code>application/json</code></td>
	</tr>
</table>


		初始化成功返回数据格式为JSON，返回数据包括:最近联系人ContactList、用户信息User、初始心跳数据SyncKey
		从这一步开始就需要进行和微信的心跳数据通讯


## 4、心跳和接收消息


		为什么要将接收信息和心跳放到一起进行处理呢，因为微信的心跳数据在有新消息的时候会发生变化，每当有新消息发送过来会伴随着新的心跳数据，最初的心跳数据将会被舍弃，如果长时间不更改心跳数据或者使用旧的心跳数据会造成登录超时或者非法操作等错误，严重者甚至封号，这个过程是保证你的机器人长时间稳定运行的重要环节


### (1)心跳


		注:SynckeyHttpHeader根据不同版本的微信号使用的路径不一样对应的版本地址为
			v1:https://webpush.wx.qq.com/cgi-bin/mmwebwx-bin/
			v2:https://webpush.wx2.qq.com/cgi-bin/mmwebwx-bin/,
			synckey的值为初始化中的Synckey的key和value拼接成为:key1_value1|key2_value2|....|keyN_valueN
<table>
	<tr>
	<td>URL</td>
	<td><code>SynckeyHttpHeader+synccheck?pass_ticket=&skey=&sid=&synckey=&uin=&deviceid=&r=^timestamp&_=timestamp</code></td>
	</tr>
	<tr>
	<td>Method</td>
	<td><code>POST</code></td>
	</tr>
	<tr>
	<td>Param</td>
	<td><code>{Baserequest:""}</code></td>
	</tr>
	<tr>
	<td>Content-type</td>
	<td><code>application/json</code></td>
	</tr>
</table>


		返回retcode和selector，如果retcode不为0即心跳失败,selector 标志消息类型，0为没有消息，第二次心跳额synckey不变，不为0就代表有新的消息，调用消息获取接口获取新的消息,心跳数据发生变化


### (2)消息获取


<table>
	<tr>
	<td>URL</td>
	<td><code>httpHeader+webwxsync?pass_ticket=&skey=&sid=</code></td>
	</tr>
	<tr>
	<td>Method</td>
	<td><code>POST</code></td>
	</tr>
	<tr>
	<td>Param</td>
	<td><code>{Baserequest:"",SyncKey{List:{key:"",value:""},Count:count(List)},rr:^timestamp}</code></td>
	</tr>
	<tr>
	<td>Content-type</td>
	<td><code>application/json</code></td>
	</tr>
</table>


		此接口返回最新的消息AddMsgList和新的心跳数据Synckey
		后面逐步完善和更新代码和注解，目前php代码比较完善!




## 捐赠


<center><img src="https://i.imgur.com/zMFLzt9.jpg" width="200" align=center /></center>
