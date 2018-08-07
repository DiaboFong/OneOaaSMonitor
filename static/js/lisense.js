/**
 * Created by sunshine on 2017/8/22.
 */
var api_config = {
    license: "/vendor"
}

var lisenseApp ={
    init:function () {
        lisenseApp.submit();
    },

    submit:function () {
        $('#present').on('click',function () {
            var lisense_email = $('#lisense_email').val();
            var lisense_vendor = $('#lisense_vendor').val();
            var lisense_machinecode = $('#lisense_machinecode').val();
            var emailReg = /\w+[@]{1}\w+[.]\w+/;
            var machinecodeReg = /^([A-Za-z0-9]){32,32}$/;
            var vendorReg =/^([A-Za-z0-9_-]){46,50}$/;

            if(!lisense_email) {
                $('#lisense_email').next().show();
            }else if (!emailReg.test(lisense_email)){
                $('#lisense_email').next().show();
                return false;
            } else {
                $('#lisense_email').next().hide();
            }


            if(!lisense_machinecode) {
                $('#lisense_machinecode').next().show();
            } else if(!machinecodeReg.test(lisense_machinecode)){
                $('#lisense_machinecode').next().show();
                return false;
            } else {
                $('#lisense_machinecode').next().hide();
            }

            if(!lisense_vendor) {
                $('#lisense_vendor').next().show();
            }else if (!vendorReg.test(lisense_vendor)){
                $('#lisense_vendor').next().show();
            } else {
                $('#lisense_vendor').next().hide();
            }

            if (!lisense_email =='' && !lisense_vendor == '' && !lisense_machinecode == ''){
                $.ajax({
                    url: api_config.license,
                    type: 'post',
                    data: {
                        email:lisense_email,
                        vendor:lisense_vendor,
                        machinecode:lisense_machinecode
                    },
                    success: function(data) {
                        if(data.code == 200){
                            $('.download-content').hide();
                            $('.version-content').append('<h1 style="color: #fff;text-align: center">申请成功，LicenseKey已发送至您的邮箱，请注意查收！</h1>');
                        }else {
                            $('.download-content').hide();
                            $('.version-content').append('<h1 style="color: #fff;text-align: center">data.msg</h1>');
                        }
                    },error:function (error) {
                        alert(error);
                    }
                })
            }
        });
    },
}

$(document).ready(function() {
    lisenseApp.init();
})