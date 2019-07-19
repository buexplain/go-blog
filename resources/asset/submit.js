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
        //新增 window window.parent 兼容 支持debug 调试
        if(typeof callback != 'string') {
            throw 'submit._callUserFunc param error: callback is must string';
        }
        var key = callback.split('.');
        var that = window;
        var func = null;
        var obj = window;
        var l = key.length;
        for(var i=0; i<l; i++) {
            if(!(key[i] in obj)) {
                break;
            }
            obj = obj[key[i]];
            if(i === (l - 2)) {
                that = obj;
            }
            if(i === (l - 1)) {
                func = obj;
            }
        }
        var result = undefined;
        if(func !== null && typeof func == 'function') {
            if(func.length === 0) {
                result = func.apply(that);
            }else{
                var args = (arguments.length === 1 ? [arguments[0]] : Array.apply(null, arguments));
                if(args.length > 1) {
                    result = func.apply(that, args.slice(1, args.length));
                }else {
                    result = func.apply(that);
                }
            }
        }else {
            throw 'submit._callUserFunc not found function: window.'+callback;
        }
        return result;
    },

    /**
     * 以表单方式进行提交
     * @param DOMObject
     */
    form: function(DOMObject) {
        var that = this;

        var data = that._getAttrData(DOMObject);

        var url = data['url'];
        delete data['url'];

        var _method = data['_method'];
        delete data['_method'];
        if(_method === undefined || _method === "") {
            _method = 'get';
        }

        var target = data['target'];
        delete data['target'];

        var form = document.createElement('form');

        form.id = 'form-' + (new Date()).getTime();

        form.action = url;

        if(_method === 'put' || _method === 'patch') {
            form.method = 'post';
        }else if(_method === 'delete') {
            form.method = 'get';
        }else {
            form.method = _method;
        }
        if(_method !== 'get' && _method !== 'post') {
            data['_method'] = _method;
        }

        if(target) {
            form.target = target;
        }

        for(var i in data) {
            var input = document.createElement('input');
            input.setAttribute('name', i);
            input.setAttribute('type', 'hidden');
            input.setAttribute('value', data[i]);
            form.appendChild(input);
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

        var url = data['url'];
        delete data['url'];

        var _method = data['_method'];
        delete data['_method'];
        if(_method === undefined || _method === "") {
            _method = 'get';
        }

        var method = _method;
        if(_method !== 'get' && _method !== 'post') {
            data['_method'] = _method;
        }
        if(_method === 'put' || _method === 'patch') {
            method = 'post';
        }else if(_method === 'delete') {
            method = 'get';
        }

        var content_type = data['content_type']; //注意 dom 树的属性不支持大小写
        delete data['content_type'];
        if(content_type === undefined || content_type === "") {
            content_type = 'application/x-www-form-urlencoded';
        }

        var success = data['success'];
        delete data['success'];

        var error = data['error'];
        delete data['error'];

        layui.use(['jquery'], function() {
            var $ = layui.jquery;
            $.ajax({
                type: method,
                url: url,
                data: data,
                contentType:content_type,
                success: function (data) {
                    if(success !== undefined && success.length > 0) {
                        that._callUserFunc(success, DOMObject, data);
                    }else {
                        console.log(data);
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
        });
    },

    /**
     * 提交前寻问
     * @param DOMObject
     * @param type string form ajax
     * @param tips string 提示文字
     */
    confirm: function (DOMObject, type, tips) {
        var that = this;

        if(!tips) {
            tips = '此操作不可撤销，你确定执行吗？';
        }

        layui.use(['layer'], function() {
            var layer = layui.layer;
            layer.confirm(tips, {icon: 3, title:'提示'}, function(index) {
                    layer.close(index);
                    if(type === 'form') {
                        that.form(DOMObject);
                    }else if(type === 'ajax') {
                        that.ajax(DOMObject);
                    }else {
                        throw 'submit.confirm param error';
                    }
                },function (index) {
                    layer.close(index);
                }
            );
        });
    }
};