#!/usr/bin/python
import os
import asyncio
import concurrent.futures
import threading
import json
import logging

from quart import Quart, request

import xssbot
from config import MAX_WORKER_COUNT, REQUIRED_FIELDS

app = Quart(__name__)


@app.before_serving
async def app_init():
    loop = asyncio.get_event_loop()
    await xssbot.xssbot(loop)


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
    app.run("0.0.0.0", port=5000)
