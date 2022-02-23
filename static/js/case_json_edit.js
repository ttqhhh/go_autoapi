
check_sel = '            <div class="layui-inline" style="margin: 0;width: 100px">\n' +
            '                <select name="check_name" id="json_head" >\n' +
            '                </select>\n' +
            '            </div>\n'

pre = '        <div class="json-nb">\n' +
    '            <label class="layui-form-label">检查点</label>\n' +
    '            <div class="layui-inline data_block">\n' +
    '            <div class="layui-inline" style="margin: 0;width: 100px">\n' +
    '                <select name="check_name" id="json_head" >\n' +
    '                </select>\n' +
    '            </div>\n' +
    '            </div>\n' +
    '\n' +
    '            <div class="layui-inline" style="margin:0">\n' +
    '                <button class="layui-btn layui-btn-normal layui-btn-sm" type="button" id="right_add" style="margin: 0">＋</button class="layui-btn layui-btn-normal">\n' +
    '                <button class="layui-btn layui-btn-normal layui-btn-sm" type="button" id="left_rm" style="margin: 0">－</button>\n' +
    '            </div>\n' +
    '\n' +
    '            <div class="layui-inline" style="margin: 0;width: 100px">\n' +
    '                <select name="check_type" id="check_type">\n' +
    '                    <option value="eq" selected>等于</option>\n' +
    '                    <option value="neq">不等于</option>\n' +
    '                    <option value="in">包含于</option>\n' +
    '                    <option value="exist">存在此字段</option>\n' +
    '                    <option value="lt">小于</option>\n' +
    '                    <option value="gt">大于</option>\n' +
    '                    <option value="lte">小于等于</option>\n' +
    '                    <option value="gte">大于等于</option>\n' +
    '                    <option value="isTrue">为真</option>\n' +
    '                    <option value="isFalse">为假</option>\n' +
    '                </select>\n' +
    '            </div>\n' +
    '\n' +
    '            <div class="layui-inline" style="margin: 0">\n' +
    '                <input type="text" name="value" id="value" lay-verify="title" autocomplete="off" placeholder="输入验证值" class="layui-input">\n' +
    '            </div>\n' +
    '\n' +
    '            <div class="layui-inline" style="margin: 0;width: 150px">\n' +
    '                <select name="value_type" id="value_type">\n' +
    '                    <option value="请选择数据类型">请选择数据类型</option>\n' +
    '                    <option value="number">number</option>\n' +
    '                    <option value="string">string</option>\n' +
    '                </select>\n' +
    '            </div>\n' +
    '            <div class="layui-inline" style="margin:0">\n' +
    '                <button class="layui-btn layui-btn-danger layui-btn-sm" type="button" id="down_add" style="margin: 0">↓</button>\n' +
    '                <button class="layui-btn layui-btn-danger layui-btn-sm" type="button" id="up_rm" style="margin:0">↑</button>\n' +
    '            </div>\n' +
    '        </div>'


layui.use(['form', 'layedit', 'laydate'], function() {
    var form = layui.form
        , layer = layui.layer
        , $ = layui.jquery
        , layedit = layui.layedit
        , laydate = layui.laydate;


});


