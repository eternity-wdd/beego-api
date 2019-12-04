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

- 请求参数示例：
```
  {"area":"朝阳区","latitude":"39.82","longitude":"118.45"}      格式根据访问的接口进行调整, 数值可以参考
```

- 其他：
```
  其中model是使用beego自带工具自动生成的, 使用方法详见：https://www.cnblogs.com/lz0925/p/11910025.html， 有问题可以留言或者在博客中加我VX

  项目中使用的地图接口用的我自己的key，如果过期，可以通知我，或者自行替换自己的key。
  项目中使用的天气接口本来是公司的，我将返回的数据保存了，通过我的服务器返回数据（固定）供大家使用。
```
