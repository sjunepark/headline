/**
 * IIFE(Immediately Invoked Function Expression)
 */
(function(thebell, $, undefined)
{
	thebell.detector = 
	{
		ui : "",
		ip : "",
		pa : "",
		wb : "",
		os : "",
		ua : "",
		send : function()
		{
			try
			{
				$.ajax(
				{
					
					type     : "POST",
					dataType : "jsonp",
					url      : "https://www.thebell.co.kr/log/log.asp",
					data     : {"ui":this.ui, "pa":this.pa , "wb":this.wb, "os":this.os}, 
					success   :  function(response)
	    			{
						//alert("zzz");
					}
				});
			}
			catch(e) {}
		},
		call : function(x, y)
		{
			var self = this;
			var $offset = $("iframe[name='frame_tabmenu']").offset();
			x = x + $offset.left + ($(window).width() >= 1050 ? ($(window).width() - 1050) / 2 * -1 : 0);
			y = y + $offset.top;
			self.send(x, y);
		},
		init : function()
		{
		
			var self = this;
			self.dr = document.referrer; //url
			self.pa = location.pathname;
			self.pr = document.URL.substring(document.URL.indexOf(self.pa) + self.pa.length + 1);
			self.ui = "";

			function getArgs(){
				var args   = new Object();
				var query = location.search.substring(1)   //? 제거
				var pairs = query.split("&");              //& 
				
				for (var i = 0; i < pairs.length; i++) {
					var pos2 = pairs[i].indexOf('=');          // "name = value" 을 찾는다.
					if (pos2 == -1) continue;                 
					var argname = pairs[i].substring(0, pos2); 
					var value = pairs[i].substring(pos2+1);     
					value = decodeURIComponent(escape(value)); // UTF-8 DECODING ENCODING => unescape()         
					//value = decodeURIComponent(value);          
					args[argname] = value;                      
				 }
				 return args;
			}
			var args  = getArgs(); //URL에서 args를 파싱
			self.co   = args.code || ''; //정의된 전달인자가 있으면 사용하고 없으면 기본값 사용
			self.al   = args.all || '';
	
			/**
			 * 브라우저 종류 및 운영체제 정보 확인
			 */
			self.ua = navigator.userAgent.toUpperCase();
			self.ua2 = navigator.appName;

			if (self.ua.indexOf("MSIE 10") != -1)
			{
				self.wb = "IE10";	
			}
			else if (self.ua.indexOf("MSIE 9") != -1)
			{
				self.wb = "IE9";	
			}
			else if (self.ua.indexOf("MSIE 8") != -1)
			{
				self.wb = "IE8";
			}
			else if (self.ua.indexOf("MSIE 7") != -1)
			{
				self.wb = "IE7";
			}
			else if (self.ua.indexOf("MSIE 6") != -1)
			{
				self.wb = "IE6";
			}
			else if (self.ua.indexOf("MSIE") != -1)
			{
				self.wb = "";
			}
			else if (self.ua.indexOf("EDG") != -1)
			{
				self.wb = "EG";
			}
			else if (self.ua.indexOf("FIREFOX") != -1)
			{
				self.wb = "FI";
			}
			else if (self.ua.indexOf("WHALE") != -1)
			{
				self.wb = "WH";
			}			
			else if (self.ua.indexOf("CHROME") != -1)
			{
				self.wb = "CH";
			}
			else if (self.ua.indexOf("OPERA") != -1)
			{
				self.wb = "OP";
			}
			else if (self.ua.indexOf("SAFARI") != -1)
			{
				self.wb = "SA";
			}
			else if (self.ua.indexOf("MAC") != -1)
			{
				self.wb = "MA";
			}
			else if(self.ua2 == 'Netscape' && self.ua.search('TRIDENT') != -1)
			{		
				self.wb = "IE11";
			}else{
				self.wb = "";
			}
			
			if (self.ua.indexOf("WINDOWS") != -1)
			{
				self.os = "WI";		
			}
			else if (self.ua.indexOf("IPHONE") != -1)
			{
				self.os = "IP";
			}
			else if (self.ua.indexOf("IPOD") != -1)
			{
				self.os = "ID";				
			}
			else if (self.ua.indexOf("ANDROID") != -1)
			{
				self.os = "AN";				
			}
			else if (self.ua.indexOf("MAC OS") != -1)
			{
				self.os = "MA";
			}
			else if (self.ua.indexOf("BLACKBERRY") != -1)
			{
				self.os = "BL";
			}
			else if (self.ua.indexOf("LINUX") != -1)
			{
				self.os = "LI";
			}
			else
			{
				self.os = "";
			}
			
			if (_uuii != undefined)
			{
				for (i = 0;i < _uuii.length;i++)
				{
					var u = _uuii.pop(i);
					if (u[0] == "account")
						self.ui = u[1];
				}
			}
			if (_uuip != undefined)
			{
				for (i = 0;i < _uuip.length;i++)
				{
					var u = _uuip.pop(i);
					if (u[0] == "ip")
						self.ip = u[1];
				}
			}
			self.send();

			/*
			if (self.pa == "/front/index.asp")
			{
				$(document).bind("mousedown", function(e)
				{
					if (e.button == 0)
					{
						var x, y;
						if ($(window).width() >= 1050)
						{
							x = e.pageX - ($(window).width() - 1050) / 2;
						}
						else
						{
							x = e.pageX;
						}
						y = e.pageY;
						
						if (x > 0)
						{
							self.send(x, y);
						}
					}
				});
			}
			else if (self.pa == "/front/tabmenu.asp")
			{
				$(document).bind("mousedown", function(e)
				{
					if (e.button == 0)
					{
						try
						{
							window.parent.thebell.detector.call(e.pageX, e.pageY);
						}
						catch (e) {}
					}
				});
			}
			else if (self.pa == "/front/datainfo/goodmorning_popup_new.asp")
			{
				$(document).bind("mousedown", function(e)
				{
					if (e.button == 0)
					{
						self.send(e.pageX, e.pageY);
					}
				});			
			}
			*/
			/*
			기타메뉴
			*/
			/*
			if (self.pa != "/front/tabmenu.asp")
			$.ajax(
			{
				type     : "POST",
				dataType : "jsonp",
				url      : "https://www.thebell.co.kr/log/log.asp",
				data     : {"wb":self.wb, "os":self.os, "pr":self.pr, "dr":self.dr, "pa":self.pa, "ua":self.ua, "ui":self.ui, "co":self.co, "al":self.al}, 
				error   : function (data, status, err)
				{
					try
					{
						var img = new Image(1, 1);
						img.src = "https://www.thebell.co.kr/log/log.asp?wb="+self.wb+"&os="+self.os+"&pr="+self.pr+"&dr="+self.dr+"&pa="+self.pa+"&ua="+escape(self.ua)+"&co="+self.co+"&al="+self.al;
					} catch (e) {alert(e);}
				}
			});
			
			
			// head line
			try{
				if(opener==undefined) self.op=self.dr;			
				else self.op=opener.location.href;
			}catch(ex){}		
				
			if (self.pa == "/front/allheadline/headline.asp")
			$.ajax(
			{
				type     : "POST",
				dataType : "jsonp",
				url      : "https://www.thebell.co.kr/log/headline.asp",
				data     : { "dr":self.dr, "pa":self.pa, "op":self.op, "ui":self.ui, "pr":self.pr}, 
				error   : function (data, status, err)
				{
					try
					{
						var img = new Image(1, 1);
						img.src = "https://www.thebell.co.kr/log/headline.asp?dr="+self.dr+"&op="+self.op+"&ui="+self.ui;
					} catch (e) {alert(e);}
				}
			});	
			
			
			// good morning			

			if ((self.pa == "/front/datainfo/goodmorning_popup_new.asp")||
				(self.pa == "/front/datainfo/goodmorning_new.asp"))						
			$.ajax(
			{
				type     : "POST",
				dataType : "jsonp",
				url      : "https://www.thebell.co.kr/log/good.asp",
				data     : { "dr":self.dr, "pa":self.pa, "op":self.op, "ui":self.ui, "pr":self.pr}, 
				error   : function (data, status, err)
				{
					try
					{
						var img = new Image(1, 1);
						img.src = "https://www.thebell.co.kr/log/headline.asp?dr="+self.dr+"&op="+self.op+"&ui="+self.ui;
					} catch (e) {alert(e);}
				}
			});	

			*/

		}
	};
	thebell.detector.init();
})(window.thebell = window.thebell || {}, jQuery);
