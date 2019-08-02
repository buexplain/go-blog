//初始化页面骨架
var skeleton = null;
$(function () {
    layui.use(['skeleton', 'jquery'], function() {
        var node = [{"name":"我的桌面","id":1,"pid":0,"url":"/backend/home"}, {"name":"菜单管理","id":2,"pid":0,"url":"/backend/menu"}];
        skeleton = layui.skeleton('left-nav', 'top-tab');
        //初始化菜单栏
        skeleton.menu.init(node, 0, 'id', 'pid', 'name', 'url');
        //打开第一个节点
        skeleton.menu.open(node[0].id);
        return;
        $.get('/admin/skeleton/node/json?isNav?=1', {isNav:1}, function (json) {
            if(json.code === 0) {
                skeleton = layui.skeleton('left-nav', 'top-tab');
                skeleton.nav.set(json.data, 0, 'ID', 'Pid', 'Name', 'URL');
                skeleton.nav.open(json.data[0].ID);
            }else {
                alert(json.msg);
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
        skeleton.tab.add("修改密码", "/admin/user/forget", id);
    }
    skeleton.tab.change(id);
}