/**
 * Created by CW on 2017/9/7.
 */
var login = {
    init:function () {
        login.submit();
    },
    submit:function () {
        $("#login_login").on( 'click',function(){
            var loginUsername =  $("#login_username").val();
            var loginPassword =  $("#login_password").val();
            if(!loginUsername){
                $("#login_username").next("span").css("display","block");
                return;
            };
            if(!loginPassword){
                $("#login_username").next("span").css("display","none");
                $("#login_password").next("span").css("display","block");
                $("#login_password").next("span").html("请填写秘密!");
                return;
            };
            $.ajax({
                url:'/api/login',
                type: "post",
                data:{
                    'username':loginUsername,
                    'password':loginPassword,
                },
                success: function(data) {
                    console.log(data);
                    if(data.code == 200){
                        window.location.href = "/license";
                    }else {
                        $("#login_password").next("span").css("display","block");
                        $("#login_password").next("span").html("用户名或密码错误，请重新填写！");
                    }
                },
                error:function (error) {
                    console.log("登录失败！");
                }
            });
        })
    },
}

$(document).ready(function() {
    login.init();
    $("body").keydown(function() {
        if (event.keyCode == "13") {//keyCode=13是回车键
            $("#login_login").click();
        }
    });
})