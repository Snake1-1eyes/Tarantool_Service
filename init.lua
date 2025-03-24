box.cfg({listen="0.0.0.0:3301"})
-- Создаём пользователя для подключения
box.schema.user.create('storage', {password='passw0rd', if_not_exists=true})
-- Даём все-все права
box.schema.user.grant('storage', 'super', nil, nil, {if_not_exists=true})
box.once('init', function()
    s = box.schema.space.create('kv')
    s:format({
        {name = 'key', type = 'string'},
        {name = 'value', type = 'any'}
    })
    s:create_index('primary', {
        type = 'hash',
        parts = {'key'}
    })
end)
require('msgpack').cfg{encode_invalid_as_nil = true}