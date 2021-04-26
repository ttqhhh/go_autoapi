<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>首页</title>
    <meta name="renderer" content="webkit">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
    <link rel="stylesheet" href="//unpkg.com/layui@2.6.5/dist/css/layui.css">
    <!-- 注意：如果你直接复制所有代码到本地，上述css路径需要改成你本地的 -->
</head>
<body class="layui-layout-body" style="border-bottom: solid 5px #1aa094;">
<div class="layui-layout layui-layout-admin" style="height: 1100px; width: 2000px">
    <div class="layui-header header header-demo">
        <div class="layui-main">
            <div class="admin-login-box">
                <!-- 展示最右app图标 -->
                <a class="logo" style="left: 0;" href="https://h5.izuiyou.com/favicon.ico">
                    <span style="font-size: 22px; color: yellow">自动化测试平台</span>
                </a>
                <div class="admin-side-toggle">
                    <i class="fa fa-bars" aria-hidden="true"></i>
                </div>
                <div class="admin-side-full">
                    <i class="fa fa-life-bouy" aria-hidden="true"></i>
                </div>
            </div>
{{/*            <ul class="layui-nav admin-header-item">*/}}
{{/*                <li class="layui-nav-item">*/}}
{{/*                    <a href="javascript:;">清除缓存</a>*/}}
{{/*                </li>*/}}
{{/*                <li class="layui-nav-item">*/}}
{{/*                    <a href="javascript:;">浏览网站</a>*/}}
{{/*                </li>*/}}
{{/*                <li class="layui-nav-item" id="video1">*/}}
{{/*                    <a href="javascript:;">视频</a>*/}}
{{/*                </li>*/}}
{{/*                <li class="layui-nav-item">*/}}
{{/*                    <a href="javascript:;" class="admin-header-user">*/}}
{{/*                        <img src="img/0.jpg" />*/}}
{{/*                        <span>layui</span>*/}}
{{/*                    </a>*/}}
{{/*                    <dl class="layui-nav-child">*/}}
{{/*                        <dd>*/}}
{{/*                            <a href="javascript:;"><i class="fa fa-user-circle" aria-hidden="true"></i> 个人信息</a>*/}}
{{/*                        </dd>*/}}
{{/*                        <dd>*/}}
{{/*                            <a href="javascript:;"><i class="fa fa-gear" aria-hidden="true"></i> 设置</a>*/}}
{{/*                        </dd>*/}}
{{/*                        <dd id="lock">*/}}
{{/*                            <a href="javascript:;">*/}}
{{/*                                <i class="fa fa-lock" aria-hidden="true" style="padding-right: 3px;padding-left: 1px;"></i> 锁屏 (Alt+L)*/}}
{{/*                            </a>*/}}
{{/*                        </dd>*/}}
{{/*                        <dd>*/}}
{{/*                            <a href="login.html"><i class="fa fa-sign-out" aria-hidden="true"></i> 注销</a>*/}}
{{/*                        </dd>*/}}
{{/*                    </dl>*/}}
{{/*                </li>*/}}
{{/*            </ul>*/}}
{{/*            <ul class="layui-nav admin-header-item-mobile">*/}}
{{/*                <li class="layui-nav-item">*/}}
{{/*                    <a href="login.html"><i class="fa fa-sign-out" aria-hidden="true"></i> 注销</a>*/}}
{{/*                </li>*/}}
{{/*            </ul>*/}}
{{/*        </div>*/}}
    </div>

    <div class="layui-side layui-bg-black" style="background-color: #1E9FFF">
        <div class="layui-side-scroll" style="background-color: #FFF">
            <!-- 左侧导航区域（可配合layui已有的垂直导航） -->
            <ul class="layui-nav layui-nav-tree left_menu_li"  lay-filter="test">
                <li class="layui-nav-item first-item layui-this">
                    <a href="javascript:;"><i class="layui-icon" style="margin-right:10px;">&#xe68e;</i>首页</a>
                </li>
                <li class="layui-nav-item">
                    <a class="" href="javascript:;">台账管理</a>
                    <dl class="layui-nav-child">
                        <dd><a target="myFrameName" href="table.html">计量器具台账</a></dd>
                        <dd><a target="myFrameName" href="javascript:;">计量标准台账</a></dd>
                        <dd><a target="myFrameName" href="javascript:;">机构资质台账</a></dd>
                        <dd><a target="myFrameName" href="javascript:;">批量修改台账</a></dd>
                        <dd><a target="myFrameName" href="form.html">任务凭证</a></dd>
                    </dl>
                </li>
                <li class="layui-nav-item">
                    <a href="javascript:;">溯源计划</a>
                    <dl class="layui-nav-child">
                        <dd><a href="javascript:;">编制年度计划</a></dd>
                        <dd><a href="javascript:;">编制月度计划</a></dd>
                        <dd><a href="">监督抽查</a></dd>
                    </dl>
                </li>
                <li class="layui-nav-item"><a href="">检定管理</a></li>
                <li class="layui-nav-item"><a href="">资质管理</a></li>
                <li class="layui-nav-item"><a href="">到期提醒</a></li>
                <li class="layui-nav-item"><a href="">查询统计</a></li>
                <li class="layui-nav-item"><a href="">系统维护</a></li>
                <li class="layui-nav-item"><a href="">考核兑现</a></li>
            </ul>
        </div>
    </div>

    <div class="layui-body">
        <div class="iframe-mask" id="iframe-mask" style="display: none;"></div>
        <!-- 内容主体区域 -->
{{/*        <iframe class="admin-iframe" name="myFrameName" src="main.html"></iframe>*/}}
        <iframe class="admin-iframe" name="myFrameName" src="/service/index" style="width: 100%; height: 810px"></iframe>
    </div>

    <div class="layui-footer">
        <!-- 底部固定区域 -->
        © layui.com - 底部固定区域
    </div>
</div>
<script src="resource/layui.js"></script>
<script src="js/index.js"></script>
</body>
</html>