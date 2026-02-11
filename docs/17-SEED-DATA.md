# Seed Data Specification
## Демо-данные для разработки и тестирования

---

## 1. Назначение

Seed data нужен для:
- Разработки (видеть реальные данные при работе над UI)
- Тестирования (проверять фильтры, сортировку, пагинацию)
- Демо (показывать инвесторам / партнёрам)
- Скриншотов (App Store)

---

## 2. Команда запуска

```bash
make seed          # Заполняет dev БД демо-данными
make seed-clean    # Удаляет все seed данные
make seed-courts   # Только корты Астаны (для production)
```

---

## 3. Данные

### Users (20 пользователей)

| # | Имя | Пол | NTRP | Рейтинг | Игр | Район |
|---|-----|-----|------|---------|-----|-------|
| 1 | Алмас Бекмухамедов | M | 3.5 | 1320 | 35 | Есильский |
| 2 | Сара Муканова | F | 3.0 | 1180 | 22 | Алматинский |
| 3 | Марат Тулегенов | M | 2.5 | 1050 | 12 | Сарыаркинский |
| 4 | Дина Касымова | F | 4.0 | 1450 | 45 | Есильский |
| 5 | Ержан Абдыкаримов | M | 3.0 | 1200 | 28 | Байконурский |
| 6 | Айгуль Нурланова | F | 2.5 | 1080 | 8 | Алматинский |
| 7 | Тимур Садвакасов | M | 4.5 | 1620 | 50 | Есильский |
| 8 | Камила Оразова | F | 3.5 | 1300 | 30 | Сарыаркинский |
| 9 | Данияр Жумабеков | M | 3.0 | 1150 | 18 | Нуринский |
| 10 | Мадина Ахметова | F | 2.0 | 980 | 5 | Байконурский |
| 11-20 | (аналогично, разнообразие уровней и районов) | Mix | 1.5-5.0 | 800-1800 | 0-60 | Mix |

**Телефоны для dev:** +77070000001 ... +77070000020
**OTP для всех в dev:** 1234

**Superadmin:** user #1 (Алмас), role=superadmin

### Communities (3 сообщества)

| # | Название | Тип | Участников | Доступ | Verified |
|---|---------|-----|-----------|--------|----------|
| 1 | NTC Astana | club | 12 | closed | true |
| 2 | Astana Tennis League | league | 15 | closed | true |
| 3 | Weekend Tennis | group | 8 | open | false |

**Роли:**
- NTC Astana: user #1 = owner, user #4 = admin, user #7 = moderator
- Astana Tennis League: user #4 = owner, user #1 = admin
- Weekend Tennis: user #5 = owner

### Events (10 ивентов)

| # | Сообщество | Тип | Статус | Участников | Когда |
|---|-----------|-----|--------|-----------|-------|
| 1 | NTC | find_partner | open | 1/2 | Завтра 18:00 |
| 2 | NTC | organized_game | filling | 3/4 | Послезавтра 10:00 |
| 3 | ATL | tournament | full | 8/8 | Суббота 09:00 |
| 4 | ATL | training | open | 2/6 | Воскресенье 11:00 |
| 5 | Weekend | find_partner | completed | 2/2 | Вчера 19:00 |
| 6 | Weekend | organized_game | completed | 4/4 | 3 дня назад |
| 7 | NTC | tournament | in_progress | 4/4 | Сегодня 10:00 |
| 8 | ATL | find_partner | cancelled | 0/2 | — |
| 9 | NTC | organized_game | open | 1/4 | Через 3 дня |
| 10 | Weekend | training | open | 0/8 | Через неделю |

### Matches (15 завершённых матчей)

Между участниками completed ивентов. Результаты разнообразные:
- Singles и doubles
- 2 и 3 сета
- С тай-брейками и без
- Все confirmed (рейтинг уже пересчитан)

### Rating History

30 записей — по 2 на матч (winner + loser). Показывает динамику рейтинга за 3 месяца для графика в профиле.

### Chat (5 чатов)

| # | Тип | Участники | Сообщений |
|---|-----|----------|-----------|
| 1 | community | NTC Astana (все 12) | 20 |
| 2 | community | ATL (все 15) | 15 |
| 3 | event | Event #3 participants | 10 |
| 4 | personal | user #1 ↔ user #4 | 8 |
| 5 | personal | user #2 ↔ user #5 | 5 |

Сообщения: микс из текстовых (приветствия, обсуждение времени, результаты). На русском.

### Posts (8 постов)

| # | Сообщество | Автор | Тип | Likes |
|---|-----------|-------|-----|-------|
| 1 | NTC | user #1 | текст | 5 |
| 2 | NTC | user #4 | текст + фото | 8 |
| 3 | ATL | user #4 | текст | 3 |
| 4 | ATL | system | match_result | 12 |
| 5 | Weekend | user #5 | текст | 2 |
| 6-8 | mix | mix | mix | 0-6 |

### Badges (assigned)

| User | Badge | Дата |
|------|-------|------|
| user #7 (Тимур) | "Ветеран" (50 игр) | 2 мес. назад |
| user #4 (Дина) | "Победная серия" (5 побед подряд) | 1 мес. назад |
| user #1 (Алмас) | "Первый шаг" (1 игра) | 3 мес. назад |
| user #1 (Алмас) | "На высоте" (рейтинг 1300+) | 1 мес. назад |

### Notifications (20 штук)

Для user #1 (основной тестовый аккаунт):
- 5 непрочитанных
- 15 прочитанных
- Типы: match_result, event_reminder, new_message, badge_earned, community_joined

### Friends

| User | Friends |
|------|---------|
| user #1 | user #4, user #5, user #7 |
| user #2 | user #5, user #8 |
| user #4 | user #1, user #7, user #8 |

### Courts (12 кортов Астаны — для production тоже)

| # | Название | Тип | Координаты | Адрес |
|---|---------|-----|-----------|-------|
| 1 | NTC Astana (National Tennis Centre) | крытый | 51.1280, 71.4306 | пр. Кабанбай батыра, 42 |
| 2 | Tennis Park Astana | открытый | 51.1350, 71.4200 | ул. Сыганак, 60 |
| 3 | Pro Tennis Club | крытый | 51.1450, 71.4700 | пр. Мәңгілік Ел, 29 |
| 4 | Mega Silk Way Tennis | крытый | 51.0900, 71.4100 | ТРЦ Mega, корп. 3 |
| 5 | Central Stadium Courts | открытый | 51.1280, 71.4320 | стадион Астана Арена |
| 6 | Достык Tennis | открытый | 51.1500, 71.4500 | парк Достық |
| 7 | EXPO Tennis Courts | крытый | 51.0870, 71.4170 | территория EXPO |
| 8 | Keruen Tennis | крытый | 51.1330, 71.4280 | ТРЦ Керуен, -1 этаж |
| 9 | Botanic Garden Courts | открытый | 51.1200, 71.3900 | Ботанический сад |
| 10 | Left Bank Sport Complex | крытый | 51.1400, 71.4600 | Левый берег, спорткомплекс |
| 11 | Saryarka Tennis | открытый | 51.1600, 71.4300 | район Сарыарка |
| 12 | Bayterek Park Courts | открытый | 51.1282, 71.4307 | набережная Ишим |

**Примечание:** Координаты приблизительные. Уточнить перед production deploy.

---

## 4. Реализация

```go
// scripts/seed/main.go

func main() {
    cfg := config.Load()
    db := database.Connect(cfg.DatabaseURL)

    log.Println("Seeding users...")
    users := seedUsers(db, 20)

    log.Println("Seeding communities...")
    communities := seedCommunities(db, users)

    log.Println("Seeding events...")
    events := seedEvents(db, communities, users)

    log.Println("Seeding matches...")
    seedMatches(db, events, users)

    log.Println("Seeding chats...")
    seedChats(db, communities, events, users)

    log.Println("Seeding posts...")
    seedPosts(db, communities, users)

    log.Println("Seeding courts...")
    seedCourts(db)

    log.Println("Seeding badges...")
    seedBadges(db, users)

    log.Println("Seeding notifications...")
    seedNotifications(db, users[0]) // for main test user

    log.Println("Done! Seed complete.")
}
```
