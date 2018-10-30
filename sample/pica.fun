// It's a full demo for pica.

// Init vars
name = 'demo'
description = 'This is a demo for pica.'
version = '0.0.1'
author = 'jeremaihloo@gmail.com<jeremaihloo>'
baseUrl = 'http://localhost:8080'

headers = {
  'Content-Type' = 'application/json'
}

// Apis format: [method] [path] [description]

// GET /api/users api_get_users 获取用户列表
headers.Authorization = 'slfjadslkfjlasdjfjas=='
headers['Content-Type'] = 'application/json'
assert(status==0)

// POST /api/users api_create_user 新建用户
post = {
  // 用户名
  name = 'test'
  // 密码
  age = 10
}