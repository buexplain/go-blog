layui.use(['element', 'jquery', 'layer', 'form'], function() {
	var element = layui.element;
	var $ = layui.jquery;
	var layer = layui.layer;
	var form = layui.form;

	/* 汉堡按钮 开始 */
		$(".hamburger-btn").on('click', function() {
			$(".hamburger-nav").toggle();
		});
	/* 汉堡按钮 结束 */

	/* 登录 开始*/
		form.on('submit(j-f-signIn)', function(data) {
			layer.msg("登录成功，页面即将刷新："+JSON.stringify(data.field), function() {
				window.location.reload()
			});
			return false;
		});
	/* 登录 结束*/

	/*注册 开始*/
		form.on('submit(j-f-register)', function(data) {
			layer.msg("注册成功，页面即将刷新："+JSON.stringify(data.field), function() {
				window.location.reload()
			});
			return false;
		});
	/*注册 结束*/

	/*忘记密码 开始*/
	form.on('submit(j-f-forget)', function(data) {
		layer.msg("修改成功，页面即将刷新："+JSON.stringify(data.field), function() {
			window.location.reload()
		});
		return false;
	});
	/*忘记密码 结束*/

	/*修改密码 开始*/
	form.on('submit(j-f-changePassword)', function(data) {
		layer.msg("修改成功，页面即将刷新："+JSON.stringify(data.field), function() {
			window.location.reload()
		});
		return false;
	});
	/*修改密码 结束*/
	
});
