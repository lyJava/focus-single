// 全局管理对象
gf = {
    // 刷新验证码
    reloadCaptcha: function() {
        $("img.captcha").attr("src","/captcha?v="+Math.random());
    },
}

// 用户模块
gf.user = {
    // 退出登录
    logout: function () {
        swal({
            title:   "注销登录",
            text:    "您确定需要注销当前登录状态吗？",
            icon:    "warning",
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
    delete: function (id,url) {
        url = url || "/"
        swal({
            title:   "删除内容",
            text:    "您确定要删除该内容吗？",
            icon:    "warning",
            buttons: ["取消", "确定"]
        }).then((value) => {
            if (value) {
                jQuery.ajax({
                    type: 'DELETE',
                    url : '/content/delete',
                    data: {
                        id: id
                    },
                    sync: true,
                    success: function (data) {
                        swal({
                            title:   "删除完成",
                            text:    "3秒后自动跳转到",
                            icon:    "success",
                            timer:   2000,
                            buttons: false
                        }).then((value) => {
                            window.location.href = url;
                        })
                    }
                });
            }
        });
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
            url : '/interact/zan',
            data: {
                id:   id,
                type: type
            },
            sync: true,
            success: function (r, status) {
                if (r.code <= 0) {
                    let number = $(elem).find('.number').html()
                    $(elem).find('.number').html(parseInt(number)+1)
                    $(elem).find('.icon').removeClass('icon-zan').addClass('icon-zan-done')
                } else {
                    swal({
                        text:   r.message,
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
            url:  '/interact/zan',
            data: {
                id:   id,
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
        var type = $(elem).attr("data-type")
        if ($(elem).find('.icon').hasClass('icon-cai-done')) {
            this.cancelCai(elem, type, id)
        } else {
            this.cai(elem, type, id)
        }
    },
    // 踩
    cai: function (elem, type, id) {
        jQuery.ajax({
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
        });
    },
    // 取消踩
    cancelCai: function (elem, type, id) {
        jQuery.ajax({
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
        });
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
        required:    "这是必填字段",
        remote:      "请修正此字段",
        email:       "请输入有效的电子邮件地址",
        url:         "请输入有效的网址",
        date:        "请输入有效的日期",
        dateISO:     "请输入有效的日期 (YYYY-MM-DD)",
        number:      "请输入有效的数字",
        digits:      "只能输入数字",
        creditcard:  "请输入有效的信用卡号码",
        equalTo:     "你的输入不相同",
        extension:   "请输入有效的后缀",
        maxlength:   $.validator.format("最多可以输入 {0} 个字符"),
        minlength:   $.validator.format("最少要输入 {0} 个字符"),
        rangelength: $.validator.format("请输入长度在 {0} 到 {1} 之间的字符串"),
        range:       $.validator.format("请输入范围在 {0} 到 {1} 之间的数值"),
        max:         $.validator.format("请输入不大于 {0} 的数值"),
        min:         $.validator.format("请输入不小于 {0} 的数值")
    });
})

/**
 * 创建弹窗(异步)
 *
 * @param title
 * @param text
 * @param icon
 * @param dangerMode
 * @returns {Promise<unknown>}
 */
function newSwal(title, text, icon, dangerMode) {
    return new Promise((resolve, reject) => {
        swal({
            title: title,
            text: text,
            icon: icon,
            buttons: ["取消", "确定"],
            dangerMode: dangerMode
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
function swalSingleBtn(title, text, icon, button, dangerMode) {
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

function deleteReplyFunc(url, id, successCallback, errorCallback) {
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

function ajaxPromise(url, type, param, message) {
    return new Promise((resolve, reject) => {
        jQuery.ajax({
            type: type,
            url: url,
            data: param,
            success: function (r) {
                if (r.code === 0) {
                    resolve(message);
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

function deleteConfirmation(url, id) {
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
                    swal({text: message, icon: "success", button: "确定"}).then(
                        function () {
                            location.reload(); // 刷新页面同步回复统计
                        });
                },
                function (message) {
                    swal({text: message, icon: "warning", button: "确定"});
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
async function deleteConfirmationPromise(url, id) {
    await newSwal("删除回复", "您确定删除回复吗？", "warning", true).then(async (value) => {
        if (value) {
            await ajaxPromise(url, "DELETE", {id}, "删除成功")
                .then(async (message) => {
                    await swalSingleBtn("", message, "success", "确定", false).then(() => {
                        // 刷新页面同步回复统计
                        location.reload();
                    });
                })
                .catch((message) => {
                    swal({text: message, icon: "warning", button: "确定"});
                });
        } else {
            console.log("删除回复取消了");
        }
    });
}

/**
 * 采纳回复
 *
 * @param url     请求URL
 * @param id      内容ID
 * @param replyId 回复ID
 * @param msg     成功提示信息
 * @returns {Promise<void>}
 */
async function adoptReply(url, type, id, replyId, msg) {
    await ajaxPromise(url, type, {id, replyId}, msg)
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
async function submitReply(url, type, param, jBtnEle, message) {
    await ajaxPromise(url, type, param, message)
        .then(async (message) => {
            jBtnEle.removeAttr('disabled');
            await swalSingleBtn("", message, "success", "确定", false)
                .then(() => {
                    // 刷新页面同步回复统计
                    location.reload();
                });
        }).catch((message) => {
            swal({text: message, icon: "warning", button: "确定"});
        });
}