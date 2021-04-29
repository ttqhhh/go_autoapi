
check_sel = '            <div class="layui-inline">\n' +
            '                <select name="service_name">\n' +
            '                    <option selected>第几层</option>\n' +
            '                </select>\n' +
            '            </div>\n'

pre = '        <div class="json-nb">\n' +
    '            <label class="layui-form-label">检查点01</label>\n' +
    '            <div class="layui-inline data_block">\n' +
        '            <div class="layui-inline">\n' +
        '                <select name="service_name">\n' +
        '                    <option selected>第一层</option>\n' +
        '                </select>\n' +
        '            </div>\n' +
    '            </div>\n' +
    '\n' +
    '            <div class="layui-inline">\n' +
    '                <button class="layui-btn layui-btn-normal" type="button" id="right_add">＋</button>\n' +
    '            </div>\n' +
    '\n' +
    '            <div class="layui-inline">\n' +
    '                <button class="layui-btn layui-btn-normal" type="button" id="left_rm">－</button>\n' +
    '            </div>\n' +
    '\n' +
    '            <div class="layui-inline">\n' +
    '                <select name="service_name">\n' +
    '                    <option selected>eq</option>\n' +
    '                    <option >in</option>\n' +
    '                    <option >need</option>\n' +
    '                </select>\n' +
    '            </div>\n' +
    '\n' +
    '            <div class="layui-inline">\n' +
    '                <input type="text" name="author" lay-verify="title" autocomplete="off" placeholder="输入验证值" class="layui-input">\n' +
    '            </div>\n' +
    '\n' +
    '            <div class="layui-inline">\n' +
    '                <select name="service_name">\n' +
    '                    <option selected>number</option>\n' +
    '                    <option >string</option>\n' +
    '                </select>\n' +
    '            </div>\n' +
    '            <div class="layui-inline">\n' +
    '                <button class="layui-btn layui-btn-danger" type="button" id="down_add">↓</button>\n' +
    '            </div>\n' +
    '            <div class="layui-inline">\n' +
    '                <button class="layui-btn layui-btn-danger" type="button" id="up_rm">↑</button>\n' +
    '            </div>\n' +
    '        </div>'


layui.use(['form', 'layedit', 'laydate'], function() {
    var form = layui.form
        , layer = layui.layer
        , $ = layui.jquery
        , layedit = layui.layedit
        , laydate = layui.laydate;

    $(document).on('click', '#right_add', function () {
        $(this).parent().parent().find(".data_block").append(check_sel)
        form.render()
    });

    $(document).on('click', '#left_rm', function () {
        $(this).parent().parent().find(".data_block").children().last().remove()
        form.render()
    });

    $(document).on('click', '#down_add', function () {
        // alert("down_add")
        $(".json-height").append(pre)
        form.render()
    });

    $(document).on('click', '#up_rm', function () {
        // alert("down_add")
        $(this).parent().parent().remove()
        form.render()
    });
});