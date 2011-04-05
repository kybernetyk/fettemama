#!/usr/bin/env python
''' 
    blog.py:
        mongodb backed 'blog'.
    (c) Leon Szpilewski 2011
        http://nntp.pl
    License: GPL v3
'''

import pymongo
import datetime
import textwrap

class Blog:
    def __init__(self):
        self.conn = pymongo.Connection()
        self.db = self.conn.blog
        self.posts = self.db.posts
        self.comments = self.db.comments
        self.admin = self.db.admin
        self.state = 'idle'
    #
    def __del__(self):
        self.conn.end_request()

    def get_all_news(self, max_num):
        cursor = self.posts.find()
        cursor.sort('id', pymongo.DESCENDING)
        cursor.limit(max_num)
        return cursor
    #
    def get_news_for_day(self, date):
        today = datetime.datetime(day = date.day, month = date.month, year = date.year)
        day = datetime.timedelta(days=1)
        tomorrow = today + day
        qry = {'date': {"$gte" : today,
                        "$lt" : tomorrow}}
        cursor = self.posts.find(qry)
        cursor.sort('id', pymongo.DESCENDING)
        return cursor

    def get_news_item(self, item_id):
        cursor = self.posts.find({'id':item_id})
        if cursor.count() == 0:
            return None
        return cursor[0]
    #

    def get_comments_for_item_id(self, item_id):
        cursor = self.comments.find({'post_id' : int(item_id)})
        cursor.sort('id', pymongo.DESCENDING)
        return cursor

    def get_free_post_id(self):
        cursor = self.posts.find()
        cursor.sort('id', pymongo.DESCENDING)
        cursor.limit(1)
        id = 0
        if cursor.count() != 0:
            x = cursor[0]
            x.setdefault('id', 0)
            id = int(x['id'])
        return id+1

    def get_free_comment_id(self):
        cursor = self.comments.find()
        cursor.sort('id', pymongo.DESCENDING)
        cursor.limit(1)
        id = 0
        if cursor.count() != 0:
            x = cursor[0]
            x.setdefault('id', 0)
            id = int(x['id'])
        return id+1

    def check_password(self, testpwd):
        cursor = self.admin.find({'admin_pass':testpwd})
        if cursor.count() == 1:
            return True
        return False

    def new_post(self, commandline):
        password = commandline[1]
        if not self.check_password(password):
            return 'wrong password.\n'
        self.state = 'posting'
        self.current_post = []
        return "enter post. enter $end to end input and save post.\n" + 8 * "0123456789" + "\n"

    def append_post(self, content):
        if content.strip() == '$end':
            id = self.get_free_post_id()
            post = ""
            for line in self.current_post:
                line_prefix = ""
                wrap_width = 70

                #a prepending tab/> indicates indented quotes
                if line[0] == '\t' or line[0] == '>':
                    line = line[1:]     #i'd kill for char* here
                    line_prefix = "\t> "
                    wrap_width = 60
                
                for l in textwrap.wrap(line, wrap_width):
                    post += line_prefix
                    post += l
                    post += '\n'
            self.state = 'idle'
            d = {'id' : id,
                 'content' : post.strip(),
                 'date' : datetime.datetime.now()}
            self.posts.insert(d)
            return "added new post with id " + str(id) + "\n"
        self.current_post.append(content)
        return ''

    def add_comment(self, commandline):
        comment = commandline[3:]
        comment = ' '.join(comment)
        post_id = commandline[1]
        id = self.get_free_comment_id()
        author = commandline[2]
        d = {'post_id':int(post_id),
             'content':comment,
             'author':author,
             'id' : id}
        self.comments.insert(d)
        return "comment added\n"


    def render_news_excerpt(self, item):
        ret = 'Post #' + str(item['id']) + '\n'
        ret += '\t' + item['excerpt'] + '\n\n'
        return ret

    def render_news_item(self, item_id, show_comments = False):
        item = self.get_news_item(int(item_id))
        if item == None:
            ret = "couldn't open post!\n"
            return ret
        date = ""
        if 'date' in item:
            date = ", " + str(item['date'])
        ret = 'Post #' + str(item['id']) + date + '\n'
        for line in item['content'].splitlines(True):
            ret += "\t"+line
        ret += '\n'
        if show_comments:
            ret += "\n" + self.render_comments(item_id)
        else:
            comments = self.get_comments_for_item_id(item['id'])
            ret += "\ncomments: " + str(comments.count())  + "\n"
        return ret

    def render_comments(self, item_id):
        comments = self.get_comments_for_item_id(item_id)
        if comments.count() == 0:
            return "no comments for this posts\n"
        ret = "comments for post # " + str(item_id) + ":\n"
        for comment in comments:
            comment.setdefault('id', 0)
            ret += '\t*[' + comment['author'] + '] ' + comment['content'] + '\n'
        return ret

    def render_h_delimiter(self):
        return 80 * "-" + "\n\n"

    def render_latest_news(self, howmany):
        howmany = int(howmany)
        all_news = self.get_all_news(howmany)
        ret_string = "showing latest %s posts ...\n" % (howmany)
        ret_string += len(ret_string) * "=" + "\n\n"
        for news_item in all_news:
            ret_string += self.render_news_item(news_item['id'])
            ret_string += self.render_h_delimiter()
            howmany -= 1
            if howmany == 0:
                break
        return ret_string

    def render_todays_news(self):
        today = datetime.datetime.now()
        news = self.get_news_for_day(today)
        ret_string = "showing %s posts for today ...\n" % (str(news.count()))
        ret_string += len(ret_string) * "=" + "\n\n"
        for news_item in news:
            ret_string += self.render_news_item(news_item['id'])
            ret_string += self.render_h_delimiter()
        return ret_string

    def render_help(self):
        ret_string = 'fettemama.org v0.1\nimplemented commands:\n\n';
        ret_string += 'help\t- this list\n'
        ret_string += 'news [num]\t- shows num latest posts. num default = 5\n'
        ret_string += "today\t- shows today's posts\n"
        ret_string += 'read <post id>\t-shows the complete post with the given id\n'
        ret_string += 'comments <post id>\t- reads comments for a post\n'
        ret_string += 'comment <post id> <your_name> <comment>\t- adds a comment\n'
        ret_string += "rst\t- close connection\n"
        ret_string += 'version\t- shows version information\n'
        return ret_string

    def render_version(self):
        ret_string = 'fettemama.org blog system version v0.1\n\t(c) don vito 2011\n\twritten in python 2.6\n\tuses mongodb for data storage\n\n'
        return ret_string

    def process_input(self, input):
        tmp = input.split(' ')
				itms = []
				for itm in tmp:
					i = itm.strip()
					itms.append(i)

        if len(itms) == 0:
            return 'continue', "\n"

        if self.state == 'posting':
            return 'continue', self.append_post(input);

        if itms[0] == 'news':
            num = 5
            if len(itms) == 2:
                num = itms[1]
            return 'continue', self.render_latest_news(num)

        if itms[0] == 'today':
            return 'continue', self.render_todays_news()

        if itms[0] == 'help':
            return 'continue', self.render_help()

        if itms[0] == 'read':
            if len(itms) != 2:
                return 'continue', "syntax: read <post_id>\n"
            return 'continue', self.render_news_item((itms[1]), True)

        if itms[0] == 'comments':
            if len(itms) != 2:
                return 'continue', "syntax: comments <post_id>\n"
            return 'continue', self.render_comments(itms[1]);

        if itms[0] == 'comment':
            if len(itms) < 4:
                return  'continue', 'syntax: comment <post_id> <your_nick> <comment>\n'
            return 'continue', self.add_comment(itms)

        if itms[0] == 'post':
            if len(itms) < 2:
                return 'continue', 'syntax: post <admin_password>\n'
            return 'continue', self.new_post(itms)

        if itms[0] == 'version':
            return 'continue', self.render_version()
        if itms[0] == 'rst':
            return 'close', "bye\n"
        return 'continue', "error. command not recognized!\n"

    def render_prompt(self):
        if self.state == 'posting':
            return "input > "
        return "\n#: "

if __name__ == "__main__":
    blog = Blog()
    print blog.process_input("news")
    print blog.render_prompt()
#
