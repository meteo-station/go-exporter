#!/bin/bash
# Отправляет тестовые метрики в mosquitto
# Использование: ./scripts/send_meteo.sh [device_id] [host] [port]
#
# Пример: ./scripts/send_meteo.sh esp-01
# Пример: ./scripts/send_meteo.sh esp-01 192.168.1.10 1883

DEVICE_ID="${1:-01}"
HOST="${2:-localhost}"
PORT="${3:-1883}"

TOPIC="esp-meteo-station/${DEVICE_ID}/data"
PAYLOAD='{"bme688_t":26.89217,"bme688_p":749.6011,"bme688_h":65.80383,"bme688_eco2":1339.564,"bme688_evo2_acc":3,"bme688_evoc":3.355276,"bme688_evoc_acc":3,"bme688_gas_perc":50.68876,"bme688_gas_perc_acc":3,"bme688_iaq":187.701,"bme688_iaq_acc":3,"bme688_iaq_stat":133.9564,"bme688_iaq_stat_acc":3,"bme688_stab_stat":1,"bme688_run_in_stat":1}'

echo "→ topic:   ${TOPIC}"
echo "→ payload: ${PAYLOAD}"
echo "→ broker:  ${HOST}:${PORT}"

mosquitto_pub -h "${HOST}" -p "${PORT}" -t "${TOPIC}" -m "${PAYLOAD}"
