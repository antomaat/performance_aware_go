

def lexer(string):
    tokens = []

    while len(string):
        json_string, string = lex_string(string)
        if json_string is not None:
            tokens.append(json_string)
            continue
        json_number, string = lex_number(string)
        if json_number is not None:
            tokens.append(json_number)
            continue
        if string[0] in " ":
            string = string[1:]
        if string[0] in [',', ':', '[', ']', '{', '}']:
            tokens.append(string[0])
            string = string[1:]
    return tokens

def lex_string(string):
    json_string = ''
    if string[0] == '"':
        string = string[1:]
    else:
        return None, string

    for c in string:
        if c == '"':
            return json_string, string[len(json_string)+1:]
        else:
            json_string += c

def lex_number(string):
    json_number = ''

    number_characters = [str(d) for d in range(0, 10)] + ['-', 'e', '.']

    for c in string:
        if c in number_characters:
            json_number += c
        else:
            break;
    rest = string[len(json_number):]

    if not len(json_number):
        return None, string
    if '.' in json_number:
        return float(json_number), rest

    return int(json_number), rest

def parse_array(tokens):
    json_array = []

    t = tokens[0]
    if t == ']':
        return json_array, tokens[1:]

    while True:
        json, tokens = parse(tokens)
        json_array.append(json)

        t = tokens[0]
        if t == ']':
            return json_array, tokens[1:]
        elif t != ',':
            raise Exception('Expected comma after object in array')
        else:
            tokens = tokens[1:]

def parse_object(tokens):
    json_object = {}

    t = tokens[0]
    if t == '}':
        return json_object, tokens[1:]

    while True:
        json_key = tokens[0]
        if type(json_key) is str:
            tokens = tokens[1:]
        else:
            raise Exception('Expected string key, got: {}'.format(json_key))

        if tokens[0] != ':':
            raise Exception('Expected colon after key in object, got: {}'.format(t))

        json_value, tokens = parse(tokens[1:])

        json_object[json_key] = json_value

        t = tokens[0]
        if t == '}':
            return json_object, tokens[1:]
        elif t != ',':
            raise Exception('Expected comma after pair in object, got: {}'.format(t))

        tokens = tokens[1:]

def parse(tokens):
    t = tokens[0]
    if t == '[':
        return parse_array(tokens[1:])
    if t == '{':
        return parse_object(tokens[1:])
    else:
        return t, tokens[1:]

def call():
    tokens = lexer('{"pairs": [{"x0": 2}]}')

    parser = parse(tokens)[0]
    print(parser['pairs'][0]['x0'])

call()
    
