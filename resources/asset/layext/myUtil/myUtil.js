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
        /**
         * 根据文件后缀解析编辑器配置
         * @param ext 文件后缀
         */
        codeMirrorOptionMap: function (ext) {
            var option = {mode:''};
            var depends = {css:[], js:[]};
            //默认为clike的配置
            option.mode = 'text/x-csrc';
            depends.js = [
                'addon/hint/show-hint.js',
                'mode/clike/clike.js'
            ];
            depends.css = ['addon/hint/show-hint.css'];

            //判断是否为clike文件
            if(['h', 'c', 'cc', 'm', 'mm', 'cpp', 'java', 'scala', 'kt', 'ceylon'].indexOf(ext) !== -1) {
                if(['cpp'].indexOf(ext) !== -1) {
                    option.mode = 'text/x-c++src';
                }
                if(['m', 'mm'].indexOf(ext) !== -1) {
                    option.mode = 'text/x-objectivec';
                }
                if(['scala'].indexOf(ext) !== -1) {
                    option.mode = 'text/x-scala';
                }
                if(['kt'].indexOf(ext) !== -1) {
                    option.mode = 'text/x-kotlin';
                }
                if(['ceylon'].indexOf(ext) !== -1) {
                    option.mode = 'text/x-ceylon';
                }
            }
            //判断是否为php文件
            if(['php'].indexOf(ext) !== -1) {
                option.mode = 'application/x-httpd-php';
                depends.js = [
                    'mode/htmlmixed/htmlmixed.js',
                    'mode/xml/xml.js',
                    'mode/javascript/javascript.js',
                    'mode/css/css.js',
                    'mode/clike/clike.js',
                    'mode/php/php.js'
                ];
                depends.css = [];
            }
            //判断是否为go文件
            if(['go'].indexOf(ext) !== -1) {
                option.mode = 'text/x-go';
                depends.js = [
                    'mode/go/go.js'
                ];
                depends.css = [];
            }
            //判断是否为python文件
            if(['py'].indexOf(ext) !== -1) {
                option.mode = {name: "python",version: 3,singleLineStringErrors: false};
                depends.js = [
                    'mode/python/python.js'
                ];
                depends.css = [];
            }
            //判断是否为css文件
            if(['css'].indexOf(ext) !== -1) {
                option.mode = 'text/css';
                depends.js = [
                    'mode/css/css.js',
                    'addon/hint/show-hint.js',
                    'addon/hint/css-hint.js'
                ];
                depends.css = ['addon/hint/show-hint.css'];
            }
            if(['less'].indexOf(ext) !== -1) {
                option.mode = 'text/x-less';
                depends.js = [
                    'mode/css/css.js'
                ];
                depends.css = [];
            }
            if(['scss'].indexOf(ext) !== -1) {
                option.mode = 'text/x-scss';
                depends.js = [
                    'mode/css/css.js'
                ];
                depends.css = [];
            }
            if(['gss'].indexOf(ext) !== -1) {
                option.mode = 'text/x-gss';
                depends.js = [
                    'mode/css/css.js',
                    'addon/hint/show-hint.js',
                    'addon/hint/css-hint.js'
                ];
                depends.css = ['addon/hint/show-hint.css'];
            }
            //判断是否为html文件
            if(['html'].indexOf(ext) !== -1) {
                option.mode = {
                    name: "htmlmixed",
                    scriptTypes: [
                        {matches: /\/x-handlebars-template|\/x-mustache/i,mode: null},
                        {matches: /(text|application)\/(x-)?vb(a|script)/i,mode: "vbscript"}
                    ]
                };
                option['selectionPointer'] = true;
                depends.js = [
                    'addon/selection/selection-pointer.js',
                    'mode/xml/xml.js',
                    'mode/javascript/javascript.js',
                    'mode/css/css.js',
                    'mode/vbscript/vbscript.js',
                    'mode/htmlmixed/htmlmixed.js'
                ];

                depends.css = [];
            }
            //判断是否为js文件
            if(['js'].indexOf(ext) !== -1) {
                option.mode = 'text/javascript';
                option['continueComments'] = "Enter";
                option['extraKeys'] = {"Ctrl-Q": "toggleComment"};
                depends.js = [
                    'addon/comment/continuecomment.js',
                    'addon/comment/comment.js',
                    'mode/javascript/javascript.js'
                ];
                depends.css = [];
            }
            if(['ts'].indexOf(ext) !== -1) {
                option.mode = 'text/typescript';
                depends.js = [
                    'mode/javascript/javascript.js'
                ];
                depends.css = [];
            }
            if(['json'].indexOf(ext) !== -1) {
                option.mode = 'application/ld+json';
                option['autoCloseBrackets'] = true;
                option['lineWrapping'] = true;
                depends.js = [
                    'addon/comment/continuecomment.js',
                    'addon/comment/comment.js',
                    'mode/javascript/javascript.js'
                ];
                depends.css = [];
            }
            //判断是否为vue文件
            if(['vue'].indexOf(ext) !== -1) {
                option.mode = {
                    name: "vue"
                };
                option['selectionPointer'] = true;
                depends.js = [
                    "addon/mode/overlay.js",
                    "addon/mode/simple.js",
                    "addon/selection/selection-pointer.js",
                    "mode/xml/xml.js",
                    "mode/javascript/javascript.js",
                    "mode/css/css.js",
                    "mode/coffeescript/coffeescript.js",
                    "mode/sass/sass.js",
                    "mode/pug/pug.js",
                    "mode/handlebars/handlebars.js",
                    "mode/htmlmixed/htmlmixed.js",
                    "mode/vue/vue.js"
                ];
                depends.css = [];
            }
            //判断是否为sql文件
            if(['sql'].indexOf(ext) !== -1) {
                option.mode = 'text/x-mariadb';
                option['extraKeys'] = {"Ctrl-Space": "autocomplete"};
                depends.js = [
                    'mode/sql/sql.js',
                    'addon/hint/show-hint.js',
                    'addon/hint/sql-hint.js'
                ];
                depends.css = ['addon/hint/show-hint.css'];
            }
            //判断是否为toml文件
            if(['toml'].indexOf(ext) !== -1) {
                option.mode = {name: "toml"};
                depends.js = [
                    'mode/toml/toml.js',
                ];
            }
            return {option:option, depends:depends};
        }
    };
    exports(MOD_NAME, myUtil);
});