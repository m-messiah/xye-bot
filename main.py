# coding=utf-8
from os import environ
import urllib
from flask import Flask, request, jsonify
import re
from random import randint
app = Flask(__name__)
app.config['DEBUG'] = False

try:
    from bot_token import BOT_TOKEN
except ImportError:
    BOT_TOKEN = environ["TOKEN"]

URL = "https://api.telegram.org/bot%s/" % BOT_TOKEN
MyURL = "https://xye-bot.appspot.com"

NON_LETTERS = re.compile(ur'[^а-яё \-]+', flags=re.UNICODE)
PREFIX = re.compile(u"^[бвгджзйклмнпрстфхцчшщ]+", flags=re.UNICODE)

DELAY = {}

def error():
    return 'Hello World! I am XyE_bot (https://telegram.me/xye_bot)'


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
                    return jsonify(
                        method="sendMessage",
                        chat_id=chat['id'],
                        text=u"Привет! Я бот-хуебот.\n"
                             u"Я буду хуифицировать "
                             u"некоторые из твоих фраз"
                    )
                else:
                    if chat['id'] not in DELAY:
                        DELAY[chat['id']] = randint(0, 4)
                    else:
                        DELAY[chat['id']] -= 1
                    if DELAY[chat['id']] == 0:
                        response = huify(text)
                        del DELAY[chat['id']]
                        if response:
                            return jsonify(
                                method="sendMessage",
                                chat_id=chat['id'],
                                text=response
                            )
                    
            return jsonify(result="SKIP", text="SKIP")
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
