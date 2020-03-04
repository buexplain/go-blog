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

	/* 加载标签列表 */
	var tagHTML = $("#j-tag");
	$.get('/index-tag?tagID='+tagHTML.attr('data-tagID'), function (data) {
		tagHTML.html(data);
	});

	/* 加载归档列表 */
	var fileHTML = $("#j-file");
	$.get('/index-place?place='+fileHTML.attr('data-place'), function (data) {
		fileHTML.html(data);
	});
});
