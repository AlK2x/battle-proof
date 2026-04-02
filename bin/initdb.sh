#!/bin/bash

set -euo pipefail

sed -e "s|\${DB_USER_APP_PASSWORD}|${DB_USER_APP_PASSWORD}|g" \
    ./database/mysql/init/1_init_user_db.sql.tpl \
    > ./database/mysql/init/1_init_user_db.sql

sed -e "s|\${DB_JAM_APP_PASSWORD}|${DB_JAM_APP_PASSWORD}|g" \
    ./database/mysql/init/2_init_jam_db.sql.tpl \
    > ./database/mysql/init/2_init_jam_db.sql

sed -e "s|\${DB_BATTLE_APP_PASSWORD}|${DB_BATTLE_APP_PASSWORD}|g" \
    ./database/mysql/init/3_init_battle_db.sql.tpl \
    > ./database/mysql/init/3_init_battle_db.sql