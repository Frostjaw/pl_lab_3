# pl_lab_3
Programming languages lab_3
Программа запускается со входными параметрами : имя_программы port - для сервера, имя_программы ip:port количество_подключений - для клиента. Где port - номер порта (от 0 до 65535), ip:port - ip(при работе на одной машине сервера и клиента - 127.0.0.1) и порт сервера для подключения (любые валидные ip и порты), количество_подключений - количество соединений клиента.
Клиент и сервер осуществляют передачу и валидацию ключей из 10 шагов (на последнем шаге новый ключ не отправляется). В консоль выводится текущий ключ, полученный ключ, статус проверки и отправленный ключ. Для читабельности логов запуск клиентов происходит без горутины.
