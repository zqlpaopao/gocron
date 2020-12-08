# cron安装

1. 安装mysql

   ```
   docker run -p 12345:3306 --name=mysql -v /Users/zhangsan/Desktop/workspace-app/mysql/conf:/etc/mysql/conf.d -v /Users/zhangsan/Desktop/workspace-app/mysql/logs:/logs -v /Users/zhangsan/Desktop/workspace-app/mysql/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 -d --net es_elastic --ip 172.19.0.6 mysql:5.7.25
   ```

2. 安装corn

   ```
   docker run --name gocron --link mysql:db -p 5920:5920 --net es_elastic -d ouqg/gocron
   ```

3. 访问

   ```
   127.0.0.1:5920
   ```

4. 配置mysql

   ```
   本地连接docker的mysql
   数据名 crm_cron
   字符集 utf8mb4
   排序规则 utf8mb4_general_ci
   ```

   ![image-20201208105503888](cronweb使用/image-20201208105503888.png)

5. 



## 源码安装

- <font color=red size=5x>docker安装好mysql，继续下面步骤</font>

1. 下载源码

   ```
   git clone https://github.com/zqlpaopao/gocron.git
   ```

2. 启动程序

   ```
   go run cmd/gocron/gocron.go web
   ```

3. 访问地址

   ```
   127.0.0.1:5921
   ```

4. 配置信息

   ![配置cron](cronweb使用/image-20201208112625076.png)

5. 点击安装跳转

   ![image-20201208112956554](cronweb使用/image-20201208112956554.png)

6. 登陆

   ![image-20201208113030462](cronweb使用/image-20201208113030462.png)



7. 此时数据库详情

![image-20201208113100508](cronweb使用/image-20201208113100508.png)

# 使用

## 简单任务配置

秒 分 时 天 月 周

![image-20201208172842394](cronweb使用/image-20201208172842394.png)

<font color=green size=5x>查看执行日志</font>

![image-20201208173043744](cronweb使用/image-20201208173043744.png)

## 任务节点配置

先配置任务节点

![image-20201208173348524](cronweb使用/image-20201208173348524.png)

查看任务 新增

![image-20201208173420495](cronweb使用/image-20201208173420495.png)

