window.onload = function() {
	//제이쿼리 데이터픽커 설정
	$.datepicker.regional['ko']= {
		closeText:'닫기',
		prevText:'이전달',
		nextText:'다음달',
		currentText:'오늘',
		monthNames:['1월(JAN)','2월(FEB)','3월(MAR)','4월(APR)','5월(MAY)','6월(JUM)','7월(JUL)','8월(AUG)','9월(SEP)','10월(OCT)','11월(NOV)','12월(DEC)'],
		monthNamesShort:['1월','2월','3월','4월','5월','6월','7월','8월','9월','10월','11월','12월'],
		dayNames:['일','월','화','수','목','금','토'],
		dayNamesShort:['일','월','화','수','목','금','토'],
		dayNamesMin:['일','월','화','수','목','금','토'],
		weekHeader:'Wk',
		dateFormat:'yy-mm-dd',
		firstDay:0,
		isRTL:false,
		showMonthAfterYear:true,
		showOn: "both",
		buttonImage: "https://image.thebell.co.kr/thebell10/img/icon-cal.png", 
		buttonImageOnly: true, 
		yearSuffix:''
	};
	$.datepicker.setDefaults($.datepicker.regional['ko']);
	
	$("#eDate").datepicker();
	$("#sDate").datepicker();

};


$(document).ready(function(){
	/**************************************
	** 공통
	**************************************/
	/* 달력 */
	/*$("#eDate").val($.datepicker.formatDate($.datepicker.ATOM, new Date()));

	$(".dateBox .today").bind("click",function(){
		$("#eDate").val($.datepicker.formatDate($.datepicker.ATOM, new Date()));
	});*/
	//$("#eDate").val($.datepicker.formatDate($.datepicker.ATOM, new Date()));
	//$(".dateBox .today").bind("click",function(){$("#eDate").val($.datepicker.formatDate($.datepicker.ATOM, new Date()));});


	/* GNB-user */
	$(".userM li .icon").bind("click",function(){
		if($(this).parent().hasClass("on")){
			$(this).parent().removeClass("on");
		}else{
			$(".userM li").removeClass("on");
			$(this).parent().addClass("on");
		}
	});

	/* 전체 메뉴 */
	$(".topMenu .all .icon").bind("click",function(){
		if($(this).hasClass("on")){
			$(this).removeClass("on");
			$(".topMenu .allmenuBox").hide();
		}else{
			$(this).addClass("on");
			$(".topMenu .allmenuBox").show();
		}
	});

	/* 전체메뉴 닫기 */
	$(".topMenu .allmenuBox .closeBox .btn").bind("click",function(){
		$(".topMenu .all .icon").removeClass("on");
		$(".topMenu .allmenuBox").hide();
	});

	/* 셀렉트 박스 */
    var select = $(".selectBox select");
    select.change(function(){
        var select_name = $(this).children("option:selected").text();
        $(this).siblings("label").text(select_name);
    });
	
	/* GNB 검색 */
	$(".gnbWrap .search .close").bind("click",function(){
		$(".gnbWrap .search").removeClass("on");
	});

	/**************************************
	** asideBox
	**************************************/
	/* best clicks */
	 /*bestSlide = 
		$('.bestSlide').bxSlider({
		pagerCustom: '.bestBox .bestP',
		nextSelector: '.bestBox .slider-next',
		prevSelector: '.bestBox .slider-prev',
		onSlideBefore:function(){
			 var current = bestSlide.getCurrentSlide();
			 if(current == 3){
				$(".bestP li").css("left","-110px");
			 }else if(current < 4){
				 $(".bestP li").css("left","0px");
			 }
		 }
	 });
    */
	/* primary issue *
	 primaySlide = 
		 $('.primaySlide').bxSlider({
		 pagerSelector:'.primaryBox .slideP',
		 nextSelector: '.primaryBox .slider-next',
		 prevSelector: '.primaryBox .slider-prev',
		 onSlideBefore:function(){
			 var current = primaySlide.getCurrentSlide();
			 $(".primaryBox .tit li").hide();
			$(".primaryBox .tit li").eq(current).show();
		 }
	 });

	 /* other issue *
	 otherSlide = 
		 $('.otherSlide').bxSlider({
		 pagerSelector:'.otherBox .slideP',
		 nextSelector: '.otherBox .slider-next',
		 prevSelector: '.otherBox .slider-prev',
		 onSlideBefore:function(){
			 var current = otherSlide.getCurrentSlide();
			$(".otherBox .tit li").hide();
			$(".otherBox .tit li").eq(current).show();
		 }
	 });
	 */

	 /**************************************
	 ** 서브 메인
	 **************************************/
	 /* 서브매인 북마크 *
	 $('.bookmarkSlide').bxSlider({
		minSlides: 4,
		maxSlides: 4,
		moveSlides: 4,
		slideWidth: 230,
		slideMargin:20,
		pagerSelector:'.bookMarkBox .slideP',
		nextSelector: '.bookMarkBox .slider-prev',
		prevSelector: '.bookMarkBox .slider-next'
	 });
	*/

	tabList(".storiList",".storiView",1,"mouseover");	

	/* 스토리 리스트 */
	$(".storiList li").bind("mouseover",function(e){
		var thisY = $(this).offset().top - 337;
		$(".storiList .checkArrow").css("top",thisY+"px");
	});


	/**************************************
	** 검색 리스트
	**************************************/
	$(".searchOption .option").bind("click",function(){
		$(this).hide();
		$(".searchOption").addClass("on");
	});
	$(".searchOption .cancel").bind("click",function(){
		$(".searchOption .option").show();
		$(".searchOption").removeClass("on");
	});

	/* 검색 리스트 */
	//tabList(".searchResult .searchResult-tabList",".searchResult .searchResult-tabView",1,"click");

	/* 딜 */
	//tabList(".dilListBox .tabList",".dilListBox .tabView",1,"click");

	
	/**************************************
	** ALL HEADLINE - 올 헤드라인
	**************************************/
	$(".headLineBox .sel").bind("click",function(){
		$(".headLineBox").toggleClass("on");
	});
	headlineChange();
	setInterval(headlineChange,5000);

	/**************************************
	** FREE 최신뉴스
	**************************************/
	$(".closeBox .btn").bind("click",function(){
		//$(".realtimeNews").toggleClass("on");
	});
	headlineChangeFree();
	setInterval(headlineChangeFree,5000);
	
	 $(window).resize(function(){
		 if($(".allH .headClose").hasClass("on") == false){
			 var winH = window.innerHeight - 120;
		 }else{
			 var winH = window.innerHeight - 5;
		 }
		 $(".headLineContent .asideBox").css("height",winH+"px");
	 }).resize();

	 $(".allH .headClose").bind("click",function(){
		 if($(this).hasClass("on") == false){
			$(this).addClass("on");
			$(".allH.headerBox").css({"top":"-114px"});
			$(".headLineContent").css("margin-top","5px");
			$(this).text("메뉴 보이기");
			$(window).resize();
		 }else{
			$(this).removeClass("on");
			$(".allH.headerBox").css({"top":"0px"});
			$(".headLineContent").css("margin-top","119px");
			$(this).text("메뉴 숨기기")
			$(window).resize();
		 }
	 });

	 /**************************************
	 ** 테이블 Q & A
	 **************************************/
	 $(".dfTable.qna .listT").bind("click",function(){
		if($(this).hasClass("on") == true){
			$(this).removeClass("on");
			$(this).next().removeClass("on");
		}else{
			$(".dfTable.qna .viewT").removeClass("on");
			$(".dfTable.qna .listT").removeClass("on");
			$(this).next().addClass("on");
			$(this).addClass("on");
		}
	 });

	 /**************************************
	 ** 팝업 
	 **************************************/
	 $(".popContent .contentView").css("height", window.innerHeight - 82+"px");
	 $(".popContent .close").bind("click",function(){
		window.close();
	 });

	 /**************************************
	 ** 이슈
	 **************************************/
	 tabList(".issueContentBox .iconList",".issueContentBox .issueListView",1,"click");

});
//기사뷰 이미지 레이어팝업창 띄우기- 광고 제외 :not(.ADVIMG)
$(document).on('click','.viewSection img:not(.ADVIMG)',function(e){

	//선행될 작업 img figure 체크하여 태그치환
	e.preventDefault();
	if($(this).parent().is("figure")){
		var temphtml = $(this).parent().clone();
			//temphtml = temphtml.wrap('<a href="#" onclick="return false" class="close"></a>');
		
		$(".window").html("");
		$(".window").html(temphtml);
		$(".image").wrap('<a href="#" onclick="return false" class="close"></a>' );
		$(".gnbBox").removeClass('on');
		$(".bannerBox").hide();
		wrapWindowByMask();


	}else{
		if($(this).is("#newsimg")){
			var temphtml = $(this).parent().parent().parent().parent().clone();
			$(".window").html("");
			$(".window").html(temphtml);
			$(".table_LSize,.table_MSize,.table_SSize,.table_RSize").wrap('<a href="#" onclick="return false" class="close"></a>' );
			$(".window table").removeClass('table_left table_right');
			$(".gnbBox").removeClass('on');
			$(".bannerBox").hide();
			wrapWindowByMask();

		}else{
			//alert("img");
			var temphtml = $(this).clone();
			$(".window").html("");
			$(".window").html(temphtml);
			$(".window img").wrap('<a href="#" onclick="return false" class="close"></a>' );
			$(".gnbBox").removeClass('on');
			$(".bannerBox").hide();
			wrapWindowByMask();
		}
		
	}
	
	//닫기 버튼을 눌렀을 때
	$('.window .close').click(function (e) {  
		//링크 기본동작은 작동하지 않도록 한다.
		//e.preventDefault();  
		$('#mask, .window').hide();  
		$(".gnbBox").not(".allH").addClass('on')
		$(".bannerBox").show(); // 광고 다시 노출

	});       

	//검은 막을 눌렀을 때
	$('#mask').click(function () {  
		$(this).hide();  
		$('.window').hide();  
		$(".gnbBox").not(".allH").addClass('on')
		$(".bannerBox").show(); // 광고 다시 노출
	}); 

}).css("cursor", "pointer");

/**************************************
** 탭 리스트
**************************************/
function tabList(list, view, n, action) {
	var firstN = n-1;
	$(view + " > ul > li").hide();
	$(view + " > ul > li:eq("+firstN+")").show();
	$(list + "> ul > li").removeClass("on");
	$(list + "> ul > li").eq(firstN).addClass("on");


	$(list + "> ul > li").bind(action,function(){
		var maxN = $(list).find("li").length;
		var thisN = $(this).index();
		if (list==".tabNewsSection .tabNewsList4"){
			var todayDate = new Date(); 
			todayDate.setHours( todayDate.getHours() + 1 ); 
			setCookie("tabNewsList4",thisN+1, todayDate);
		}
		if($(this).hasClass("off") == false){

			$(list).find("li").each(function(e){
				if($(this).hasClass("off") == true){
				}else if($(this).hasClass("hide") == true){
				}else{
					$(this).removeClass("on");
				}
			});
			$(this).addClass("on");
			$(view + " > ul > li").hide();
			$(view + "> ul > li:eq("+thisN+")").show();
		}
	});

};

function setCookie(name, value, expirehours){
	var todayDate = new Date(); 
	//todayDate.setHours( todayDate.getHours() + expirehours ); 
	document.cookie = name + "=" + escape( value ) + "; path=/; expires=" + todayDate.toGMTString() + ";" 
}



/**************************************
** 팝업
**************************************/
function popWin(url, w, h, sb) {
 var newWin;
 var setting = "width="+w+", height="+h+", top=5, left=20, scrollbars="+sb;
 newWin = window.open (url, "", setting);
 newWin.focus();
}

/**************************************
** 팝업
**************************************/
function popheadline(url, w, h, sb) {
 var newWin;
 var setting = "width="+w+", height="+h+", top=5, left=20, resizable=1, scrollbars="+sb;
 newWin = window.open (url, "", setting);
 newWin.focus();
}


/**************************************
** 헤드라인
**************************************/
function headlineChange(){
	var li = $(".headLineBox li");
	li.eq(li.length-1).addClass("on");
	var n = $(".headLineBox li.on").index();
	if(n == li.length-1){
		li.removeClass("on");
		li.eq(0).addClass("on")
	}else{
		li.removeClass("on");
		li.eq(n+1).addClass("on");	
	}
	$(".headLineBox .tit a .txt").text($(".headLineBox li.on a").text());
	$(".headLineBox .tit a").attr("href",$(".headLineBox li.on a").attr("href"));
	$(".headLineBox .tit .time").text($(".headLineBox li.on .time").text());
};

/**************************************
** 프리 헤드라인
**************************************/
function headlineChangeFree(){
	var li = $(".realtimeNews li");
	li.eq(li.length-1).addClass("on");
	var n = $(".realtimeNews li.on").index();
	if(n == li.length-1){
		li.removeClass("on");
		li.eq(0).addClass("on")
	}else{
		li.removeClass("on");
		li.eq(n+1).addClass("on");	
	}

	$(".realtimeNews .titBox a").attr("href",$(".realtimeNews li.on a").attr("href"));
	$(".realtimeNews .titBox .tit").text($(".realtimeNews li.on a").text());
	
};

/*
	기사 이미지 레이어팝업처리전 배경 
*/
function wrapWindowByMask(){
	 
	//화면의 높이와 너비를 구한다.
	var maskHeight = $(document).height();  
	var maskWidth = $(window).width();  

	//마스크의 높이와 너비를 화면 것으로 만들어 전체 화면을 채운다.
	$("#mask").css({"width":maskWidth,"height":maskHeight});  

	//애니메이션 효과 - 일단 0초동안 까맣게 됐다가 60% 불투명도로 간다.

	$("#mask").fadeIn(0);      
	$("#mask").fadeTo("slow",0.6);    

	$(".window").show();
	$(".window").each(function () {
		var left = ( $(window).scrollLeft() + ($(window).width() - $(this).width()) / 2 );
		var top = ( $(window).scrollTop() + ($(window).height() - $(this).height()) / 2 );

		if(top<0) top = 0;
		if(left<0) left = 0;

		$(this).css({"left":left, "top":top});
	});


}