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

//渲染markdown样式
var preview = document.getElementById('j-preview');
//高亮显示代码部分
Vditor.highlightRender({
	enable:true,
	style:'monokai',
	lineNumber:true
}, preview);
//为 element 中的代码块添加复制按钮
Vditor.codeRender(preview);
//转换 preview 中的文本为数学公式
Vditor.mathRender(preview);
//转换 preview 中 class 为 className 的元素为流程图/时序图/甘特图
Vditor.mermaidRender(preview);
//图表渲染
Vditor.chartRender(preview);
//五线谱渲染
Vditor.abcRender(preview);
//为特定链接分别渲染为视频、音频、嵌入的 iframe
Vditor.mediaRender(preview);
//对使用 Lute 渲染的数学公式进行渲染
Vditor.mathRender(preview);
