//初始化页面骨架
var skeleton = null;
$(function () {
    layui.use(['skeleton', 'jquery'], function() {
        try {
            skeleton = layui.skeleton('left-nav', 'top-tab');
            //初始化菜单栏
            skeleton.menu.init(myMenu, 'ID', 'Pid', 'Name', 'URL');
            //打开第一个节点
            skeleton.menu.open(myMenu[0].ID);
        }catch (e) {
            console.log('初始化后台骨架失败：'+e);
        }

    });
});

/**
 * 忘记密码
 */
function forget() {
    var id  = '2019-07-29';
    if(skeleton.tab.has(id) === false) {
        skeleton.tab.add("修改密码", "/backend/home/user/forget", id);
    }
    skeleton.tab.change(id);
}


/**
 * 退出登录
 */
function signOut() {
    window.location.href = '/backend/sign?_method=delete';
}