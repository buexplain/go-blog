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

	/* 加载部件 */
	var tagHTML = $("#j-tag");
	var fileHTML = $("#j-file");
	$.get('/index-widget', {tagID:tagHTML.attr('data-tagID'), place:fileHTML.attr('data-place')}, function (json) {
		if(json.code === 0) {
			tagHTML.html(json.data.tag);
			fileHTML.html(json.data.place);
		}else {
			tagHTML.html(json.message);
		}
	});
});
