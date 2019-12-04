# beego-api

- 一个使用beego写的API

- 支持Api日志

- 支持Swagger注解文档

### 使用说明

- 导入数据表到mysql中（建议不低于MySQL5.7）

```
  数据库文件：beego_api.sql
```

- 修改数据库配置
```
  conf文件夹下的 app.conf文件
  
  dbuser = root
  dbpassword = root
  dbhost = localhost
  dbport = 3306
  dbname = beego_api
  dbprefix = bg_
```

- 安装所有依赖：

```
  go get -v ./...   
```
- 运行：  
```
  bee run -downdoc=true -gendoc=true
```

- 访问：
```
localhost:8080/swagger,
```
