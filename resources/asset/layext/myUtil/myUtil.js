layui.define([], function(exports) {
    var MOD_NAME = 'myUtil';
    var myUtil = {
        queryString: {
            /**
             * 替换querystring中的参数
             */
            update: function(key, value, uri) {
                if(uri === undefined) {
                    uri = window.location.href;
                }
                var reg = new RegExp("([?&])" + key.replace(']', '\\]').replace('[', '\\[') + "=.*?(&|$)", "i");
                var separator = uri.indexOf('?') !== -1 ? "&" : "?";
                if (uri.match(reg)) {
                    return uri.replace(reg, '$1' + key + "=" + value + '$2');
                }else {
                    return uri + separator + key + "=" + value;
                }
            },
            /**
             * 获取uri中的参数
             */
            get: function (key, def, uri) {
                if(def === undefined) {
                    def = null;
                }
                if(uri === undefined) {
                    uri = window.location.href;
                }
                key += '=';
                var index = uri.indexOf(key);
                if(index === -1) {
                    return def;
                }
                uri = uri.substr(index + key.length, uri.length);
                index = uri.indexOf('&');
                if(index === -1) {
                    return uri;
                }
                return uri.substr(0, index);
            }
        },
        /**
         * 根据当前url生成分页url
         * @param targetPage 目标页码
         * @param limit 每页大小
         * @returns {string}
         */
        createPageUrl: function (targetPage, limit) {
            var url = this.queryString.update('page', targetPage);
            url = this.queryString.update('limit', limit);
            return url;
        }
    };
    exports(MOD_NAME, myUtil);
});