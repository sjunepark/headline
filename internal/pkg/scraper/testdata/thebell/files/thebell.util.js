$.fn.formatnumber = function()
{
	return this.each(function()
	{
		if (this.value != "")
		{
			var rxSplit = new RegExp('([0-9])([0-9][0-9][0-9][,.])');
			var arrNumber = this.value.split('.');
			
			arrNumber[0] += '.';
			
			do
	    	{
	    	    arrNumber[0] = arrNumber[0].replace(rxSplit, '$1,$2');
	    	}
	    	while (rxSplit.test(arrNumber[0]));
	    	
			if (arrNumber.length > 1) 
	    	{
	    	    this.value = arrNumber.join('');
	    	}
	    	else 
	    	{
	    	    this.value = arrNumber[0].split('.')[0];
	    	}
		}		
	});	
}

function setCookie (name, value, expires)
{
    document.cookie = name + "=" + escape (value) + "; path=/; expires=" + expires.toGMTString();  
}

function getCookie(Name) 
{
    var search = Name + "="
    
    if (document.cookie.length > 0) 
    { 
    	offset = document.cookie.indexOf(search)
        if (offset != -1) 
        { // 쿠키가 존재하면
            offset += search.length
            // set index of beginning of value
            end = document.cookie.indexOf(";", offset)
            // 쿠키 값의 마지막 위치 인덱스 번호 설정
            if (end == -1)
            {
                end = document.cookie.length
            }
            return unescape(document.cookie.substring(offset, end))
        }
    }
    return "";
}

function buyProc(mid)
{
	window.open("/front/event/eventMagazineAgree.asp?mdiv=1&free_flag=Y&fID=" + mid,"", "width=615, height=800,left=300,top=10, toolbar=no, location=no, directories=no, status=no, menubar=no, resizable=no, scrollbars=yes, copyhistory=no")
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
		if(confirm("더벨의 Premium 유료 컨텐츠 입니다.\n이 서비스의 이용을 위해선 별도의 결제가 필요합니다.\n이용정보안내 페이지로 이동하시겠습니까?"))
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
function do_print(newskey) 
{
	window.open('/front/NewsPrint.asp?key=' + newskey, '', 'width=740, height=800, resizable=yes, scrollbars=yes');
}
function noticeview(key)
{
	top.location.href="/front/free/community/board_view.asp?board_no=1&id="+key;
}

function nomember()
{
	if(confirm("더벨의 Premium 유료 컨텐츠 입니다.\n이 서비스의 이용을 위해선 별도의 결제가 필요합니다.\n이용정보안내 페이지로 이동하시겠습니까?"))
	{
		this.location.href="/front/free/member/memberChargedInfo.asp";
	}
	else
	{

	}
}

function poplinkwarp02()
{	
	popup04.showopen04(560,432);	// 개인정보취급방침
}

function poplinkwarp03()
{
	popup03.showopen03(560,482);	// 서비스이용약관
}

function premium_mainpop() //더벨플러스 신규
{	
	//window.open('/premium/premium_main.asp', 'premium', 'menubar=no, scrollbars=no, statusbar=no, width=1100, height=800');
	window.open('/plus/index.asp', 'premium', 'location=yes, resizable=yes, menubar=yes, scrollbars=yes, status=yes, titlebar=yes, toolbar=yes, width=1400, height=800');
}

function premium_oldpop() //더벨플러스 구
{	
	window.open('/premium/premium_main.asp', 'premium', 'menubar=no, scrollbars=no, statusbar=no, width=1100, height=800');
	//window.open('/plus/index.asp', 'premium', 'location=yes, resizable=yes, menubar=yes, scrollbars=yes, status=yes, titlebar=yes, toolbar=yes, width=1400, height=800');
}


function premium_info()
{
	window.open('/premium/premium_info.asp','premiuminfo','width=800 height=400');
	//alert("thebell plus이용을 위해선 프리미엄 상품에 가입해주셔야 합니다.\n전환가입안내:02-724-4101");
    //this.location.href="/front/datainfo/dataService.asp?code=09";
	/*alert("thebell plus이용을 위해선 프리미엄 상품에 가입해주셔야 합니다.\n전환가입안내:02-724-4101");
	this.location.href="/front/free/service/plus.asp";*/
	//window.open('/premium/', 'premium', 'menubar=no, scrollbars=no, statusbar=no, width=1100, height=800');
}

function premium_loginpop()
{
	window.open('/premium/', 'premium', 'menubar=no, scrollbars=no, statusbar=no, width=1100, height=800');
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
			$('.articletxt').css('font-size', newFontSize + "px");
			//$('.articletxt').each(function(idx){
			//	$(this).css('font-size', newFontSize + "px");
			//})
			break;

		case 'minus':
			var currentFontSize = $('#article_main').css('font-size');
			var currentFontSizeNum = parseFloat(currentFontSize, 10);
			var newFontSize = currentFontSizeNum - 1;
			$('#article_main').css('font-size', newFontSize + "px");
			$('.articletxt').css('font-size', newFontSize + "px");
			//$('.articletxt').each(function(idx){
			//	$(this).css('font-size', newFontSize + "px");
			//})
			break;
	}
}
/*
	전자결제시 이용자 유의사항
*/
function go_payinfopop()
{

	window.open("/pg/payinfo.asp","", "width=615, height=615,left=0,top=0, toolbar=no, location=no, directories=no, status=no, menubar=no, resizable=no, scrollbars=yes, copyhistory=no")
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

function win_YK(filename,winhow) 
{
	var params = filename.split("&");
	var name = params[0].split("=")[1];
	var path = params[1].split("=")[1];
	
	var image = new Image();
	image.onload = function()
	{
		var width  = image.width ;
		var height = image.height;
		
		if (width == 0) 
			width = 600;
		else
			width += 20;
		if (height == 0) 
			height = 480
		else
			height += 20;
		
		var window_option = "";
		if (width > 1000 || height > 700)
		{
			window_option = ",scrollbars=yes";
		}
		var photo_win = window.open(filename + "&width:"+width, "WIN_SUB", "width=" + width + ",height=" + height + window_option);
		photo_win.focus();
	};
	image.src = "https://image.thebell.co.kr/news/photo/" + path.substring(0, 4) + "/" + path.substring(4, 6) + "/" + path.substring(6, 8) + "/" + name;
}
function setCode() 
{
	var mypagecode_win = window.open("/member/MyNewspop.asp","", "width=800, height=539,left=300,top=100, toolbar=no, location=no, directories=no, status=no, menubar=no, resizable=no,		scrollbars=no, copyhistory=no");
	mypagecode_win.focus();
}


jQuery.extend({

	// Validation 관련 함수--------------------------------------
	isEmpty: function(str){
		if(str==null || $.trim(str) =="" ) return true;
		return false;
	},		

	isEmptyArray: function(arr){
		var i;
		 for (i in arr) {
				return false;
		}
		return true;
	},	
	
	isEmail: function(str){
		var check1 = /(@.*@)|(\.\.)|(@\.)|(\.@)|(^\.)/; 
		var check2 = /^[a-zA-Z0-9\-\.\_]+\@[a-zA-Z0-9\-\.]+\.([a-zA-Z]{2,4})$/;
		if ( !check1.test(str) && check2.test(str) ) { 
			return true;
		} else {
			return false;
		}
	},
	
	isEmailList: function(str){		
		var strList, arrList, strEmail, idx;
		if($.isEmpty(str)) return false;
		else{
			strList = str;
			strList = $.replaceStr(strList,";",",");			
			
			arrList = strList.split(",");
			
			for(var i=0; i< arrList.length; i++)
			{
				strEmail = arrList[i];	
				strEmail = $.trim(strEmail);	
				
				idx = strEmail.indexOf("<");
				if(idx>=0) strEmail = strEmail.substring(idx);
				
				idx = strEmail.indexOf(" ");
				if(idx>=0) strEmail = strEmail.substring(idx);
				
				strEmail = $.replaceStr(strEmail," ","");
				strEmail = $.replaceStr(strEmail,"<","");
				strEmail = $.replaceStr(strEmail,">","");		
				if(!$.isEmail(strEmail))  return false;
			}
			return true;	
		}
	},	
	
	isNumeric: function(str){
		if ($.isEmpty(str)) return false;
		if( isNaN(str*1) ) return false;
		else return true;
	},
	
	isDate: function(yyyymmdd){
		var yyyy,mm,dd;
	
		if($.getByte(yyyymmdd)!=8) return false;
		yyyy = yyyymmdd.substring(0,4);
		mm   = yyyymmdd.substring(4,6);
		dd   = yyyymmdd.substring(6,8);
		var m = parseInt(mm,10) - 1;
		var d = parseInt(dd,10);
	
		var end = new Array(31,28,31,30,31,30,31,31,30,31,30,31);
		if ((yyyy % 4 == 0 && yyyy % 100 != 0) || yyyy % 400 == 0) {
			end[1] = 29;
		}
	
		return (d >= 1 && d <= end[m]);
	},

	isKorean: function(str){
		var retValue = true;
	
		if (  str==null || str ==""  ) 	return false;
		else {
			for( var i = 0; i < str.length; i++ ) {
				var chr = escape(str.charAt(i));					//입력된 값의 하나하나를 아스키(ASCII) 값으로 변환시킨 후...
	
				if ( chr.length == 1 ) return false;					//영문의 경우 아스키값이 1자리니까...
				else if ( chr.indexOf("%u")	 != -1 )  continue;
				else if ( chr.indexOf("%") != -1 ) return false;		//"~"와 같은 특수문자의 경우 아스키값이 3자리니까...
			}
		}
		return true;
	},
		
	isAlphabet: function(str){
		var temp1;
		var len = str.length;
		if (  str==null || str ==""  ) 	return false;
		
		for(l=0;l<len;l++){
		     temp1 = str.charAt(l);
		     if (escape(temp1).length >= 4) return false;
	             if ( (temp1<'a' || temp1 > 'z') && (temp1 <'A' || temp1 >'Z') ) return false;
	       }
	     return true;
	},	
	
	isResNo: function(str){
		if ($.isEmpty(str)) return false;
		
		var str_serial1=str.substring(0,6)
		var str_serial2=str.substring(6,13)
	
		if (str_serial1.length != 6){
			return false;
		}else if (str_serial2.length != 7){
			return false;
		}else if(isNaN(str_serial1) || isNaN(str_serial2)){
			return false;
		}else{
			var a1=str_serial1.substring(0,1)
			var a2=str_serial1.substring(1,2)
			var a3=str_serial1.substring(2,3)
			var a4=str_serial1.substring(3,4)
			var a5=str_serial1.substring(4,5)
			var a6=str_serial1.substring(5,6)
	
			var check_digit=a1*2+a2*3+a3*4+a4*5+a5*6+a6*7
	
			var b1=str_serial2.substring(0,1)
			var b2=str_serial2.substring(1,2)
			var b3=str_serial2.substring(2,3)
			var b4=str_serial2.substring(3,4)
			var b5=str_serial2.substring(4,5)
			var b6=str_serial2.substring(5,6)
			var b7=str_serial2.substring(6,7)
	
			var check_digit=check_digit+b1*8+b2*9+b3*2+b4*3+b5*4+b6*5
	
			check_digit = check_digit%11
			check_digit = 11 - check_digit
			check_digit = check_digit%10
	
			if (str_serial1.substring(2,3) > 1){
				return false;
			}else	if (str_serial1.substring(4,5) > 3){
				return false;
			}else	if (str_serial2.substring(0,1) > 4 || str_serial2.substring(0,1) == 0){
				return false;
			}else	if (check_digit != b7){
				return false;
			}else{
				return true;
			}
		}
	},
	
	isPhone:function(str) {
	
	   var RetValue = true;
	   var Count;
	   var PermitChar =
	         "0123456789-";
	
	   for (var i = 0; i < str.length; i++) {
	      Count = 0;
	      for (var j = 0; j < PermitChar.length; j++) {
	         if(str.charAt(i) == PermitChar.charAt(j)) {
	            Count++;
	            break;
	         }
	      }
	
	      if (Count == 0) {
	         RetValue = false;
	         break;
	      }
	   }
	   return RetValue;
	},

	// 숫자 관련 함수--------------------------------------
	formatNumber: function(str){
    	var txtNumber = '' + str;
    	var rxSplit = new RegExp('([0-9])([0-9][0-9][0-9][,.])');
    	var arrNumber = txtNumber.split('.');
    	arrNumber[0] += '.';
    	do
    	{
    	    arrNumber[0] = arrNumber[0].replace(rxSplit, '$1,$2');
    	}
    	while (rxSplit.test(arrNumber[0]));
    	
    	if (arrNumber.length > 1) 
    	{
    	    return arrNumber.join('');
    	}
    	else 
    	{
    	    return arrNumber[0].split('.')[0];
    	}
	},
	
	// 휴대폰 format
	mobileNumber: function(str){
    	var txtNumber = '' + str;
    	var rxSplit = new RegExp('([0-9])([0-9][0-9][0-9][,.])');
    	var arrNumber = txtNumber.split('.');
    	arrNumber[0] += '.';
    	do
    	{
    	    arrNumber[0] = arrNumber[0].replace(rxSplit, '$1,$2');
    	}
    	while (rxSplit.test(arrNumber[0]));
    	
    	if (arrNumber.length > 1) 
    	{
    	    return arrNumber.join('');
    	}
    	else 
    	{
    	    return arrNumber[0].split('.')[0];
    	}
	},	
	
	// 스트링 관련 함수--------------------------------------
	getByte: function(str){
		var byteSize = 0;
		var retValue = "";
	
		if ( str==null || str ==""  ) {
			return 0;
		} else {
			for( var i = 0; i < str.length; i++ ) {
				var chr = escape(str.charAt(i));		//입력된 값의 하나 하나를 아스키(ASCII) 값으로 변환시킨 후...
	
				if ( chr.length == 1 ) {					//영문의 경우 아스키값이 1자리니까...
					byteSize ++;
					retValue += str.charAt(i);
				} else if ( chr.indexOf("%u") != -1 ) {		//한글의 경우"%"u로 시작하니까...
					byteSize += 2;
				} else if ( chr.indexOf("%") != -1 ) {		//"~"와 같은 특수문자의 경우 아스키값이 3자리니까...
					byteSize += chr.length/3;
					retValue += str.charAt(i);
				}
			}
		}
		return byteSize;
	},	
	
	replaceStr: function(str, strSrc, strDest){
		var strTmp  = str;
	   	str = "";
	
		while ( strTmp.indexOf(strSrc) >= 0 ) {
			str +=  strTmp.substring(0,strTmp.indexOf(strSrc))+strDest;
			strTmp = strTmp.substring(strTmp.indexOf(strSrc)+strSrc.length);
		}
		str = str + strTmp;
		return str;
	},	
	
	addStr: function(str, start, strAdd){	
   		str = str.substring(0,start)+strAdd+str.substring(start);	
   		return str;
	},


	indexOf: function(str, charFor){
		for (i=0; i < str.length; i++)
		{
			if (charFor == $.mid(str, i, charFor.length)) {return i; break;}
		}
		return -1;
	},
	
	lastIndexOf: function(str, charFor){
		for (i = str.length-1; i>=0; i--)
		{
			if (charFor == $.mid(str, i, charFor.length)){return i;break;}
		}
		return -1;
	},
	
	mid: function(str, start, len){
		if (start < 0 || len < 0) return "";
		var iEnd, iLen = String(str).length;
		if (start + len > iLen)
			iEnd = iLen;
		else
			iEnd = start + len;
		return String(str).substring(start,iEnd);
	},
		
	// 이미지 관련 함수-----------------------------------
	// 사용예<img src="/images/logo.gif" onload="javascript:imgResize(this, true, 400);" onerror="javascript:imgResize(this, false, 400);">
	imgResize: function(imgObj, bDisplay, maxWidth)
	{
		if(maxWidth==null || $.trim(maxWidth) =="" ) maxWidth= 500;
		
		if(bDisplay)
			if(imgObj.width > maxWidth)	imgObj.width = maxWidth;
		else
			imgObj.style.display = 'none'; //** 안보이게 스타일 시트로 처리
	},	
	
	// flash(파일주소, 가로, 세로, 배경색, 윈도우모드, 변수, 경로)
	flash:function(URL,w,h,bg,win,vars,base)
	{
		var s=
		"<object classid='clsid:d27cdb6e-ae6d-11cf-96b8-444553540000' codebase='http://fpdownload.macromedia.com/pub/shockwave/cabs/flash/swflash.cab#version=10,0,0,0' width='"+w+"' height='"+h+"' align='middle'>"+
		"<param name='allowScriptAccess' value='always' />"+
		"<param name='movie' value='"+url+"' />"+
		"<param name='wmode' value='"+win+"' />"+
		"<param name='menu' value='false' />"+
		"<param name='quality' value='high' />"+
		"<param name='FlashVars' value='"+vars+"' />"+
		"<param name='bgcolor' value='"+bg+"' />"+
		"<param name='base' value='"+base+"' />"+
		"<embed src='"+url+"' base='"+base+"' wmode='"+win+"' menu='false' quality='high' bgcolor='"+bg+"' width='"+w+"' height='"+h+"' align='middle' type='application/x-shockwave-flash' pluginspage='http://www.macromedia.com/go/getflashplayer' />"+
		"</object>";
		document.write(s);
	},

	
	
	// 날짜관련 함수--------------------------------------
	getToDay: function(tocken){
		if(tocken==null) tocken="";
		var toDt  = new Date();
		var yearStr  = toDt.getFullYear();
		var monthStr = toDt.getMonth()+1;
		var dateStr  = toDt.getDate();
		if (monthStr<10)  monthStr = "0"+monthStr;
		if (dateStr<10)  dateStr= "0"+dateStr;
		return 	yearStr +tocken+ monthStr+tocken + dateStr;
	},
	
	addYearStr: function(yyyymmdd,val){
		var tarDay = new Date();
		tarDay.setYear(yyyymmdd.substring(0,4)*1+val);
		tarDay.setMonth(yyyymmdd.substring(4,6)*1-1);
		tarDay.setDate(yyyymmdd.substring(6,8));
		var yearStr  = tarDay.getFullYear();
		var monthStr = tarDay.getMonth()+1;
		var dateStr  = tarDay.getDate();
		if (monthStr<10)  monthStr = "0"+monthStr;
		if (dateStr<10)  dateStr= "0"+dateStr;
		return 	yearStr +""+ monthStr+"" + dateStr;
	},
	
	addMonStr: function(yyyymmdd,val){
		var tarDay = new Date();
		tarDay.setYear(yyyymmdd.substring(0,4)*1);
		tarDay.setMonth(yyyymmdd.substring(4,6)*1-1+val);
		tarDay.setDate(yyyymmdd.substring(6,8));
		var yearStr  = tarDay.getFullYear();
		var monthStr = tarDay.getMonth()+1;
		var dateStr  = tarDay.getDate();
		if (monthStr<10)  monthStr = "0"+monthStr;
		if (dateStr<10)  dateStr= "0"+dateStr;
		return 	yearStr +""+ monthStr+"" + dateStr;
	},
	
	addDateStr: function(yyyymmdd,val){
		var tarDay = new Date();
		tarDay.setYear(yyyymmdd.substring(0,4)*1);
		tarDay.setMonth(yyyymmdd.substring(4,6)*1-1);
		tarDay.setDate(yyyymmdd.substring(6,8)*1+val);
		var yearStr  = tarDay.getFullYear();
		var monthStr = tarDay.getMonth()+1;
		var dateStr  = tarDay.getDate();
		if (monthStr<10)  monthStr = "0"+monthStr;
		if (dateStr<10)  dateStr= "0"+dateStr;
		return 	yearStr +""+ monthStr+"" + dateStr;
	},	
	
	duringMonth: function(f_yyyymmdd, t_yyymmdd){
		var ret_Val= 0;
		var monNum = 1;
		var v_fromDt = f_yyyymmdd;
		if (v_fromDt>t_yyymmdd) return -1;
		if (v_fromDt==t_yyymmdd) return 0;
	

		while( v_fromDt<t_yyymmdd)
		{
			ret_Val++;								
			v_fromDt=$.addMonStr(v_fromDt,1);
		}
		return ret_Val;
	},
	
	duringDate: function(f_yyyymmdd, t_yyymmdd){
		var ret_Val=0;
		var monNum = 1;
		var v_fromDt = f_yyyymmdd;
				
		if (v_fromDt>t_yyymmdd) return -1;
		if (v_fromDt==t_yyymmdd) return 0;

		while( v_fromDt<t_yyymmdd)
		{
			ret_Val++;								
			v_fromDt=$.addDateStr(v_fromDt,1);
		}
		return ret_Val;
	},	
	
	// 현재 일자부터 yyyymmdd 까지 일수
	duringDays: function(yyyymmdd) {
		toDay = new Date();
		tarDay = new Date();
		tarDay.setYear(yyyymmdd.substring(0,4));
		tarDay.setMonth(yyyymmdd.substring(4,6)*1-1);
		tarDay.setDate(yyyymmdd.substring(6,8));
		betweenDate = parseInt((tarDay-toDay)/86400000);
		return betweenDate;
	},	
	// Byte 짜르기
	byteCutLength:function (str, len){
		var byteSize = 0;
	
		if ( str==null || str ==""  ) {
			return str;
		} else {
			for( var i = 0; i < str.length; i++ ) 
			{
				var chr = escape(str.charAt(i));		//입력된 값의 하나 하나를 아스키(ASCII) 값으로 변환시킨 후...
	
				if ( chr.length == 1 ) {					//영문의 경우 아스키값이 1자리니까...
					byteSize ++;
				} else if ( chr.indexOf("%u") != -1 ) {		//한글의 경우"%"u로 시작하니까...
					byteSize += 2;
				} else if ( chr.indexOf("%") != -1 ) {		//"~"와 같은 특수문자의 경우 아스키값이 3자리니까...
					byteSize += chr.length/3;
				}
				if (byteSize > len) return str.substring(0,i) + "...";
			}
		}
		return str;
	},	
	// 변환 관련 함수--------------------------------------
	urlEncode: function(str){
	  var lsRegExp = /\+/g;
	  return escape(String(str).replace(lsRegExp, " ")); 
	},
	
	urlDecode: function(str){
	  var lsRegExp = /\+/g;
	  return unescape(String(str).replace(lsRegExp, " ")); 
	},
	
	getAjaxRequest:function(uri, divLayer, param) 
	{
		var loading = "<div id='loading_"+ divLayer +"' style='position:relative;display:none;z-index:9999;'><img src=\"https://image.thebell.co.kr/m/loading.gif\" /></div>";		
		var targetTop= $('#'+divLayer).offset().top+ (($('#'+divLayer).height()==0)? 40 : $('#'+divLayer).height()/2);
	
	 	var h=32;
	 	var w=32;
		$border = $("<div/>").attr({"id":'loading_'+ divLayer}).css(
		{
			"position"	  : "absolute",
			"border"	  : "0px solid #000000",			
			"z-index"	  : "9000",
			"left"		  : "50%", 
			"top"		  : targetTop, 
			"margin-top"  : ((h + 2) / 2 * -1) + "px",
			"margin-left" : ((w + 2) / 2 * -1) + "px",
			"width"		  : (w) + "px", 
			"height"	  : (h) + "px"
		});
		$("#"+divLayer).append($border);
		$img = $("<img/>").attr({"src": "https://image.thebell.co.kr/m/loading.gif"}).css({"width":(w) + "px","height":(h) + "px"});
		$border.append($img);

		$.ajax({
		    url 	  : uri,
		    type  	  : "post",
		    data 	  : param,
		    dataType  : "html",
		    async     : "false",
		    error	  : function(xhr, textStatus) {	 if(loading) {$('#loading_'+ divLayer).remove();}},
		    complete  : function() { if(loading) {$('#loading_'+ divLayer).remove();} },	  	    
		    success   : function(HTML) {	

		    	$("#"+divLayer).each(function(){
		    		$(this).find('script').remove();
		    		$(this).find('style').remove();
		    		$(this).find("#bodyLayer").remove();
		    		$(this).empty();
		    	});
				$bodyLayer = $("<div/>").attr({"id":'bodyLayer'}).css(
				{
					"width"	  : "100%",
					"height"  : "100%"
				});
				$("#"+divLayer).html($bodyLayer);	
							
				$bodyLayer.html(HTML);
			}
		});
	}

})
$(document).on("keyup",".numclass",function(){
	$(this).val($(this).val().replace(/[^0-9]/gi,""));
});