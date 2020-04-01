layui.use(['jquery'], function() {
	var $ = layui.jquery;
	/* 异步加载标签与归档 开始 */
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
	/* 异步加载标签与归档 结束 */
});
