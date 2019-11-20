layui.define([], function(exports) {
    var MOD_NAME = 'myUtil';
    var myUtil = {
        /**
         * 替换querystring中的参数
         * @param uri
         * @param key
         * @param value
         * @returns {*}
         */
        updateQueryStringParameter: function(uri, key, value) {
            var reg = new RegExp("([?&])" + key.replace(']', '\\]').replace('[', '\\[') + "=.*?(&|$)", "i");
            var separator = uri.indexOf('?') !== -1 ? "&" : "?";
            if (uri.match(reg)) {
                return uri.replace(reg, '$1' + key + "=" + value + '$2');
            }else {
                return uri + separator + key + "=" + value;
            }
        },
        /**
         * 根据当前url生成分页url
         * @param targetPage 目标页码
         * @param limit 每页大小
         * @returns {string}
         */
        createPageUrl: function (targetPage, limit) {
            var url = window.location.href;
            url = this.updateQueryStringParameter(url, 'page', targetPage);
            url = this.updateQueryStringParameter(url, 'limit', limit);
            return url;
        }
    };
    exports(MOD_NAME, myUtil);
});