# 基于百度上传组件进行封装的layui上传组件

## 文档
百度上传组件文档：http://fex.baidu.com/webuploader/doc/index.html#WebUploader_Uploader_events

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


```