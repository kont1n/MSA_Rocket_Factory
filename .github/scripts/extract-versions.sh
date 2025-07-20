#!/bin/bash

# Скрипт для извлечения версий инструментов из Taskfile.yml
# Использование: .github/scripts/extract-versions.sh

# Путь к Taskfile.yml
TASKFILE="Taskfile.yml"

# Проверка наличия файла
if [ ! -f "$TASKFILE" ]; then
  echo "Ошибка: Файл $TASKFILE не найден" >&2
  exit 1
fi

# Извлечение всех переменных из секции vars
echo "Извлекаем переменные из Taskfile.yml:"

# Определяем начало и конец секции vars
VARS_START=$(grep -n "^vars:" "$TASKFILE" | cut -d: -f1)
if [ -z "$VARS_START" ]; then
  echo "Ошибка: секция vars не найдена в $TASKFILE" >&2
  exit 1
fi

VARS_START=$((VARS_START + 1))

# Ищем следующую секцию после vars или конец файла
NEXT_SECTION=$(tail -n +$VARS_START "$TASKFILE" | grep -n "^[a-z]" | head -1 | cut -d: -f1)
if [ -n "$NEXT_SECTION" ]; then
  VARS_END=$((VARS_START + NEXT_SECTION - 2))
else
  VARS_END=$(wc -l < "$TASKFILE")
fi

# Извлекаем все строки из секции vars
VARS_SECTION=$(sed -n "${VARS_START},${VARS_END}p" "$TASKFILE")

# Инициализируем ассоциативный массив для хранения переменных
declare -A VARS

# Извлекаем имя и значение каждой переменной
while IFS= read -r line; do
  # Пропускаем пустые строки и строки с комментариями
  if [[ "$line" =~ ^[[:space:]]*$ || "$line" =~ ^[[:space:]]*# ]]; then
    continue
  fi

  # Извлекаем имя и значение
  if [[ "$line" =~ ^[[:space:]]*([A-Z_0-9]+):\ *\'([^\']*)\' ]]; then
    var_name=${BASH_REMATCH[1]}
    var_value=${BASH_REMATCH[2]}
    VARS["$var_name"]="$var_value"
    echo "- $var_name: ${VARS[$var_name]}"
  elif [[ "$line" =~ ^[[:space:]]*([A-Z_0-9]+):\ *\"([^\"]*)\" ]]; then
    var_name=${BASH_REMATCH[1]}
    var_value=${BASH_REMATCH[2]}
    VARS["$var_name"]="$var_value"
    echo "- $var_name: ${VARS[$var_name]}"
  elif [[ "$line" =~ ^[[:space:]]*([A-Z_0-9]+):\ *(.*) ]]; then
    var_name=${BASH_REMATCH[1]}
    var_value=${BASH_REMATCH[2]}
    VARS["$var_name"]="$var_value"
    echo "- $var_name: ${VARS[$var_name]}"
  fi
done <<< "$VARS_SECTION"

# Находим список модулей
if [ -n "${VARS[MODULES]}" ]; then
  MODULES="${VARS[MODULES]}"
  echo "- найдены модули: $MODULES"
else
  # Если не найдено в vars, пытаемся найти в другом месте (для обратной совместимости)
  MODULES=$(sed -n 's/.*MODULES: \(.*\)/\1/p' "$TASKFILE" | head -1)
  echo "- модули (из старого формата): $MODULES"
fi

# Установка переменных GitHub Actions
if [ -n "$GITHUB_ENV" ]; then
  echo "Устанавливаем переменные в GITHUB_ENV:"
  # Экспортируем все переменные
  for var_name in "${!VARS[@]}"; do
    echo "$var_name=${VARS[$var_name]}" >> $GITHUB_ENV
    echo "  $var_name -> GITHUB_ENV"
  done
  # Для совместимости добавляем MODULES отдельно, если оно не в vars
  if [ -z "${VARS[MODULES]}" ] && [ -n "$MODULES" ]; then
    echo "MODULES=$MODULES" >> $GITHUB_ENV
    echo "  MODULES -> GITHUB_ENV"
  fi
fi

if [ -n "$GITHUB_OUTPUT" ]; then
  echo "Устанавливаем переменные в GITHUB_OUTPUT:"
  # Экспортируем все переменные
  for var_name in "${!VARS[@]}"; do
    echo "$var_name=${VARS[$var_name]}" >> $GITHUB_OUTPUT
    echo "  $var_name -> GITHUB_OUTPUT"
  done
  # Для совместимости добавляем MODULES отдельно, если оно не в vars
  if [ -z "${VARS[MODULES]}" ] && [ -n "$MODULES" ]; then
    echo "MODULES=$MODULES" >> $GITHUB_OUTPUT
    echo "  MODULES -> GITHUB_OUTPUT"
  fi
fi