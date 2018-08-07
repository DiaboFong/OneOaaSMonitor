/**
 * Created by CW on 2017/9/7.
 */
//user manage
var userManage = {
    vmTable:[],
    init:function () {
        userManage.datatables("#user");
        userManage.checked();
        userManage.licenseAccredit();
        userManage.pageChange();
        userManage.userModal();
    },
    //初始化表格
    datatables:function (tableId) {
        var language = {
            search: '搜索：',
            lengthMenu: "每页显示 _MENU_ 记录",
            zeroRecords: "没找到相应的数据！",
            info: "分页 _PAGE_ / _PAGES_",
            infoEmpty: "暂无数据！",
            infoFiltered: "(从 _MAX_ 条数据中搜索)",
            paginate: {
                first: '<<',
                last: '>>',
                previous: '上一页',
                next: '下一页',
            }
        }
        var userManageRecord = $(tableId).DataTable({
            paging: true, //隐藏分页
            ordering: true, //开启排序
            order: [0,"desc"],
            info: false, //隐藏左下角分页信息
            searching: false, //开启搜索
            pageLength: 10, //每页显示几条数据
            lengthChange: false, //不允许用户改变表格每页显示的记录数
            language: language, //汉化
            sLoadingRecords:true,
            deferRender: true,
            scrollX: true, //水平滚动条
            scrollCollapse: true,
            scroller: true,
            ajax: {
                url: '/users/OneS/data',
                type: 'GET',
            },
            columns: [
                {
                    visible: false,
                    data: 'userid',
                },
                {
                    title: "<input id=\"user_selected\" type=\"checkbox\">", "defaultContent": "<input type=\"checkbox\">",
                    bSortable: false,
                },
                {
                    title : '用户名',
                    data: 'username',
                    width:'149px',
                    bSortable: false,
                },
                // {
                //     title : '密码',
                //     data: 'password2',
                //     width:'150px',
                //     bSortable: false,
                // },
                {
                    title : '邮　箱',
                    data: 'email',
                    width:'200px',
                    bSortable: false,
                },
                {
                    title : '公　司 ',
                    data: 'company',
                    width:'200px',
                    bSortable: false,
                },
                {
                    title : '工　作',
                    data: 'work',
                    width:'150px',
                    bSortable: false,
                },
                {
                    title : '电　话',
                    data: 'phone',
                    width:'150px',
                    bSortable: false,
                },
                {
                    title : '机器码',
                    data: 'machinecode',
                    bSortable: false,
                },
                {
                    title : '许可证',
                    data: 'vendor_num',
                    width: '400px',
                    bSortable: false,
                },
                // {
                //    title: '操作',
                //    data: "id",
                //    width: '150px',
                //    bSortable: false,
                //    "render": function (data, type, row,meta) {
                //       var syncHtml = '<button class="btn btn-xs btn-primary" onclick="">删除</button>';
                //       return syncHtml;
                //    }
                // },
            ],
        });
        this.vmTable = userManageRecord;
        //设置button
        new $.fn.dataTable.Buttons( userManageRecord, {
            buttons: [
                {
                    extend: 'excelHtml5',
                    bom: "utf-8",
                    text: '导出Excel',
                },
            ],
        });
        //将button放置到底部
        var tableContainer = userManageRecord.buttons().container();
        tableContainer.appendTo(
            userManageRecord.table().container()
        );
        return userManageRecord;
    },
    // 选中未授权用户
    checked:function () {
        // 选中当前页未授权行
        $("#user_selected").on('click', function () {
            var trNum = $("#user tbody tr");
            if(!$("#user_selected").is(":checked")){
                $("#user tbody").find("input").prop("checked",false);
                $("#user tbody").find("tr").removeClass("selected");
            }else {
                $("#user tbody").find("input").prop("checked",true);
                $("#user tbody").find("tr").addClass("selected");
            }
            for(var i=0; i<trNum.length; i++){
                var vendorNum = $(trNum[i]).children("td").eq(7).html();
                // 判断是否已经授权
                if( vendorNum !== ""){
                    $("#user tbody tr:nth-child("+(i+1)+")").find("input").prop("checked",false);
                    $("#user tbody tr:nth-child("+(i+1)+")").removeClass("selected");
                }
            }
        });
        // 选中当前行
        $('#user').unbind('#user_selected').on('click','tbody tr', function () {
            var vendorNum = $(this).children("td").eq(7).html();
            if(vendorNum === ""){
                $(this).toggleClass("selected");
                if($(this).hasClass("selected")){
                    $(this).find("input").prop("checked",true);
                }else{
                    $(this).find("input").prop("checked",false);
                    if(!$("#user tbody").find("tr").hasClass("selected")) {
                        $("#user_selected").prop("checked", false);
                    }
                }
            }else{
                $(this).find("input").prop("checked",false);
            }
        });
    },
    // 证书授权
    licenseAccredit:function () {
        $("#license_accredit").on('click',function(){
            var str = "";
            var row_selected = userManage.vmTable.rows('.selected').data('userid');
            if( row_selected.length == 0){
                $(".modal-body p").html("请选择授权用户！");
                $(".modal-body p").css("color","red");
                $("#user_model").parent(".shade").fadeIn(400);
                return
            }
            for(var i=0; i<row_selected.length; i++){
                str += row_selected[i].userid+",";
            }
            $.ajax({
                url: "/users/grant",
                type: "get",
                data:{
                    id: str,
                },
                success: function(data) {
                    if(data.code == 200){
                        $(".modal-body p").html(data.msg);
                        $(".modal-body p").css("color","green");
                        $("#user_model").parent(".shade").fadeIn(400);
                    }else {
                        $(".modal-body p").html(data.msg);
                        $(".modal-body p").css("color","red");
                        $("#user_model").parent(".shade").fadeIn(400);
                    }
                },
                error:function (error) {
                    $(".modal-body p").html(error);
                    $(".modal-body p").css("color","red");
                    $("#user_model").parent(".shade").fadeIn(400);
                }
            })
        })
    },
    // 翻页去除选中项
    pageChange:function () {
        $('#user').on('page.dt',function () {
            var trNum = $("#user tbody tr");
            for(var i=0; i<trNum.length; i++){
                $("#user tbody tr:nth-child("+(i+1)+")").removeClass("selected");
            }
        });
    },
    // 模态框关闭事件
    userModal:function () {
        $("#user_model").on('click','.close',function () {
            $("#user_model").parent(".shade").fadeOut(300);
            $('#user').DataTable().ajax.url("/users/OneS/data").load();
        });
        $("#user_model").on('click','.btn-close',function () {
            $("#user_model").parent(".shade").fadeOut(300);
            $('#user').DataTable().ajax.url("/users/OneS/data").load();
        });
    }
};
$(document).ready(function(){
    userManage.init();
})