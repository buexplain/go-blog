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

	/* 如果不是详情页面，则加载部件 */
	if(window.location.href.indexOf('article') === -1) {
		var tagHTML = $("#j-tag");
		var placeHTML = $("#j-place");
		$.get('/index-widget', {categoryID:$("#j-categoryID").attr("data-categoryID"), tagID:tagHTML.attr('data-tagID'), place:placeHTML.attr('data-place')}, function (json) {
			if(json.code === 0) {
				tagHTML.html(json.data.tag);
				placeHTML.html(json.data.place);
			}else {
				tagHTML.html(json.message);
			}
		});
	}
});
