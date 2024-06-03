// 全局管理对象
gf = {
    // 刷新验证码
    reloadCaptcha: function () {
        $("img.captcha").attr("src", "/captcha?v=" + Math.random());
    },
}

// 用户模块
gf.user = {
    // 退出登录
    logout: function () {
        swal({
            title: "注销登录",
            text: "您确定需要注销当前登录状态吗？",
            icon: "warning",
            buttons: ["取消", "确定"]
        }).then((value) => {
            if (value) {
                window.location.href = "/user/logout";
            }
        });
    },
}

// 内容模块
gf.content = {
    // 删除内容
    delete: async function (id, url, title) {
        // url = url || "/"
        await personalContentDelete({id}, title, "删除成功", url);
        /* swal({
            title: "删除内容",
            text: "您确定要删除该内容吗？",
            icon: "warning",
            buttons: ["取消", "确定"]
        }).then((value) => {
            if (value) {
                jQuery.ajax({
                    type: 'DELETE',
                    url: '/content/delete',
                    data: {
                        id: id
                    },
                    sync: true,
                    success: function (data) {
                        swal({
                            title: "删除完成",
                            text: "3秒后自动跳转到",
                            icon: "success",
                            timer: 2000,
                            buttons: false
                        }).then((value) => {
                            window.location.href = url;
                        });
                    }
                });
            }
        }); */
    },

}

gf.personal = {
    // 删除我的信息
    deleteMessage: async function (id) {
        await myMessageDelete(id);
    }
}

// 互动模块
gf.interact = {
    // 检查赞
    checkZan: function (elem, id) {
        var type = $(elem).attr("data-type")
        if ($(elem).find('.icon').hasClass('icon-zan-done')) {
            this.cancelZan(elem, type, id)
        } else {
            this.zan(elem, type, id)
        }
    },
    // 赞
    zan: function (elem, type, id) {
        jQuery.ajax({
            type: 'PUT',
            url: '/interact/zan',
            data: {
                id: id,
                type: type
            },
            sync: true,
            success: function (r, status) {
                if (r.code <= 0) {
                    let number = $(elem).find('.number').html()
                    $(elem).find('.number').html(parseInt(number) + 1)
                    $(elem).find('.icon').removeClass('icon-zan').addClass('icon-zan-done')
                } else {
                    swal({
                        text: r.message,
                        button: "确定"
                    })
                }
            }
        });
    },
    // 取消赞
    cancelZan: function (elem, type, id) {
        jQuery.ajax({
            type: 'DELETE',
            url: '/interact/zan',
            data: {
                id: id,
                type: type
            },
            sync: true,
            success: function (r, status) {
                if (r.code <= 0) {
                    let number = $(elem).find('.number').html()
                    $(elem).find('.number').html(parseInt(number) - 1)
                    $(elem).find('.icon').removeClass('icon-zan-done').addClass('icon-zan')
                } else {
                    swal({
                        text: r.message,
                        button: "确定"
                    })
                }
            }
        });
    },
    // 检查是执行踩还是取消踩
    checkCai: function (elem, id) {
        const type = $(elem).attr("data-type");
        if ($(elem).find('.icon').hasClass('icon-cai-done')) {
            this.cancelCai(elem, type, id);
        } else {
            this.cai(elem, type, id);
        }
    },
    // 踩
    cai: function (elem, type, id) {
        /*jQuery.ajax({
            type: 'PUT',
            url : '/interact/cai',
            data: {
                id:   id,
                type: type
            },
            sync: true,
            success: function (r, status) {
                if (r.code <= 0) {
                    let number = $(elem).find('.number').html()
                    $(elem).find('.number').html(parseInt(number)+1)
                    $(elem).find('.icon').removeClass('icon-cai').addClass('icon-cai-done')
                } else {
                    swal({
                        text:   r.message,
                        button: "确定"
                    })
                }
            }
        });*/
        const param = {
            id: id,
            type: type
        }
        clickCai("/interact/cai", "PUT", param, "", elem);
    },
    // 取消踩
    cancelCai: function (elem, type, id) {
        /*jQuery.ajax({
            type: 'DELETE',
            url:  '/interact/cai',
            data: {
                id:   id,
                type: type
            },
            sync: true,
            success: function (r, status) {
                if (r.code <= 0) {
                    let number = $(elem).find('.number').html()
                    $(elem).find('.number').html(parseInt(number) - 1)
                    $(elem).find('.icon').removeClass('icon-cai-done').addClass('icon-cai')
                } else {
                    swal({
                        text: r.message,
                        button: "确定"
                    })
                }
            }
        });*/
        const param = {
            id: id,
            type: type
        }
        cancelCai("/interact/cai", "DELETE", param, "", elem);
    }
}

jQuery(function ($) {
    // 为必填字段添加提示
    $('.required').prepend('&nbsp;<span class="icon iconfont red">&#xe71b;</span>');

    // 回车搜索
    $("#search").keydown(function (e) {
        if (e.keyCode == 13) {
            gf.search();
            e.preventDefault();
        }
    });

    // 分页高亮
    let pageItem = $("ul.pagination li.page-item")
    if (pageItem.length > 4) {
        pageItem.each(function (index, element) {
            if (index < 2 || index > pageItem.length - 3) {
                return
            }
            if ($(element).attr("class").indexOf("disabled") > -1) {
                $(element).removeClass("disabled").addClass("active");
                return
            }
        });
    }

    $.extend($.validator.messages, {
        required: "这是必填字段",
        remote: "请修正此字段",
        email: "请输入有效的电子邮件地址",
        url: "请输入有效的网址",
        date: "请输入有效的日期",
        dateISO: "请输入有效的日期 (YYYY-MM-DD)",
        number: "请输入有效的数字",
        digits: "只能输入数字",
        creditcard: "请输入有效的信用卡号码",
        equalTo: "你的输入不相同",
        extension: "请输入有效的后缀",
        maxlength: $.validator.format("最多可以输入 {0} 个字符"),
        minlength: $.validator.format("最少要输入 {0} 个字符"),
        rangelength: $.validator.format("请输入长度在 {0} 到 {1} 之间的字符串"),
        range: $.validator.format("请输入范围在 {0} 到 {1} 之间的数值"),
        max: $.validator.format("请输入不大于 {0} 的数值"),
        min: $.validator.format("请输入不小于 {0} 的数值")
    });
})

/**
 * 个人中心我的消息删除
 *
 * @param param 请求参数
 */
const myMessageDelete = async (param) => {
    console.log("个人中心消息删除", param);
    await newSwal("提示", `您确定删除回复信息吗`, `warning`, ["取消", "确定"], true).then(async val => {
        if (val) {
            await ajaxPromise(`/reply/delete`, "DELETE", {id: param}, undefined).then(resp => {
                if (resp.code === 0) {
                    swalSingleBtn(undefined, "删除成功", "success", "确定", false).then(()=> {
                        location.reload();
                    });
                } else {
                    swalSingleBtn(undefined, resp.message, "warning", "确定", false);
                }
            });
        } else {
            console.log("取消了内容回复删除操作");
        }
    });
}

/**
 * 个人中心内容删除
 * 
 * @param param       请求参数
 * @param title       内容标题
 * @param message     提示信息
 * @param redirectUrl 跳转URL
 */
const personalContentDelete = async (param, title, message, redirectUrl) => {
    console.log("个人中心内容删除参数", param, redirectUrl);
    await newSwal("提醒", title ? `你确定删除该内容【${title}】吗？`: `你确定删除该信息吗`, `warning`, ["取消", "确定"], true).then(async val => {
        if (val) {
            await ajaxPromise(`/content/delete`, "DELETE", param, message).then(resp => {
                swal({
                    title: resp,
                    text: "2秒后自动刷新当前页面",
                    icon: "success",
                    timer: 2000,
                    buttons: false
                }).then(() => {
                    window.location.href = redirectUrl || "/";
                });
            });
        } else {
            console.log("取消了内容删除操作");
        }
    });
}

/**
 * 点击踩
 *
 * @param url     请求URL
 * @param type    请求类型
 * @param param   请求参数
 * @param message 提示信息
 * @param ele     html元素选择器
 * @returns {Promise<void>}
 */
const clickCai = async (url, type, param, message, ele) => {
    await ajaxPromise(url, type, param, message).then(resp => {
        if (resp.code <= 0) {
            const num = $(ele).find('.number').html();
            $(ele).find('.number').html(parseInt(num) + 1);
            $(ele).find('.icon').removeClass('icon-cai').addClass('icon-cai-done');
        } else {
            swal({
                text: r.message,
                button: "确定"
            });
        }
    });
}

/**
 * 取消踩
 *
 * @param url     请求URL
 * @param type    请求类型
 * @param param   请求参数
 * @param message 提示信息
 * @param ele     html元素选择器
 * @returns {Promise<void>}
 */
const cancelCai = async (url, type, param, message, ele) => {
    await ajaxPromise(url, type, param, message).then(resp => {
        if (resp.code <= 0) {
            const num = $(ele).find('.number').html();
            $(ele).find('.number').html(parseInt(num) - 1);
            $(ele).find('.icon').removeClass('icon-cai-done').addClass('icon-cai');
        } else {
            swal({
                text: resp.message,
                button: "确定"
            });
        }
    });
}


/**
 * 创建弹窗(异步)
 *
 * @param title      标题
 * @param text       文本信息
 * @param icon       消息类型
 * @param buttons    按钮数组
 * @param dangerMode 是否红色提醒
 * @returns {Promise<unknown>}
 */
const newSwal = (title, text, icon, buttons, dangerMode) => {
    return new Promise((resolve, reject) => {
        swal({
            title: title,
            text: text,
            icon: icon,
            //buttons: ["取消", "确定"],
            buttons: buttons,
            dangerMode: dangerMode,
        }).then((value) => {
            resolve(value);
        }).catch((error) => {
            reject(error);
        });
    });
}

/**
 * 创建弹窗(单按钮～异步)
 *
 * @param title      标题
 * @param text       信息文本
 * @param icon       图标类型
 * @param button     按钮文本
 * @param dangerMode 是否显示危险(为true按钮会变成红色)
 * @returns {Promise<unknown>}
 */
const swalSingleBtn = (title, text, icon, button, dangerMode) => {
    return new Promise((resolve, reject) => {
        swal({
            title: title,
            text: text,
            icon: icon,
            button: button,
            dangerMode: dangerMode
        }).then((value) => {
            resolve(value);
        }).catch((error) => {
            reject(error);
        });
    });
}

const deleteReplyFunc = (url, id, successCallback, errorCallback) => {
    jQuery.ajax({
        type: 'DELETE',
        url: url,
        data: {
            id: id,
        },
        success: function (r) {
            if (r.code === 0) {
                successCallback("删除成功");
            } else {
                errorCallback(r.message);
            }
        },
        error: function (xhr, status, error) {
            errorCallback("请求失败：" + error);
        }
    });
}


/**
 * ajax的Promise方式
 * 
 * @param {*} url      请求URL
 * @param {*} type     请求类型
 * @param {*} param    请求参数
 * @param {*} message  响应信息
 * @returns 
 */
const ajaxPromise = (url, type, param, message) => {
    return new Promise((resolve, reject) => {
        $.ajax({
            type: type,
            url: url,
            data: param,
            success: function (r) {
                if (r.code === 0 && message) {
                    resolve(message);
                } else if (r.code === 0 && !message) {
                    resolve(r);
                } else {
                    reject(r.message);
                }
            },
            error: function (xhr, status, error) {
                reject("请求失败：" + error);
            }
        });
    });
}

const deleteConfirmation = (url, id) => {
    swal({
        title: "删除回复",
        text: "您确定删除回复吗？",
        icon: "warning",
        buttons: ["取消", "确定"],
        dangerMode: true,
    }).then((value) => {
        if (value) {
            deleteReplyFunc(url, id,
                function (message) {
                    swal({ text: message, icon: "success", button: "确定" }).then(
                        function () {
                            location.reload(); // 刷新页面同步回复统计
                        });
                },
                function (message) {
                    swal({ text: message, icon: "warning", button: "确定" });
                }
            );
        } else {
            console.log("删除回复取消了");
        }
    });
}

/**
 * 删除回复
 *
 * @param url 请求URL
 * @param id 回复ID
 * @returns {Promise<void>} 返回异步
 */
const deleteConfirmationPromise = async (url, id) => {
    await newSwal("删除回复", "您确定删除回复吗？", "warning", ["取消", "确定"], true).then(async (value) => {
        if (value) {
            await ajaxPromise(url, "DELETE", { id }, "删除成功")
                .then(async (message) => {
                    await swalSingleBtn("", message, "success", "确定", false).then(() => {
                        // 刷新页面同步回复统计
                        location.reload();
                    });
                })
                .catch((message) => {
                    swal({ text: message, icon: "warning", button: "确定" });
                });
        } else {
            console.log("删除回复取消了");
        }
    });
}

const loadReplyData = async (url, type, param) => {
    return await ajaxPromise(url, type, param, "");
}

/**
 * 采纳回复
 *
 * @param url     请求URL
 * @param id      内容ID
 * @param type    请求类型
 * @param replyId 回复ID
 * @param msg     成功提示信息
 * @returns {Promise<void>}
 */
const adoptReply = async (url, type, id, replyId, msg) => {
    await ajaxPromise(url, type, { id, replyId }, msg)
        .then(async (message) => {
            await swalSingleBtn("", message, "success", "确定", false).then(() => {
                // 刷新页面同步回复统计
                location.reload();
            });
        }).catch((message) => {
            swalSingleBtn("", message, "warning", "确定", false);
        });
}

/**
 * 提交评论
 *
 * @param url     请求URL
 * @param type    请求类型
 * @param param   请求参数
 * @param jBtnEle 提交按钮(jQuery元素)
 * @param message 提示信息
 * @returns {Promise<void>}
 */
const submitReply = async (url, type, param, jBtnEle, message) => {
    await ajaxPromise(url, type, param, message)
        .then(async (message) => {
            jBtnEle.removeAttr('disabled');
            await swalSingleBtn("", message, "success", "确定", false)
                .then(() => {
                    // 刷新页面同步回复统计
                    location.reload();
                });
        }).catch((message) => {
            swalSingleBtn("", message, "warning", "确定", false);
        });
}

/**
 * ajax的表单提交(promise)
 *
 * @param form jQuery表单元素对象
 * @returns {Promise<unknown>}
 */
const ajaxSubmitPromise = (form) => {
    return new Promise((resolve, reject) => {
        form.ajaxSubmit({
            dataType: "json",
            success: function (r) {
                if (r.code === 0) {
                    resolve(r);
                } else {
                    reject(r);
                }
            },
            error: function (xhr, status, error) {
                reject("请求失败：" + error);
            }
        });
    });
}

/**
 * 文章发布表单提交
 *
 * @param jForm       jQuery表单对象
 * @param contentType 内容类型
 * @param title       弹窗标题
 * @returns {Promise<void>}
 */
const submitContentForm = async (jForm, contentType, title) => {
    console.log("内容表单提交参数", contentType, title)
    await ajaxSubmitPromise(jForm).then(async (resp) => {
        await newSwal(title + "成功", resp.message, "success", ["继续" + title, "查看详情"], false).then((val) => {
            if (val) {
                window.location.href = `/${contentType}/${resp.data.contentId}`;
            } else {
                window.location.reload();
            }
        }).catch((resp) => {
            swalSingleBtn("", resp.message, "warning", "确定", false);
        });
    });
}

/**
 * 内容修改表单提交
 *
 * @param jForm       jQuery表单对象
 * @param contentType 内容类型
 * @param contentId   内容ID
 * @param title       弹窗标题
 * @returns {Promise<void>}
 */
const submitContentUpdateForm = async (jForm, contentType, contentId) => {
    console.log("内容修改表单参数", contentType, contentId);
    await ajaxSubmitPromise(jForm).then(async (resp) => {
        await newSwal("操作成功", "内容保存成功", "success", ["继续继续修改", "查看详情"], false).then((val) => {
            if (val) {
                console.log("返回数据", resp.data);
                window.location.href = `/${contentType}/${contentId}`;
            } else {
                window.location.reload();
            }
        }).catch((resp) => {
            swalSingleBtn("", resp.message, "warning", "确定", false);
        });
    });
}

/**
 * 删除内容
 *
 * @param url   请求URL
 * @param title 内容标题
 * @param msg   成功提示信息
 * @param ask   是否询问
 * @returns {Promise<void>}
 */
const deleteReply = async (url, title, msg, ask) => {
    console.log("删除内容参数：", url, title);
    if (ask) {
        await newSwal("提示", `你确定要删除该【${title}】吗？`, "error", ["取消", "确定"], true).then(async res => {
            if (res) {
                await ajaxPromise(url, "delete", undefined, msg)
                    .then(async (message) => {
                        await swalSingleBtn("", message, "success", "确定", false).then(() => {
                            // 刷新页面同步回复统计
                            location.reload();
                        });
                    }).catch(async (message) => {
                        await swalSingleBtn("", message, "warning", "确定", false);
                    });
            } else {
                console.log("取消了操作");
            }
        });

    } else {
        await ajaxPromise(url, "delete", null, msg)
            .then(async (message) => {
                await swalSingleBtn("", message, "success", "确定", false).then(() => {
                    // 刷新页面同步回复统计
                    location.reload();
                });
            }).catch(async (message) => {
                await swalSingleBtn("", message, "warning", "确定", false);
            });
    }
}

/**
 * 富文本编辑器初始化
 *
 * @param height        高度
 * @param placeholder   提示信息
 * @param content       富文本内容
 * @param uploadMaxSize 文件上传最大字节
 * @returns {编辑器}
 */
const vditorInit = (height, placeholder, content, uploadMaxSize) => {
    // 编辑器初始化
    return new Vditor('vditor', {
        cdn: '/plugin/vditor/',
        theme: 'classic',
        height: height,
        icon: 'ant',
        mode: 'wysiwyg',
        cache: {
            enable: false
        },
        placeholder: placeholder,
        value: content,
        toolbar: [
            "emoji",
            "headings",
            "bold",
            "italic",
            "strike",
            "link",
            "|",
            "list",
            "ordered-list",
            "check",
            "outdent",
            "indent",
            "|",
            "quote",
            "line",
            "code",
            "inline-code",
            "insert-before",
            "insert-after",
            "|",
            "upload",
            "table",
            "|",
            "undo",
            "redo",
            "|",
            "fullscreen",
            "edit-mode",
            "outline",
            "preview"
        ],
        upload: {
            accept: 'image/*',     // 附件格式
            url: '/file',          // 上传路径
            linkToImgUrl: '/file', // 粘贴图片上传
            max: uploadMaxSize,    // 20 * 1024 * 1024, // 最大上传文件大小（10MB）
            headers: {
                "X-Requested-With": "XMLHttpRequest"
            },
            filename(name) {
                return name.replace(/[^(a-zA-Z0-9\u4e00-\u9fa5\.)]/g, '').replace(/[\?\\/:|<>\*\[\]\(\)\$%\{\}@~]/g, '').replace('/\\s/g', '')
            },
            // 格式化上传返回
            format(file, response) {
                const { code, data, message } = JSON.parse(response)
                return JSON.stringify({ message, code, data: { errFiles: [], succMap: { "image.png": data.url } } })
            }
        },
        preview: {
            maxWidth: 1920
        },
    });
}