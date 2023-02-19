<h1>go工具集合</>

## httputil
httputil模块实现了http接口链式调用以及自动解析返回数据
eg: res := httputil.New("http://127.0.0.1", http.MethodGet, nil).Do().ParseResponseBody(data)