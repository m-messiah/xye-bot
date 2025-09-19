// Xye Bot - JavaScript version with wasmer.io Edge compatibility
import { suggestions } from './suggestions.js';

// Constants
const BOT_NAME = 'xye_bot';
const DELAY_LIMIT = 4;
const WORDS_AMOUNT = 1;
const HUIFIED_PREFIX = 'ху';
const VOWELS = 'оеаяуюыі';
const RULES_VALUES = 'еяюиії';
const APPLY_UA_RULES = 'ії';

// Vowel transformation rules
const VOWELS_RULES = { 'о': 'е', 'а': 'я', 'у': 'ю', 'ы': 'и' };
const VOWELS_RULES_UA = { 'о': 'е', 'а': 'я', 'у': 'ю', 'ы': 'и', 'и': 'і', 'і': 'ї' };

// Regular expressions
const PREFIX_TO_SKIP_RE = /^[бвгджзйклмнпрстфхцчшщьъ]+/;
const NON_LETTERS_RE = /[^а-яёії-]+/g;
const ONLY_DASHES_RE = /^-*$/;
const XAXAXA_RE = /^[ах]+$/;

// Telegram API types
class TelegramResponse {
    constructor(chatId, text, method = 'sendMessage', replyToId = null, parseMode = '', disableWebPagePreview = true) {
        this.chat_id = chatId;
        this.text = text;
        this.method = method;
        if (replyToId) this.reply_to_message_id = replyToId;
        if (parseMode) this.parse_mode = parseMode;
        this.disable_web_page_preview = disableWebPagePreview;
    }
}

class RequestInfo {
    constructor(updateMessage) {
        this.updateMessage = updateMessage;
        this.cacheId = updateMessage.chat.id.toString();
        this.lastResponse = null;
    }

    answer(message, parseMode = '') {
        this.lastResponse = new TelegramResponse(this.updateMessage.chat.id, message, 'sendMessage', this.updateMessage.message_id, parseMode);
    }

    isAnswerNeeded() {
        const randDelay = Math.floor(Math.random() * DELAY_LIMIT) + 1;
        return this.updateMessage.message_id % randDelay === 0;
    }

    modify() {
        return huify(this.updateMessage.text, WORDS_AMOUNT);
    }

    handleCommand() {
        const command = this.getCommand();
        return command.handle();
    }

    getCommand() {
        const commandName = getCommandName(this.updateMessage.text);
        switch (commandName) {
            case '/start':
                return new CommandStart(this);
            case '/stop':
                return new CommandStop(this);
            case '/help':
                return new CommandHelp(this);
            default:
                return new CommandNotFound(this);
        }
    }
}

// Command classes
class CommandStart {
    constructor(request) {
        this.request = request;
    }

    handle() {
        this.request.answer('Я вежливый Хуебот\nЯ буду хуифицировать некоторые из ваших фраз', '');
        return null;
    }
}

class CommandStop {
    constructor(request) {
        this.request = request;
    }

    handle() {
        this.request.answer('Я не умею обрабатывать эту команду. Удалите меня из чата, если надоел', '');
        return null;
    }
}

class CommandHelp {
    constructor(request) {
        this.request = request;
    }

    handle() {
        this.request.answer(
            'Бот использует вежливый режим и отвечает на случайные сообщения\n\n' +
            'Если хочется запустить своего бота, код на [GitHub](https://github.com/m-messiah/xye-bot)\n' +
            '[Donate](https://www.paypal.com/donate?hosted_button_id=7KT6MFDSHPXL6)',
            'MarkdownV2'
        );
        return null;
    }
}

class CommandNotFound {
    constructor(request) {
        this.request = request;
    }

    handle() {
        throw new Error('команда не найдена');
    }
}

// Utility functions
function createTelegramResponse(chatId, text, replyToId = null, parseMode = '') {
    return new TelegramResponse(chatId, text, 'sendMessage', replyToId, parseMode);
}

function getCommandName(text) {
    if (!text.startsWith('/')) return '';
    
    const commandName = text.split(' ')[0];
    const commandParts = commandName.split('@');
    
    if (commandParts.length > 1 && commandParts[1] === BOT_NAME) {
        return commandParts[0];
    }
    
    return commandName;
}

function newRequestFromBody(body) {
    try {
        const update = JSON.parse(body);
        let updateMessage = update.message;
        
        if (!updateMessage) {
            if (!update.edited_message) {
                throw new Error('no message in update');
            }
            updateMessage = update.edited_message;
        }
        
        return new RequestInfo(updateMessage);
    } catch (error) {
        throw error;
    }
}

// Stop list functions
function inStopList(word) {
    return isXAXAXA(word);
}

function isXAXAXA(word) {
    return XAXAXA_RE.test(word);
}

// Huification functions
function isUAWord(word) {
    return APPLY_UA_RULES.split('').some(letter => word.includes(letter));
}

function getRules(word) {
    return isUAWord(word) ? VOWELS_RULES_UA : VOWELS_RULES;
}

function isHuifyApplicable(word) {
    // Skip words with only dashes
    if (ONLY_DASHES_RE.test(word)) {
        return { applicable: false };
    }
    
    const postfix = word.replace(PREFIX_TO_SKIP_RE, '');
    
    // Skip already huified words
    if (postfix.length < 6 || (word.startsWith(HUIFIED_PREFIX) && RULES_VALUES.includes(postfix.substring(2, 4)))) {
        return { applicable: false };
    }
    
    // Skip words from stop list
    if (inStopList(word)) {
        return { applicable: false };
    }
    
    return { applicable: true, postfix };
}

function huifyWord(postfix, rules) {
    const firstTwo = postfix.substring(0, 2);
    
    if (rules[firstTwo]) {
        const nextTwo = postfix.substring(2, 4);
        if (!VOWELS.includes(nextTwo)) {
            return HUIFIED_PREFIX + rules[firstTwo] + postfix.substring(2);
        }
        if (rules[nextTwo]) {
            return HUIFIED_PREFIX + rules[nextTwo] + postfix.substring(4);
        }
        return HUIFIED_PREFIX + postfix.substring(2);
    }
    
    return HUIFIED_PREFIX + postfix;
}

function tryHuifyWord(text) {
    const word = text.toLowerCase().replace(NON_LETTERS_RE, '');
    
    // Special cases
    if (word === 'бот') {
        return { result: 'хуебот', huified: true };
    }
    
    if (word === 'путин') {
        return { result: 'хуйло', huified: true };
    }
    
    const { applicable, postfix } = isHuifyApplicable(word);
    if (applicable) {
        return { result: huifyWord(postfix, getRules(word)), huified: true };
    }
    
    return { result: word, huified: false };
}

function tryHuify(text, amount) {
    const words = text.trim().split(/\s+/);
    
    if (words.length < 1 || words.length - amount > 4) {
        return '';
    }
    
    const candidateWords = words.length > amount ? words.slice(-amount) : words;
    const answer = [];
    let isHuified = false;
    
    for (const word of candidateWords) {
        const { result, huified } = tryHuifyWord(word);
        if (result.length > 0) {
            answer.push(result);
        }
        isHuified = isHuified || huified;
    }
    
    return isHuified ? answer.join(' ') : '';
}

function huify(text, amount) {
    const huified = tryHuify(text, amount);
    if (!huified) {
        return '';
    }
    
    const randomSuggestion = suggestions[Math.floor(Math.random() * suggestions.length)];
    return randomSuggestion + huified;
}

// Main handler function for wasmer
async function handler(request) {
    if (request.method !== 'POST') {
        return new Response('OK', { status: 200 });
    }
    
    try {
        const body = await request.text();
        const requestInfo = newRequestFromBody(body);
        
        try {
            const commandResult = requestInfo.handleCommand();
            if (commandResult === null && requestInfo.lastResponse) {
                // Command was handled successfully
                return new Response(JSON.stringify(requestInfo.lastResponse), {
                    status: 200,
                    headers: { 'Content-Type': 'application/json' }
                });
            }
        } catch (error) {
            // Command not found or error, continue to huification
        }
        
        if (!requestInfo.isAnswerNeeded()) {
            return new Response('', { status: 200 });
        }
        
        const output = requestInfo.modify();
        if (!output) {
            return new Response('', { status: 200 });
        }
        
        const telegramResponse = createTelegramResponse(requestInfo.updateMessage.chat.id, output, requestInfo.updateMessage.message_id, '');
        return new Response(JSON.stringify(telegramResponse), {
            status: 200,
            headers: { 'Content-Type': 'application/json' }
        });
    } catch (error) {
        console.error('Error processing request:', error);
        return new Response('Internal Server Error', { status: 500 });
    }
}

addEventListener("fetch", (fetchEvent) => {
    fetchEvent.respondWith(handler(fetchEvent.request));
});
