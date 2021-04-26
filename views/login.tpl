<!DOCTYPE html>

<html>
<head>
    <title>首页</title>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <link rel="shortcut icon"
          href="https://h5.izuiyou.com/favicon.ico"
          type="image/x-icon"/>
    <link rel="stylesheet" type="text/css" href="../static/layui/css/layui.css" media="all"/>

    <style type="text/css">
    </style>
</head>

<body>
<script src="/static/layui/layui.js"></script>
<script src="https://code.jquery.com/jquery-3.1.1.min.js"></script>
<script>
    function login() {
        var username = $("[name='user_name']").val();
        var password = $("[name='password']").val();
        var data = {user_name: username, password: password}
        $.ajax({
            url: '/auto/login',
            type: 'POST',
            contentType: 'application/json;charset=utf-8',
            data: JSON.stringify(data),
            success: function (resp) {
                code = resp.code;
                if (code == 200) {
                    // 登录成功
                    window.location.href = "/service/index";
                } else {
                    // 登录失败
                    layer.alert(resp.msg);
                }
            }
        })
    }

    function forget() {
        layer.alert("客观别急，等会儿我再实现~~~");
    }
</script>
{{/*<div></div>*/}}
<div id="background" style="position:absolute;z-index:-1;width:100%;height:100%;top:0px;left:0px;">

    <!-- 背景图 -->
    <img style="z-index: 1; position: absolute" src="https://static.ixiaochuan.cn/planck/home_background.81871fe3ab.jpg"
         width="100%" height="100%"/>
    <!-- 登录框 -->
    <div style="z-index: 2; position: absolute; height: 100%; width: 100%">
        {{/*        <div style="text-align: center; vertical-align: center; margin: 0px auto">*/}}
        <div style="width: 100px; height: 100px; position: absolute; left: 35%; top: 40%; margin: -50px 0 0 -50px;">
            <form class="layui-form">
                <div class="layui-form-item" style="width: 500px">
                    <label class="layui-form-label">用户名</label>
                    <div class="layui-input-block" style="width: 350px">
                        <input type="text" name="user_name" required lay-verify="required" placeholder="请输入用户名"
                               autocomplete="off" class="layui-input">
                    </div>
                </div>
                <div class="layui-form-item" style="width: 500px">
                    <label class="layui-form-label">密码</label>
                    <div class="layui-input-block" style="width: 350px">
                        <input type="password" name="password" required lay-verify="required" placeholder="请输入密码"
                               autocomplete="off" class="layui-input">
                        <div class="layui-form-mid layui-word-aux" style="color: red">同Confluence账号和密码</div>
                    </div>
                </div>
                <div class="layui-form-item" style="width: 500px">
                    <div class="layui-input-block">
                        <button type="button" id="login_btn" class="layui-btn" onclick="login()">登&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;录</button>
                        <button type="button" id="forget_btn" class="layui-btn layui-btn-warm" onclick="forget()">
                            忘记密码?
                        </button>
                    </div>
                </div>
            </form>
        </div>
    </div>
</div>
</body>
</html>
