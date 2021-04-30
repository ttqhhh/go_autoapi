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

/** 此处入参为一个校验点数组对象、一个能匹配上校验点的json对象（冒烟响应json）
 * 支持的验证符号有：eq、need、in、lt、gt、lte、gte
 * @param checkpoints
 * [
        {"node":"data.code","checkType":"eq", "value":1, "valueType":"number"},
        {"node":"data.name","checkType":"eq", "value":"wahaha", "valueType":"string"},
        {"node":"data.sex","checkType":"in", "value":"2,3,4,5", "valueType":"number"}
    ]
 * @param json
 * var json = {
        "data":{
            "code":3,
            "name": "验证generate",
            "sex": "男的",
            "relations":[
                {
                    "father":"张三",
                    "mother":"小红"
                },
                {
                    "father":"李四",
                    "mother":"小美"
                },
            ]
        }
    }
 * @returns {{}}
 * {
    "$.data.code": {
        "eq": 1
    },
    "$.data.name": {
        "eq": "wahahah"
    },
    "$.data.code": {
        "in": "2,3,4,5"
    }
}
 */
function generateJsonPath(checkpoints, json) {
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

    for (let i = 0; i < checkpoints.length; i++) {
        // 根据node值，判断每层中是否有数组情况
        var finalNode = "$";
        var checkpoint = checkpoints[i]
        var node = checkpoint.node
        var innerJson = json;
        var nodes = node.split(".");
        for (let j = 0; j < nodes.length; j++) {
            var innerNode = nodes[j];
            if (Array.isArray(innerJson)) {
                innerJson = innerJson[0][innerNode];
                finalNode += "[0]" + "." + innerNode;
            } else {
                innerJson = innerJson[innerNode];
                finalNode += "." + innerNode;
            }
        }
        var key = "$." + node;
        var value = data[key];
        delete data[key];
        data[finalNode] = value;
    }

    // var innerJson;
    // var nodes = node.split(".")
    // for (let j = 0; j < nodes.length; j++) {
    //     var innerNode = nodes[i];
    //     innerJson = json[innerNode]
    //     if (Array.isArray(innerJson)) {
    //
    //     }
    // }
    //
    result["code"] = 1;
    result["msg"] = null;
    result["data"] = data;
    return result;
}

/***
 * 解析jsonpath，入参为一个json对象
 * @param checkpoints
 * {
            "$.data.code": {
                "eq": 1
            },
            "$.data.name": {
                "eq": "wahaha"
            },
            "$.data.sex": {
                "in": "2,3,4,5"
            },
            "$.data.relations[0].mother[0].inblood": {
                "in": "2,3,4,5"
            }
    }
 *
 * @returns {{}}
 * [
        {"node":"data.code","checkType":"eq", "value":1, "valueType":"number"},
        {"node":"data.name","checkType":"eq", "value":"wahaha", "valueType":"string"},
        {"node":"data.sex","checkType":"in", "value":"2,3,4,5", "valueType":"number"}
    ]
 */
function analysisJsonPath(jsonpath) {
    // jsonpath = jsonpath["data"]; // 用于测试
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

/*******************************************************************单测函数*********************************************************************/

function generateJsonPathTest() {
    var param = [
        {"node":"data.code","checkType":"eq", "value":1, "valueType":"number"},
        {"node":"data.name","checkType":"eq", "value":"wahaha", "valueType":"string"},
        {"node":"data.sex","checkType":"in", "value":"2,3,4,5", "valueType":"string"},
        {"node":"data.relations.mother","checkType":"in", "value":"2,3,4,5", "valueType":"string"}
    ]
    var json = {
        "data":{
            "code":3,
            "name": "验证generate",
            "sex": "男的",
            "relations":[
                {
                    "father":"张三",
                    "mother":"小红"
                },
                {
                    "father":"李四",
                    "mother":"小美"
                },
            ]
        }
    }
    var result = generateJsonPath(param, json);
    console.log("generateJsonPath处理结果："+JSON.stringify(result))
    return result;
}

function analysisJsonPathTest() {
    // var param = generateJsonPathTest()
    var param = {
            "$.data.code": {
                "eq": 1
            },
            "$.data.name": {
                "eq": "wahaha"
            },
            "$.data.sex": {
                "in": "2,3,4,5"
            },
            "$.data.relations[0].mother[0].inblood": {
                "in": "2,3,4,5"
            }
    }
    console.log("analysisJsonPath处理结果："+JSON.stringify(analysisJsonPath(param)));
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