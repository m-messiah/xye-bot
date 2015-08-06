# coding=utf-8
from os import environ
import urllib
from flask import Flask, request, jsonify
from google.appengine.api import urlfetch
import re
app = Flask(__name__)
app.config['DEBUG'] = True

try:
    from bot_token import BOT_TOKEN
except ImportError:
    BOT_TOKEN = environ["TOKEN"]

URL = "https://api.telegram.org/bot%s/" % BOT_TOKEN
MyURL = "https://xye-bot.appspot.com"

NON_LETTERS = re.compile(ur'[^а-яё \-]+', flags=re.UNICODE)
PREFIX = re.compile(u"^[бвгджзйклмнпрстфхцчшщ]+", flags=re.UNICODE)


def error():
    return 'Hello World! I am XyE_bot (https://telegram.me/xye_bot)'


def send_reply(response, chat_id):
    response = {
        'chat_id': chat_id,
        'text': response.encode("utf8")
    }
    app.logger.debug("SENT\t%s", response)
    payload = urllib.urlencode(response)
    if response['text'] == '':
        return
    o = urlfetch.fetch(URL + "sendMessage",
                       payload=payload,
                       method=urlfetch.POST)
    app.logger.debug(str(o.content))


@app.route('/', methods=['POST', 'GET'])
def index():
    if request.method == 'GET':
        return error()
    else:
        if 'Content-Type' not in request.headers:
            return error()
        if request.headers['Content-Type'] != 'application/json':
            return error()
        try:
            update = request.json
            message = update['message']
            chat = message['chat']
            text = message.get('text')
            if text:
                app.logger.debug("MESSAGE FROM\t%s",
                                 chat['username'] if 'username' in chat
                                 else chat['id'])
                if text == "/start" or text == "/help":
                    send_reply(u"Привет! Я бот-хуеплёт.\n"
                               u"Я буду хуифицировать некоторые из твоих фраз",
                               chat["id"])
                else:
                    response = huify(text)
                    if response:
                        send_reply(response, chat['id'])
            return jsonify(result="OK", text="Accepted")
        except Exception as e:
            app.logger.warning(str(e))
            return jsonify(result="Fail", text=str(e))


def huify(text):
    vowels = {u'о', u'е', u'а', u'я', u'у', u'ю'}
    rules = {u'о': u'е', u'а': u'я', u'у': u'ю'}
    words = text.split()
    if len(words) > 3:
        return None
    word = NON_LETTERS.sub(u"", words[-1].lower())
    postfix = PREFIX.sub(u"",  word)
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
