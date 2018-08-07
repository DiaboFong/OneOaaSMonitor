/**
 * Created by sunshine on 2017/7/13.
 */

var api_config = {
    send_sms_code: '/api/send_sms_code',
    apply: "/api/apply"
}

var downloadApp = {

    init: function() {

        // 提交表单
        $('#submit').on('click', function() {
            downloadApp.submit();
        });

        // 提交表单
        $('#sendingVerificationCode').on('click', function() {
            downloadApp.verificationCode();
        });

        // 表单信息验证
        downloadApp.formInfoProof();
    },

    // 表单信息验证
    formInfoProof: function() {
        var date = new Date();
        var nowTime = parseInt(date.valueOf());
        var startTime = parseInt(window.localStorage.startTimes);
        var passTime = parseInt(Math.round(nowTime - startTime) / 1000);
        var remainTime = 120 - passTime;
        if (passTime < 119) {
            $('#sendingVerificationCode').addClass('btn-sent');
            $('#sendingVerificationCode').attr('disabled', 'disabled');
            $('#sendingVerificationCode').text(remainTime + '秒后可重新获取');

            var timeStorage = '';
            var nums = remainTime;

            function sendCode() {
                $('#sendingVerificationCode').text(nums + '秒后可重新获取');
                timeStorage = setInterval(doLoop, 1000); //一秒执行一次
            }

            function doLoop() {
                nums--;
                if (nums > 0) {
                    $('#sendingVerificationCode').text(nums + '秒后可重新获取');
                } else {
                    clearInterval(timeStorage); //清除js定时器
                    $('#sendingVerificationCode').attr('disabled', false);
                    $('#sendingVerificationCode').removeClass('btn-sent');
                    $('#sendingVerificationCode').text('发送验证码');
                    nums = 120; //重置时间
                }
            }
            sendCode();
        }

        $('#companyName').on('change', function() {
            if (!$('#companyName').val()) {
                $('#companyName').next().show();
            } else {
                $('#companyName').next().hide();
            }
        });

        $('#berth').on('change', function() {
            if (!$('#berth').val()) {
                $('#berth').next().show();
            } else {
                $('#berth').next().hide();
            }
        });

        $('#userName').on('change', function() {
            if (!$('#userName').val()) {
                $('#userName').next().show();
            } else {
                $('#userName').next().hide();
            }
        });

        $('#email').on('change', function() {
            var reg = /\w+[@]{1}\w+[.]\w+/;
            var vals = $('#email').val();
            if (vals == "") {
                $('#email').next().show();
                $('#email').next().text('邮箱不能为空！');
            } else if (!reg.test(vals)) {
                $('#email').next().show();
                $('#email').next().text('邮箱格式错误！');
            } else if (vals.indexOf("@qq.") > 0 || vals.indexOf("@sina.") > 0 || vals.indexOf("@163.") > 0 || vals.indexOf("@google.") > 0 || vals.indexOf("@foxmail.") > 0 || vals.indexOf("@gmail.") > 0 || vals.indexOf("@hotmail.") > 0 || vals.indexOf("@126.") > 0) {
                $('#email').next().show();
                $('#email').next().text('邮箱必须为企业邮箱');
            } else {
                $('#email').next().hide();
            }
        });

        $('#phone').on('change', function() {
            var phone = /(1[3-9]\d{9}$)/;
            if (!phone.test($('#phone').val())) {
                $('#phone').parent().find('.error-show').show();
            } else {
                $('#phone').parent().find('.error-show').hide();
            }
        });

        $('#verificationCode').on('change', function() {
            if (!$('#verificationCode').val()) {
                $('#verificationCode').next().show();
            } else if ($('#verificationCode').val().length < 6) {
                $('#verificationCode').next().show();
            } else {
                $('#verificationCode').next().hide();
            }
        });
    },

    // 提交表单
    submit: function() {
        var companyName = $('#companyName').val();
        var berth = $('#berth').val();
        var userName = $('#userName').val();
        var email = $('#email').val();
        var phone = $('#phone').val();
        var verificationCode = $('#verificationCode').val();
        var emailReg = /\w+[@]{1}\w+[.]\w+/;
        var phoneReg = /(1[3-9]\d{9}$)/;

        // 公司验证
        if (!companyName) {
            $('#companyName').next().show();
            return false;
        }

        // 职位验证
        if (!berth) {
            $('#berth').next().show();
            return false;
        }

        // 名称验证
        if (!userName) {
            $('#userName').next().show();
            return false;
        }


        // 邮箱验证
        if (!email) {
            $('#email').next().show();
            return false;
        } else if (!emailReg.test(email)) {
            $('#email').next().show();
            return false;
        } else {
            $('#email').next().hide();
        }

        // 手机号码验证
        if (!phone) {
            $('#phone').parent().find('.error-show').show();
            return false;
        } else if (!phoneReg.test(phone)) {
            $('#phone').parent().find('.error-show').show();
            return false;
        } else {
            $('#phone').parent().find('.error-show').hide();
        }

        // 验证码验证
        if (!verificationCode) {
            $('#verificationCode').next().show();
            return false;
        }

        $.ajax({
            url: api_config.apply,
            type: 'post',
            data: {
                username: userName,
                email: email,
                phone: phone,
                company: companyName,
                work: berth,
                smscode: verificationCode,
                machinecode: $('#machinecode').val()
            },
            success: function(data) {
                if (data.code == 200) {
                    $('.download-success').show();
                    $('.version-content').hide();
                    $('.download-success').text(data.msg);
                    $('.oneoaas-banner-bg2').css('margin-bottom','0px');
                } else {
                    $('.download-form-submit').find('div.error-show').show();
                    $('.download-form-submit').find('div.error-show').text(data.msg);
                }
            }
        })
    },

    // 发送验证码手机号码验证
    verificationCode: function() {
        var phone = /(1[3-9]\d{9}$)/;
        if (!phone.test($('#phone').val())) {
            $('#phone').parent().find('.error-show').show();
        } else {
            $.ajax({
                url: api_config.send_sms_code,
                type: 'post',
                data: { phone: $('#phone').val() },
                success: function(data) {
                    if (data.code == 200) {
                        var date = new Date();
                        window.localStorage.startTimes = date.valueOf();
                        $('#sendingVerificationCode').addClass('btn-sent');
                        $('#sendingVerificationCode').attr('disabled', 'disabled');

                        var timeStorage = '';
                        var nums = 120;

                        function sendCode() {
                            $('#sendingVerificationCode').text(nums + '秒后可重新获取');
                            timeStorage = setInterval(doLoop, 1000); //一秒执行一次
                        }

                        function doLoop() {
                            nums--;
                            if (nums > 0) {
                                $('#sendingVerificationCode').text(nums + '秒后可重新获取');
                            } else {
                                clearInterval(timeStorage); //清除js定时器
                                $('#sendingVerificationCode').attr('disabled', false);
                                $('#sendingVerificationCode').removeClass('btn-sent');
                                $('#sendingVerificationCode').text('发送验证码');
                                nums = 120; //重置时间
                            }
                        }
                        sendCode();
                    } else {
                        $('.download-form-submit').find('div.error-show').show();
                        $('.download-form-submit').find('div.error-show').text(data.msg);
                    }
                }
            })
        }
    },

}

$(document).ready(function() {
    downloadApp.init();
})