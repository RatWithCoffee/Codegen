# Кодогенератор
Ищет методы структуры в файле api.go с меткой apigen (например, apigen:api {"url": "/user/create", "auth": true, "method": "POST"}) и генерирует для них http-обертки, валидацию параметров, заполнение структуры параметрами метода.

У полей структур есть метки apivalidator, по которым генерируем код проверяет поля:
* `required` - поле не должно быть пустым (не должно иметь значение по-умолчанию)
* `paramname` - если указано, то берется из параметра с этим именем, иначе - `lowercase` от имени
* `enum` - "одно из"
* `default` - если указано и приходит пустое значение (значение по-умолчанию) - устанавливается то, что написано
  в `default`
* `min` - >= X для типа `int`, для строк `len(str)` >=
* `max` - <= X для типа `int`

Также в сгенерированном коде есть подобие авторизации (для методов с "auth": true), которое просто проверяет в загловке наличие фиксиравнного значения
