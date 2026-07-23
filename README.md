# Как зайти на сервер?

## Лаунчер
Подойдёт **ЛЮБОЙ** лаунчер. Лично от себя могу порекоммендовать [Prism Launcher](https://prismlauncher.org/download/). Если по неведомой причине вы не можете подключиться к аккаунту Mojang с купленной лицензией, её проверку можно отключить в конфиге. Для этого нужно изменить файл `PrismLauncher/accounts.json` вручную или через команду:
> Для Windows: (Терминал можно открыть, нажав Win+R и введя 'cmd')
> ```bat
> echo {"accounts": [{"entitlement": {"canPlayMinecraft": true,"ownsMinecraft": true},"type": "MSA"}],"formatVersion": 3} > %appdata%/PrismLauncher/accounts.json
> ```
> Для Linux:
> ```sh
> echo '{"accounts": [{"entitlement": {"canPlayMinecraft": true,"ownsMinecraft": true},"type": "MSA"}],"formatVersion": 3}' > ~/.local/share/PrismLauncher/accounts.json
> ```
> Для Linux (если Prism установлен через Flatpak):
> ```sh
> echo '{"accounts": [{"entitlement": {"canPlayMinecraft": true,"ownsMinecraft": true},"type": "MSA"}],"formatVersion": 3}' > ~/.var/app/org.prismlauncher.PrismLauncher/data/PrismLauncher/accounts.json
> ```
> Для macOS:
> ```sh
> echo '{"accounts": [{"entitlement": {"canPlayMinecraft": true,"ownsMinecraft": true},"type": "MSA"}],"formatVersion": 3}' > ~/.var/app/org.prismlauncher.PrismLauncher/data/PrismLauncher/accounts.json
> ```

## Версия и моды
Сервер работает на ядре **Fabric** и версии Minecraft **26.2**. Далее будет предоставлен список модов для игры:

### Серверные моды
Установка этих модов на клиент **НЕ ТРЕБУЕТСЯ**, но они работают на сервере.
**Fabric-моды**
* 2032 World Height - увеличивает максимальную высоту строительства до 2032 блоков
* AppleSkin - поддержка функционирования этого же мода для клиента
* Auth - авторизация по паролю
* Chunky - инструмент для генерации чанков
* Invisible Frames - возможность сделать рабку невидимой через Shift+ПКМ
* Skin Restorer - подтягивает скины по нику с серверов Mojang и [Ely.by](https://ely.by/)
* Clumps, Ferrite Core, Lithium - оптимизация производительности сервера
* Fabric API - зависимость для вышеперечисленных модов
**Дата-паки**
* Minecart Improvements - официальный датапак для улучшения вагонеток (у нас их максимум скорости поднят в 4 раза)

### Обязательные моды
Обязательных модов, необходимых для игре на сервере в текущей сборке **НЕТ**.

### Рекоммендуемые моды
Эти моды не обязательны, чтобы зайти на сервер, но улучшают игровой опыт. Вы можете поплнить этот список удобными вам клиентскими модами или исключить ненужные для вашей сборки.
* AppleSkin - улучшение интерфейса здоровья и голода
* Continuity - соединённые текстуры стекла и некоторых других строительных блоков (требуется активировать встроенный текстурпак)
* Iris - поддержка шейдеров
* LambDynamicLights - доработка освещения
* Mod Menu - меню настройки модов
* Xaero's Minimap - миникарта и маркеры места смерти
* Xaero's World Map - полноэкранная карта мира
* Zoomify - зум на клавишу C
* Ferrite Core, Lithium, Sodium - оптимизация производительности клиента
* Fabric API, Fabric Language Kotlin, Placeholder API, YetAnotherConfigLib - зависимости для вышеперечисленных модов

### Установка модов
Моды можно скопировать из приложенного архива в папку модов вашей сборки или вручную можно с сайтов [Modrinth](https://modrinth.com/discover/mods) и [CurseForge](https://www.curseforge.com/minecraft). Там же можно скачать ресурспаки, шейдеры и т.д. В Prism Launcher всё это можно сделать сразу через интерфейс программы.

## Подключение к серверу
К серверу можно подключиться по одному из двух доменов:
* `rwsheld.joinserver.xyz` - прямое подключение с минимальным пингом
* `rwsheld.joinserver.ru` - подключение через прокси для обхода пранков от РКН

Если у вас не работает один из доменов, обязательно попробуйте другой - бывает даже такое что с прокси не работает, а без него всё в порядке.

# Общение и поддержка
Общаться на тематику сервера и сопутствующего события можно в каналах категории `🟩 Minecraft` на Rain World Rus. Если у вас остались вопросы или возникли проблемы, можете обратиться к стаффу сервера:
* @spiderhead71
* @mr_crepp
* @marussell 