layui.define([], function(exports) {
    var MOD_NAME = 'myUtil';
    var myUtil = {
        /**
         * 根据当前url生成分页url
         * @param targetPage 目标页码
         * @param limit
         * @returns {string}
         */
        createPageUrl: function (targetPage, limit) {
            var url = window.location.href;
            if(window.location.search === "") {
                url = window.location.pathname+'?page='+targetPage+'&limit='+limit;
            }else{
                if(url.indexOf('page=') !== -1) {
                    url = url.replace(/page=\d+/, 'page='+targetPage);
                }else {
                    url += '&page='+targetPage;
                }
                if(url.indexOf('limit=') !== -1) {
                    url = url.replace(/limit=\d+/, 'limit='+limit);
                }else {
                    url += '&limit='+limit;
                }
            }
            return url;
        }
    };
    exports(MOD_NAME, myUtil);
});