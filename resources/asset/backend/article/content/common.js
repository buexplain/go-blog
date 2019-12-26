/**
 * 内容
 */
class Content {
    constructor(data) {
        data = data || null;
        this.data = data;
    }

    online(id) {
        if(this.data === null) {
            return false;
        }
        return  this.data.Content.Online === id;
    }

    getCategory() {
        if(this.data === null || this.data.Category === null || this.data.Category.length === 0) {
            return false;
        }
        return this.data.Category[0].ID;
    }

    hasTag(id) {
        if(this.data === null || this.data.Tag === null) {
            return false;
        }
        return this.data.Tag.some(function (v, k) {
            return v.ID === id
        });
    }

    getBody() {
        if(this.data === null) {
            return '';
        }
        return this.data.Content.Body;
    }

    static getInstance(contentID) {
        return new Promise(function(resolve, reject) {
            if(!contentID) {
                return resolve(new Content(null));
            }
            $.ajax({
                url: "/backend/article/content/show/"+contentID,
                async:true,
                data: {},
                type: "get",
                success: function (json) {
                    if (json.code !== 0) {
                        reject(json.message);
                    } else {
                        resolve(new Content(json.data));
                    }
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    reject(textStatus+errorThrown);
                }
            });
        });
    }
}

/**
 * 所有的标签
 */
class TagList {
    constructor(data) {
        data = data || null;
        this.data = data;
    }

    process (content) {
        var result = [];
        this.data.forEach(function (v, k) {
            result.push({name: v.Name, value: v.ID, selected: content.hasTag(v.ID)});
        });
        return result;
    }

    static getInstance() {
        return new Promise(function(resolve, reject) {
            $.ajax({
                url: "/backend/article/content/tag",
                async:true,
                data: {},
                type: "get",
                success: function (json) {
                    if (json.code !== 0) {
                        reject(json.message);
                    } else {
                        resolve(new TagList(json.data));
                    }
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    reject(textStatus+errorThrown);
                }
            });
        });
    }
}

/**
 * 所有的分类
 */
class CategoryList {
    constructor(data) {
        data = data || null;
        this.data = data;
    }

    process () {
        var result = [];
        this.data.forEach(function (v, k) {
            result.push({
                "id":v.ID,
                "title": v.Name,
                "parentId": v.Pid
            });
        });
        return result;
    };

    static getInstance() {
        return new Promise(function(resolve, reject) {
            $.ajax({
                url: "/backend/article/content/category/-1",
                async:true,
                data: {},
                type: "get",
                success: function (json) {
                    if (json.code !== 0) {
                        reject(json.message);
                    } else {
                        resolve(new CategoryList(json.data));
                    }
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    reject(textStatus+errorThrown);
                }
            });
        });
    }
}

/**
 * 上下线
 */
class Online {
    constructor() {
        this.data = [{ID:1, Name:'上线'}, {ID:2, Name:'下线'}];
    }

    process (content) {
        var result = [];
        this.data.forEach(function (v, k) {
            result.push({name: v.Name, value: v.ID, selected: content.online(v.ID)});
        });
        return result;
    };

    static getInstance() {
        return new Promise(function(resolve, reject) {
            resolve(new Online());
        });
    }
}

/**
 * 图片懒加载
 */
const lazyLoadImage = () => {
    const loadImg = (it) => {
        const testImage = document.createElement('img');
        testImage.src = it.getAttribute('data-src');
        testImage.addEventListener('load', () => {
            it.src = testImage.src;
            it.style.backgroundImage = 'none';
            it.style.backgroundColor = 'transparent';
        });
        it.removeAttribute('data-src')
    };

    if (!('IntersectionObserver' in window)) {
        document.querySelectorAll('img').forEach((data) => {
            if (data.getAttribute('data-src')) {
                loadImg(data)
            }
        });
        return false
    }

    if (window.imageIntersectionObserver) {
        window.imageIntersectionObserver.disconnect();
        document.querySelectorAll('img').forEach(function (data) {
            window.imageIntersectionObserver.observe(data)
        });
    } else {
        window.imageIntersectionObserver = new IntersectionObserver((entries) => {
            entries.forEach((entrie) => {
                if ((typeof entrie.isIntersecting === 'undefined'
                    ? entrie.intersectionRatio !== 0
                    : entrie.isIntersecting) && entrie.target.getAttribute('data-src')) {
                    loadImg(entrie.target)
                }
            })
        });
        document.querySelectorAll('img').forEach(function (data) {
            window.imageIntersectionObserver.observe(data)
        });
    }
};

class Upload {
    constructor(data) {
        data = data || null;
        this.data = data;
    }

    getMarkdownText () {
        if(this.data === null) {
            return '';
        }
        var images = ['jpeg', 'gif', 'jpg', 'png', 'bmp'];
        images.is = function(url) {
            for(var i in this) {
                if(url.substr(url.length - this[i].length, this[i].length).toLocaleLowerCase() === this[i]) {
                    return true;
                }
            }
            return false;
        };
        var tmp = '['+this.data.Name+'](/'+this.data.Path+' "'+this.data.Name+'")';
        if(images.is(this.data.Path)) {
            tmp = '!'+tmp;
        }
        return tmp;
    };

    static getInstance(formData) {
        return new Promise(function(resolve, reject) {
            $.ajax({
                url: "/backend/article/content/upload",
                async:true,
                data: formData,
                type: "post",
                processData: false,
                contentType: false,
                success: function (json) {
                    if (json.code !== 0) {
                        reject(json.message);
                    } else {
                        resolve(new Upload(json.data));
                    }
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    reject(textStatus+errorThrown);
                }
            });
        });
    }
}

class Uploads
{
    constructor(data) {
        data = data || null;
        this.data = data;
    }

    getMarkdownText () {
        if(this.data === null) {
            return '';
        }
        var tmp = '';
        for(var i in this.data) {
            if(tmp !== '') {
                tmp += '\r\n\r\n';
            }
            tmp += this.data[i].getMarkdownText();
        }
        return tmp;
    }

    static getInstance(formDataArr) {
        var arr = [];
        for (var i in formDataArr) {
            arr.push(Upload.getInstance(formDataArr[i]));
        }
        return Promise.all(arr).then(function(values) {
            return new Promise(function(resolve, reject) {
                resolve(new Uploads(values));
            });
        });
    }
}

//更多操作
layui.use(['jquery'], function () {
    var $ = layui.jquery;
    var more = $("#j-more");
    more.hide();
    var moreBtn = $("#j-more-btn");
    moreBtn.on('click', function () {
        if(more.css('display') === 'none') {
            more.show();
            moreBtn.removeClass('layui-icon-down').addClass('layui-icon-up');
        }else{
            more.hide();
            moreBtn.removeClass('layui-icon-up').addClass('layui-icon-down');
        }
    });
});
