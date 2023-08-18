# kursRates
web сервис который по запросу забирает данные из публичного API национального банка и сохраняет данные в локальную датабазу TEST

# Описание
Здесь я использовал фрэймворк Gorilla, потому что он упрощает извелечение параметров из URL и использовать их в дальнейших операциях. А также был использован Go MySQL Driver потому что он имеет большую аудиторию пользователей, позволяет эффективно работать с базами данных.

# Endpoints
GET /currency/save/{date} - получает данные с API Нац. банка с заданной датой и сохраняет их в базе данных Test
GET /currency/{date}/{code} - возвращает определенный курс валют с заданной датой и кодом в формате JSON
GET /currency/{date} - возвращает все курсы валют с заданной датой в формате JSON

