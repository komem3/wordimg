#!/bin/bash

gcloud functions deploy SlackEmojiGen --runtime go113 --trigger-http --project $1
