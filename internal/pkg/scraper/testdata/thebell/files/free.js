function go(page)
{
	var forms = document.frm;
	forms.page.value = page;
	forms.submit();
}


function serach()
{
	var forms = document.frm;
	forms.submit();
}

function buyProcPdf(mid)
{
	window.open("/event/eventMagazineAgree.asp?mdiv=4&free_flag=Y&fID=" + mid,"", "width=615, height=800,left=300,top=10, toolbar=no, location=no, directories=no, status=no, menubar=no, resizable=no, scrollbars=yes, copyhistory=no")
}
function open_popup(url, wid, hei, scroll, winName) 
{
	var url = url;
	var posi = "width="+ wid +",height="+hei+",toolbar=no,location=no,status=no,menubar=no,top=10,left=50,scrollbars=" + scroll +",resizable=no" ;
	if (winName == "")
	{
		winName = popup;
	}
	var pop = window.open(url,winName,posi);
	return pop;
}
function f_paymentSite(pId, pLevel){
	if(pId == "" || pLevel == "0"){
		if(confirm("더벨의 유료 컨텐츠 입니다.\n이 서비스의 이용을 위해선 별도의 결제가 필요합니다.\n이용정보안내 페이지로 이동하시겠습니까?"))
		{
			this.location.href="/front/free/member/memberChargedInfo.asp";
		}
		else
		{

		}
	}else{
		window.location.href="/front/index.asp"
	}
}

function newsview(key)
{
	top.location.href="/front/free/contents/news/article_view.asp?key="+key;
}

function noticeview(key)
{
	top.location.href="/front/free/community/board_view.asp?board_no=1&id="+key;
}
function mediaview(key)
{
	top.location.href="/front/free/community/board_view.asp?board_no=6&id="+key;
}

function nomember()
{
	if(confirm("더벨의 유료 컨텐츠 입니다.\n이 서비스의 이용을 위해선 별도의 결제가 필요합니다.\n이용정보안내 페이지로 이동하시겠습니까?"))
	{
		this.location.href="/free/company/ChargedInfo.asp?lcode=21";
	}
}

function poplinkwarp02(){	
	popup04.showopen04(560,432);	// 개인정보취급방침
}

function poplinkwarp03(){
	popup03.showopen03(560,482);	// 서비스이용약관
}

function poplinkwarp05(){
	popup05.showopen05(560,432);	// 청소년보호
}

function premium_mainpop()
{
	//window.open('/premium/premium_main.asp', 'premium', 'menubar=no, scrollbars=no, statusbar=no, width=1100, height=800');
	window.open('/plus/index.asp', 'premium', 'location=yes, resizable=yes, menubar=yes, scrollbars=yes, status=yes, titlebar=yes, toolbar=yes, width=1400, height=800');
}

function premium_info()
{
	window.open('/premium/premium_info.asp','premiuminfo','width=800 height=400');
	/*alert("thebell plus이용을 위해선 프리미엄 상품에 가입해주셔야 합니다.\n전환가입안내:02-724-4101");
	this.location.href="/front/free/service/plus.asp";*/
	//window.open('/premium/', 'premium', 'menubar=no, scrollbars=no, statusbar=no, width=1100, height=800');
}

function premium_loginpop()
{
	nomember();
	//window.open('/premium/premium_info.asp','premiuminfo','width=800 height=400');
	//window.open('/premium/', 'premium', 'menubar=no, scrollbars=no, statusbar=no, width=1100, height=800'); 구 안내창
}

function quick()
{
	this.location.href="/front/mypage/myPageKeyword.asp";
}

function scrap()
{
	this.location.href="/front/mypage/myPageScrap.asp?lmenu=Scrap";
}

function mynews()
{
	this.location.href="/front/mypage/myPageNews.asp";
}

function calendar()
{
	this.location.href="/front/datainfo/Calendar_new.asp?code=09";
}

function report()
{
	this.location.href="/front/datainfo/ReportsFeed.asp?code=09";
}

function pop_samsung()
{
	window.open("/front/popup/pop_samsung.html","samsung","width=398,height=600")
}

var originalFontSize = $('#article_main').css('font-size');

function changeTextStyle(type)
{
	
	switch(type)
	{
		case 'plus':
			var currentFontSize = $('#article_main').css('font-size');
			var currentFontSizeNum = parseFloat(currentFontSize, 10);
			var newFontSize = currentFontSizeNum + 1;
			$('#article_main').css('font-size', newFontSize + "px");
			break;

		case 'minus':
			var currentFontSize = $('#article_main').css('font-size');
			var currentFontSizeNum = parseFloat(currentFontSize, 10);
			var newFontSize = currentFontSizeNum - 1;
			$('#article_main').css('font-size', newFontSize + "px");
			break;
	}
}
function go_url(gb,mode,year)
{
	document.location.href="leagueTable.asp?code=09&mode="+mode+"&global="+gb+"&year="+year;
}
function catchKeyPress(evt, id, act){
	evt = evt || window.event;
	var keyCode = (( window.netscape ) ? evt.which : evt.keyCode);
	if(keyCode==13){
		if(act=='click'){
			document.getElementById(id).click();
		} else if(act=='focus'){
			document.getElementById(id).focus();
		} else if(act=='submit'){
			document.getElementById(id).submit();
		}
	}
}
