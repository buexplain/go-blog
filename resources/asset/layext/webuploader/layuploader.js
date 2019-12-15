layui.define(['webuploader', 'table', 'layer', 'jquery', 'element'], function(exports) {
    layui.link('/resources/asset/layext/webuploader/webuploader.css');
    var webuploader = layui.webuploader;
    var $ = layui.jquery;
    var table = layui.table;
    var layer = layui.layer;
    var element = layui.element;

    /**
     * 列表
     */
    class Lists {
        constructor(option) {
            var that = this;
            this.upload = null;
            this.lay_filter = webuploader.Base.guid();
            option = option || {};
            this.option = option;
            if(!this.option.elem) {
                throw "not found option: elem";
            }
            var elemObj = null;
            if(typeof this.option.elem === "string") {
                elemObj = $(this.option.elem);
            }else {
                //jquery对象
                if(this.option.elem.context) {
                    elemObj = this.option.elem;
                }else {
                    //dom对象
                    elemObj = $(this.option.elem);
                }
            }
            if(elemObj === null) {
                throw "not found option: elem obj";
            }
            //设置filter属性，用于后续表格的右侧按钮的监听
            elemObj.attr('lay-filter', this.lay_filter);
            //设置id属性，用户重置表格行高
            if(this.option.id === undefined) {
                if(elemObj.attr('id')) {
                    //存在id，记录id
                    this.option.id = elemObj.attr('id');
                }else {
                    //不存在id，设置id
                    this.option.id = 'j-'+this.lay_filter;
                    elemObj.attr('id', this.option.id);
                }
            }
            //重置条数
            this.option.limit = 1000;
            this.option.data = [];
            this.option.cols = [[
                {field: 'name', title: '名称', fixed: 'left',style:'text-align:center;', templet: function (data) {
                    if(data.src) {
                        var html = '<span><img src="'+data.src+'" title="'+data.name+'">';
                        html += '<br>'+data.name;
                        html += '</span>';
                        return html;
                    }
                    return data.name;
                }}
                ,{field: 'size', title: '大小'}
                ,{field: 'progress', title: '进度',templet: function (data) {
                    var html = '';
                    html += '<div class="layui-progress layui-progress-big" lay-showPercent="true" lay-filter="j-progress'+that.lay_filter+data.id+'">';
                    html += '<div class="layui-progress-bar" lay-percent="'+data.progress+'%" style="width: '+data.progress+'%;">';
                    html += '<span class="layui-progress-text">'+data.progress+'%</span>';
                    html += '</div>';
                    html += '</div>';
                    return html;
                }}
                ,{field: 'status', title: '状态', templet: function (data) {
                    var html = '<span id="j-status'+that.lay_filter+data.id+'">'+data.status+'</span>';
                    return html;
                }}
                ,{fixed: 'right', title:'操作', templet:function (data) {
                    var btn = '';
                    btn += '<button class="layui-btn layui-btn-sm layui-btn-danger" id="j-del'+that.lay_filter+data.id+'" lay-event="del">删除</button>';
                    return btn;
                }}
            ]];
            that.option.done = function () {
                //如果有缩略图，则设置为自动高度
                for (var i in that.option.data) {
                    if(that.option.data[i].src !== '') {
                        var css = {
                            'height': 'auto'
                        };
                        var tmp = $.find('div[lay-id="'+that.option.id+'"]');
                        if(tmp.length > 0) {
                            $(tmp[0]).find('.layui-table-cell').css(css);
                        }else {
                            $(".layui-table-cell").css(css);
                        }
                        break;
                    }
                }
            };
            this.instance = null;
            table.on('tool('+this.lay_filter+')', function(obj) {
                var data = obj.data;
                if(obj.event === 'del') {
                    var tmp = [];
                    for(var i in that.option.data) {
                        if(that.option.data[i].id !== data.id) {
                            tmp.push(that.option.data[i]);
                        }
                    }
                    that.option.data = tmp;
                    obj.del();
                    that.upload._removeFile(data.id);
                }
            });
        }

        _setUpload (upload) {
            this.upload = upload;
        }

        _addFile (id, name, src, size, status) {
            this.option.data.push({id:id, name:name, src:src, size:webuploader.Base.formatSize(size), progress:0, status:status});
            if(this.instance === null) {
                this.instance = table.render(this.option);
            }else {
                this.instance.reload(this.option);
            }
        }

        _delBtnLock(id) {
            var btn_id = '#j-del'+this.lay_filter+id;
            var o = $(btn_id);
            if(o.length > 0) {
                o.attr('disabled', 'disabled');
                o.addClass('layui-btn-disabled');
            }
        }

        _delBtnUnLock(id) {
            var btn_id = '#j-del'+this.lay_filter+id;
            var o = $(btn_id);
            if(o.length > 0) {
                o.removeAttr('disabled');
                o.removeClass('layui-btn-disabled');
            }
        }

        /**
         * 设置状态进度条
         * @param id
         * @param progress
         */
        progress(id, progress) {
            for(var i in this.option.data) {
                if(this.option.data[i].id === id) {
                    this.option.data[i].progress = progress;
                    var progress_id = 'j-progress'+this.lay_filter+id;
                    element.progress(progress_id, this.option.data[i].progress+'%');
                    break;
                }
            }
        }

        /**
         * 设置状态描述
         * @param id
         * @param status
         */
        status(id, status) {
            for(var i in this.option.data) {
                if(this.option.data[i].id === id) {
                    this.option.data[i].status = status;
                    var status_id = 'j-status'+this.lay_filter+id;
                    document.getElementById(status_id).innerHTML = this.option.data[i].status;
                    break;
                }
            }
        }
    }

    /**
     * 上传
     */
    class Upload {
        constructor(option) {
            var that = this;
            this.lists = null;
            option = option || {};
            this.option = option;
            //给开始上传按钮绑定上传事件
            if(!that.option.send_btn) {
                throw "not found option: send_btn";
            }
            that.send_btn = $(that.option.send_btn);
            if(that.send_btn.length === 0) {
                throw "not found send button: " + that.option.send_btn;
            }else if(that.send_btn[0].nodeName.toLocaleLowerCase() !== 'button') {
                throw "send button must be button label: " + that.option.send_btn;
            }
            if(that.send_btn.length >1 ) {
                that.send_btn = that.send_btn.eq(0);
            }
            that.send_btn.on('click', function () {
                that._send();
            });
            this.instance = webuploader.create(this.option);
            //添加文件
            this.instance.on('fileQueued', function(file) {
                if (file.getStatus() === webuploader.Base.File.Status.INVALID) {
                    that.lists._addFile(file.id, file.name, '', file.size, '文件错误: '+file.statusText);
                }else {
                    //文件预览
                    var ratio = window.devicePixelRatio || 1, thumbnailWidth = 110 * ratio, thumbnailHeight = 110 * ratio;
                    that.instance.makeThumb(file, function(error, src) {
                        if (error) {
                            that.lists._addFile(file.id, file.name, '', file.size, '等待上传');
                        }else {
                            that.lists._addFile(file.id, file.name, src, file.size, '等待上传');
                        }
                    }, thumbnailWidth, thumbnailHeight);
                }
            });
            //监听错误
            this.instance.on('error', function(handler) {
                if(handler ==='F_EXCEED_SIZE') {
                    layer.msg('文件超出上传大小', {icon: 2});
                }else if(handler==='Q_TYPE_DENIED') {
                    layer.msg('此类文件不允许上传', {icon: 2});
                }else if(handler === 'F_DUPLICATE') {
                    layer.msg('文件已经在队列中', {icon: 1});
                }
            });
            //当开始上传流程时触发
            this.instance.on('startUpload', function () {
                //锁定上传按钮
                that._sendBtnLock();
            });
            //某个文件开始上传前触发，一个文件只会触发一次。
            this.instance.on('uploadStart', function (file) {
                //锁定删除按钮
                that.lists._delBtnLock(file.id);
            });
            //不管成功或者失败，文件上传完成时触发。
            this.instance.on('uploadComplete', function (file) {
                //解锁删除按钮
                that.lists._delBtnUnLock(file.id);
            });
            //更新上传进度条
            this.instance.on('uploadProgress', function(file, percentage) {
                percentage = (percentage * 100).toFixed(0);
                if(percentage > 2) {
                    //始终减一，只有当服务端返回结果后再更新进度条为100，或者设置其它信息
                    percentage = percentage - 1;
                }
                that.lists.progress(file.id, percentage);
            });
            //上传失败
            this.instance.on('uploadError', function(file, reason) {
                if(file.skipped) {
                    //跳过的文件不做处理
                    return;
                }
                var status = file.getStatus();
                var message = reason === undefined ? "上传错误: "+status : "上传错误: "+status+' '+reason;
                that.lists.status(file.id, message);
                //重置上传状态
                if(status === webuploader.Base.File.Status.ERROR) {
                    file.setStatus(webuploader.Base.File.Status.INITED);
                }
            });
            //当所有文件上传结束时触发
            this.instance.on( 'uploadFinished', function() {
                //释放上传按钮
                that._sendBtnUnLock();
            });
        }

        _setLists (lists) {
            this.lists = lists;
        }

        _removeFile(id) {
            if(this.instance.getFile(''+id)) {
                this.instance.removeFile(id);
            }
        }

        _send() {
            if(this.instance.getFiles().length === 0) {
                layer.msg('没有可上传的文件', {icon: 2});
                return;
            }
            this.instance.upload();
        }

        _sendBtnLock()
        {
            if(this.send_btn.attr('data-old-text') === undefined) {
                this.send_btn.attr('data-old-text', this.send_btn.text());
            }
            this.send_btn.attr('disabled', 'disabled');
            this.send_btn.text('正在上传');
            this.send_btn.addClass('layui-btn-disabled');
        }

        _sendBtnUnLock() {
            this.send_btn.text(this.send_btn.attr('data-old-text'));
            this.send_btn.removeAttr('disabled');
            this.send_btn.removeClass('layui-btn-disabled');
        }
    }

    /**
     * 工厂
     */
    class Factory {
        constructor(upload, lists) {
            this.upload = new Upload(upload);
            this.lists = new Lists(lists);
            this.upload._setLists(this.lists);
            this.lists._setUpload(this.upload);
        }
    }

    exports('layuploader', {factory:Factory, webUploader:webuploader});
});