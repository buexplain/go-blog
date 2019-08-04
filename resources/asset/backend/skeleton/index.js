//初始化页面骨架
var skeleton = null;
$(function () {
    layui.use(['skeleton', 'jquery'], function() {
        $.get('/backend/skeleton/all', {}, function (json) {
            if(json.code === 0) {
                var node = json.data;
                skeleton = layui.skeleton('left-nav', 'top-tab');
                //初始化菜单栏
                skeleton.menu.init(node, 0, 'ID', 'Pid', 'Name', 'URL');
                //打开第一个节点
                skeleton.menu.open(node[0].ID);
            }else {
                alert(json.message);
            }
        });
    });
});

/**
 * 忘记密码
 */
function forget() {
    var id  = '2019-07-29';
    if(skeleton.tab.has(id) === false) {
        skeleton.tab.add("修改密码", "/backend/user/forget", id);
    }
    skeleton.tab.change(id);
}


/**
 * 退出登录
 */
function signOut() {
    window.location.href = '/backend/sign?_method=delete';
}