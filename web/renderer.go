package main

import (
	"fmt"
	"time"
	"html"
)

func RenderHeader() string {
	return `<html>
    	<head>
    		<meta http-equiv="content-type" content="text/html; charset=UTF-8">
    		<link rel="alternate" type="application/rss+xml" title="RSS FEED AFFE" href="/rss.xml">
    		<title>fefemama.org - THE BEST BLOG IN THE UNIVERSE WRITTEN IN Go :]</title>

    	<script type="text/javascript">
    	  var _gaq = _gaq || [];
	      _gaq.push(['_setAccount', 'UA-705689-1']);
    	  _gaq.push(['_trackPageview']);
    	  (function() {
    	    var ga = document.createElement('script'); ga.type = 'text/javascript'; ga.async = true;
    		ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';
    		var s = document.getElementsByTagName('script')[0]; s.parentNode.insertBefore(ga, s);
    	  })();
		</script>
    	</head>

    	<body>
    		<h1><a href="/">Fefemama</a></h1>
    		<b>I love the smell of nopslides in the morning ...</b><br>
				<p align=right>Fragen? <a href="/faq.html">Antworten!</a></p>
    	 `
}

func RenderFooter() string {
	return `
    <div align=right>Proudly made without PHP, Java, Perl, MySQL and Postgres
    </div>
    </body>
    </html>
    `
}

func RenderPost(post *BlogPost, withComments bool) string {
    post.Comments, _ = Db.GetComments(post.Id)
    
	s := "<li>"
	s += fmt.Sprintf("<a href='/post?id=%d'>[%d]</a> ", post.Id, len(post.Comments))
	//s += strings.Replace(post.Content, "\n", "<br>", -1)
	s += post.Content
	s += "</li>"

	if withComments {
		s += "Comments:<ul>"
		for _, comment := range post.Comments {
			//comment := html.EscapeString(comment.Content)
			comment := comment.Content
			comment = strings.Replace(comment, "<", "(", -1)
			comment = strings.Replace(comment, ">", ")", -1)
		
			//author := html.EscapeString(comment.Author)
			author := comment.Author
			author = strings.Replace(author, "<", "(", -1)
			author = strings.Replace(author, ">", ")", -1)
			
			s += fmt.Sprintf("<li>[%s] %s</li>", author, comment)
		}
		s += "</ul>"
		s += "<p><a href='/comment.html'>Willst du einen Kommentar hinterlassen?</a></p>"
	}

	return s
}

func RenderPosts(posts *[]BlogPost) string {
	s := ""

	var cur_date time.Time
	for _, p := range *posts {
		post_date := time.SecondsToLocalTime(p.Timestamp)

		if !(cur_date.Day == post_date.Day && cur_date.Month == post_date.Month && cur_date.Year == post_date.Year) {
			cur_date = *post_date
			if len(s) > 0 {
				s += "</ul>"
			}
			s += "<h3>"
			s += cur_date.Format("Mon Jan _2 2006")
			s += "</h3><ul>"
		}

		s += RenderPost(&p, false)
		//s += "<br>"
	}
	s += "</ul>"

	return s
}
