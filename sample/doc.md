
# demo

> This is a demo for pica.

Version: 0.0.1
Author: jeremaihloo@gmail.com jeremaihloo

## Init Scope

### Headers:

| --- | --- | -- |
| name| value | description |

| Accept | [* /*] | - |

| Accept-Language | [en-US,en;q=0.8] | - |

| Cache-Control | [max-age=0] | - |

| Connection | [keep-alive] | - |

| Content-Type | [application/json] | - |

| Referer | [http://www.baidu.com/] | - |

| User-Agent | [Pica Api Test Client/0.0.1 https://github.com/jeremaihloo/pica] | - |


## API


### GET /api/users 获取用户列表


#### Query
| --- | --- | -- |
| name | type | description |



#### Response
Headers:
| --- | --- | -- |
| name| value | description |

| Content-Length | [343] | - |

| Content-Type | [application/json; charset=utf-8] | - |

| Date | [Thu, 28 Jun 2018 05:52:58 GMT] | - |



Body:
```

Json:
{
  "items": [
    {
      "age": 23,
      "name": "jeremaihloo"
    },
    {
      "age": 10,
      "name": "test"
    },
    {
      "age": 10,
      "name": "test"
    },
    {
      "age": 10,
      "name": "test"
    },
    {
      "age": 10,
      "name": "test"
    },
    {
      "age": 10,
      "name": "test"
    },
    {
      "age": 10,
      "name": "test"
    },
    {
      "age": 10,
      "name": "test"
    },
    {
      "age": 10,
      "name": "test"
    },
    {
      "age": 10,
      "name": "test"
    },
    {
      "age": 10,
      "name": "test"
    },
    {
      "age": 10,
      "name": "test"
    },
    {
      "age": 10,
      "name": "test"
    }
  ]
}



```

### POST /api/users 新建用户


#### Query
| --- | --- | -- |
| name | type | description |


#### Body
| --- | --- | -- |
| name | type | description |

| Content-Length | [343] | - |

| Content-Type | [application/json; charset=utf-8] | - |

| Date | [Thu, 28 Jun 2018 05:52:58 GMT] | - |



#### Response
Headers:
| --- | --- | -- |
| name| value | description |

| Content-Length | [24] | - |

| Content-Type | [application/json; charset=utf-8] | - |

| Date | [Thu, 28 Jun 2018 05:52:58 GMT] | - |



Body:
```

Json:
{
  "age": 10,
  "name": "test"
}



```

