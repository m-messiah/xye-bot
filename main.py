# coding=utf-8
import logging
import webapp2
from webapp2_extras import json
from re import compile, UNICODE
from random import randint

NON_LETTERS = compile(ur'[^а-яё \-]+', flags=UNICODE)
PREFIX = compile(u"^[бвгджзйклмнпрстфхцчшщьъ]+", flags=UNICODE)

DELAY = {}


class MainPage(webapp2.RequestHandler):
    def show_error(self):
        self.response.headers['Content-Type'] = 'application/json'
        self.response.write(json.encode({
            'result': "Info",
            "name": 'Hello World! I am XyE_bot (https://telegram.me/xye_bot)'
        }))

    def get(self):
        return self.show_error()

    def post(self):
        if 'Content-Type' not in self.request.headers:
            return self.show_error()
        if 'application/json' not in self.request.headers['Content-Type']:
            return self.show_error()
        try:
            update = json.decode(self.request.body)
        except Exception:
            return self.show_error()
        response = {'text': None}
        if 'message' in update:
            message = update['message']
            if 'chat' in message and 'text' in message:
                text = message['text']
                chat_id = message['chat']['id']
                logging.debug(message)
                if "/start" in text or "/help" in text:
                    response = {
                        'method': "sendMessage",
                        'chat_id': chat_id,
                        'text': u"Привет! Я бот-хуебот.\n"
                                u"Я буду хуифицировать "
                                u"некоторые из твоих фраз"}
                else:
                    if chat_id not in DELAY:
                        DELAY[chat_id] = randint(0, 4)
                    else:
                        DELAY[chat_id] -= 1
                    if DELAY[chat_id] == 0:
                        del DELAY[chat_id]
                        response = {
                            'method': "sendMessage",
                            'chat_id': chat_id,
                            'text': huify(text)}
        self.response.headers['Content-Type'] = 'application/json'
        self.response.write(json.encode(response if response['text'] else {}))


def huify(text):
    vowels = {u'о', u'е', u'а', u'я', u'у', u'ю'}
    rules = {u'о': u'е', u'а': u'я', u'у': u'ю'}
    words = text.split()
    if len(words) > 3:
        return None
    word = NON_LETTERS.sub(u"", words[-1].lower())
    postfix = PREFIX.sub(u"",  word)
    if word == u"бот":
        return u'хуебот'
    if len(postfix) < 3:
        return None

    if postfix[0] in rules:
        if postfix[1] not in vowels:
            return u"ху%s%s" % (rules[postfix[0]], postfix[1:])
        else:
            if postfix[1] in rules:
                return u"ху%s%s" % (rules[postfix[1]], postfix[2:])
            else:
                return u'ху%s' % postfix[1:]
    else:
        return u"ху%s" % postfix

app = webapp2.WSGIApplication([('/', MainPage)])

if __name__ == '__main__':
    app.run()
