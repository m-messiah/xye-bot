import { test, before } from 'node:test';
import assert from 'node:assert/strict';

let fetchHandler;

before(async () => {
    globalThis.addEventListener = (_type, listener) => {
        fetchHandler = listener;
    };
    await import('../src/index.js');
});

function send(body) {
    let pending;
    fetchHandler({
        request: new Request('http://localhost/', { method: 'POST', body }),
        respondWith: (promise) => {
            pending = promise;
        },
    });
    return pending;
}

function post(update) {
    return send(JSON.stringify(update));
}

// message_id 12 is divisible by 1..4, so isAnswerNeeded() always passes.
const chat = { id: 1 };
const messageId = 12;

test('text message gets a huified reply', async () => {
    const res = await post({ message: { message_id: messageId, chat, text: 'программирование' } });
    assert.equal(res.status, 200);
    const body = await res.json();
    assert.match(body.text, /хуограммирование$/);
});

for (const [name, extra] of [
    ['sticker', { sticker: {} }],
    ['photo', { photo: [] }],
    ['voice', { voice: {} }],
    ['new_chat_members', { new_chat_members: [] }],
    ['pinned_message', { pinned_message: {} }],
]) {
    test(`${name} message without text returns 200, not 500`, async () => {
        const res = await post({ message: { message_id: messageId, chat, ...extra } });
        assert.equal(res.status, 200);
        assert.equal(await res.text(), '');
    });
}

test('edited message without text returns 200', async () => {
    const res = await post({ edited_message: { message_id: messageId, chat, sticker: {} } });
    assert.equal(res.status, 200);
});

test('/start command still answers', async () => {
    const res = await post({ message: { message_id: messageId, chat, text: '/start' } });
    assert.equal(res.status, 200);
    assert.equal((await res.json()).method, 'sendMessage');
});

test('malformed JSON still returns 500', async () => {
    const res = await send('not json');
    assert.equal(res.status, 500);
});
