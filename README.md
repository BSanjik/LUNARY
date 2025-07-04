# LUNARY
LUNARY MobileApp

Обязательно к прочтению!!!!!
# 1. Клонируем репозиторий (если ещё не клонировал)
git clone https://github.com/BSanjik/LUNARY.git


# 2. Обновляем локальную ветку main
git checkout main
git pull origin main

# 3. Создаем новую ветку для своей задачи (замени task-name на название задачи)
git checkout -b task-name

# 4. Делаем изменения в коде, добавляем файлы
git add .
git commit -m "Краткое описание изменений"

# 5. Отправляем ветку на сервер
git push -u origin task-name

# 6. Заходим на GitHub и создаем Pull Request из ветки task-name в main

# 7. После одобрения Pull Request, ветка main обновится через GitHub