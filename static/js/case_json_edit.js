
check_sel = '            <div class="layui-inline">\n' +
            '                <select name="service_name" id="json_head">\n' +
            '                </select>\n' +
            '            </div>\n'

pre = '        <div class="json-nb">\n' +
    '            <label class="layui-form-label">检查点01</label>\n' +
    '            <div class="layui-inline data_block">\n' +
    '            <div class="layui-inline">\n' +
    '                <select name="service_name" id="json_head">\n' +
    '                    <option selected>新增第一层</option>\n' +
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

    /**
     * 解析json树,并返回相关子层级的所有节点node
     *
     * 入参：json对象（不要传字符串）、path数组对象
     */
    function analysisJson(json, path) {
        if (path !== "undefined" && path != null) {
            for (let i = 0; i < path.length; i++) {
                json = json[path[i]];
            }
        }
        // 获取相关层级的所有子节点
        return Object.keys(json);
    }
    /** 测试数据 **/
        // var obj = {a:1, 'b':'foo', c:[false,'false',null, 'null', {d:{e:1.3e5,f:'1.3e5'}}]};
    var obj = {a:1, 'b':'foo', c:{d:{e:1.3e5,f:'1.3e5'}}};
    // analysisJson(obj,["c","d","e"])
    /**
     * 给第一个下拉框填上默认的第一组json数据
     */

    $.each(analysisJson(obj), function(i,v){
        $("#json_head").append('<option value='+v+'>'+ v+'</option>')
    });

    $(document).on('click', '#right_add', function () {
        $(this).parent().parent().find(".data_block").append(check_sel)
        var arr =new Array();
        const sel_val = $(this).parent().parent().find(".data_block").find("option:selected")
        sel_val.each(function() {
            arr.push($(this).text())
        });
        alert(arr)
        var last_json = $(this).parent().parent().find(".data_block").children().last().find("#json_head")
        $.each(analysisJson(obj, arr), function(i,v){
            last_json.append('<option value='+v+'>'+ v+'</option>')
        });
        form.render()
    });

    $(document).on('click', '#left_rm', function () {
        $(this).parent().parent().find(".data_block").children().last().remove()
        form.render()
    });

    $(document).on('click', '#down_add', function () {
        // alert("down_add")
        const these = $(this);
        these.parent().parent().after(pre)
        // alert($(this).parent().parent().next().find("#json_head").text())
        $.each(analysisJson(obj), function(i,v){
            these.parent().parent().next().find("#json_head").append('<option value='+v+'>'+ v+'</option>')
        });
        form.render()
        form.render("select")
    });

    $(document).on('click', '#up_rm', function () {
        // alert("down_add")
        $(this).parent().parent().remove()
        form.render()
    });

    /** 返回数据映射填充 **/

    var js_list =  [
        {"node":"data.code","checkType":"eq", "value":"1", "valueType":"number"},
        {"node":"data.name","checkType":"eq", "value":"wahaha", "valueType":"string"},
        {"node":"data.sex","checkType":"in", "value":"2,3,4,5", "valueType":"number"}
    ]

    function make_response_to_select(json){
        $.each(json, function (i,v){
            $(".json-height").append(
                '        <div class="json-nb">\n' +
                '            <label class="layui-form-label">检查点01</label>\n' +
                '            <div class="layui-inline data_block">\n'
            )
            const arr = v.node.split(".");
            const temp = new Array();
            $.each(arr, function (i,v){
                $(".json-height").append(
                    '            <div class="layui-inline">\n' +
                    '                <select name="check_name" id="json_head">\n' +
                    '                   <option selected value='+v+'>'+v+'</option>\n')
                // 获取到正常的数据？？？
                if(i===0){
                    $.each(analysisJson(obj), function(i,v){
                        $(".json-height").append('<option value='+v+'>'+ v+'</option>\n')
                    });
                }else{
                    $.each(analysisJson(obj, temp), function(i,v){
                        $(".json-height").append('<option value='+v+'>'+ v+'</option>\n')
                    });
                }
                temp.push(v)
                $(".json-height").append(
                    '                </select>\n' +
                    '            </div>\n'
                )
            });
            $(".json-height").append(
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
                '                    <option value='+v.checkType+'>'+v.checkType+'</option>\n' +
                '                    <option >eq</option>\n' +
                '                    <option >in</option>\n' +
                '                    <option >need</option>\n' +
                '                </select>\n' +
                '            </div>\n' +
                '\n' +
                '            <div class="layui-inline">\n' +
                '                <input type="text" value='+v.value+' name="author" lay-verify="title" autocomplete="off"  class="layui-input">\n' +
                '            </div>\n' +
                '\n' +
                '            <div class="layui-inline">\n' +
                '                <select name="service_name">\n' +
                '                    <option selected>'+v.valueType+'</option>\n' +
                '                    <option >number</option>\n' +
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
            )
        });
    }

    make_response_to_select(js_list)

});


