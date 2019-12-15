# 基于百度上传组件进行封装的layui上传组件

## 文档
百度上传组件文档：http://fex.baidu.com/webuploader/doc/index.html#WebUploader_Uploader_events

## 升级注意
webuploader.js做了如下更改：
1. 支持layui模块的导出
2. formData参数支持闭包获取

## 使用方式

```JavaScript
//配置组件地址
layui.config({
    base: '/resources/asset/'
}).extend({
    //百度上传组件
    webuploader:'layext/webuploader/webuploader',
    //layupload上传组件
    layuploader:'layext/webuploader/layuploader'
});

layui.use(['layuploader', 'jquery'], function () {
    var layUploader = layui.layuploader;
    var $ = layui.jquery;
    var fc = null;

    layUploader.webUploader.Uploader.register({
        'before-send-file':'checkMd5'
    }, {
        checkMd5: function (file) {
            //秒传md5校验
            var owner = this.owner;
            var deferred = layUploader.webUploader.Deferred();
            fc.lists.status(file.id, '正在校验md5');
            owner.md5File(file, 0 ,file.size).progress(function(percentage) {
                percentage = (percentage * 100).toFixed(0);
                if(percentage > 2) {
                    percentage = percentage - 1;
                }
                fc.lists.progress(file.id, percentage);
            }).then(function(md5) {
                $.get('/backend/article/attachment/check/'+md5, function (json) {
                    if(json.code === 0) {
                        owner.skipFile(file);
                        fc.lists.progress(file.id, 100);
                        fc.lists.status(file.id, '秒传成功');
                        deferred.reject();
                    }else {
                        fc.lists.progress(file.id, 0);
                        fc.lists.status(file.id, '正在上传');
                        deferred.resolve();
                    }
                });
            });
            return deferred.promise();
        }
    });

    fc = new layUploader.factory({
        //开始上传按钮的id
        send_btn:'#j-send_btn',
        pick: {
            //选择文件按钮的id
            id: '#j-webUploader',
            label: '<i class="layui-icon layui-icon-upload"></i>选择文件'
        },
        //服务器地址
        server: '/backend/article/attachment/upload',
        //额外的表单信息
        formData:{
            _token:$("input[name='_token']").val(),
            folder:function () {
                return $("#j-folder").val();
            }
        }
    }, {
        //上传列表的表格的id（数据表格）
        elem: '#j-webUploader-table',
    });

    //当文件上传成功时触发
    fc.upload.instance.on( 'uploadSuccess', function(file, response) {
        if(response.code === 0) {
            fc.lists.progress(file.id, 100);
            fc.lists.status(file.id, '上传成功');
        }else {
            fc.lists.progress(file.id, 0);
            fc.lists.status(file.id, response.message);
        }
    });

    //全部上传完成时候触发
    fc.upload.instance.on('uploadFinished', function() {
        // window.location.reload();
    });
});
```