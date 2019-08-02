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
    //树形数据迭代帮助类
    treeHelper: 'backend/skeleton/treeHelper',
    //表格树组件
    treeTable:'layext/treeTable/treeTable'
});