/**
 * 提交对象
 */
var submit = {
    /**
     * 获取 DOMObject 上以 data- 开头的属性
     * @param DOMObject
     * @private
     */
    _getAttrData: function (DOMObject) {
        var attrs = DOMObject.attributes;
        var data = {};
        for(var i in attrs) {
            if(attrs[i].nodeName === undefined || attrs[i].nodeName.indexOf('data-') !== 0) {
                continue;
            }
            data[attrs[i].nodeName.substr(5)] = attrs[i].nodeValue;
        }
        return data;
    },

    /**
     *  从window对象中找到 callback 字符串描述的函数，然后进行调用
     *  示例：_callUserFunc('location.reload') 等价于 window.location.reload()
     * @param callback string
     * @returns {undefined}
     * @private
     */
    _callUserFunc: function(callback) {
        if(typeof callback != 'string') {
            throw 'submit._callUserFunc param error: callback is must string';
        }
        var scanCtx = function(rootObj, funcChain) {
            var key = funcChain.split('.');
            var lastParent = rootObj;
            var func = null;
            var obj = rootObj;
            var l = key.length;
            for(var i=0; i<l; i++) {
                if(i === 0 && key[0] === 'window') {
                    continue;
                }
                if(i === 1 && key[0] === 'window' && key[1] === 'parent') {
                    continue;
                }
                if(!(key[i] in obj)) {
                    break;
                }
                obj = obj[key[i]];
                if(i === (l - 2)) {
                    lastParent = obj;
                }
                if(i === (l - 1)) {
                    func = obj;
                }
            }
            return [lastParent, func];
        };
        var ctx = scanCtx(window, callback);
        if(ctx[1] === null) {
            ctx = scanCtx(window.parent, callback);
        }
        if(ctx[1] === null || typeof ctx[1] !== 'function') {
            throw 'submit._callUserFunc not found function: '+callback;
        }
        var args = Array.apply(null, arguments);
        args.shift();
        return ctx[1].apply(ctx[0], args);
    },

    /**
     * 弹出信息
     * @private
     */
    _alert: function() {
        var args = arguments;
        if(args.length === 0) {
            return;
        }
        try {
            if(window.layui !== undefined) {
                layui.use(['layer'], function() {
                    var layer = layui.layer;
                    layer.alert(...args);
                });
            }else if(window.layer !== undefined) {
                layer.alert(...args);
            }else {
                alert(args[0]);
            }
        }catch (e) {
            console.log(e);
        }
    },

    /**
     * 以表单方式进行提交
     * @param DOMObject
     */
    form: function(DOMObject) {
        var that = this;

        var data = that._getAttrData(DOMObject);

        var debug = data['debug'];
        delete data['debug'];
        if(debug === undefined || debug === "" || debug === "false") {
            debug = false;
        }else {
            debug = true;
        }

        var url = data['url'];
        delete data['url'];
        if(url === undefined || url === "") {
            throw 'submit.form: invalid url';
        }

        var method = data['method'];
        delete data['method'];
        if(method === undefined || method === "") {
            method = 'get';
        }

        var target = data['target'];
        delete data['target'];

        var form = document.createElement('form');
        form.id = 'form-' + (new Date()).getTime();
        form.action = url;

        if(method === 'put' || method === 'patch') {
            //改为post请求
            form.method = 'post';
        }else if(method === 'delete') {
            //改为get请求
            form.method = 'get';
        }else {
            form.method = method;
        }
        //模拟请求
        if(method !== 'get' && method !== 'post') {
            data['_method'] = method;
        }

        if(target) {
            form.target = target;
        }

        //请求之前预处理数据回调
        var prepare = data['call-prepare-data'];
        delete data['call-prepare-data'];
        if(prepare !== undefined && prepare.length > 0) {
            try {
                data = that._callUserFunc(prepare, data);
            }catch (e) {
                that._alert(e.toString(), {icon:0});
                return;
            }
        }

        for(var i in data) {
            var input = document.createElement('input');
            input.setAttribute('name', i);
            input.setAttribute('type', 'hidden');
            input.setAttribute('value', data[i]);
            form.appendChild(input);
        }

        if(debug) {
            console.log(debug);
            return;
        }

        var body = document.getElementsByTagName('body')[0];

        body.append(form);
        form.submit();

        try {
            body.removeChild(form);
        }catch (e) {
            console.log(e);
        }
    },

    /**
     * 以ajax方式进行提交
     * @param DOMObject
     */
    ajax: function(DOMObject) {
        var that = this;

        var data = that._getAttrData(DOMObject);

        var debug = data['debug'];
        delete data['debug'];
        if(debug === undefined || debug === "" || debug === "false") {
            debug = false;
        }else {
            debug = true;
        }

        var url = data['url'];
        delete data['url'];
        if(url === undefined || url === "") {
            throw 'submit.ajax: invalid url';
        }

        var method = data['method'];
        delete data['method'];
        if(method === undefined || method === "") {
            method = 'GET';
        }

        var content_type = data['content_type']; //注意 dom 树的属性不支持大小写
        delete data['content_type'];
        if(content_type === undefined || content_type === "") {
            content_type = 'application/x-www-form-urlencoded';
        }

        //成功回调
        var success = data['call-success'];
        delete data['call-success'];

        //错误回调
        var error = data['call-error'];
        delete data['call-error'];

        //请求之前预处理数据回调
        var prepare = data['call-prepare-data'];
        delete data['call-prepare-data'];
        if(prepare !== undefined && prepare.length > 0) {
            try {
                data = that._callUserFunc(prepare, data);
            }catch (e) {
                that._alert(e.toString(), {icon:0});
                return;
            }
        }

        if(debug) {
            console.log(data);
            return;
        }

        var request = function(jQuery) {
            jQuery.ajax({
                type: method,
                url: url,
                data: data,
                contentType:content_type,
                success: function (data) {
                    if(success !== undefined && success.length > 0) {
                        that._callUserFunc(success, DOMObject, data);
                    }else {
                        //默认的回调逻辑
                        if(typeof data === 'object' && data.code !== undefined && (data.msg !== undefined || data.message !== undefined)) {
                            //假设返回结构是一个 {code:1,message:"xx"} 或 {code:1,msg:"xx"} 的对象
                            var message = '';
                            if(data.msg !== undefined) {
                                message = data.msg;
                            }else {
                                message = data.message;
                            }
                            if(message !== '') {
                                if(data.code === 0) {
                                    //操作成功
                                    that._alert(message, {icon: 1});
                                }else {
                                    //操作失败
                                    that._alert(message, {icon: 2});
                                }
                            }
                        }else if(typeof data === 'string') {
                            //返回的是字符串，不为空字符串则弹出
                            if(data !== '') {
                                that._alert(data);
                            }
                        }else {
                            //未知返回，直接弹出
                            that._alert(data);
                        }
                    }
                },
                error: function (jqXHR) {
                    if(error !== undefined) {
                        that._callUserFunc(error, DOMObject, jqXHR);
                    }else {
                        console.log(jqXHR);
                    }
                }
            });
        };

        if(window.layui !== undefined) {
            layui.use(['jquery'], function() {
                var $ = layui.jquery;
                request($);
            });
        }else {
            request(window.$);
        }
    },

    /**
     * 提交前寻问
     * @param DOMObject
     * @param type string form ajax
     * @param tips string 提示文字
     */
    confirm: function (DOMObject, type, tips) {
        var that = this;

        if(!type) {
            type = 'ajax';
        }

        if(!tips) {
            tips = '此操作不可撤销，你确定执行吗？';
        }

        var request = function () {
            if(type === 'form') {
                that.form(DOMObject);
            }else if(type === 'ajax') {
                that.ajax(DOMObject);
            }else {
                throw 'submit.confirm param error';
            }
        };

        var layerConfirm = function (layer) {
            layer.confirm(tips, {icon: 3, title:'提示'}, function(index) {
                    layer.close(index);
                    request();
                },function (index) {
                    layer.close(index);
                }
            );
        };

        if(window.layui !== undefined) {
            layui.use(['layer'], function() {
                var layer = layui.layer;
                layerConfirm(layer);
            });
        }else if(window.layer !== undefined) {
            layerConfirm(window.layer);
        }else {
            if(window.confirm(tips)) {
                request();
            }
        }
    }
};