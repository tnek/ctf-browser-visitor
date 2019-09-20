#!/usr/bin/python
import os
import asyncio
import threading
import json
import logging

from quart import Quart, request
from hypercorn.asyncio import serve

import xssbot

logging.basicConfig(level=logging.DEBUG)


app = Quart(__name__)
REQUIRED_FIELDS = {"url"}


@app.route("/visit", methods=["GET", "POST"])
async def visit():
    form = request.args
    job = form.get("job", None)
    if not job:
        return "{}"

    config = json.loads(job)
    logging.info(config)
    if not all(field in config for field in REQUIRED_FIELDS):
        return '{"status":"fail"}'

    await xssbot.queue_job(config)
    return '{"status":"ok"}'


if __name__ == "__main__":
    app.run("0.0.0.0", 5000)
