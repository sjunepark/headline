(function(){
	
	window.premium = 
	{
		version : "0.9",
		name    : "프리미엄 서비스"	
	};
	
	var is_premium = false;
	
	premium.page = 
	{
		init : function(bool_premium)
		{
			if (bool_premium == undefined)
			{
				is_premium = false;
			}
			else
			{
				is_premium = bool_premium;
			}
		},
		
		bind : function(scope)
		{
			if (is_premium == false)
			{
				return;	
			}
			
			if (scope == undefined)
			{
				scope = "";
			}

			$(scope + " .premium").bind("click", function()
			{
				var data_id    = $(this).attr("data-id");
				var data_kind  = $(this).attr("data-kind");
				var data_form  = $(this).attr("data-form");
				
				var params 	  = (data_form != undefined && data_form != "") ? $("form[name='" + data_form + "']").serialize() : "";
				
				var doc_win;
				
				

				if (typeof(parent.opener) == "object")
					try
					{
						doc_win = parent.opener.document.location;
					}
					catch (e)
					{
						doc_win = self.location;
					}
				else
					doc_win = self.location;
				
				if (data_kind == "leaguetable")				// 리그테이블
				{
					doc_win.href = "/datainfo/LTList.asp?Code=0902&com=" + data_id + "&" + params;
				}
				else if (data_kind == "company")	// 발행사
				{
					doc_win.href = "/datainfo/CompanyInfo.asp?Code=0902&com=" + data_id;				
				}
				else if (data_kind == "agent")		// 주관사
				{
					doc_win.href = "/datainfo/AgentInfo.asp?Code=0902&com=" + data_id;										
				}
				else if (data_kind == "deal")
				{		
					doc_win.href = "/datainfo/DealInfo.asp?Code=0902&seq=" + data_id;	
				}
				else if (data_kind == "djmain")
				{
					doc_win.href = "/front/djib/djibMain.asp?Code=0902&code=08";
				}
				else if (data_kind == "djview")
				{
					doc_win.href = "/front/newsview.asp?key=" + data_id;
				}
				else if (data_kind == "newslist")
				{
					doc_win.href = "/front/newslist.asp?code=" + data_id;
				}
				else if (data_kind == "bp")
				{
					doc_win.href = "/datainfo/bp.asp?Code=0902";
				}
				else if (data_kind == "bp_detail")
				{
					doc_win.href = "/datainfo/BPList.asp?Code=0902&period=" + data_id;
				}
				else if (data_kind == "djlist")
				{
					if (data_id == "0801")
						doc_win.href = "/front/djib/djibBChoice.asp";
					else if (data_id == "0802")
						doc_win.href = "/front/djib/djibCommentary.asp";
					else if (data_id == "0803")
						doc_win.href = "/front/djib/djibTopNews.asp";
					else if (data_id == "0804")
					{
						ibwirePop();	
					}
				}
			}).css("cursor", "pointer").each(function()
			{
				var data_color = $(this).attr("data-color");
				var data_kind  = $(this).attr("data-kind");
				
				if (data_kind != "djmain" && data_kind != "djview")
				{
					$(this).css("color", (data_color == undefined || data_color == "") ? "#586CFB" : data_color);
				}
			}).show();
		}
	};
	
	
	premium.info =
	{
		init : function()
		{
			$(".premium_content")
			.bind("mouseover", function()
			{
				$(this).attr("src", "https://image.thebell.co.kr/images/premium/btn_pmcont_01_on.png");
			})
			.bind("mouseout", function()
			{
				$(this).attr("src", "https://image.thebell.co.kr/images/premium/btn_pmcont_01.png");
			})
			.bind("click", function()
			{
				var from = $(this).attr("data-from");
				if (from == undefined)
					from = "";
				var premium_win = window.open("/premium/premium_info.asp?from="+from, "premiumwin", "width=837, height=600, scrollbars=yes");
				premium_win.focus();
			})
			.css("cursor", "pointer")
			.show();

			$(".premium_content_new")
			.bind("mouseover", function()
			{
				$(this).attr("src", "https://image.thebell.co.kr/images/premium/btn_pmcont_01_on.png");
			})
			.bind("mouseout", function()
			{
				$(this).attr("src", "https://image.thebell.co.kr/images/premium/btn_pmcont_01.png");
			})
			.bind("click", function()
			{
				var from = $(this).attr("data-from");
				if (from == undefined)
					from = "";
				var premium_win = window.open("/premium/premium_info.asp?from=deal", "premiumwin", "width=837, height=600, scrollbars=yes");
				premium_win.focus();
			})
			.css("cursor", "pointer")
			.show();
		}
	}
})();
