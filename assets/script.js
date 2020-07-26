(function () {

    let autoSaveHandle;
    const autoSave = function () {
        saveArticle(1, 0);
        let now = new Date();
        now = now.toLocaleString();
        $('.alert-info').html("Last AutoSave:" + now).fadeIn(300)
    }

    function refreshContent() {
        let content = $(".article-input-content").val();
        if (content == "") {
            content = "View content here!"
        }
        $('.article-view').html(marked(content));
    }

    if ($(".article-input-content").length > 0) {
        autoSaveHandle = setInterval(refreshContent, 500);
    }

    if ($(".article-view").length > 0) {
        //自动保存草稿
        // setInterval(autoSave, 10000);

    }


    //打开或者关闭预览
    const swichPerview = function () {
        $("#viewSwitch").click(function () {
            if ($(".markdown-content").hasClass("full-width")) {
                $(".markdown-content").removeClass("full-width");
                $(".article-view").show();
            } else {
                $(".markdown-content").addClass("full-width");
                $(".article-view").hide();
            }
        })
    }();

    /**
     * 管理排序按钮
     */
    const obSelect = function () {
        let obBtns = $("#ob-select");
        obBtns.children().each(function (e) {
            if ($(this).attr("obstring") == obBtns.attr("nowob")) {
                $(this).addClass("btn-primary");
                let sign = $(this).attr("obstatus") == "asc" ? " ↑ " : " ↓ ";
                $(this).html($(this).html() + sign);
            }
        });
        obBtns.children().on("click", function () {
            let path = window.location.pathname;
            let obString = $(this).attr("obstring");
            let obStatus = $(this).attr("obstatus") == "asc" ? "desc" : "asc";
            let p = getRequestParam("p");
            path = addURLParam(path, "obstring", obString);
            path = addURLParam(path, "obstatus", obStatus);
            if (p != "") {
                path = addURLParam(path, "p", p);
            }
            window.location.href = path;
        })
    }();

    const drawNav = function () {
        let totalRows = $(".page-nav").attr("totalRows") * 1;
        let curPage = $(".page-nav").attr("curPage") * 1;
        let perPage = 20;
        let totalPages = Math.ceil(totalRows / perPage);
        let navArr = Array();
        let uri = location.href;
        if (totalPages <= 20) {
            for (let i = 1; i <= totalPages; i++) {
                if (curPage === i) {
                    navArr.push($("<li class='active'><a class='curPage' href='#'> " + i + " </a></li>"));
                } else {
                    uri = addURLParam(uri, "page", i);
                    let liEle = "<li ><a href='" + uri + "'>" + i + "</a></li>";
                    navArr.push($(liEle));
                }
            }
        } else {
            //渲染最前面五页
            for (let i = 1; i <= 5; i++) {
                if (curPage === i) {
                    navArr.push($("<li class='active'><a class='curPage' href='#'> " + i + " </a></li>"));
                } else {
                    uri = addURLParam(uri, "page", i);
                    let liEle = "<li ><a href='" + uri + "'>" + i + "</a></li>";
                    navArr.push($(liEle));
                }
            }
            // 渲染中间部分
            if (curPage > 4 && curPage < (totalPages - 3)) {

                if (curPage > 8) {
                    let middlePage = Math.ceil((curPage) / 2) == 5 ? 6 : Math.ceil((curPage) / 2);
                    uri = addURLParam(uri, "page", middlePage)
                    let liEle = "<li ><a href='" + uri + "'> ... </a></li>";
                    navArr.push($(liEle));
                }

                let forStart = curPage > 8 ? curPage - 2 : 6;
                let forEnd = curPage < (totalPages - 6) ? (curPage + 2) : totalPages - 5;
                for (let i = forStart; i <= forEnd; i++) {
                    if (curPage === i) {
                        navArr.push($("<li class='active'><a class='curPage' href='#'> " + i + " </a></li>"));
                    } else {
                        uri = addURLParam(uri, "page", i);
                        let liEle = "<li ><a href='" + uri + "'>" + i + "</a></li>";
                        navArr.push($(liEle));
                    }
                }
                if (curPage < (totalPages - 7)) {
                    let middlePage = Math.ceil(curPage + (totalPages - curPage) / 2);
                    uri = addURLParam(uri, "page", middlePage)
                    let liEle = "<li ><a href='" + uri + "'> ... </a></li>";
                    navArr.push($(liEle));
                }
            }
            // 渲染最后五页
            for (let i = (totalPages - 4); i <= totalPages; i++) {
                if (curPage == i) {
                    navArr.push($("<li class='active'><a class='curPage' href='#'> " + i + " </a></li>"));
                } else {
                    uri = addURLParam(uri, "page", i);
                    let liEle = "<li ><a href='" + uri + "'>" + i + "</a></li>";
                    navArr.push($(liEle));
                }
            }
        }
        //添加前一页和后一页
        if (curPage == 1) {
            navArr.unshift($("<li class='disabled'><a href='#'> ← </a></li>"));
        } else {
            uri = addURLParam(uri, "page", curPage - 1);
            let liEle = "<li ><a href='" + uri + "'> ← </a></li>";
            navArr.unshift($(liEle));
        }
        if (curPage == totalPages) {
            navArr.push($("<li class='disabled'><a href='#'> → </a></li>"));
        } else {
            uri = addURLParam(uri, "page", curPage * 1 + 1);
            let liEle = "<li ><a href='" + uri + "'> → </a></li>";
            navArr.push($(liEle));
        }
        $(".pager").append(navArr)
    }();

    /**
     * 向url追加参数,如果原来参数存在,则会移除后重新添加
     * @param url
     * @param name
     * @param value
     * @returns {*}
     */
    function addURLParam(url, name, value) {
        if (url.indexOf(name) != -1) {
            let tname = name;
            if (url.indexOf("?" + name) != -1) {
                //兼容处理?name的情况
                tname = "?" + name;
            }
            if (url.indexOf("&" + name) != -1) {
                tname = "&" + name;
            }
            sStr = url.split(tname);
            if (sStr[1].indexOf("&") == -1) {
                //nothing after
                url = sStr[0];
            } else {
                url = sStr[0] + sStr[1].substring(sStr[1].indexOf("&"))
            }
        }
        url += (url.indexOf("?") == -1 ? "?" : "&");
        url += encodeURI(name) + "=" + encodeURI(value);
        return url;
    }

    /**
     * 获取制定的参数
     * @returns {Object}
     * @constructor
     */
    function getRequestParam(param) {
        let url = location.href;
        let paraString = url.substring(url.indexOf("?") + 1, url.length).split("&");
        let paraObj = {}
        for (let i = 0; j = paraString[i]; i++) {
            paraObj[j.substring(0, j.indexOf("=")).toLowerCase()] = j.substring(j.indexOf("=") + 1, j.length);
        }
        let returnValue = paraObj[param.toLowerCase()];
        if (typeof (returnValue) == "undefined") {
            return "";
        } else {
            return returnValue;
        }
    }

    const CtriI = function () {
        //捕获  ctrl+i 操作
        $(window).keydown(function (e) {
            if (e.keyCode === 73 && e.ctrlKey) {
                e.preventDefault();
                $('#myModal').modal("show");
            }
        });
    }();

    const resetHeight = function () {
        let bodyHeight = $("body").outerHeight();
        $(".contents").height(bodyHeight - 71);
        $(".list-content , .albums-list").height(bodyHeight - 71 - 100)
        $(".article-input-content").height(bodyHeight - 71 - 100 - 60);
        $(".article-view").height(bodyHeight - 71 - 100 - 60 - 8);
        let totalrows = $(".articles-content-tbody").children().length;
        listHeight = (4.8 * totalrows) + "%";
        $(".articles-list").height(listHeight)

    }();

    //获取剪贴板中的图片
    const getPasteImg = function () {
        document.addEventListener('paste', function (event) {
            const items = event.clipboardData && event.clipboardData.items;
            let file = null;
            if (items && items.length) {
                // 检索剪切板items
                for (let i = 0; i < items.length; i++) {
                    if (items[i].type.indexOf('image') !== -1) {
                        file = items[i].getAsFile();
                        break;
                    }
                }
            }
            let name = window.btoa(Math.random() * 200 + file.name + file.size)
            console.log(name)
            name = name.substring(0, 10);
            uploadImg(file, name);
        });
    }();
    /**
     * 定义插入url的方法
     */
    const insertUrl = function () {
        $("#doUrlInsert").on("click", function (e) {
            let urlName = $('#urlName').val();
            let urlLink = $('#urlLink').val();
            if (urlLink == "") {
                return false;
            }
            if (urlName == "") {
                urlName = "Clikc Here"
            }
            let contentBox = $(".article-input-content")
            let linkText = "[" + urlName + "](" + urlLink + ") \n";
            contentBox.val(contentBox.val() + linkText);
            $('#myModal').modal('hide');
        });
    }();

    /**
     * 通过上传来插入图片到内容中
     */
    const inputUpload = function () {
        $("#doImgUpload").on("click", function (e) {
            let inputFile = $('#inputFile')[0].files[0];
            let fileName = $('#fileName').val();
            if (inputFile == undefined) {
                return false;
            }
            if (fileName == "") {
                fileName = inputFile.name
            }
            var res = uploadImg(inputFile, fileName);
            console.log(res);
            let contentBox = $(".article-input-content");
            var markDownStr = "![" + fileName + "](" + res.viewUrl + ") \n";
            contentBox.val(contentBox.val() + markDownStr);
            $('#myModal').modal('hide');
        });
    }();
    /**
     * 上传文章头图
     */
    const headImageUpload = function(){
        $("#headImageUpload").on("click", function (e) {
            let inputFile = $('#headImage')[0].files[0];
            if (inputFile == undefined) {
                return false;
            }
            fileName = inputFile.name
            var res = uploadImg(inputFile, fileName);
            $(".pic-view").attr("src",res.viewUrl);
            $("#headImageUrl").val(res.viewUrl);
            // console.log(res);
        });
    }();

    /**
     * file 为文件对象
     * @param file
     * @param name
     */
    function uploadImg(file, name) {
        var response;
        let data = new FormData();
        data.append("file", file);
        data.append("filename", name);
        $.ajax({
            type: "post",
            url: "/admin/api/upload",
            data: data,
            contentType: false,
            async:false,
            //设置之后multipart/form-data
            processData: false,
            // 默认情况下会对发送的数据转化为对象 不需要转化的信息
            success: function (res) {
                let obj = $(".save-notice");
                if (res.errNo != 0) {
                    obj.removeClass("alert-info").addClass("alert-warning").html("上传错误" + res.errMsg).show(300).delay(3000).hide(300);
                } else {
                    obj.removeClass("alert-warning").addClass("alert-info").html("上传成功").show(300).delay(3000).hide(300);
                    response = res;
                }
            },
        });
        return response;
    }

    const changeSwitch = function () {
        $(".albums-select").on("click", function (e) {
            if (!$(this).hasClass("btn-primary")) {
                $(".albums-select").removeClass("btn-primary");
                $(this).addClass("btn-primary");
                $("#isPublic").val($(this).attr("targetval"))
            }
        })
    }();
    /**
     * 保存专辑信息
     */
    const saveAlbum = function () {
        $(".save-album").on("click", function () {
            let form = document.forms.namedItem("albumEditForm");
            let data = new FormData(form);
            if (data.get("albumName") == "") {
                $(".albums-notice").show();
                return false;
            }
            $.ajax({
                type: "post",
                url: "/admin/albums/save",
                data: data,
                contentType: false,
                //设置之后multipart/form-data
                processData: false,
                // 默认情况下会对发送的数据转化为对象 不需要转化的信息
                success: function (res) {
                    $(".albums-notice").hide();
                    if (res.errNo == 0) {

                        $('#albumEdit').modal('hide');
                        window.location.reload()
                    } else {
                        //#TODO 错误处理
                        console.log("错误发生了");
                    }

                },
            });
        });
    }();

    const editAlbum = function () {
        $(".editAlbum").on("click", function (e) {
            $(".albums-notice").hide();
            let id = $(this).attr("albumId");
            $.ajax({
                type: "get",
                url: "/admin/albums/info/" + id,
                success: function (res) {
                    if (res.errNo == 0) {
                        $("#albumId").val(res.data.id);
                        $("#albumName").val(res.data.albumName);
                        $("#isPublic").val(res.data.isPublic);
                        if (res.data.isPublic == 1) {
                            $(".albums-public").click();
                        } else {
                            $(".albums-private").click();
                        }
                        $('#albumEdit').modal('show');
                    } else {
                        //#todo 错误处理
                        console.log('#TODO ,获取数据失败了');
                    }

                },
            });
        });
    }();

    const pubClick = function () {
        $(".save-draft").on("click", function () {
            saveArticle(1, 0);
            $('.alert-info').html("save success！").fadeIn(300).delay(3000).fadeOut(400);
        });

        $(".save-publish").on("click", function () {
            clearInterval(autoSaveHandle);
            saveArticle(0, 0);
            $('.alert-info').html("publish success! return to article list。").fadeIn(300).delay(300).fadeOut(100, function () {
                window.location.href = "/admin/articles?obstring=mtime&obstatus=desc";
            });

        });

        $(".save-publish-single").on("click", function () {
            clearInterval(autoSaveHandle);
            saveArticle(0, 1);
            $('.alert-info').html("publish success! return to article list。").fadeIn(300).delay(300).fadeOut(100, function () {
                window.location.href = "/admin/articles?obstring=mtime&obstatus=desc";
            });
        });

    }();

    const tagSelect = function (){
        $(document).on("ready",function(){
            var currentTagsIdStr = $("#tagsValue").val();
            if (currentTagsIdStr != undefined) {
                var currentTags = currentTagsIdStr.split(",")
            }
            $(".tag-content").children().each(function(){
               var id = $(this).attr("tagId");
               if ($.inArray(id,currentTags) != -1) {
                   $(this).toggleClass("label-primary").toggleClass("label-default");
               }
            });
        });
        let selectedCount = 0;
        $(document).on("click",".tag-label",function(e){
            if ($(".tag-content").children(".label-primary").length <= (4-selectedCount)){
                //可以继续选择的情况
                $(this).toggleClass("label-primary").toggleClass("label-default");

            } else if($(this).hasClass("label-primary")) {
                //已经超过5个标签了
                $(this).toggleClass("label-primary").toggleClass("label-default");
            }
        });
        $(".saveTags").on("click",function(){
            var tagsValues = Array();
            var tagsString = Array();
            $(".tag-content").children(".label-primary").each(function(){
                tagsValues.push($(this).attr("tagId"));
                tagsString.push($(this).html());
            });
            $("#tagsValue").val(tagsValues.join(","));
            $("#tagsShow").val(tagsString.join(","));
        });
    }();




    function saveArticle(draft = 1, single = 0) {
        var data = new Object();
        data['content'] = $(".article-input-content").val();
        data['title'] = $(".title-input").val();
        if(data['title'] == "") {
            data['title'] = "temp article";
        }
        data['id'] = $(".articleId").val();
        data['pubStatus'] = $(".pubStatus").val();
        data['albumId'] = $(".albumId").val();
        data['tags'] = $("#tagsValue").val();
        if (draft == 1) {
            //存草稿
            data['pubStatus'] = 0;
        } else {
            data['pubStatus'] = 1;
        }
        if (single == 1) {
            //独立页面
            data['independPage'] = 1;
        }
        $.post({
            url: "/admin/api/articlesave",
            data: data,
            success: function (res) {
                $(".articleId").val(res.articleId);
            },
        });
    }

    /**
     * 删除文章的方法
     */
    const removeArticle = function(){
        $(".articleRemove").on("click",function(){
            console.log($(this).attr("removeId"));
            $(".removeBtn").attr("articleid",$(this).attr("removeId"));
        });
        $(".removeBtn").on("click",function(){
            $.get({
                url: "/admin/api/articleremove?id="+$(this).attr("articleid"),
                success: function (res) {
                    console.log(res);
                    $('#remove').modal('hide');
                    window.location.reload();
                },
            });
        })
    }();


    const articleSetting = function(){
        $(".setting,.settingsave").on("click",function(){
            var settingPanle = $(".setting-body").toggle()
        })
    }();


    /**
     * 添加tag的方法
     */
    const tagAdd = function(){
        $(".tagInputs").on("submit",function(){
            var tagName = $("#newTagName").val();
            console.log(tagName)
            //添加的标签不能为空
            if (tagName == "" || tagName == undefined) {
                return false;
            }
            $(".tag-content").children().each(function(){
                //不能添加已有标签
                if (tagName == $(this).html()) {
                    return false;
                }
            });
            $.post({
                url: "/admin/api/tagadd",
                data: {"tagName":tagName},
                success: function (res) {
                    tagInfo = $('<span class="label label-default tag-label" tagId="'+res.data.tagId+'">'+res.data.tagName+'</span>');
                    $(".tag-content").append(tagInfo);
                    $("#newTagName").val("");
                },
            });
            return false;
        })
    }();

    //文件结尾
})();
