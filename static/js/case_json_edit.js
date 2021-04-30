
check_sel = '            <div class="layui-inline">\n' +
            '                <select name="check_name" id="json_head">\n' +
            '                </select>\n' +
            '            </div>\n'

pre = '        <div class="json-nb">\n' +
    '            <label class="layui-form-label">检查点01</label>\n' +
    '            <div class="layui-inline data_block">\n' +
    '            <div class="layui-inline">\n' +
    '                <select name="check_name" id="json_head">\n' +
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
    '                <select name="check_type" id="check_type">\n' +
    '                    <option value="eq" selected>eq</option>\n' +
    '                    <option value="in">in</option>\n' +
    '                    <option value="need">need</option>\n' +
    '                </select>\n' +
    '            </div>\n' +
    '\n' +
    '            <div class="layui-inline">\n' +
    '                <input type="text" name="value" id="value" lay-verify="title" autocomplete="off" placeholder="输入验证值" class="layui-input">\n' +
    '            </div>\n' +
    '\n' +
    '            <div class="layui-inline">\n' +
    '                <select name="value_type" id="value_type">\n' +
    '                    <option value="number">number</option>\n' +
    '                    <option value="string">string</option>\n' +
    '                </select>\n' +
    '            </div>\n' +
    '            <div class="layui-inline">\n' +
    '                <button class="layui-btn layui-btn-danger" type="button" id="down_add">↓</button>\n' +
    '            </div>\n' +
    '            <div class="layui-inline">\n' +
    '                <button class="layui-btn layui-btn-danger" type="button" id="up_rm">↑</button>\n' +
    '            </div>\n' +
    '        </div>'


function analysisJsonPath(jsonpath) {
    // jsonpath = jsonpath["data"]; // 用于测试
    alert(jsonpath)
    var result = new Array();
    var keys = Object.keys(jsonpath)
    for (let i = 0; i < keys.length; i++) {
        var key = keys[i];
        var item = jsonpath[key];
        var innerKeys = Object.keys(item);
        var checkType = innerKeys[0];
        var checkValue = item[checkType];

        // 生成checkpoint
        var checkpoint = {};
        // 对key做截取，去掉'$.'后为json的node节点value
        var node = key.slice(2);
        // 将node中的'[0]'剔除掉
        while (node.indexOf('[0]')!=-1) {
            var index = node.indexOf('[0]')
            var temp = node.substring(0, index) + node.substring(index + 3, node.length);
            node = temp;
        }
        checkpoint["node"] = node;
        checkpoint["checkType"] = checkType;
        checkpoint["value"] = checkValue;
        var valueType = "string"
        if (typeof(checkValue) == "number") {
            valueType = "number";
        }
        checkpoint["valueType"] = valueType;

        // 将解析出来的校验点push到数组中
        result.push(checkpoint);
    }
    return result;
}




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
        if (path != "undefined" && path != null) {
            for (let i = 0; i < path.length; i++) {
                if (json == undefined) {
                    return new Array();
                }
                if (Array.isArray(json)) {
                    if (json.length > 0) {
                        // 当json为一个数组时，取数组第一层
                        json = json[0];
                    } else {
                        // 当json数组中没有数据时，再去取深层节点肯定报错，所以，直接返回
                        return new Array();
                    }
                }
                json = json[path[i]];
            }
        }
        // 获取相关层级的所有子节点
        // 首先验证undefined
        if (json == undefined) {
            return new Array();
        }
        // alert("是否数组？？" + Array.isArray(json) + "\n\n" + json);
        var isArray = Array.isArray(json);
        if (isArray) {
            if (json.length > 0) {
                // 确保数组中有第一个元素
                json = json[0];
            } else {
                // 返回空数组
                return new Array()
            }
        }
        var keys = Object.keys(json);
        return keys;
    }
    /** 测试数据 **/
        // var obj = {a:1, 'b':'foo', c:[false,'false',null, 'null', {d:{e:1.3e5,f:'1.3e5'}}]};
    // var obj = {"a":1, "b":"foo'", "c":{"d":{"e":"1.3e5","f":"1.3e5'"}}};
    // analysisJson(obj,["c","d","e"])
    const obj = get_response()
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
        // alert(arr)
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

    function make_response_to_select(json){
        $.each(json, function (i,v){
            let html_head =
                '        <div class="json-nb">\n' +
                '            <label class="layui-form-label">检查点</label>\n' +
                '            <div class="layui-inline data_block">\n'
            const arr = v.node.split(".");
            const temp = new Array();
            let html_body = ""
            $.each(arr, function (i,v){
                 html_body = html_body +
                     '            <div class="layui-inline">\n' +
                    '                <select name="check_name" id="json_head">\n' +
                    '                   <option selected value='+v+'>'+v+'</option>\n'
                // 获取到正常的数据？？？
                if(i===0){
                    $.each(analysisJson(obj), function(i,v){
                        html_body  = html_body  + '<option value=' + v + '>' + v + '</option>\n'
                    });
                }else{
                    // alert(temp)
                    // alert(analysisJson(obj, temp))
                    $.each(analysisJson(obj, temp), function(i,v){
                        html_body  = html_body  + '<option value='+v+'>'+ v+'</option>\n'
                    });
                }
                temp.push(v)
                html_body = html_body + '</select>\n</div>\n'
            });

            let html_butt =
                '            <div class="layui-inline">\n' +
                '                <button class="layui-btn layui-btn-normal" type="button" id="right_add">＋</button>\n' +
                '            </div>\n' +
                '\n' +
                '            <div class="layui-inline">\n' +
                '                <button class="layui-btn layui-btn-normal" type="button" id="left_rm">－</button>\n' +
                '            </div>\n' +
                '\n' +
                '            <div class="layui-inline">\n' +
                '                <select name="check_type" id="check_type">\n' +
                '                    <option checked value=' + v.checkType + '>' + v.checkType + '</option>\n' +
                '                    <option value="eq">eq</option>\n' +
                '                    <option value="in">in</option>\n' +
                '                    <option value="need">need</option>\n' +
                '                </select>\n' +
                '            </div>\n' +
                '\n' +
                '            <div class="layui-inline">\n' +
                '                <input type="text" id="value" value=' + v.value + ' name="value" lay-verify="title" autocomplete="off"  class="layui-input">\n' +
                '            </div>\n' +
                '\n' +
                '            <div class="layui-inline">\n' +
                '                <select name="value_type" id="value_type">\n' +
                '                    <option checked value='+v.valueType+'>' + v.valueType + '</option>\n' +
                '                    <option value="number">number</option>\n' +
                '                    <option value="string">string</option>\n' +
                '                </select>\n' +
                '            </div>\n' +
                '            <div class="layui-inline">\n' +
                '                <button class="layui-btn layui-btn-danger" type="button" id="down_add">↓</button>\n' +
                '            </div>\n' +
                '            <div class="layui-inline">\n' +
                '                <button class="layui-btn layui-btn-danger" type="button" id="up_rm">↑</button>\n' +
                '            </div>\n' +
                '        </div>'

            $(".json-height").append(html_head + html_body + '</div>\n' +html_butt)
        });
    }

    function get_response(){
        var request_param = $("#request_param").val()
        var request_url = $("#api_url").val()
        let body = ""
        $.ajax({
            type: 'POST',
            contentType: "application/x-www-form-urlencoded",
            dataType: "json",
            url: "/auto/perform_smoke",
            async: false,
            timeout: 500000,
            data: {
                "api_url":request_url,
                "parameter":request_param
            },
            success:function (data){
                // alert(data.data.body)
                body = JSON.parse(data.data.body)
            }
        });
        return body
    }

    const js_list = $("#check_point").val()
    const res = analysisJsonPath(JSON.parse(js_list))
    make_response_to_select(res)

});


