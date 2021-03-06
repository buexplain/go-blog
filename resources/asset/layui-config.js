//对layui进行全局配置
layui.config({
    base: '/resources/asset/'
}).extend({
    //后台骨架
    skeleton: 'backend/skeleton/skeleton',
    //后台骨架菜单栏部分
    skeletonMenu: 'backend/skeleton/skeletonMenu',
    //后台骨架切换卡部分
    skeletonTab: 'backend/skeleton/skeletonTab',
    //本项目自定义工具集
    myUtil:'layext/myUtil/myUtil',
    //表格树组件
    treeTable:'layext/treeTable/treeTable',
    //树组件
    dtree:'layext/dtree/dtree',
    //select下拉选项卡
    xmSelect:'layext/xm-select/xm-select',
    //百度上传组件
    webuploader:'layext/webuploader/webuploader',
    //layupload上传组件
    layuploader:'layext/webuploader/layuploader',
    //复制组件
    clipboard:'layext/clipboard/clipboard'
});