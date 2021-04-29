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


/**
 * 此处入参为一个校验点数组对象
 * 支持的验证符号有：eq、need、in、lt、gt、lte、gte
 *
    [
        {"node":"data.code","checkType":"eq", "value":"1", "valueType":"number"},
        {"node":"data.name","checkType":"eq", "value":"wahaha", "valueType":"string"},
        {"node":"data.sex","checkType":"in", "value":"2,3,4,5", "valueType":"number"}
    ]
 *
 */
function generateJsonPath(checkpoints) {
    var result = {};
    var data = {};
    // 对checkpoint进行循环遍历，生成相应的jsonpath
    for (let i = 0; i < checkpoints.length; i++) {
        var checkpoint = checkpoints[i]
        var node = checkpoint.node
        var checkType = checkpoint.checkType
        var value = checkpoint.value
        var valueType = checkpoint.valueType

        if (valueType == "string"){
            // nothing to do
        } else if (valueType == "number"){
            if (isNaN(value)) {
                var msg = " "+node + " 节点检验值不是数字类型"
                result["code"] = -1
                result["msg"] = msg;
                result["data"] = null;
                return result;
            }
            value = parseFloat(value)
        }
        var checkMap = {};
        checkMap[checkType] = value
        data["$."+node] = checkMap
    }
    result["code"] = 1;
    result["msg"] = null;
    result["data"] = data;
    return result;
}

// json解析函数的验证函数
function analysisJsonTest() {
    var param = {a:1, b:'foo', c:[false,'false',null, 'null', {d:{e:1.3e5,f:'1.3e5'}}]}
    var param = {
        "a": 1,
        "b": "foo",
        "c": [{
            //         "d": {
            //             "e": 1.3e5,
            //             "f": "1.3e5"
            //         },
            //         "x":{
            //             "e": 1.3e5,
            //             "f": "1.3e5"
            //         }
            //     }, {
            //         "g": {
            //             "h": 1.3e5,
            //             "i": "1.3e5"
            //         }
        }]
    }
    // alert(analysisJson(param, ["aaa"]))
    // alert(analysisJson(param, ["aaa", "ccc"]))
    // analysisJson(param, ["c"])
    analysisJson(param, ["c", "d"])
}