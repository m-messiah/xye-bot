# coding=utf-8
import sys
# local
sys.path.insert(1, '/usr/local/google_appengine')
sys.path.insert(1, '/usr/local/google_appengine/lib/yaml/lib')
# travis
sys.path.insert(1, 'google_appengine')
sys.path.insert(1, 'google_appengine/lib/yaml/lib')

sys.path.insert(1, './lib')
from unittest import TestCase
import webapp2
from webapp2_extras import json
from main import app, huify


class TestApp(TestCase):
    def test_show_error(self):
        request = webapp2.Request.blank("/")
        response = request.get_response(app)
        self.assertEqual(response.status_int, 200)
        self.assertIn("application/json", response.headers['Content-Type'])
        self.assertDictEqual(
            json.decode(response.body),
            {"name": 'Hello World! I am XyE_bot (https://telegram.me/xye_bot)',
             "result": "Info"})

    def test_get(self):
        request = webapp2.Request.blank("/")
        response = request.get_response(app)
        self.assertEqual(response.status_int, 200)
        self.assertIn("application/json", response.headers['Content-Type'])
        self.assertDictEqual(
            json.decode(response.body),
            {"name": 'Hello World! I am XyE_bot (https://telegram.me/xye_bot)',
             "result": "Info"})

    def test_bad_post(self):
        request = webapp2.Request.blank("/")
        request.method = "POST"
        response = request.get_response(app)
        self.assertEqual(response.status_int, 200)
        self.assertIn("application/json", response.headers['Content-Type'])
        self.assertDictEqual(
            json.decode(response.body),
            {"name": 'Hello World! I am XyE_bot (https://telegram.me/xye_bot)',
             "result": "Info"})

    def test_json_empty_post(self):
        request = webapp2.Request.blank("/")
        request.method = "POST"
        request.headers["Content-Type"] = "application/json"
        response = request.get_response(app)
        self.assertEqual(response.status_int, 200)
        self.assertIn("application/json", response.headers['Content-Type'])
        self.assertDictEqual(
            json.decode(response.body),
            {"name": 'Hello World! I am XyE_bot (https://telegram.me/xye_bot)',
             "result": "Info"})

    def test_json_start_post(self):
        request = webapp2.Request.blank("/")
        request.method = "POST"
        request.headers["Content-Type"] = "application/json"
        request.body = json.encode({
            'update': 1,
            'message': {
                u'date': 1450696897,
                u'text': u'/start',
                u'from': {
                    u'username': u'm_messiah',
                    u'first_name': u'Maxim',
                    u'last_name': u'Muzafarov',
                    u'id': 3798371
                },
                u'message_id': 1,
                u'chat': {
                    u'type': u'user',
                    u'id': 3798371,
                    u'username': u'm_messiah',
                    u'first_name': u'Maxim',
                    u'last_name': u'Muzafarov',
                }
            }
        })
        response = request.get_response(app)
        self.assertEqual(response.status_int, 200)
        self.assertIn("application/json", response.headers['Content-Type'])
        self.assertDictEqual(
                json.decode(response.body),
                {
                    'method': 'sendMessage',
                    'text': u"Привет! Я бот-хуебот.\n"
                            u"Я буду хуифицировать "
                            u"некоторые из твоих фраз",
                    'chat_id': 3798371,
                }
        )

    def test_huify_rus(self):
        self.assertEqual(huify(u'привет'), u"хуивет")
        self.assertEqual(huify(u'привет бот'), u"хуебот")
        self.assertEqual(huify(u'доброе утро'), u"хуютро")
        self.assertEqual(huify(u'ты пьяный'), u"хуяный")
        self.assertEqual(huify(u'были'), u"хуили")
        self.assertEqual(huify(u'китайцы'), u"хуитайцы")

    def test_huify_huified(self):
        self.assertEqual(huify(u'хуитайцы'), None)
        self.assertEqual(huify(u'хуютро'), None)
        self.assertEqual(huify(u'хутор'), u"хуютор")

    def test_huify_non_rus(self):
        self.assertEqual(huify(u'hello'), None)
        self.assertEqual(huify(u'123'), None)

    def test_huify_dash(self):
        self.assertEqual(huify(u'когда-то'), u"хуегда-то")
        self.assertEqual(huify(u'semi-drive'), None)
        self.assertEqual(huify(u'шалтай-болтай'), u"хуялтай-болтай")
        self.assertEqual(huify(u'https://www.edx.org/by-sec-li-mitx-3'), None)

    def test_huify_url(self):
        self.assertEqual(huify(u'сайт.рф'), u"хуяйтрф")
        self.assertEqual(huify(u'http://www.ru'), None)
