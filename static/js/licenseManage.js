/**
 * Created by CW on 2017/9/7.
 */
//license manage
var licenseManage = {
    vmTable:[],
    init:function(){
        licenseManage.dataTables("#license");
        licenseManage.licenseApplyFor();
        licenseManage.licensePresent();
        licenseManage.LicenseModal();
    },
    //初始化表格
    dataTables:function (tableId) {
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
        var licenseManageRecord = $(tableId).DataTable({
            paging: true, //隐藏分页
            ordering: true, //开启排序
            order: [0,"desc"],
            info: false, //隐藏左下角分页信息
            searching: false, //开启搜索
            pageLength : 10, //每页显示几条数据
            lengthChange: false, //不允许用户改变表格每页显示的记录数
            language: language, //汉化
            sLoadingRecords:true,
            deferRender: true,
            scrollX: true, //水平滚动条
            scrollCollapse: true,
            scroller: true,
            ajax: {
                url: '/users/list/license',
                type: 'GET',
            },
            columns: [
                {
                    visible: false,
                    data: 'id',
                },
                {
                    title : '申请时间',
                    data: 'apply_date',
                    width: '200px',
                    bSortable: false,
                    "render":function(data, type, row, meta){
                        return (new Date(data)).Format("yyyy-MM-dd hh:mm:ss");
                    }
                },
                {
                    title : 'license 期限(天)',
                    data: 'duration',
                    width: '150px',
                    bSortable: false,
                },
                {
                    title : 'license到期时间 ',
                    data: 'expire_date',
                    width: '200px',
                    bSortable: false,
                    "render":function(data, type, row, meta){
                        return (new Date(data)).Format("yyyy-MM-dd hh:mm:ss");
                    }
                },
                {
                    title : '代理商编号',
                    data: 'user.vendor_num',
                    width: '200px',
                    bSortable: false,
                },
                {
                    title : '机器码',
                    data: 'user.machinecode',
                    width: '200px',
                    bSortable: false,
                },
                {
                    title : '授权码',
                    data: 'license_key',
                    width: '400px',
                    bSortable: false,
                },
                {
                    title : '用户邮箱',
                    data: 'user.email',
                    width: '200px',
                    bSortable: false,
                },
            ],
        });
        licenseManage.vmTable = licenseManageRecord;
        //设置button
        new $.fn.dataTable.Buttons( licenseManageRecord, {
            buttons: [
                {
                    extend: 'excelHtml5',
                    bom: "utf-8",
                    text: '导出Excel',
                    // exportOptions:{"columns":[1,2,3,4,5,6,7,8,9,10] } //导出指定列
                },
            ],
        });
        //将button放置到底部
        var tableContainer = licenseManageRecord.buttons().container();
        tableContainer.appendTo(
            licenseManageRecord.table().container()
        );
        //日期格式化
        Date.prototype.Format = function(fmt) {
            var o = {
                "M+": this.getMonth() + 1,
                "d+": this.getDate(),
                "h+": this.getHours(),
                "m+": this.getMinutes(),
                "s+": this.getSeconds(),
                "q+": Math.floor((this.getMonth() + 3) / 3), //季度
                "S": this.getMilliseconds() //毫秒
            };
            if (/(y+)/.test(fmt)) {
                fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
            }
            for (var k in o) {
                if (new RegExp("(" + k + ")").test(fmt)) {
                    fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
                }
            }
            return fmt;
        } ;
        return licenseManage;
    },
    // 证书申请
    licenseApplyFor: function () {
        $("#license_apply_for").on('click', function () {
            // 如果代理商编号存在，则填充代理商编号
            var lisenseVendor= $("#license tbody tr").children("td").eq(3).html();
            if(lisenseVendor){
                $('#lisense_vendor').val(lisenseVendor);
            }

            $("#apply_for_license").fadeIn(400);
        })
    },
    // 证书申请
    licensePresent: function () {
        $("#license_present").on('click', function () {
            var lisenseEmail = $('#lisense_email').val();
            var lisenseVendor = $('#lisense_vendor').val();
            var lisenseMachinecode = $('#lisense_machinecode').val();
            var emailReg = /\w+[@]{1}\w+[.]\w+/;
            var machinecodeReg = /^([A-Za-z0-9]){32,32}$/;
            var vendorReg =/^([A-Za-z0-9_-]){46,50}$/;

            if(!lisense_email) {
                $('#lisense_email').next().show();
            }else if (!emailReg.test(lisenseEmail)){
                $('#lisense_email').next().show();
                return false;
            } else {
                $('#lisense_email').next().hide();
            }
            
            if(!lisense_machinecode) {
                $('#lisense_machinecode').next().show();
            } else if(!machinecodeReg.test(lisenseMachinecode)){
                $('#lisense_machinecode').next().show();
                return false;
            } else {
                $('#lisense_machinecode').next().hide();
            }

            if(!lisense_vendor) {
                $('#lisense_vendor').next().show();
            }else if (!vendorReg.test(lisenseVendor)){
                $('#lisense_vendor').next().show();
            } else {
                $('#lisense_vendor').next().hide();
            }

            $("#apply_for_license").fadeOut(300);

            $.ajax({
                url: "/vendor",
                type: 'post',
                data: {
                    email: lisenseEmail,
                    vendor: lisenseVendor,
                    machinecode: lisenseMachinecode,
                },
                success: function (data) {
                    if (data.code == 200) {
                        $(".modal-body p").html(data.msg);
                        $(".modal-body p").css("color", "green");
                        $("#license_model").parent(".shade").fadeIn(400);
                        $("#apply_for_license form")[0].reset();
                    } else {
                        $(".modal-body p").html(data.msg);
                        $(".modal-body p").css("color", "red");
                        $("#license_model").parent(".shade").fadeIn(400);
                    }
                },
                error: function (error) {
                    $(".modal-body p").html(error);
                    $(".modal-body p").css("color", "red");
                    $("#license_model").parent(".shade").fadeIn(400);
                }
            })
        })
    },
    // 模态框关闭事件
    LicenseModal:function () {
        $("#license_model").on('click','.close',function () {
            $("#license_model").parent(".shade").fadeOut(300);
            $('#license').DataTable().ajax.url("/users/list/license").load();
        });
        $("#license_model").on('click','.btn-close',function () {
            $("#license_model").parent(".shade").fadeOut(300);
            $('#license').DataTable().ajax.url("/users/list/license").load();
        });
        // 关闭 LicenseKey form 表单
        $("#apply_for_license").on('click','.close',function () {
            $("#apply_for_license").fadeOut(300);
        });
    }
};
$(document).ready(function(){
    licenseManage.init();
});