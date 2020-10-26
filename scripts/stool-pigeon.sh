#!/bin/bash

# MAIN_LOCATION=$GRANDMA_MAIN_LOCATION
# ADDITIONAL_LOCATION=$GRANDMA_ADDITIONAL_LOCATION
# USER_NAME=$GRANDMA_USER_NAME
MAIN_LOCATION=$1
ADDITIONAL_LOCATION=$2
USER_NAME=$3
DEVICE=$4

EXTENSION="jpg"

if [ -z "$MAIN_LOCATION" ]; then
  echo 'Error: please set MAIN_LOCATION environment variable'
  exit 1
fi

if [ -z "$USER_NAME" ]; then
  echo 'Error: please set USER_NAME environment variable'
  exit 1
fi

if [ -z "$ADDITIONAL_LOCATION" ]; then
  echo 'Error: please set ADDITIONAL_LOCATION environment variable'
  exit 1
fi

NOW_TIME="$(date +%s)"
if [ -z "$NOW_TIME" ]; then
  echo "Error: now time check fail"
  exit 2
fi

TIMEZONE="$(date +%z)"
TIMEZONE=${TIMEZONE: 2:1}
TIMEZONE=$(($TIMEZONE*60*60))

NOW_TIME=$(($NOW_TIME+$TIMEZONE))

PHOTO_NAME=/tmp/${USER_NAME}_${MAIN_LOCATION}_${ADDITIONAL_LOCATION}_${NOW_TIME}.${EXTENSION}

FSWEB_OUT="$(fswebcam -r 3264x2448 --jpeg 100 -d ${DEVICE} ${PHOTO_NAME} 2>&1 | grep Writing)"
if [ -z "$FSWEB_OUT" ]; then
  echo "Error: Capturing frame fail"
  exit 2
fi

PIC_SEND_STRING='pic=@'${PHOTO_NAME}
CURL_OUT="$(curl -F ${PIC_SEND_STRING} https://huezer.xyz/upload 2>&1 | grep Successfully)"
if [ -z "$CURL_OUT" ]; then
  echo "Error: sending fail"
  exit 3
fi

rm ${PHOTO_NAME}
