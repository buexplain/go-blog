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
            get: function (key, def, url) {
                if (def === undefined) {
                    def = null;
                }
                key += '=';
                let index = url.indexOf(key);
                if (index === -1) {
                    return def;
                }
                url = url.substr(index + key.length, url.length);
                let index_and = url.indexOf('&');
                let index_jing = url.indexOf('#');
                if (index_and !== -1 && index_jing !== -1) {
                    if (index_jing < index_and) {
                        index = index_jing;
                    } else {
                        index = index_and;
                    }
                } else if (index_and !== -1) {
                    index = index_and;
                } else if (index_jing !== -1) {
                    index = index_jing;
                } else {
                    index = -1;
                }

                if (index === -1) {
                    return url;
                }
                return url.substr(0, index);
            }
        },
        /**
         * 判断是否为图片文件
         * @param file_name
         * @returns {boolean}
         */
        isImage: function(file_name) {
            if(typeof file_name !== 'string') {
                throw "error file name: "+file_name;
            }
            var images = ['jpeg', 'gif', 'jpg', 'png', 'bmp'];
            for(var i in images) {
                if(file_name.substr(file_name.length - images[i].length, images[i].length).toLocaleLowerCase() === images[i]) {
                    return true;
                }
            }
            return false;
        },
        /**
         * 判断文件是否可编辑
         * @param file_name
         * @returns {boolean}
         */
        isEditable: function(file_name) {
            if(typeof file_name !== 'string') {
                throw "error file name: "+file_name;
            }
            var names = ["c", "cpp", "php", "java", "go", "py", "css", "html", "js", "vue", "txt"];
            for(var i in names) {
                if(file_name.substr(file_name.length - names[i].length, names[i].length).toLocaleLowerCase() === names[i]) {
                    return true;
                }
            }
            return false;
        },
        /**
         * 根据当前url生成分页url
         * @param targetPage 目标页码
         * @param limit 每页大小
         * @returns {string}
         */
        createPageUrl: function (targetPage, limit) {
            var url = this.queryString.update('page', targetPage);
            return this.queryString.update('limit', limit, url);
        },
        /**
         * 字节友好显示
         * @param size
         * @returns {string}
         */
        renderBytes: function (size) {
            if(size === undefined || size === null || parseFloat(size) <= 0) {
                return "0 Bytes";
            }
            size = parseFloat(size);
            var unitArr = ["Bytes","KB","MB","GB","TB","PB","EB","ZB","YB"];
            var index = Math.floor(Math.log(size)/Math.log(1024));
            var new_size = (size/Math.pow(1024,index)).toFixed(2);
            return new_size +' '+ unitArr[index];
        },
        /**
         * 对无法进行字符串拼接的字符串变量进行掩码处理
         */
        maskStr: function (str) {
            return {
                str: str,
                toString: function () {
                    return btoa(unescape(encodeURIComponent(this.str)));
                }
            }

        },
        /**
         * 对已经掩码处理的字符串进行解码操作
         */
        unMaskStr: function (str) {
            return {
                str: str,
                toString: function () {
                    return decodeURIComponent(escape(atob(this.str)));
                }
            }
        },
    };
    exports(MOD_NAME, myUtil);
});