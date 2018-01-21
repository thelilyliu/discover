#!/bin/sh
curl -X POST -F "images_file=@moment.jpg" \
                "https://gateway-a.watsonplatform.net/visual-recognition/api/v3/classify?api_key={8d7aced8efa9ce11cca985d203dce5989cc20148}&version=2016-05-20" -o ./response.json