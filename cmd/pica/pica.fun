// It's a full demo for pica.

// Init vars
baseUrl = 'http://localhost:9000'
headers = {
    'Content-Type' = 'application/json'
}

// Apis format: [method] [path] [description]

// Get /api/users 获取用户列表
headers = {
    'Authorization': 'slfjaslkfjlasdjfjas=='
}
must(json.a == 2)

// Post /api/users 新建用户
post = {
    'a' = 'b'
}
must(json.a == 2)