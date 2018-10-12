name = '{{.Name}}'
description = '{{.Description}}'
author = '{{.Author}}'
version = '{{.Version}}'

baseUrl = '{{.BaseUrl}}'

headers = {
    'Content-Type' = 'application/json'
}

// Apis format: [method] [path] [description]

// GET /api/users 获取用户列表
headers['Authorization'] = 'slfjaslkfjlasdjfjas=='

// POST /api/users 新建用户
post = {
    // 用户名
    'name' = 'test'
    // 密码
    'age' = 10
}