#!/bin/sh
set -e
sleep 10
# until psql $POSTGRES_URL -c '\l'; do
#   >&2 echo "Postgres is unavailable - sleeping"
#   sleep 1
# done

# >&2 echo "Postgres is up - continuing"

if [ "$DJANGO_MANAGEPY_MIGRATE" = 'on' ]; then
    /venv/bin/python manage.py migrate --noinput
fi

if [ "$DJANGO_MANAGEPY_COLLECTSTATIC" = 'on' ]; then
    /venv/bin/python manage.py collectstatic --noinput
fi

exec "$@"
