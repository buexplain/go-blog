layui.use(['jquery'], function() {
	var $ = layui.jquery;
	/* 更新点击量 开始 */
	var hits = $("#j-hits");
	var contentID = hits.attr('data-contentID');
	var hitsKey = "hitsKey"+contentID;
	var param = {
		isIncr: window.localStorage.getItem(hitsKey) === null,
	};
	$.get('/article-hits/'+hits.attr('data-contentID'), param, function (json) {
		if(json.code === 0) {
			hits.append(json.data);
			window.localStorage.setItem(hitsKey, json.data);
		}
	});
	/* 更新点击量 结束 */
});
