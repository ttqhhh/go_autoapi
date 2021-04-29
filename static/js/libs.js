/**
 * 解析json树,并返回相关子层级的所有节点node
 *
 * 入参：json对象（不要传字符串）、path数组对象
 */
function analysisJson(json, path) {
    if (path != "undefined" && path != null) {
        for (let i = 0; i < path.length; i++) {
            json = json[path[i]];
        }
    }
    // 获取相关层级的所有子节点
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
                result["code"] = -1;
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